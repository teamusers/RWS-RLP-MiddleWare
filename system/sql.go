package system

import (
	"fmt"
	"net/url"
	"time"

	"lbe/config"
	"lbe/log"
	"lbe/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Ewriter implements the Printf method required by GORM's logger.
type Ewriter struct {
	mlog *logrus.Logger
}

// Printf formats the log message and writes it using Logrus.
func (m *Ewriter) Printf(format string, v ...interface{}) {
	logstr := fmt.Sprintf(format, v...)
	m.mlog.Info(logstr)
}

// NewWriter returns a new instance of Ewriter using the project's logger.
func NewWriter() *Ewriter {
	return &Ewriter{mlog: log.GetLogger()}
}

// DB is the global MSSQL database instance.
var DB *gorm.DB

func init() {
	// Fetch configuration from YAML.
	cfg := config.GetConfig()

	// Check if the application should start.
	if cfg.AllStart == 0 {
		return
	}

	// Log loaded DB config for debugging.
	log.Infof("DB cfg → host=%q, port=%d, user=%q", cfg.Database.Host, cfg.Database.Port, cfg.Database.User)

	// Customize GORM logger.
	newLogger := logger.New(
		NewWriter(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Build the DSN using net/url to handle special characters.
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(cfg.Database.User, cfg.Database.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.Port),
	}
	params := u.Query()
	params.Set("database", cfg.Database.DBName)
	params.Set("packet size", "4096")
	u.RawQuery = params.Encode()
	dsn := u.String()

	// Create a masked DSN for logs (hide password).
	safeURL := *u
	safeURL.User = url.UserPassword(cfg.Database.User, "***")
	safeDSN := safeURL.String()
	log.Info("Using DSN: " + safeDSN)

	// Open a connection to MSSQL using GORM.
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Assign the connection to the global variable.
	DB = db

	// Optional: Check if a specific table exists.
	if db.Migrator().HasTable(&model.SysChannel{}) {
		fmt.Println("Table sys_channel exists in MSSQL.")
	} else {
		fmt.Println("Table sys_channel does not exist in MSSQL!")
	}
}

// GetDb returns the global MSSQL database instance.
func GetDb() *gorm.DB {
	return DB
}
