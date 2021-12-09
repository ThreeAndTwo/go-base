package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Interface interface{}

type Config struct {
	URI                     string            `mapstructure:"url"`
	AuthMechanism           string            `mapstructure:"auth_mechanism"`
	AuthMechanismProperties map[string]string `mapstructure:"auth_mechanism_properties"`
	AuthSource              string            `mapstructure:"auth_source"`
	Username                string            `mapstructure:"username"`
	Password                string            `mapstructure:"password"`
	PasswordSet             bool              `mapstructure:"password_set"`
	TimeOut                 time.Duration     `mapstructure:"timeout"`
}

type Mongo struct {
	config *Config
	Client *mongo.Client
	exitCh chan struct{}
}

func New(c *Config) (*mongo.Client, error) {
	if err := c.check(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.TimeOut)
	defer cancel()

	credential := options.Credential{
		AuthMechanism:           c.AuthMechanism,
		AuthMechanismProperties: c.AuthMechanismProperties,
		AuthSource:              c.AuthSource,
		Username:                c.Username,
		Password:                c.Password,
		PasswordSet:             c.PasswordSet,
	}

	clientOpts := options.Client().ApplyURI(c.URI).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpts)
	return client, err
}

func (c *Config) check() error {
	if "" == c.URI || "" == c.Username || "" == c.Password {
		return fmt.Errorf("config error, plz check your config")
	}

	if 0 == c.TimeOut {
		c.TimeOut = 10 * time.Second
	}

	return nil
}

func (m *Mongo) GetConfig() *Config {
	return m.config
}

func (m *Mongo) Close() error {
	close(m.exitCh)
	return m.Client.Disconnect(context.TODO())
}
