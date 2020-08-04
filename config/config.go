package config

type (
	Config struct {
		Port string `yaml:"port"`
		Pubsub Pubsub `yaml:"pubsub"`
		Sample Sample `yaml:"sample"`
	}

	Pubsub struct {
		ProjectId string `yaml:"project_id"`
		TopicId string `yaml:"topic_id"`
	}

	Sample struct {
		BaseUrl string `yaml:"baseUrl"`
	}
)
