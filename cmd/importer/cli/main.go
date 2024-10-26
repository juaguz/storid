package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/juaguz/storid/cmd/importer/internal"
	"github.com/juaguz/storid/internal/accounts/transactions/importer"
	"github.com/juaguz/storid/internal/platform/config"
	"github.com/juaguz/storid/internal/platform/dispatcher"
	"github.com/juaguz/storid/internal/platform/filereaders"
	"go.uber.org/fx"
)

type ImportHandler struct {
	importer *importer.FileImporter
	filePath string
}

func NewImportHandler(importer *importer.FileImporter, filePath string) *ImportHandler {
	return &ImportHandler{
		importer: importer,
		filePath: filePath,
	}
}

func determineFileReaderMode(mode string) fx.Option {
	switch mode {
	case "s3":
		return fx.Provide(
			func(cfg *config.Config) *s3.Client {
				return cfg.S3Config.Client
			},
			func(client *s3.Client) importer.FileReader {
				bucket := os.Getenv("S3_BUCKET_NAME")
				return filereaders.NewS3FileReader(client, bucket)
			},
		)
	case "local":
		return fx.Provide(
			func() importer.FileReader {
				return filereaders.NewLocalFileReader()
			},
		)
	default:
		log.Fatalf("Unknown mode: %s", mode)
		return nil
	}
}

func (h *ImportHandler) RunImport(ctx context.Context) {
	if err := h.importer.Import(ctx, h.filePath); err != nil {
		log.Fatalf("Error running import: %v", err)
	}
	fmt.Println("Import completed successfully")
}

func main() {
	filePath := flag.String("file", "", "Path to the file to import")
	mode := flag.String("mode", "local", "Choose file reader mode: s3 or local")
	flag.Parse()

	app := fx.New(
		internal.NewApp(),
		//fx.Provide(
		//	func() *filereaders.LocalFileReader {
		//		return filereaders.NewLocalFileReader()
		//	},
		//	fx.Annotate(
		//		filereaders.NewLocalFileReader,
		//		fx.As(new(importer.FileReader)),
		//	),
		//),
		determineFileReaderMode(*mode),
		fx.Provide(
			//this is a flag to determine if the app is running as a CLI or not
			func() bool {
				return false
			},
			fx.Annotate(
				dispatcher.NewSimpleEventDispatcher,
				fx.As(new(importer.EventDispatcher)),
				fx.As(new(dispatcher.EventDispatcher)),
			),
		),
		fx.Provide(func(importer *importer.FileImporter) *ImportHandler {
			return NewImportHandler(importer, *filePath)
		}),
		fx.Invoke(func(handler *ImportHandler) {
			handler.RunImport(context.Background())
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}

	defer app.Stop(context.Background())
}
