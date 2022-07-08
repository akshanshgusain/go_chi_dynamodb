package routes

import (
	"net/http"
	"time"
)

type Config struct {
	timeout time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Cors(next *http.Handler) {

}

func (c *Config) SetTimeout(timeoutSeconds int) *Config {
	c.timeout = time.Duration(timeoutSeconds) * time.Second
	return c
}

func (c *Config) GetTimeout() time.Duration {
	return c.timeout
}
