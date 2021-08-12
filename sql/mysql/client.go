package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type Config struct {
	user            string        `mapstructure:"user"`
	pass            string        `mapstructure:"pass"`
	host            string        `mapstructure:"host"`
	db              string        `mapstructure:"db"`
	maxIdleConn     int           `mapstructure:"max_idle_conn"`
	maxOpenConn     int           `mapstructure:"max_open_conn"`
	maxLifeTimeConn time.Duration `mapstructure:"max_lifetime_conn"`
	maxIdleTimeConn time.Duration `mapstructure:"max_idletime_conn"`
}

type MySQL struct {
	config *Config
	client *gorm.DB
	exitCh chan struct{}
}

func New(config *Config) (*MySQL, error) {
	if err := config.check(); err != nil {
		return nil, fmt.Errorf("invalid config:%s", err)
	}

	url := config.user + ":" + config.pass + "@tcp(" + config.host + ")/" + config.db +
		"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(config.maxIdleConn)
	db.DB().SetMaxOpenConns(config.maxOpenConn)
	db.DB().SetConnMaxIdleTime(config.maxIdleTimeConn)
	db.DB().SetConnMaxLifetime(config.maxLifeTimeConn)

	client := &MySQL{
		config: config,
		client: db,
		exitCh: make(chan struct{}),
	}
	return client, nil
}

func (config *Config) check() error {
	if "" == config.user || "" == config.pass || "" == config.host || "" == config.db {
		return fmt.Errorf("config error")
	}

	if 0 == config.maxIdleConn || 0 == config.maxOpenConn {
		config.maxIdleConn = 80
		config.maxOpenConn = 80
	}

	if config.maxIdleConn != config.maxOpenConn {
		config.maxOpenConn = config.maxIdleConn
	}

	if config.maxLifeTimeConn == time.Duration(0) {
		config.maxLifeTimeConn = time.Duration(300)
	}

	if config.maxIdleTimeConn == time.Duration(0) {
		config.maxIdleTimeConn = time.Duration(300)
	}

	return nil
}
