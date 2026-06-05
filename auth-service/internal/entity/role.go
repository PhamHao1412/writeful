package entity

import "fmt"

type Role struct {
	BaseEntity

	Code string `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name string `json:"name"`
}

func (Role) TableName() string {
	return fmt.Sprintf("%vroles", SchemaName())

}
