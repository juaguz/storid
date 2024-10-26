package internal

import (
	"context"

	"github.com/juaguz/storid/internal/accounts/balances"
	"github.com/juaguz/storid/internal/accounts/transactions/importer"
	"github.com/juaguz/storid/internal/accounts/transactions/repositories"
	"github.com/juaguz/storid/internal/platform/config"
	"github.com/juaguz/storid/internal/platform/db"
	"github.com/juaguz/storid/internal/platform/dispatcher"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewApp() fx.Option {
	return fx.Options(
		fx.Provide(
			zap.NewProduction,
			config.LoadConfig,
			func(cfg *config.Config) *db.Config {
				return &db.Config{
					Host:     cfg.DBConfig.Host,
					Port:     cfg.DBConfig.Port,
					User:     cfg.DBConfig.User,
					Password: cfg.DBConfig.Password,
					Database: cfg.DBConfig.Database,
				}
			},
			db.NewDB,
			fx.Annotate(
				balances.NewBalanceRefresher,
				fx.As(new(dispatcher.EventHandler))),
			fx.Annotate(
				repositories.NewTransactionRepository,
				fx.As(new(importer.TransactionRepository)),
			),
			repositories.NewTransactionRepository,
			importer.NewFileImporter,
		),
		fx.Invoke(registerHandlers),
	)
}

func registerHandlers(d dispatcher.EventDispatcher, handler dispatcher.EventHandler) {
	d.Register(context.Background(), importer.EventImported, handler)
}
