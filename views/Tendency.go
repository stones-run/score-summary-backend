package views

import (
	"github.com/gin-gonic/gin"
	"score-summary-backend/database"
	"score-summary-backend/models"
	"score-summary-backend/response"
)

type Tendency struct {
}

func (t *Tendency) Tendency(c *gin.Context) {
	var scores []models.Score
	studentNum := c.Param("studentNum")

	db := database.DB.Where("student_num = ? ", studentNum).Find(&scores)
	if response.DBError(c, db) {
		return
	}
	if db.RowsAffected == 0 {
		response.F(c, "数据缺失")
		return
	}
}
