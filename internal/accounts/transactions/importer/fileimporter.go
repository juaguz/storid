package importer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/juaguz/storid/internal/accounts/transactions/dto"
	"github.com/juaguz/storid/internal/platform/currencies"
	"go.uber.org/zap"
)

const (
	//chunkLimit the amount of lines to import at once
	// this would be capped by the DB limit, ideally the repository should handle this
	chunkLimit = 1000

	//numWorkers the amount of workers to use
	numWorkers = 10

	//EventImported name of the event when the file is done
	EventImported = "Imported"
)

type EventDispatcher interface {
	Dispatch(ctx context.Context, event string, payload []byte)
}

type TransactionRepository interface {
	Create(transaction []*dto.Transaction) error
}

type FileReader interface {
	Open(ctx context.Context, filePath string) (io.ReadCloser, error)
}

type FileImporter struct {
	Logger                *zap.Logger
	FileReader            FileReader
	TransactionRepository TransactionRepository
	Dispatcher            EventDispatcher
}

func NewFileImporter(logger *zap.Logger, fileReader FileReader, transactionRepo TransactionRepository, dispatcher EventDispatcher) *FileImporter {
	return &FileImporter{
		Logger:                logger,
		FileReader:            fileReader,
		TransactionRepository: transactionRepo,
		Dispatcher:            dispatcher,
	}
}

func (fi *FileImporter) Import(ctx context.Context, filePath string) error {
	return fi.processFile(ctx, filePath)
}

// process file can be a standalone function to be used in other places
func (fi *FileImporter) processFile(ctx context.Context, filePath string) error {
	f, err := fi.FileReader.Open(ctx, filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// Leer y descartar la cabecera
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("error reading header: %w", err)
	}

	var wg sync.WaitGroup
	chunkChan := make(chan [][]string, numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				if err := fi.createRecords(chunk); err != nil {
					fmt.Printf("error creating records: %v\n", err)
				}
			}
		}()
	}

	var chunk [][]string
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading record: %w", err)
		}

		chunk = append(chunk, record)
		if len(chunk) == chunkLimit {
			chunkChan <- chunk
			chunk = nil
		}
	}

	if len(chunk) > 0 {
		chunkChan <- chunk
	}

	// Close the channel and wait for all workers to finish
	close(chunkChan)
	wg.Wait()

	fi.Dispatcher.Dispatch(ctx, EventImported, nil)

	return nil
}

func (fi *FileImporter) createRecords(records [][]string) error {
	transactions := make([]*dto.Transaction, 0, len(records))

	for _, record := range records {
		transaction, err := fi.parseRecord(record)
		if err != nil {
			return err
		}
		transactions = append(transactions, transaction)
	}

	return fi.TransactionRepository.Create(transactions)
}

func (fi *FileImporter) parseRecord(record []string) (*dto.Transaction, error) {
	if len(record) != 4 {
		return nil, fmt.Errorf("invalid record: %v", record)
	}

	// Parsear amount
	amount, err := currencies.StringToCents(record[2])
	if err != nil {
		return nil, fmt.Errorf("error parsing amount: %w", err)
	}

	dateStr := record[1]
	currentYear := time.Now().Year()
	date, err := fi.parseDate(currentYear, dateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	id := record[0]
	accountID, err := strconv.Atoi(record[3])
	if err != nil {
		return nil, fmt.Errorf("error parsing account ID: %w", err)
	}

	operationType := dto.Credit
	if amount < 0 {
		operationType = dto.Debit
	}

	transaction := &dto.Transaction{
		ExternalID: id,
		Date:       date,
		Amount:     amount,
		AccountID:  uint(accountID),
		Type:       operationType,
	}
	return transaction, nil
}

func (fi *FileImporter) parseDate(currentYear int, dateStr string) (time.Time, error) {
	fullDateStr := fmt.Sprintf("%d/%s", currentYear, dateStr)
	layout := "2006/1/02"

	date, err := time.Parse(layout, fullDateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %w", err)
	}
	return date, nil
}
