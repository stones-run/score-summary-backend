package response

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Cd int

const (
	Success Cd = iota // 0
	Fail              // 1
)

type response struct {
	Cd   Cd          `json:"cd"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// S 成功的response
func S(c *gin.Context, data interface{}) {
	R(c, data, "", Success)
}

// F 失败的response
func F(c *gin.Context, msg string) {
	R(c, nil, msg, Fail)
}

// DBError 数据库操作是否失败的判定
func DBError(c *gin.Context, result *gorm.DB) bool {
	if result.Error != nil {
		R(c, nil, result.Error.Error(), Fail)
		return true
	}
	return false
}

// Auto 自动返回
func Auto(c *gin.Context, result *gorm.DB, data interface{}) {
	if !DBError(c, result) {
		S(c, data)
	}
}

// R 返回response
func R(c *gin.Context, data interface{}, msg string, cd Cd) {
	response := response{Cd: cd, Data: data, Msg: msg}
	c.JSON(http.StatusOK, response)
}
