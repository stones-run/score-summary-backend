package routers

import (
	"github.com/gin-gonic/gin"
	"score-summary-backend/views"
)

func LoadUsersRoutes(r *gin.Engine) {
	exam := views.Exam{}
	examsGroup := r.Group("/exam")
	{
		examsGroup.GET("", exam.List)
		examsGroup.POST("", exam.Add)
		examsGroup.DELETE("/:id", exam.Delete)
	}
	score := views.Score{}
	scoreGroup := r.Group("/score")
	{
		scoreGroup.GET("", score.GradeScore)
		scoreGroup.POST("", score.Save)
	}
}
