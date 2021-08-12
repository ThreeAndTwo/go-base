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
	Client *gorm.DB
	exitCh chan struct{}
}

func New(c *Config) (*MySQL, error) {
	if err := c.check(); err != nil {
		return nil, fmt.Errorf("invalid config:%s", err)
	}

	url := c.user + ":" + c.pass + "@tcp(" + c.host + ")/" + c.db + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(c.maxIdleConn)
	db.DB().SetMaxOpenConns(c.maxOpenConn)
	db.DB().SetConnMaxIdleTime(c.maxIdleTimeConn)
	db.DB().SetConnMaxLifetime(c.maxLifeTimeConn)

	client := &MySQL{
		config: c,
		Client: db,
		exitCh: make(chan struct{}),
	}
	return client, nil
}

func (c *Config) check() error {
	if "" == c.user || "" == c.pass || "" == c.host || "" == c.db {
		return fmt.Errorf("config error")
	}

	if 0 == c.maxIdleConn || 0 == c.maxOpenConn {
		c.maxIdleConn = 80
		c.maxOpenConn = 80
	}

	if c.maxIdleConn != c.maxOpenConn {
		c.maxOpenConn = c.maxIdleConn
	}

	if c.maxLifeTimeConn == time.Duration(0) {
		c.maxLifeTimeConn = time.Duration(300)
	}

	if c.maxIdleTimeConn == time.Duration(0) {
		c.maxIdleTimeConn = time.Duration(300)
	}

	return nil
}

func (mysql *MySQL) GetConfig() *Config {
	return mysql.config
}

func (mysql *MySQL) Close() error {
	close(mysql.exitCh)
	return mysql.Client.Close()
}