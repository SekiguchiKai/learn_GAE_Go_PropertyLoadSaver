package api

import (
	"github.com/gin-gonic/gin"
	"github.com/SekiguchiKai/learn_GAE_Go_PropertyLoadSaver/util"
	"strconv"
	"github.com/SekiguchiKai/learn_GAE_Go_PropertyLoadSaver/model"
	"github.com/pkg/errors"
	"net/http"
)

const (
	_DefaultUserLimit = 20
	_MaxUserLimit = 50
)


type userQueryParams struct {
	Limit  int
	Cursor string
}


func InitUserAPI(g *gin.RouterGroup) {
	g.GET("/user/:id", getUser)
	g.GET("/user", getUsers)
	g.POST("/user/new", createUser)
	g.PUT("/user/:id", updateUser)
	g.DELETE("/user/:id", deleteUser)
}

func getUser(c *gin.Context) {
	util.InfoLog(c, "getUser is called")

	id := getUserID(c)

	var list model.User

	s := store







}

func getUsers(c *gin.Context) {

}


func createUser(c *gin.Context) {

}

func updateUser(c *gin.Context) {

}

func deleteUser(c *gin.Context) {

}

func newUserQueryParam(c *gin.Context)(userQueryParams, error) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return userQueryParams{},err
	}

	return userQueryParams{
		Limit:limit,
		Cursor:c.Query("cursor"),
	}, nil

}

// HTTPのリクエストボディのjsonデータUserに変換
func bindUserFromJson(c *gin.Context, dst *model.User)error {
	if err := c.BindJSON(dst); err != nil {
		return err
	}

	dst.ID = getUserID(c)
	return nil
}

// IDを取得
func getUserID(c *gin.Context) string {
	return c.Param("id")
}

func validateParamsForUser(u model.User) error{
	if u.Name == "" {
		return errors.New("name is required")
	}

	if u.Address == "" {
		return errors.New("address is required")
	}

	if u.Age < 0 {
		return errors.New("age should be over 0")
	}

	return nil
}