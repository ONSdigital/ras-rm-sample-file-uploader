//+build wireinject

package inject

import (
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/file"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"cloud.google.com/go/pubsub"
	"context"

	log "github.com/sirupsen/logrus"
)

var FileProcessor = Inject()

func Inject() file.FileProcessor{
	wire.Build(NewFileProcessor, ConfigSetup, GenContext, NewPubSub)
	return file.FileProcessor{}
}

func ConfigSetup() config.Config {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("PROJECT_ID", "rm-ras-sandbox")
	viper.SetDefault("TOPIC_ID", "topic")
	config := config.Config{
		Port: viper.GetString("PORT"),
		Pubsub: config.Pubsub{
			ProjectId: viper.GetString("PROJECT_ID"),
			TopicId: viper.GetString("TOPIC_ID"),
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
		log.WithError(err).Error("Failed to create pubsub client")
	}
	log.Info("Pubsub client created")
	return client
}