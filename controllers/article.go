package controllers

import (
	"go-gorm/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleController struct {
	DB *gorm.DB
}

func NewArticleController(DB *gorm.DB) ArticleController {
	return ArticleController{DB}
}

// Create Article Handler
func (ac *ArticleController) CreateArticle(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateArticleRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newArticle := models.Article{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		UserId:    int(currentUser.ID),
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newArticle)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Article with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newArticle})
}

// Update Article Handler
func (ac *ArticleController) UpdateArticle(ctx *gin.Context) {
	articleId := ctx.Param("articleId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdateArticle
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedArticle models.Article
	result := ac.DB.First(&updatedArticle, "id = ?", articleId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No article with that title exists"})
		return
	}
	now := time.Now()
	articleToUpdate := models.Article{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		UserId:    int(currentUser.ID),
		CreatedAt: updatedArticle.CreatedAt,
		UpdatedAt: now,
	}

	ac.DB.Model(&updatedArticle).Updates(articleToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedArticle})
}

// Get Single Post Handler
func (pc *ArticleController) FindArticleById(ctx *gin.Context) {
	articleId := ctx.Param("articleId")

	var article models.Article
	result := pc.DB.First(&article, "id = ?", articleId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No Article with that title exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": article})
}

// Get All Article Handler
func (pc *ArticleController) FindArticle(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var posts []models.Article
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&posts)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}

// Delete Article Handler
func (pc *ArticleController) DeleteArticle(ctx *gin.Context) {
	articleId := ctx.Param("articleId")

	result := pc.DB.Delete(&models.Article{}, "id = ?", articleId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No article with that title exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
