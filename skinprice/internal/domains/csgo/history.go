package csgo

import (
	"time"

	"github.com/google/uuid"
)

type HistoryRecordType string

const (
	HistoryRecordTypeBuy  HistoryRecordType = "buy"
	HistoryRecordTypeSell HistoryRecordType = "sell"
)

type Record struct {
	ID         uuid.UUID
	UserSkinId uuid.UUID
	Platform   Platform
	Price      float64
	FeeAmount  float64
	TotalPrice float64
	Type       HistoryRecordType
	AppliedAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewHistoryRecord(
	UserSkinId uuid.UUID,
	Platform Platform,
	Price float64,
	FeeAmount float64,
	Type HistoryRecordType,
	AppliedAt time.Time,
) *Record {
	now := time.Now()
	r := &Record{
		ID:         uuid.New(),
		UserSkinId: UserSkinId,
		Platform:   Platform,
		Price:      Price,
		FeeAmount:  FeeAmount,
		Type:       Type,
		AppliedAt:  AppliedAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	r.RecalculateTotal()
	return r
}

func (r *Record) RecalculateTotal() {
	r.TotalPrice = r.Price * r.FeeAmount
}
