package repositories

import (
	"lh-gin/models"
	"lh-gin/tools"
)

type userLoginRecordRepository struct {
}

func NewUserLoginRecordRepository() *userLoginRecordRepository {
	return &userLoginRecordRepository{}
}

func (r *userLoginRecordRepository) AddNew(userID int64, token string, clientIP string, timeStamp int) (int64, error) {
	userLoginRecordModel := models.UserLoginRecord{}
	userLoginRecordModel.UserId = userID
	userLoginRecordModel.Token = token
	userLoginRecordModel.CreatedIp = clientIP
	userLoginRecordModel.Created = timeStamp
	lastID, err := tools.NewMysqlInstance().InsertOne(&userLoginRecordModel)
	if err != nil {
		return 0, err
	}
	return lastID, err
}
