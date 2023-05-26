package postgres

import (
	"fmt"
	"time"
)

// Config DbLife: the unit is minute
type Config struct {
	Host    string `json:"host"     required:"true"`
	User    string `json:"user"     required:"true"`
	Pwd     string `json:"pwd"      required:"true"`
	Name    string `json:"name"     required:"true"`
	Port    int    `json:"port"     required:"true"`
	Life    int    `json:"life"     required:"true"`
	MaxConn int    `json:"max_conn" required:"true"`
	MaxIdle int    `json:"max_idle" required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxConn <= 0 {
		cfg.MaxConn = 1000
	}

	if cfg.MaxIdle <= 0 {
		cfg.MaxIdle = 500
	}

	if cfg.Life <= 0 {
		cfg.Life = 2
	}
}

func (cfg *Config) getLifeDuration() time.Duration {
	return time.Minute * time.Duration(cfg.Life)
}

func (cfg *Config) dsn() string {
	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Pwd, cfg.Name, cfg.Port,
	)
}
