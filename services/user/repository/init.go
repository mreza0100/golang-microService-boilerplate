package repository

import (
	"fmt"
	"microServiceBoilerplate/configs"
	"microServiceBoilerplate/services/user/instances"
	"microServiceBoilerplate/services/user/models"

	"github.com/mreza0100/golog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewRepo(lgr *golog.Core) *instances.Repository {
	var (
		db     *gorm.DB
		readQ  *read
		writeQ *write
	)

	{
		db = getConnection(lgr)
		lgr = lgr.With("In Repository ->")
	}
	{
		readQ = &read{
			lgr: lgr.With("In Read ->"),
			db:  db,
		}
		writeQ = &write{
			lgr: lgr.With("In Read ->"),
			db:  db,
		}
		readQ.write = writeQ
		writeQ.read = readQ
	}

	return &instances.Repository{
		Read:  readQ,
		Write: writeQ,
	}
}

func getConnection(lgr *golog.Core) *gorm.DB {
	var (
		err error
		db  *gorm.DB
	)

	{
		db, err = gorm.Open(getConfigs())
		if err != nil {
			lgr.Fatal(err)
		}
	}
	{
		if err := db.AutoMigrate(&models.User{}); err != nil {
			lgr.Fatal(err)
		}
	}

	lgr.SuccessLog("Connected to DB")

	return db
}

func getConfigs() (driverConfigs gorm.Dialector, gormConfigs *gorm.Config) {
	DSN := fmt.Sprintf("host=localhost user=postgres dbname=postgres port=%v", configs.UserConfigs.DBPort)
	driverConfigs = postgres.New(postgres.Config{
		DSN: DSN,
	})

	gormConfigs = &gorm.Config{
		NamingStrategy:         schema.NamingStrategy{},
		SkipDefaultTransaction: true,
		Logger:                 logger.Default,
		// PrepareStmt:            false,
	}

	return
}
