package repositories

import (
	"lh-gin/models"
	"lh-gin/tools"
	"log"
)

type UserRepository interface {
	AddNew(user models.User) (int64, error)
	GetInfoByUsername(username string) (models.User, error)
	GetInfoByID(id int) (models.User, error)
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
	lastID, err = tools.NewMysqlInstance().InsertOne(user)
	if err != nil {
		log.Println("插入失败: ", err.Error())
		return 0, err
	}

	return lastID, nil
}

func (receiver UserManagerRepository) GetInfoByUsername(username string) (models.User, error) {
	userModel := models.User{}
	if ok, err := tools.NewMysqlInstance().Where("username=?", username).Get(&userModel); !ok {
		return userModel, err
	}
	return userModel, nil
}

func (receiver UserManagerRepository) GetInfoByID(id int) (models.User, error) {
	userModel := models.User{}
	if ok, err := tools.NewMysqlInstance().Where("id=?", id).Get(&userModel); !ok {
		return userModel, err
	}
	return userModel, nil
}
