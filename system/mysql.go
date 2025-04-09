package system

import (
	"fmt"
	"time"

	"rlp-middleware/config"
	log "rlp-middleware/log"
	model "rlp-middleware/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Ewriter struct {
	mlog *logrus.Logger
}

// Implement the Printf function required by gorm's logger
func (m *Ewriter) Printf(format string, v ...interface{}) {
	logstr := fmt.Sprintf(format, v...)
	m.mlog.Info(logstr)
}

func NewWriter() *Ewriter {
	// Use the Logger from the sys package
	return &Ewriter{mlog: log.GetLogger()}
}

var DB *gorm.DB

func init() {
	// Get configuration
	cfg := config.GetConfig()

	if cfg.AllStart == 0 {
		return
	}

	// Customize GORM logger
	newLogger := logger.New(
		NewWriter(), // Use a custom logrus logger
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound (record not found) error
			Colorful:                  false,       // Disable colorful output or Disable color printing
		},
	)

	// Construct MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// Open MySQL connection
	database, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN (Data Source Name)
		DefaultStringSize:         256,   // Default string length
		DisableDatetimePrecision:  true,  // Disable datetime precision, not supported in MySQL versions before 5.6
		DontSupportRenameIndex:    true,  // Use drop and create method when renaming indexes
		DontSupportRenameColumn:   true,  // Use `change` to rename columns, not supported in MySQL versions before 8
		SkipInitializeWithVersion: false, // Auto-configure based on version
	}), &gorm.Config{
		Logger: newLogger, // Use a custom GORM logger
	})

	// Error handling
	if err != nil {
		log.Fatal(err) // Use sys.Logger to log fatal errors
	}

	// Enable debug mode if needed:
	database = database.Debug()

	// Assign the database instance to the global variable DB
	DB = database

	if DB.Migrator().HasTable(&model.SysChannel{}) {
		fmt.Println("Table sys_channel exists.")
	} else {
		fmt.Println("Table sys_channel does not exist!")
	}
}

func GetDb() *gorm.DB {
	return DB
}
