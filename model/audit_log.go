// model/audit_log.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UserID       string    // from your auth context
	Method       string
	Path         string
	StatusCode   int
	ClientIP     string
	UserAgent    string
	RequestBody  string `gorm:"type:longtext"` // consider redaction for sensitive fields
	ResponseBody string `gorm:"type:longtext"` // optional
	LatencyMs    int64
}

func MigrateAuditLog(db *gorm.DB) error {
	return db.AutoMigrate(&AuditLog{})
}
