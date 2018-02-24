package api

import (
	"github.com/gin-gonic/gin"
	"github.com/SekiguchiKai/learn_GAE_Go_PropertyLoadSaver/util"
	"strconv"
	"github.com/SekiguchiKai/learn_GAE_Go_PropertyLoadSaver/model"
	"github.com/pkg/errors"
	"github.com/SekiguchiKai/learn_GAE_Go_PropertyLoadSaver/store"
	"net/http"
	"time"
	"context"
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

	s := store.NewUserStore(c.Request)
	var u model.User
	if exists, err := s.GetUser(id, &u); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
		return
	}else if !exists {
		util.RespondAndLog(c, http.StatusBadRequest, errors.New("Invalid id" + id).Error())
		return
	}

	c.JSON(http.StatusOK, u)

}

func getUsers(c *gin.Context) {
	util.InfoLog(c, "getUsers is called")
	params, err := newUserQueryParam(c)
	if err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	if params.Limit <= 0 {
		params.Limit = _DefaultUserLimit
	}


	var list model.UserList

	s := store.NewUserStore(c.Request)
	if err := s.GetUserList(params.Cursor, params.Limit, &list); err != nil {
		util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)

}


func createUser(c *gin.Context) {
	var params model.User
	// HTTPリクエストで受け取ったJSONを構造体にロードする
	if err := bindUserFromJson(c, &params); err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	// 構造体のバリデーションを行う
	if err := validateParamsForUser(params); err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}


	u := model.NewUser(params)
	u.UpdatedAt = time.Now().UTC()


	err := store.RunInTransaction(c.Request, func(ctx context.Context) error {

		s := store.NewUserStoreWithContext(ctx)
		// Userは一意になるようにするために、既に同じユーザーが存在するかを確認する
		if exists, err := s.ExistsUser(u.ID); err != nil {
			util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
			return err
		} else if exists {
			caution := "There is same user"
			util.RespondAndLog(c, http.StatusBadRequest, caution)
			return errors.New(caution)
		}

		// ユーザーをDatastoreに格納する
		if err := s.PutUser(u); err != nil {
			util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
			return err
		}

		return nil

	})

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, nil)
}

func updateUser(c *gin.Context) {
	var params model.User
	if err := bindUserFromJson(c, &params); err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateParamsForUser(params); err != nil {
		util.RespondAndLog(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedAt := time.Now().UTC()

	err := store.RunInTransaction(c.Request, func(ctx context.Context) error {

		s := store.NewUserStoreWithContext(ctx)

		var source model.User
		if exists, err := s.GetUser(params.ID, &source); err != nil {
			util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
			return err
		} else if !exists {
			util.RespondAndLog(c, http.StatusNotFound, "id = %s is not found", params.ID)
			return errors.New("the airport is not found")
		}


		u := model.UpdatedUser(source, params)
		u.UpdatedAt = updatedAt

		if err := s.PutUser(u); err != nil {
			util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
			return err
		}


		return nil

	})

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, nil)


}

func deleteUser(c *gin.Context) {
	id := getUserID(c)

	err := store.RunInTransaction(c.Request, func(ctx context.Context) error {

		s := store.NewUserStoreWithContext(ctx)

		var u model.User
		if exists, err := s.GetUser(id, &u); err != nil {
			util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
			return err
		} else if !exists {
			util.RespondAndLog(c, http.StatusNotFound, "id = %s is not found", id)
			return errors.New("the u is not found")
		}

		if err := s.DeleteUser(id); err != nil {
			util.RespondAndLog(c, http.StatusInternalServerError, err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, nil)

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