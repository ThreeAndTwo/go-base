package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type Interface interface {
}

type Config struct {
	User            string        `mapstructure:"user"`
	Pass            string        `mapstructure:"pass"`
	Host            string        `mapstructure:"host"`
	Db              string        `mapstructure:"db"`
	MaxIdleConn     int           `mapstructure:"max_idle_conn"`
	MaxOpenConn     int           `mapstructure:"max_open_conn"`
	MaxLifeTimeConn time.Duration `mapstructure:"max_lifetime_conn"`
	MaxIdleTimeConn time.Duration `mapstructure:"max_idletime_conn"`
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

	url := c.User + ":" + c.Pass + "@tcp(" + c.Host + ")/" + c.Db + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConn)
	db.DB().SetMaxOpenConns(c.MaxOpenConn)
	db.DB().SetConnMaxIdleTime(c.MaxIdleTimeConn)
	db.DB().SetConnMaxLifetime(c.MaxLifeTimeConn)

	client := &MySQL{
		config: c,
		Client: db,
		exitCh: make(chan struct{}),
	}
	return client, nil
}

func (c *Config) check() error {
	if "" == c.User || "" == c.Pass || "" == c.Host || "" == c.Db {
		return fmt.Errorf("config error")
	}

	if 0 == c.MaxIdleConn || 0 == c.MaxOpenConn {
		c.MaxIdleConn = 80
		c.MaxOpenConn = 80
	}

	if c.MaxIdleConn != c.MaxOpenConn {
		c.MaxOpenConn = c.MaxIdleConn
	}

	if c.MaxLifeTimeConn == time.Duration(0) {
		c.MaxLifeTimeConn = 600 * time.Second
	}

	if c.MaxIdleTimeConn == time.Duration(0) {
		c.MaxIdleTimeConn = 600 * time.Second
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
