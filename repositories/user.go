package repositories

import (
	"lh-gin/models"
	"lh-gin/utils"
	"log"
)

type UserRepository interface {
	AddNew(user models.User) (int64, error)
	GetInfoByUsername(username string) (models.User, error)
}

type UserManagerRepository struct {
}

func NewUserManagerRepository() UserRepository {
	return &UserManagerRepository{}
}
func (receiver UserManagerRepository) AddNew(user models.User) (int64, error) {
	var (
		err    error
		lastID int64
	)

	// insert db
	lastID, err = utils.NewDBMysql().InsertOne(user)
	if err != nil {
		log.Println("插入失败: ", err.Error())
		return 0, err
	}

	return lastID, nil
}

func (receiver UserManagerRepository) GetInfoByUsername(username string) (models.User, error) {
	userModel := models.User{}
	if ok, err := utils.NewDBMysql().Where("username=?", username).Get(&userModel); !ok {
		return userModel, err
	}
	return userModel, nil
}
