package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/tifarin/fullstack/api/models"
)

var users = []models.User{
	models.User{
		Username: "Latifah",
		Email:    "latifah@gmail.com",
		Password: "password",
		Age:      23,
	},
	models.User{
		Username: "Arda",
		Email:    "arda@gmail.com",
		Password: "password",
		Age:      17,
	},
}

var photos = []models.Photo{
	models.Photo{
		Title:    "Title 1",
		Caption:  "Hello world 1",
		PhotoUrl: "image.img",
	},
	models.Photo{
		Title:    "Title 2",
		Caption:  "Hello world 2",
		PhotoUrl: "image2.img",
	},
}

var comment = []models.Comment{
	models.Comment{
		Message: "sangat bagus",
	},
	models.Comment{
		Message: "kurang bagus",
	},
}

var mediasosial = []models.Media_Sosial{
	models.Media_Sosial{
		Name:           "latifah",
		SosialMediaUrl: "instagram.com/latifah",
	},
	models.Media_Sosial{
		Name:           "arda",
		SosialMediaUrl: "instagram.com/arda",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Photo{}, &models.User{}, &models.Comment{}, &models.Media_Sosial{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.Media_Sosial{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Photo{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Comment{}).AddForeignKey("photo_id", "photos(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Media_Sosial{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		photos[i].UserID = users[i].ID

		err = db.Debug().Model(&models.Photo{}).Create(&photos[i]).Error
		if err != nil {
			log.Fatalf("cannot seed photos table: %v", err)
		}
		comment[i].UserID = users[i].ID
		comment[i].PhotoID = uint32(photos[i].ID)

		err = db.Debug().Model(&models.Comment{}).Create(&comment[i]).Error
		if err != nil {
			log.Fatalf("cannot seed comment table: %v", err)
		}
		mediasosial[i].UserID = users[i].ID

		err = db.Debug().Model(&models.Media_Sosial{}).Create(&mediasosial[i]).Error
		if err != nil {
			log.Fatalf("cannot seed mediasosial table: %v", err)
		}
	}
}
