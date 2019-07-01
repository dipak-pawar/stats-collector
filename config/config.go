package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Configuration represent all configurations
type Configuration struct {
	Host     string `default:"localhost" split_words:"true"`
	Port     int    `default:"5430" split_words:"true"`
	User     string `default:"postgres" split_words:"true"`
	Database string `default:"postgres" split_words:"true"`
	Password string `default:"secret" split_words:"true"`
	SSLMode  string `default:"disable" envconfig:"POSTGRESQL_SSLMODE"`
}

func (c *Configuration) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Database,
		c.Password,
		c.SSLMode,
	)
}

func (c *Configuration) SetDatabaseName(name string) {
	c.Database = name
}

var Postgres Configuration

func init() {
	err := envconfig.Process("postgresql", &Postgres)
	if err != nil {
		log.Fatal(err.Error())
	}
}
