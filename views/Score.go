package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"score-summary-backend/database"
	"score-summary-backend/models"
	"score-summary-backend/response"
)

type Score struct {
}

func (s *Score) GradeScore(c *gin.Context) {
	var scores []models.Score
	examId := c.Query("exam_id")
	grade := c.Query("grade")
	db := database.DB.Where("exam_id = ? AND grade = ?", examId, grade).Find(&scores)
	if response.DBError(c, db) {
		return
	}
	result := make([]map[string]interface{}, len(scores))
	for i, score := range scores {
		result[i] = make(map[string]interface{})
		result[i]["student_num"] = score.StudentNum
		result[i]["student_name"] = score.StudentName
		result[i]["id"] = score.ID
		for k, v := range score.ScoreJSON {
			result[i][k] = v
		}
	}
	response.S(c, result)
}

func (s *Score) Save(c *gin.Context) {
	type update struct {
		Grade  []uint8 `json:"grade"`
		ExamId []int   `json:"exam_id"`
	}

	type list struct {
		Exam   []models.Exam  `json:"exam"`
		Score  []models.Score `json:"score"`
		Update update         `json:"update"`
	}
	var param list
	if err := c.ShouldBindJSON(&param); err != nil {
		response.F(c, err.Error())
		return
	}
	// 更新考试
	for i, exam := range param.Exam {
		db := database.DB.First(&models.Exam{}, exam.ID)
		if db.RowsAffected != 0 {
			if len(exam.Name) == 0 {
				response.F(c, "考试名称不能为空")
				return
			}
			fmt.Println("更新考试", exam.ID)
			exam.SortNum = uint(i)
			result := database.DB.Save(&exam)
			if result.Error != nil {
				response.F(c, result.Error.Error())
				return
			}
		} else {
			fmt.Println("不存在的考试", exam)
		}
	}
	// 删除
	db := database.DB.Where("grade IN (?) AND exam_id IN (?)", param.Update.Grade, param.Update.ExamId).Delete(&models.Score{})
	fmt.Println("删除数据", db.RowsAffected)
	// 存储成绩
	db = database.DB.CreateInBatches(param.Score, 20)
	response.Auto(c, db, true)
}
