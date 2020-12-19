package services

import (
	"lh-gin/constants"
	"lh-gin/models"
	"lh-gin/repositories"
	"lh-gin/requests"
	"lh-gin/tools"
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
func (this UserService) GetLogin(params requests.LoginRequest) (models.User, int) {

	//get db
	info, err := repositories.NewUserManagerRepository().GetInfoByUsername(params.Username)
	if err != nil {
		return models.User{}, constants.SERVICE_FAILED
	}
	if info.Id <= 0 {
		return models.User{}, constants.SERVICE_NO_EXIST
	}

	//validate password
	if tools.NewGenerate().GenerateMd5(params.Password) == info.Password {
		return models.User{}, constants.SERVICE_PASSWORD_ERROR
	}

	return info, constants.SERVICE_SUCCESS
}
