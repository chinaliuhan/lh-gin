package services

import (
	"errors"
	"lh-gin/models"
	"lh-gin/repositories"
	"lh-gin/requests"
	"lh-gin/utils"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

/**
新增用户
*/
func (this *UserService) AddNew(user models.User) (int64, error) {

	//insert db
	lastID, err := repositories.NewUserManagerRepository().AddNew(user)

	return lastID, err
}

/**
获取登录信息
*/
func (this UserService) GetLogin(params requests.LoginRequest) (models.User, error) {

	//get db
	info, err := repositories.NewUserManagerRepository().GetInfoByUsername(params.Username)
	if err != nil {
		return models.User{}, err
	}

	//validate password
	if utils.NewGenerate().GenerateMd5(params.Password) == info.Password {
		return models.User{}, errors.New("密码错误")
	}

	return info, nil
}
