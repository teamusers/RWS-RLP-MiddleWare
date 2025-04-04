package svs

import (
	"testing"
	"time"
)

func TestFloat(t *testing.T) {
	mysql, err := PoolsInMysql("F9RafodEquxWgiKBzyYWaqnPqdEuWcMPEVCAgJE1QS3R", "SOLANA")
	if err != nil {
		t.Error(err)
	}
	t.Log(mysql)
}
func TestKLine1min(t *testing.T) {
	mysql, err := KLine1min("6a3m2EgFFKfsFuQtP4LJJXPcAe3TQYXNyHUjjZpUxYgd", "SOLANA", time.Now().Add(-time.Hour*24*1).Format("2006-01-02 15:04:05"))
	if err != nil {
		t.Error(err)
	}
	t.Log(mysql)
}
