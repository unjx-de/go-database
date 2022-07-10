package database

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
	"time"
)

type MySQL struct {
	Url      string
	User     string
	Password string
	Database string
	ORM      *gorm.DB
}

func (config *MySQL) tryDbConnection() {
	i := 1
	total := 20
	for i <= total {
		ln, err := net.DialTimeout("tcp", config.Url, 1*time.Second)
		if err != nil {
			if i == total {
				logrus.WithField("attempt", i).Fatal("Failed connecting to database")
			}
			logrus.WithField("attempt", i).Warning("Connecting to database")
			time.Sleep(2 * time.Second)
			i++
		} else {
			_ = ln.Close()
			logrus.WithField("attempt", i).Info("Connected to database")
			i = total + 1
		}
	}
}

func (config *MySQL) initializeMySql() {
	var err error
	config.tryDbConnection()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Url,
		config.Database,
	)
	config.ORM, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to open database")
	}
}

func (config *MySQL) initializeSqLite(storageDir string) {
	var err error
	absPath := storageDir + "db.sqlite"
	config.ORM, err = gorm.Open(sqlite.Open(absPath), &gorm.Config{})
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to open database")
	}
}

func (config *MySQL) Initialize(storageDir string) {
	if config.Url == "" {
		config.initializeSqLite(storageDir)
	} else {
		config.initializeMySql()
	}
	logrus.WithField("dialect", config.ORM.Dialector.Name()).Debug("Database initialized")
}
