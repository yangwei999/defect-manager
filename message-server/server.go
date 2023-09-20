package messageserver

import (
	kafka "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/defect-manager/issue"
)

func Init(cfg *Config, handler issue.EventHandler) error {
	s := messageServer{
		handler: giteeEventHandler{
			handler:   handler,
			userAgent: cfg.UserAgent,
		},
	}

	return s.subscribe(cfg)
}

type messageServer struct {
	handler giteeEventHandler
}

func (m *messageServer) subscribe(cfg *Config) error {
	return kafka.Subscribe(cfg.GroupName, m.handler.handle, []string{cfg.Topics.DefectEvent})
}
