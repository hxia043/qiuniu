package controller

import (
	"github/hxia043/qiuniu/internal/app/models"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

func (c *Controller) RegisterRouter(engine *gin.Engine) {
	v1Group := engine.Group("/qiuniu/v1")

	v1Group.GET("/ping", models.Ping)
	v1Group.GET("/version", models.Version)
	v1Group.POST("/log", models.CollectLog)
	v1Group.GET("/list", models.ListLog)
	v1Group.GET("/download", models.DownloadLog)
	v1Group.DELETE("/clean", models.CleanLog)
}

func New() *Controller {
	return &Controller{}
}
