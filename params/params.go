package params

type Pagination struct {
	Page  int `form:"page" json:"page"`
	Size  int `form:"size" json:"size"`
	Total int64
}

type GradeScore struct {
	ExamId int   `json:"exam_id"`
	Grade  uint8 `json:"grade"`
}
