package utils

import (
	"fmt"
	"lbe/config"
	"lbe/model"
	"lbe/system"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GenerateNextRLPUserNumbering() (*model.RLPUserNumbering, error) {
	db := system.GetDb()
	now := time.Now()

	year := int64(now.Year() % 100)
	month := int64(now.Month())
	day := int64(now.Day())

	// Get latest RLP_NO
	var lastEntry model.RLPUserNumbering
	err := db.Order("rlp_no DESC").First(&lastEntry).Error
	var nextRlpNo string
	conf := config.GetConfig()
	if err == gorm.ErrRecordNotFound {
		nextRlpNo = conf.Application.RLPNumberingFormat.RLPNODefault
	} else if err != nil {
		return nil, err
	} else {
		var lastRlpNoInt uint64
		fmt.Sscanf(lastEntry.RLP_NO, "%d", &lastRlpNoInt)
		nextRlpNo = fmt.Sprintf("%011d", lastRlpNoInt+1)
	}

	// Get today's max ending number
	var todayEntry model.RLPUserNumbering
	err = db.Where("year = ? AND month = ? AND day = ?", year, month, day).
		Order("rlp_id_ending_no DESC").
		First(&todayEntry).Error

	var nextEndingNo int
	if err == gorm.ErrRecordNotFound {
		nextEndingNo = 1
	} else if err != nil {
		return nil, err
	} else {
		nextEndingNo = todayEntry.RLPIDEndingNO + 1
	}

	rlpID := fmt.Sprintf("%02d%02d%02d%05d", year, month, day, nextEndingNo)

	newRlp := &model.RLPUserNumbering{
		Year:          year,
		Month:         month,
		Day:           day,
		RLP_ID:        rlpID,
		RLP_NO:        nextRlpNo,
		RLPIDEndingNO: nextEndingNo,
	}

	if err := db.Create(newRlp).Error; err != nil {
		return nil, err
	}

	return newRlp, nil
}

func GenerateNextRLPUserNumberingWithRetry() (*model.RLPUserNumbering, error) {
	var lastErr error
	conf := config.GetConfig()
	maxAttempts := conf.Application.RLPNumberingFormat.MaxAttempts
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		newRlp, err := GenerateNextRLPUserNumbering()
		if err == nil {
			return newRlp, nil
		}

		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") {
			lastErr = err
			continue
		}

		return nil, err
	}

	return nil, fmt.Errorf("failed to generate RLP number after %d attempts: %v", maxAttempts, lastErr)
}
