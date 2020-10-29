//+build wireinject

package inject

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/file"
	logger "github.com/ONSdigital/ras-rm-sample/file-uploader/logging"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

var FileProcessor = Inject()

func Inject() file.FileProcessor {
	wire.Build(NewFileProcessor, ConfigSetup, GenContext, NewPubSub)
	return file.FileProcessor{}
}

func ConfigSetup() config.Config {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("GOOGLE_CLOUD_PROJECT", "rm-ras-sandbox")
	viper.SetDefault("PUBSUB_TOPIC", "topic")
	viper.SetDefault("SAMPLE_SERVICE_BASE_URL", "http://localhost:8080")
	config := config.Config{
		Port: viper.GetString("PORT"),
		Pubsub: config.Pubsub{
			ProjectId: viper.GetString("GOOGLE_CLOUD_PROJECT"),
			TopicId:   viper.GetString("PUBSUB_TOPIC"),
		},
		Sample: config.Sample{
			BaseUrl: viper.GetString("SAMPLE_SERVICE_BASE_URL"),
		},
	}
	return config
}

func NewFileProcessor(config config.Config, client *pubsub.Client, ctx context.Context) file.FileProcessor {
	return file.FileProcessor{Config: config, Client: client, Ctx: ctx}
}

func GenContext() context.Context {
	return context.Background()
}

func NewPubSub(config config.Config, ctx context.Context) *pubsub.Client {
	client, err := pubsub.NewClient(ctx, config.Pubsub.ProjectId)
	if err != nil {
		logger.Error("Failed to create pubsub client", zap.Error(err))
	}
	logger.Info("Pubsub client created")
	return client
}
