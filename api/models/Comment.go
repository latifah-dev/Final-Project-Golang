package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Message   string    `gorm:"size:255;not null;unique" json:"message"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	PhotoID   uint32    `gorm:"not null" json:"photo_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Comment) Prepare() {
	c.ID = 0
	c.Message = html.EscapeString(strings.TrimSpace(c.Message))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Comment) Validate() error {

	if c.Message == "" {
		return errors.New("Required Message")
	}
	if c.UserID < 1 {
		return errors.New("Required User")
	}
	if c.PhotoID < 1 {
		return errors.New("Required Photo")
	}
	return nil
}

func (c *Comment) SaveComment(db *gorm.DB) (*Comment, error) {
	var err error
	err = db.Debug().Model(&Comment{}).Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.PhotoID).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) FindAllComments(db *gorm.DB) (*[]Comment, error) {
	var err error
	comments := []Comment{}
	err = db.Debug().Model(&Photo{}).Limit(100).Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err := db.Debug().Model(&User{}).Where("id = ?", comments[i].UserID).Take(&comments[i].PhotoID).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
	}
	return &comments, nil
}

func (c *Comment) FindCommentByID(db *gorm.DB, pid uint64) (*Comment, error) {
	var err error
	err = db.Debug().Model(&Comment{}).Where("id = ?", pid).Take(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.PhotoID).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) UpdateAComment(db *gorm.DB) (*Comment, error) {

	var err error

	err = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Updates(Comment{Message: c.Message, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.PhotoID).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) DeleteAComment(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Comment{}).Where("id = ? and user_id = ?", pid, uid).Take(&Comment{}).Delete(&Comment{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Comment not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
