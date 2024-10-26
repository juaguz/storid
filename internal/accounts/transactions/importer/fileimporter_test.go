package importer

import (
	"context"
	"sync"
	"testing"

	"github.com/juaguz/storid/internal/accounts/transactions/dto"
	"github.com/juaguz/storid/internal/platform/filereaders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// normally I would use mockgen for this but for the sake of the exercise I will just create a mock
type TransactionRepositoryMock struct {
	m            sync.Mutex
	transactions []*dto.Transaction
	count        int
}

func (t *TransactionRepositoryMock) Create(transaction []*dto.Transaction) error {
	t.m.Lock()
	defer t.m.Unlock()
	t.transactions = append(t.transactions, transaction...)
	t.count = t.count + 1
	return nil
}

type EventDispatcherMock struct {
	Event string
}

func (e *EventDispatcherMock) Dispatch(ctx context.Context, event string, payload []byte) {
	e.Event = event
}

func TestFileImporter_Import(t *testing.T) {
	transactionRepository := &TransactionRepositoryMock{}
	eventDispatcher := &EventDispatcherMock{}

	fi := &FileImporter{
		zap.NewExample(),
		filereaders.NewLocalFileReader(),
		transactionRepository,
		eventDispatcher,
	}

	err := fi.Import(context.Background(), "./fixtures/random_transactions.csv")
	if err != nil {
		t.Errorf("Expected no error %v", err)
	}

	assert.Equal(t, EventImported, eventDispatcher.Event)

	assert.Len(t, transactionRepository.transactions, 100000)
}
