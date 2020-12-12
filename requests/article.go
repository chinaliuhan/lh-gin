package requests

/**
添加文章
*/
type AddArticleContentRequest struct {
	UserId   int    `form:"user_id" json:"user_id"  binding:"required,number"`
	Classify string `form:"classify" json:"classify"  binding:"required,containsrune|alphanum"`
	Title    string `form:"title" json:"title"  binding:"required,containsrune|alphanum"`
	Content  string `form:"content" json:"content"  binding:"required,containsrune|alphanum"`
}
