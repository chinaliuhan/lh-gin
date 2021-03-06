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

/**
绑定GA
*/
func (this UserService) GaBind(userID int, request *requests.GaBindRequest) (bool, int) {
	ok, err := tools.NewGoogleAuth().VerifyCode(request.GaSecret, request.GaCode)
	if !ok {
		return true, constants.SERVICE_GA_WRONG
	}
	if err != nil {
		tools.NewLogUtil().Warning("用户绑定GA失败,GA验证出现异常错误,UID:", userID, err.Error())
		return true, constants.SERVICE_FAILED
	}

	if _, err = repositories.NewUserManagerRepository().ModifyByID(userID, models.User{GaSecret: request.GaSecret}); !ok {
		return true, constants.SERVICE_FAILED
	}

	return true, constants.SERVICE_SUCCESS
}

/**
解绑GA
*/
func (this UserService) GaUnbind(userID int, request *requests.GaUnbindRequest) (bool, int) {

	info, _ := repositories.NewUserManagerRepository().GetInfoByID(userID)
	if info.Id <= 0 || info.GaSecret == "" {
		return true, constants.SERVICE_NO_EXIST
	}

	ok, err := tools.NewGoogleAuth().VerifyCode(info.GaSecret, request.GaCode)
	if !ok {
		return true, constants.SERVICE_GA_WRONG
	}
	if err != nil {
		tools.NewLogUtil().Warning("用户绑定GA失败,GA验证出现异常错误,UID:", userID, err.Error())
		return true, constants.SERVICE_FAILED
	}

	if _, err = repositories.NewUserManagerRepository().ModifyByID(userID, models.User{GaSecret: ""}); !ok {
		return true, constants.SERVICE_FAILED
	}

	return true, constants.SERVICE_SUCCESS
}
