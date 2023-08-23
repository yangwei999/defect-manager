package config

import (
	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/defect-manager/defect/infrastructure/repositoryimpl"
	"github.com/opensourceways/defect-manager/issue"
	"github.com/opensourceways/defect-manager/message-server"
)

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type Config struct {
	MessageServer messageserver.Config `json:"message_server" required:"true"`
	Kafka         kafka.Config         `json:"kafka"          required:"true"`
	Issue         issue.Config         `json:"issue"          required:"true"`

	repositoryimpl.Config
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.MessageServer,
		&cfg.Kafka,
		&cfg.Issue,
	}
}

func (cfg *Config) SetDefault() {
	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
