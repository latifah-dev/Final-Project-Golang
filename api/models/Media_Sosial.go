package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Media_Sosial struct {
	ID             uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name           string    `gorm:"size:255;not null;unique" json:"name"`
	SosialMediaUrl string    `gorm:"size:255;not null;" json:"sosialmediaurl"`
	UserID         uint32    `gorm:"not null" json:"user_id"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *Media_Sosial) Prepare() {
	m.ID = 0
	m.Name = html.EscapeString(strings.TrimSpace(m.Name))
	m.SosialMediaUrl = html.EscapeString(strings.TrimSpace(m.SosialMediaUrl))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
}

func (m *Media_Sosial) Validate() error {

	if m.Name == "" {
		return errors.New("Required Name")
	}
	if m.SosialMediaUrl == "" {
		return errors.New("Required Sosial Media Url")
	}
	if m.UserID < 1 {
		return errors.New("Required User")
	}
	return nil
}

func (m *Media_Sosial) SaveMediaSosial(db *gorm.DB) (*Media_Sosial, error) {
	var err error
	err = db.Debug().Model(&Media_Sosial{}).Create(&m).Error
	if err != nil {
		return &Media_Sosial{}, err
	}
	if m.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", m.UserID).Take(&m.SosialMediaUrl).Error
		if err != nil {
			return &Media_Sosial{}, err
		}
	}
	return m, nil
}

func (m *Media_Sosial) FindAllMediaSosial(db *gorm.DB) (*[]Media_Sosial, error) {
	var err error
	medsos := []Media_Sosial{}
	err = db.Debug().Model(&Media_Sosial{}).Limit(100).Find(&medsos).Error
	if err != nil {
		return &[]Media_Sosial{}, err
	}
	if len(medsos) > 0 {
		for i, _ := range medsos {
			err := db.Debug().Model(&User{}).Where("id = ?", medsos[i].UserID).Take(&medsos[i].SosialMediaUrl).Error
			if err != nil {
				return &[]Media_Sosial{}, err
			}
		}
	}
	return &medsos, nil
}

func (m *Media_Sosial) FindMediaSosialByID(db *gorm.DB, pid uint64) (*Media_Sosial, error) {
	var err error
	err = db.Debug().Model(&Media_Sosial{}).Where("id = ?", pid).Take(&m).Error
	if err != nil {
		return &Media_Sosial{}, err
	}
	if m.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", m.UserID).Take(&m.SosialMediaUrl).Error
		if err != nil {
			return &Media_Sosial{}, err
		}
	}
	return m, nil
}

func (m *Media_Sosial) UpdateAMediaSosial(db *gorm.DB) (*Media_Sosial, error) {

	var err error

	err = db.Debug().Model(&Photo{}).Where("id = ?", m.ID).Updates(Media_Sosial{Name: m.Name, SosialMediaUrl: m.SosialMediaUrl, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Media_Sosial{}, err
	}
	if m.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", m.UserID).Take(&m.SosialMediaUrl).Error
		if err != nil {
			return &Media_Sosial{}, err
		}
	}
	return m, nil
}

func (m *Media_Sosial) DeleteAMediaSosial(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Media_Sosial{}).Where("id = ? and user_id = ?", pid, uid).Take(&Media_Sosial{}).Delete(&Media_Sosial{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("media sosial not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
