// model/audit_log.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	ActorID      string    `gorm:"column:actor_id" json:"actor_id"`
	Method       string    `gorm:"column:method" json:"method"`
	Path         string    `gorm:"column:path" json:"path"`
	StatusCode   int       `gorm:"column:status_code" json:"status_code"`
	ClientIP     string    `gorm:"column:client_ip" json:"client_ip"`
	UserAgent    string    `gorm:"column:user_agent" json:"user_agent"`
	RequestBody  string    `gorm:"type:NVARCHAR(MAX);column:request_body" json:"request_body"`   // consider redaction for sensitive fields
	ResponseBody string    `gorm:"type:NVARCHAR(MAX);column:response_body" json:"response_body"` // optional
	LatencyMs    int64     `gorm:"column:latency_ms" json:"latency_ms"`
}

func MigrateAuditLog(db *gorm.DB) error {
	return db.AutoMigrate(&AuditLog{})
}
