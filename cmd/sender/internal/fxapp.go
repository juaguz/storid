package internal

import (
	accountrepositories "github.com/juaguz/storid/internal/accounts/repositories"
	"github.com/juaguz/storid/internal/platform/config"

	"github.com/juaguz/storid/internal/accounts/balances/repositories"
	"github.com/juaguz/storid/internal/accounts/balances/summary"
	"github.com/juaguz/storid/internal/platform/db"
	"github.com/juaguz/storid/internal/platform/notifications"
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
			func(cfg *config.Config) *notifications.SMTPConfig {
				return &notifications.SMTPConfig{
					Host:     cfg.SMTPConfig.Host,
					Port:     cfg.SMTPConfig.Port,
					Username: cfg.SMTPConfig.Username,
					Password: cfg.SMTPConfig.Password,
				}
			},
			db.NewDB,
			fx.Annotate(
				accountrepositories.NewAccountRepository,
				fx.As(new(summary.AccountRepository)),
			),
			fx.Annotate(
				notifications.NewSMTPService,
				fx.As(new(summary.EmailService)),
			),
			fx.Annotate(
				summary.NewEmailSender,
				fx.As(new(summary.Notifier)),
				fx.ResultTags(`group:"notifiers"`),
			),
			fx.Annotate(
				repositories.NewBalancesRepository,
				fx.As(new(summary.SummaryGenerator))),

			fx.Annotate(
				summary.NewSender,
				fx.ParamTags(``, `group:"notifiers"`),
			),
		),
	)
}
