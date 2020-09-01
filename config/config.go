package config

type (
	Config struct {
		Port   string `yaml:"port"`
		Pubsub Pubsub `yaml:"pubsub"`
		Sample Sample `yaml:"sample"`
	}

	Pubsub struct {
		ProjectId string `yaml:"google_cloud_project"`
		TopicId   string `yaml:"pubsub_topic"`
	}

	Sample struct {
		BaseUrl string `yaml:"baseUrl"`
	}
)
