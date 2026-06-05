package db

import (
	"content-service/internal/entity"
	"log"

	"gorm.io/gorm"
)

func AutoMigrateAndSeed(dbConn *gorm.DB) error {
	log.Println("Database AutoMigrate starting...")
	err := dbConn.AutoMigrate(
		&entity.Music{},
		&entity.Story{},
		&entity.StoryView{},
	)
	if err != nil {
		return err
	}
	log.Println("Database AutoMigrate completed successfully.")

	// Insert initial seed music tracks
	seeds := []entity.Music{
		{
			Title:    "Chúng Ta Của Tương Lai",
			Artist:   "Sơn Tùng M-TP",
			URL:      "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3",
			CoverURL: "https://images.unsplash.com/photo-1514525253161-7a46d19cd819?w=150",
			Genre:    "vpop",
		},
		{
			Title:    "Chúng Ta Của Hiện Tại",
			Artist:   "Sơn Tùng M-TP",
			URL:      "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-3.mp3",
			CoverURL: "https://images.unsplash.com/photo-1498038432885-c6f3f1b912ee?w=150",
			Genre:    "vpop",
		},
		{
			Title:    "Muộn Rồi Mà Sao Còn",
			Artist:   "Sơn Tùng M-TP",
			URL:      "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-5.mp3",
			CoverURL: "https://images.unsplash.com/photo-1501386761578-eac5c94b800a?w=150",
			Genre:    "vpop",
		},
		{
			Title:    "Blinding Lights",
			Artist:   "The Weeknd",
			URL:      "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-2.mp3",
			CoverURL: "https://images.unsplash.com/photo-1508700115892-45ecd05ae2ad?w=150",
			Genre:    "pop",
		},
		{
			Title:    "Save Your Tears",
			Artist:   "The Weeknd",
			URL:      "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-4.mp3",
			CoverURL: "https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=150",
			Genre:    "pop",
		},
		{
			Title:    "Starboy",
			Artist:   "The Weeknd ft. Daft Punk",
			URL:      "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-6.mp3",
			CoverURL: "https://images.unsplash.com/photo-1518609878373-06d740f60d8b?w=150",
			Genre:    "pop",
		},
		{
			Title:    "Sunset Lofi Study",
			Artist:   "Lofi Dreamer",
			URL:      "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-8.mp3",
			CoverURL: "https://images.unsplash.com/photo-1516450360452-9312f5e86fc7?w=150",
			Genre:    "lofi",
		},
		{
			Title:    "Sunny Days Acoustic",
			Artist:   "Guitar Forest",
			URL:      "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-10.mp3",
			CoverURL: "https://images.unsplash.com/photo-1520523839897-bd0b52f945a0?w=150",
			Genre:    "acoustic",
		},
	}

	for _, s := range seeds {
		var count int64
		dbConn.Model(&entity.Music{}).Where("title = ? AND artist = ?", s.Title, s.Artist).Count(&count)
		if count == 0 {
			if err := dbConn.Create(&s).Error; err != nil {
				log.Printf("failed to insert seed music %s: %v", s.Title, err)
			} else {
				log.Printf("Inserted seed music track: %s - %s", s.Title, s.Artist)
			}
		}
	}
	return nil
}
