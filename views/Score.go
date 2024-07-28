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

func getTotalScore(score models.Score, totalScores map[int]float32) float32 {
	totalScore, ok := totalScores[score.ID]
	if !ok {
		for _, v := range score.ScoreJSON {
			totalScore += v
		}
		totalScores[score.ID] = totalScore
	}
	return totalScore
}

func (s *Score) GradeScore(c *gin.Context) {
	var scores []models.Score
	examId := c.Param("examId")
	grade := c.Param("grade")
	db := database.DB.Where("exam_id = ? AND grade = ?", examId, grade).Find(&scores)
	if response.DBError(c, db) {
		return
	}
	if db.RowsAffected == 0 {
		response.F(c, "数据缺失")
		return
	}

	// 获取所有科目
	subjects := make([]string, 0)
	for k := range scores[0].ScoreJSON {
		subjects = append(subjects, k)
	}

	// 初始化排名信息
	type RankInfo struct {
		ClassRank int
		GradeRank int
	}
	rankMap := make(map[int]map[string]*RankInfo)
	for _, score := range scores {
		rankMap[score.ID] = make(map[string]*RankInfo)
		for _, subject := range subjects {
			rankMap[score.ID][subject] = &RankInfo{ClassRank: 1, GradeRank: 1}
		}
		rankMap[score.ID]["total"] = &RankInfo{ClassRank: 1, GradeRank: 1}
	}

	// 计算班级排名和年级排名
	totalScores := make(map[int]float32)
	// gradeCompared classCompared保存已经对比过的科目成绩
	gradeCompared := make(map[float32]bool)
	classCompared := make(map[float32]bool)

	// 计算总分排名
	for _, score1 := range scores {
		totalScore1 := getTotalScore(score1, totalScores)
		gradeCompared = make(map[float32]bool)
		classCompared = make(map[float32]bool)
		for _, score2 := range scores {
			if score1.ID != score2.ID && score1.Grade == score2.Grade {
				// 从 map 中获取总分
				totalScore2 := getTotalScore(score2, totalScores)
				// 计算班级排名和年级排名
				if totalScore2 > totalScore1 {
					if !gradeCompared[totalScore2] {
						gradeCompared[totalScore2] = true
						rankMap[score1.ID]["total"].GradeRank++
					}
					if score1.Class == score2.Class && !classCompared[totalScore2] {
						classCompared[totalScore2] = true
						rankMap[score1.ID]["total"].ClassRank++
					}
				}
			}
		}
	}
	// 计算各科排名
	for _, subject := range subjects {
		for _, score1 := range scores {
			gradeCompared = make(map[float32]bool)
			classCompared = make(map[float32]bool)
			for _, score2 := range scores {
				if score1.ID != score2.ID && score1.Grade == score2.Grade {
					// 计算班级排名和年级排名
					comparedScore := score2.ScoreJSON[subject]
					if comparedScore > score1.ScoreJSON[subject] {
						// 计算班级排名
						if !gradeCompared[comparedScore] {
							gradeCompared[comparedScore] = true
							rankMap[score1.ID][subject].GradeRank++
						}
						// 计算年级排名
						if score1.Class == score2.Class && !classCompared[comparedScore] {
							classCompared[comparedScore] = true
							rankMap[score1.ID][subject].ClassRank++
						}
					}

				}
			}
		}
	}

	// 构造返回结果
	result := make([]map[string]interface{}, len(scores))
	for i, score := range scores {
		result[i] = make(map[string]interface{})
		result[i]["student_num"] = score.StudentNum
		result[i]["student_name"] = score.StudentName
		result[i]["class"] = score.Class
		result[i]["id"] = score.ID
		for k, v := range score.ScoreJSON {
			result[i][k] = v
		}

		// 添加排名信息
		for _, subject := range subjects {
			result[i][subject+"ClassRank"] = rankMap[score.ID][subject].ClassRank
			result[i][subject+"GradeRank"] = rankMap[score.ID][subject].GradeRank
		}
		result[i]["totalClassRank"] = rankMap[score.ID]["total"].ClassRank
		result[i]["totalGradeRank"] = rankMap[score.ID]["total"].GradeRank

		// 从 map 中获取总分并添加到结果中
		result[i]["total"] = totalScores[score.ID]
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
