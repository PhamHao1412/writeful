package entity

import "github.com/google/uuid"

type Media struct {
	BaseEntity
	URL        string    `json:"url" gorm:"type:text;not null"`
	Type       string    `json:"type" gorm:"type:varchar(20);not null"`
	Provider   string    `json:"provider" gorm:"type:varchar(20);default:'cloudinary'"`
	PublicID   *string   `json:"public_id,omitempty" gorm:"type:text"`
	MimeType   *string   `json:"mime_type,omitempty" gorm:"type:text"`
	Format     *string   `json:"format,omitempty" gorm:"type:text"`
	FileSize   *int64    `json:"file_size,omitempty" gorm:"type:bigint"`
	Width      *int      `json:"width,omitempty" gorm:"type:int"`
	Height     *int      `json:"height,omitempty" gorm:"type:int"`
	UploadedBy uuid.UUID `json:"uploaded_by" gorm:"type:uuid;index"`
}

func (Media) TableName() string {
	return SchemaName() + "media"
}
func (Media) TableNameAlias(alias string) string {
	return SchemaName() + "media " + alias
}
