package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type Config struct {
	Addr        string `yaml:"addr"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	DbName      string `yaml:"dbName"`
	Prefix      string `yaml:"prefix"`
	Parameters  string `yaml:"parameters"`
	Debug       bool   `yaml:"debug"`
	Active      int    `yaml:"active"`
	Idle        int    `yaml:"idle"`
	IdleTimeout int    `yaml:"idleTimeout"`
}

//user:password@(addr)/dbname?charset=utf8&parseTime=True&loc=Local
func NewMysql(c *Config) (db *gorm.DB) {
	LogLevel := logger.Silent
	if c.Debug {
		LogLevel = logger.Info
	} else {
		LogLevel = logger.Error
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      LogLevel,
			Colorful:      true,
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", c.User, c.Password, c.Addr, c.DbName, c.Parameters)
	var err error
	config := &gorm.Config{
		Logger:               newLogger,
		DisableAutomaticPing: true,
	}
	if c.Prefix != "" {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   c.Prefix,
			SingularTable: true,
		}
	}
	db, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetConnMaxLifetime(time.Second * time.Duration(c.IdleTimeout))
	sqlDb.SetMaxOpenConns(c.Active)
	sqlDb.SetMaxIdleConns(c.Idle)
	return
}
