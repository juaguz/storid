package accounts

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/docker/docker/api/types/container"
	"github.com/juaguz/storid/internal/accounts/balances"
	"github.com/juaguz/storid/internal/accounts/balances/repositories"
	"github.com/juaguz/storid/internal/accounts/balances/summary"
	"github.com/juaguz/storid/internal/accounts/models"
	accountrepository "github.com/juaguz/storid/internal/accounts/repositories"
	"github.com/juaguz/storid/internal/accounts/transactions/importer"
	transactionrepo "github.com/juaguz/storid/internal/accounts/transactions/repositories"
	"github.com/juaguz/storid/internal/platform/currencies"
	"github.com/juaguz/storid/internal/platform/db"
	"github.com/juaguz/storid/internal/platform/dispatcher"
	"github.com/juaguz/storid/internal/platform/filereaders"
	"github.com/juaguz/storid/internal/platform/months"
	"github.com/juaguz/storid/internal/platform/notifications"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

//go:embed transactions/importer/fixtures/random_transactions.csv
var fileContent string

func setUpS3Container(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack:latest",
		ExposedPorts: []string{"4566/tcp"},
		Env: map[string]string{
			"SERVICES": "s3",
		},
		WaitingFor: wait.ForLog("Ready."),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	return container, nil
}

func setUpMailContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "rnwood/smtp4dev:latest",
		ExposedPorts: []string{"25/tcp", "80/tcp"},
		Env: map[string]string{
			"SMTP_AUTH_USERNAME": "testuser",
			"SMTP_AUTH_PASSWORD": "testpass",
			"SMTP_HELO_HOSTNAME": "localhost",
		},
		WaitingFor: wait.ForListeningPort("25/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return container, nil
}

func setUpPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
	// Creates a new postgresql container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
	}

	//starts the container
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	return postgresContainer, err
}

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	logger := zap.L()
	var postgresContainer, mailContainer, s3Container testcontainers.Container
	var err error

	g, wgCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		postgresContainer, err = setUpPostgresContainer(wgCtx)
		return err
	})

	g.Go(func() error {
		var err error
		s3Container, err = setUpS3Container(wgCtx)
		return err
	})

	g.Go(func() error {
		var err error
		mailContainer, err = setUpMailContainer(wgCtx)
		return err

	})

	if err := g.Wait(); err != nil {
		t.Fatalf("Error setting up containers: %v", err)
	}

	defer s3Container.Terminate(ctx)
	defer postgresContainer.Terminate(ctx)
	defer mailContainer.Terminate(ctx)

	// gets hosts and ports
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Error getting postgress host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Error getting postgres port: %v", err)
	}

	s3Host, err := s3Container.Host(ctx)
	if err != nil {
		t.Fatalf("Error getting s3 host: %v", err)
	}

	s3Port, err := s3Container.MappedPort(ctx, "4566")
	if err != nil {
		t.Fatalf("Error getting s3 port: %v", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", s3Host, s3Port.Port())

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatalf("Error setting AWS config: %v", err)
	}

	mailHost, err := mailContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Error getting mailhog host: %v", err)
	}

	mailPort, err := mailContainer.MappedPort(ctx, "25")
	if err != nil {
		t.Fatalf("Error getting mailhog port: %v", err)
	}

	webPort, _ := mailContainer.MappedPort(ctx, "80")
	fmt.Println("Mailhog web interface: http://localhost:" + webPort.Port())

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	bucketName := "test-bucket"
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("Create bucket: %v", err)

	}

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("random_transactions.csv"),
		Body:   strings.NewReader(fileContent),
	})

	if err != nil {
		t.Fatalf("Error saving file: %v", err)
	}

	gormDb := db.NewDB(&db.Config{
		Host:     host,
		Port:     port.Int(),
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		SSLMode:  "disable",
	}, zap.NewExample())

	// Now I can test the summary
	count := int64(0)
	gormDb.Model(&models.Account{}).Count(&count)

	if count != 20 {
		t.Errorf("Expected 20 accounts, got %d", count)
	}

	eventDispatcher := dispatcher.NewSimpleEventDispatcher(false)

	refresher := balances.NewBalanceRefresher(gormDb, logger)

	eventDispatcher.Register(context.Background(), importer.EventImported, refresher)

	transactionRepository := transactionrepo.NewTransactionRepository(gormDb)

	s3Reader := filereaders.NewS3FileReader(s3Client, bucketName)

	i := importer.NewFileImporter(zap.NewExample(), s3Reader, transactionRepository, eventDispatcher)

	err = i.Import(context.Background(), "random_transactions.csv")
	if err != nil {
		log.Fatalln(err)
	}

	gormDb.Model(&models.Transaction{}).Count(&count)
	assert.Equal(t, int64(100_000), count)

	recordedBalances := []models.Balance{}
	err = gormDb.Find(&recordedBalances).Error
	assert.NoError(t, err)
	assert.Equal(t, 20, len(recordedBalances))

	balanceRepository := repositories.NewBalancesRepository(gormDb)
	balance, err := balanceRepository.GetBalanceByAccountID(ctx, 8)
	assert.NoError(t, err)

	assert.Equal(t, -755_043, balance.TotalBalance)

	assert.Equal(t, "-7550.43", currencies.CentsToString(balance.TotalBalance))

	monthlyBalance, err := balanceRepository.GetMonthlyBalancesByAccountID(ctx, 8)
	assert.NoError(t, err)

	assert.Equal(t, 434, monthlyBalance[months.January].TransactionCount)
	assert.Equal(t, 410, monthlyBalance[months.February].TransactionCount)
	assert.Equal(t, 437, monthlyBalance[months.March].TransactionCount)
	assert.Equal(t, 404, monthlyBalance[months.April].TransactionCount)
	assert.Equal(t, 451, monthlyBalance[months.May].TransactionCount)
	assert.Equal(t, 416, monthlyBalance[months.June].TransactionCount)
	assert.Equal(t, 397, monthlyBalance[months.July].TransactionCount)
	assert.Equal(t, 514, monthlyBalance[months.August].TransactionCount)
	assert.Equal(t, 385, monthlyBalance[months.September].TransactionCount)
	assert.Equal(t, 454, monthlyBalance[months.October].TransactionCount)
	assert.Equal(t, 393, monthlyBalance[months.November].TransactionCount)
	assert.Equal(t, 400, monthlyBalance[months.December].TransactionCount)

	res, err := balanceRepository.GetSummaryBalance(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(res))

	accountRepository := accountrepository.NewAccountRepository(gormDb)

	emailService := notifications.NewSMTPService(&notifications.SMTPConfig{
		Host:     mailHost,
		Port:     mailPort.Port(),
		Username: "testuser",
		Password: "testpass",
	}, zap.NewExample())

	notifiers := []summary.Notifier{
		summary.NewEmailSender(accountRepository, emailService),
	}

	sender := summary.NewSender(balanceRepository, notifiers, zap.NewExample())

	err = sender.Send(context.Background())
	assert.NoError(t, err)
	verifyEmailReceived(t, fmt.Sprintf("http://%s:%s", mailHost, webPort.Port()), "Summary Balance")
}

type Smtp4DevMessage struct {
	CurrentPage    int `json:"currentPage"`
	FirstRowOnPage int `json:"firstRowOnPage"`
	LastRowOnPage  int `json:"lastRowOnPage"`
	PageCount      int `json:"pageCount"`
	PageSize       int `json:"pageSize"`
	Results        []struct {
		AttachmentCount int       `json:"attachmentCount"`
		DeliveredTo     string    `json:"deliveredTo"`
		From            string    `json:"from"`
		Id              string    `json:"id"`
		IsRelayed       bool      `json:"isRelayed"`
		IsUnread        bool      `json:"isUnread"`
		ReceivedDate    time.Time `json:"receivedDate"`
		Subject         string    `json:"subject"`
		To              []string  `json:"to"`
	} `json:"results"`
	RowCount int `json:"rowCount"`
}

// verifyEmailReceived checks if an email with the expected subject was sent
func verifyEmailReceived(t *testing.T, smtp4devAPI string, expectedSubject string) {
	resp, err := http.Get(fmt.Sprintf("%s/api/messages", smtp4devAPI))
	if err != nil {
		t.Fatalf("Failed to connect to smtp4dev API: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response from smtp4dev: %v", err)
	}

	var messages Smtp4DevMessage

	err = json.Unmarshal(data, &messages)
	if err != nil {
		t.Fatalf("Error parsing JSON response: %v", err)
	}

	assert.Equal(t, 20, messages.RowCount)

	// Check if we received an email with the expected subject
	found := false
	for _, message := range messages.Results {
		if message.Subject == expectedSubject {
			found = true
			fmt.Printf("Found email with subject: %s\n", message.Subject)
			break
		}
	}

	assert.True(t, found, "Expected email with subject '%s' was not received", expectedSubject)
}
