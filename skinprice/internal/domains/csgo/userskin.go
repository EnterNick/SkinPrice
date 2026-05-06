package csgo

import (
	"time"

	"github.com/google/uuid"
)

type UserSkinStatus string

const (
	UserSkinStatusOwned UserSkinStatus = "owned"
	UserSkinStatusSold  UserSkinStatus = "sold"
)

type UserSkin struct {
	ID             uuid.UUID
	SkinID         uuid.UUID
	SourcePlatform string
	Status         UserSkinStatus
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	deletedAt      *time.Time
}

func NewUserSkin(
	SkinID uuid.UUID,
	SourcePlatform string,
	notes string,
) *UserSkin {
	now := time.Now()
	return &UserSkin{
		ID:             uuid.New(),
		SkinID:         SkinID,
		SourcePlatform: SourcePlatform,
		Status:         UserSkinStatusOwned,
		Notes:          notes,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (us *UserSkin) Delete() {
	now := time.Now()
	us.deletedAt = &now
}

func (us *UserSkin) Restore() {
	us.deletedAt = nil
}

func (us *UserSkin) Sell() {
	us.Status = UserSkinStatusSold
}
