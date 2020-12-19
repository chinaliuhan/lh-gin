package repositories

import (
	"database/sql"
	"lh-gin/models"
	"lh-gin/tools"
	"log"
)

//只要任意一个struct实现了接口中的所有方法, 即认为其集成了该接口
type ArticleRepository interface {
	AddNew(article models.ArticleContent) (int64, error)
	GetInfoByUid(uid int) (models.ArticleContent, error)
}

type ArticleManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewArticleManagerRepository() ArticleRepository {

	return &ArticleManagerRepository{}
}

func (o *ArticleManagerRepository) AddNew(article models.ArticleContent) (int64, error) {
	var (
		err    error
		lastID int64
	)

	// insert db
	lastID, err = tools.NewMysqlInstance().InsertOne(article)
	if err != nil {
		log.Println("插入失败: ", err.Error())
		return 0, err
	}

	return lastID, nil
}

func (o *ArticleManagerRepository) GetInfoByUid(uid int) (models.ArticleContent, error) {
	tmp := models.ArticleContent{}
	if ok, err := tools.NewMysqlInstance().Where("user_id=?", uid).Get(&tmp); !ok {
		return tmp, err
	}
	return tmp, nil
}
