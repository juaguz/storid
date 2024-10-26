package balances

import (
	"context"
	"fmt"

	"github.com/juaguz/storid/internal/accounts/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var views = make([]string, 0, 2)

type BalanceRefresher struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewBalanceRefresher(DB *gorm.DB, log *zap.Logger) *BalanceRefresher {
	views = append(views, models.MonthlyBalanceTable, models.BalancesTable)
	return &BalanceRefresher{DB: DB, Log: log}
}

func (b *BalanceRefresher) Refresh(ctx context.Context) error {
	b.Log.Info("refreshing views", zap.Int("views", len(views)))
	for _, view := range views {
		b.Log.Info("refreshing view", zap.String("view", view))
		sql := fmt.Sprintf("REFRESH MATERIALIZED VIEW %s", view)
		if err := b.DB.Exec(sql).Error; err != nil {
			b.Log.Error("refreshing views", zap.Error(err), zap.String("view", view))
		}
	}
	return nil
}

func (b *BalanceRefresher) Handle(ctx context.Context, event string, payload []byte) {
	if err := b.Refresh(ctx); err != nil {
		b.Log.Error("error refreshing views", zap.Error(err))
	}
}
