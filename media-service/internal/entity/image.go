package entity

type Image struct {
	BaseEntity
	URL        string  `json:"url" gorm:"type:text;not null"`
	Type       string  `json:"type" gorm:"type:varchar(20);not null"`
	Provider   string  `json:"provider" gorm:"type:varchar(20);default:'cloudinary'"`
	PublicID   *string `json:"public_id,omitempty" gorm:"type:text"`
	MimeType   *string `json:"mime_type,omitempty" gorm:"type:text"`
	Format     *string `json:"format,omitempty" gorm:"type:text"`
	FileSize   *int64  `json:"file_size,omitempty" gorm:"type:bigint"`
	Width      *int    `json:"width,omitempty" gorm:"type:int"`
	Height     *int    `json:"height,omitempty" gorm:"type:int"`
	UploadedBy string  `json:"uploaded_by,omitempty" gorm:"type:uuid"`
	FileHash   *string `json:"file_hash,omitempty" gorm:"type:varchar(64);index"`
}

func (i *Image) TableName() string {
	return SchemaName() + "images"
}
