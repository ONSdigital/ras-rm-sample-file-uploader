package config

type (
	Config struct {
		Port string `yaml:"port"`
		Pubsub Pubsub `yaml:"pubsub"`
	}

	Pubsub struct {
		ProjectId string `yaml:"project_id"`
		TopicId string `yaml:"topic_id"`
	}
)
