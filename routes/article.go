package routes

import (
	"go-gorm/controllers"
	middleware "go-gorm/middlewares"

	"github.com/gin-gonic/gin"
)

type ArticleRouteController struct {
	articleController controllers.ArticleController
}

func NewRouteArticleController(articleController controllers.ArticleController) ArticleRouteController {
	return ArticleRouteController{articleController}
}

func (pc *ArticleRouteController) ArticleRoute(rg *gin.RouterGroup) {

	router := rg.Group("articles")
	router.Use(middleware.DeserializeUser())
	router.POST("/", pc.articleController.CreateArticle)
	router.GET("/", pc.articleController.FindArticle)
	router.PUT("/:articleId", pc.articleController.UpdateArticle)
	router.GET("/:articleId", pc.articleController.FindArticleById)
	router.DELETE("/:articleId", pc.articleController.DeleteArticle)
}
