package entity

type Video struct {
	BaseEntity
	ID           string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type         string `json:"type"`
	URL          string `json:"url"`
	Format       string `json:"format"`
	Duration     int    `json:"duration"`
	Size         int64  `json:"size"`
	ThumbnailURL string `json:"thumbnail_url" gorm:"column:thumbnail_url"`
}

func (i *Video) TableName() string {
	return SchemaName() + "videos"
}
