package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Photo struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Caption   string    `gorm:"size:255;not null;" json:"caption"`
	PhotoUrl  string    `gorm:"size:255;not null;" json:"photourl"`
	Author    User      `json:"author"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Photo) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Caption = html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoUrl = html.EscapeString(strings.TrimSpace(p.PhotoUrl))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Photo) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.PhotoUrl == "" {
		return errors.New("Required Photo")
	}
	if p.UserID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Photo) SavePhoto(db *gorm.DB) (*Photo, error) {
	var err error
	err = db.Debug().Model(&Photo{}).Create(&p).Error
	if err != nil {
		return &Photo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.Author).Error
		if err != nil {
			return &Photo{}, err
		}
	}
	return p, nil
}

func (p *Photo) FindAllPhotos(db *gorm.DB) (*[]Photo, error) {
	var err error
	photos := []Photo{}
	err = db.Debug().Model(&Photo{}).Limit(100).Find(&photos).Error
	if err != nil {
		return &[]Photo{}, err
	}
	if len(photos) > 0 {
		for i, _ := range photos {
			err := db.Debug().Model(&User{}).Where("id = ?", photos[i].UserID).Take(&photos[i].Author).Error
			if err != nil {
				return &[]Photo{}, err
			}
		}
	}
	return &photos, nil
}

func (p *Photo) FindPhotoByID(db *gorm.DB, pid uint64) (*Photo, error) {
	var err error
	err = db.Debug().Model(&Photo{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Photo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.Author).Error
		if err != nil {
			return &Photo{}, err
		}
	}
	return p, nil
}

func (p *Photo) UpdateAPhoto(db *gorm.DB) (*Photo, error) {

	var err error

	err = db.Debug().Model(&Photo{}).Where("id = ?", p.ID).Updates(Photo{Title: p.Title, PhotoUrl: p.PhotoUrl, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Photo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.Author).Error
		if err != nil {
			return &Photo{}, err
		}
	}
	return p, nil
}

func (p *Photo) DeleteAPhoto(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Photo{}).Where("id = ? and user_id = ?", pid, uid).Take(&Photo{}).Delete(&Photo{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Photo not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
