package services

import (
	"lh-gin/constants"
	"lh-gin/models"
	"lh-gin/repositories"
	"lh-gin/requests"
	"lh-gin/tools"
)

type chatService struct {
}

func NewChatService() *chatService {
	return &chatService{}
}

func (s *chatService) Login(params *requests.ChatLoginRequest) (models.User, int) {

	//get db
	info, err := repositories.NewUserManagerRepository().GetInfoByMobile(params.Mobile)
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
