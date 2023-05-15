package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Exam struct {
	ID      int    `gorm:"primaryKey;autoIncrement;comment:'主键'" json:"id"`
	Name    string `gorm:"not null;comment:'考试名'" json:"name"`
	SortNum uint   `gorm:"not null;comment:'排序'" json:"sort_num"`
}

func (r *Exam) TableName() string {
	return "exam"
}

type Score struct {
	ID          int            `gorm:"primaryKey;autoIncrement;comment:'主键'" json:"id"`
	ExamId      int            `gorm:"not null;comment:'考试ID'" json:"exam_id"`
	Grade       uint8          `gorm:"not null;comment:'年级'" json:"grade"`
	StudentNum  string         `gorm:"not null;comment:'学号'" json:"student_num"`
	StudentName string         `gorm:"not null;comment:'姓名'" json:"student_name"`
	Class       uint8          `gorm:"not null;comment:'班级'" json:"class"`
	Score       string         `gorm:"not null;comment:'成绩'" json:"-"`
	ScoreJSON   map[string]int `gorm:"-" json:"score"`
}

// BeforeSave 钩子函数，在保存数据之前自动调用
func (r *Score) BeforeSave(_ *gorm.DB) error {
	scoreJSON, err := json.Marshal(r.ScoreJSON)
	if err != nil {
		return err
	}
	r.Score = string(scoreJSON)
	return nil
}

// AfterFind 钩子函数，在检索数据之后自动调用
func (r *Score) AfterFind(_ *gorm.DB) error {
	err := json.Unmarshal([]byte(r.Score), &r.ScoreJSON)
	if err != nil {
		return err
	}
	return nil
}

func (r *Score) TableName() string {
	return "score"
}
