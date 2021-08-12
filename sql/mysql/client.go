package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Config struct {
	user string
	pass string
	host string
	db   string
}

type MySQL struct {
	config *Config
	client *gorm.DB
	exitCh chan struct{}
}

func New(config *Config) (*MySQL, error) {
	url := config.user + ":" + config.pass + "@tcp(" + config.host + ")/" + config.db +
		"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	client := &MySQL{
		config: config,
		client: db,
		exitCh: make(chan struct{}),
	}
	return client, nil
}

func (config *Config) Check() error {
	if "" == config.user	|| "" == config.pass || "" == config.host || "" == config.db {
		return fmt.Errorf("config error")
	}
	return nil
}