package entity

import "fmt"

type Follower struct {
	BaseEntity
	FollowerID  string `gorm:"column:follower_id"`
	FollowingID string `gorm:"column:following_id"`
}

func (Follower) TableName() string {
	return fmt.Sprintf("%sfollows", SchemaName())
}
