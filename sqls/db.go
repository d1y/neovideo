package sqls

import (
	"database/sql"
	"time"

	"d1y.io/neovideo/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// copy by https://github.com/mlogclub/simple/blob/master/sqls/db.go

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

func Open(dbConfig config.DbConfig, config *gorm.Config, models ...interface{}) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		}
	}

	if db, err = gorm.Open(sqlite.Open(dbConfig.File), config); err != nil {
		log.Errorf("opens database failed: %s", err.Error())
		return
	}

	if sqlDB, err = db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		sqlDB.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTimeSeconds) * time.Second)
		sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetimeSeconds) * time.Second)
	} else {
		log.Error(err)
	}

	if err = db.AutoMigrate(models...); nil != err {
		log.Errorf("auto migrate tables failed: %s", err.Error())
	}
	return
}

func DB() *gorm.DB {
	return db
}

func RealDb() *sql.DB {
	return sqlDB
}

func Close() {
	if sqlDB == nil {
		return
	}
	if err := sqlDB.Close(); nil != err {
		log.Errorf("Disconnect from database failed: %s", err.Error())
	}
}
