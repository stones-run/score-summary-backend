package views

import (
	"github.com/gin-gonic/gin"
	"score-summary-backend/database"
	"score-summary-backend/models"
	"score-summary-backend/response"
)

type Exam struct {
}

func (e *Exam) List(c *gin.Context) {
	type gradeData struct {
		Grade    uint8 `json:"grade"`
		IsUpload bool  `json:"is_upload"`
	}

	type ExamData struct {
		models.Exam
		Data []gradeData `json:"data"`
	}

	gradesConst := [6]uint8{1, 2, 3, 4, 5, 6}
	var exams []models.Exam
	db := database.DB.Omit("sort_num").Order("sort_num asc").Find(&exams)
	if db.Error != nil {
		response.F(c, db.Error.Error())
		return
	}
	var examData []ExamData
	for _, exam := range exams {
		var data []gradeData
		for _, gradeConst := range gradesConst {
			score := models.Score{}
			db = database.DB.Where("grade = ? and exam_id = ?", gradeConst, exam.ID).First(&score)
			if db.RowsAffected != 0 {
				data = append(data, gradeData{gradeConst, true})
			} else {
				data = append(data, gradeData{gradeConst, false})
			}
		}
		examData = append(examData, ExamData{exam, data})
	}
	response.S(c, examData)

}

func (e *Exam) Add(c *gin.Context) {
	var param models.Exam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.F(c, err.Error())
		return
	}
	result := database.DB.Where("name = ?", param.Name).Find(&param)
	if param.ID != 0 {
		response.F(c, "不能有同名的考试")
		return
	}
	result = database.DB.Create(&param)
	response.Auto(c, result, param)
}

func (e *Exam) Delete(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.First(&models.Exam{}, id)
	if result.Error == nil {
		response.F(c, "不存在的考试")
		return
	}
	result = database.DB.Delete(&models.Exam{}, id)
	if result.Error != nil {
		response.F(c, "删除考试失败")
		return
	}
	result = database.DB.Where("exam_id = ?", id).Delete(&models.Score{})
	response.Auto(c, result, true)
}
