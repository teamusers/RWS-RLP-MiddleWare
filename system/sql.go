package system

import (
	"fmt"
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

	// Build the DSN for MSSQL.
	//
	// If a static port is provided (non-zero), then use it.
	// Otherwise, if an instance name is provided, use the instance parameter.
	// The "packet size" parameter is URL-encoded as "%20".
	var dsn string
	if cfg.Database.Port != 0 {
		// Use static port mode.
		dsn = fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s&packet%%20size=4096",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)
	} else if cfg.Database.Instance != "" {
		// Use named instance mode.
		dsn = fmt.Sprintf(
			"sqlserver://%s:%s@%s?instance=%s&database=%s&packet%%20size=4096",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Instance,
			cfg.Database.DBName,
		)
	} else {
		// Fallback: use host only.
		dsn = fmt.Sprintf(
			"sqlserver://%s:%s@%s?database=%s&packet%%20size=4096",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.DBName,
		)
	}

	var safeDSN string
	if cfg.Database.Port != 0 {
		safeDSN = fmt.Sprintf("sqlserver://%s:***@%s:%d?database=%s&packet%%20size=4096",
			cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	} else if cfg.Database.Instance != "" {
		safeDSN = fmt.Sprintf("sqlserver://%s:***@%s?instance=%s&database=%s&packet%%20size=4096",
			cfg.Database.User, cfg.Database.Host, cfg.Database.Instance, cfg.Database.DBName)
	} else {
		safeDSN = fmt.Sprintf("sqlserver://%s:***@%s?database=%s&packet%%20size=4096",
			cfg.Database.User, cfg.Database.Host, cfg.Database.DBName)
	}
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

	// Optional: Check if a specific table exists (for example, model.SysChannel).
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
