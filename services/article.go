package services

import (
	"lh-gin/models"
	"lh-gin/repositories"
)

type ArticleService struct {
}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

/**
新增
*/
func (this *ArticleService) AddNew(model models.ArticleContent) (int64, error) {

	//insert db
	lastID, err := repositories.NewArticleManagerRepository().AddNew(model)

	return lastID, err
}

/**
获取详情
*/
func (this ArticleService) GetInfoByUid(uid int) (models.ArticleContent, error) {

	//get db
	info, err := repositories.NewArticleManagerRepository().GetInfoByUid(uid)
	if err != nil {
		return models.ArticleContent{}, err
	}

	return info, nil
}
