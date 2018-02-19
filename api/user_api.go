package api

import (
	"github.com/mjibson/goon"
	"net/http"
)

type UserAPI struct{}

// Battle APIをまとめる
func (ua *UserAPI) HandleBattleAPI(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	g := goon.NewGoon(r)

	switch method {
	case "GET":
		ua.GetUser(w, r, g)
	case "POST":
		ua.Post(w, r, g)
	case "PUT":
		ba.doPut(w, r, g)
	}

}
func (ua *UserAPI) getUser(w http.ResponseWriter, r *http.Request) {


}

func (ua *UserAPI) PostUser(w http.ResponseWriter, r *http.Request) {

}

func (ua *UserAPI) PutUser(w http.ResponseWriter, r *http.Request) {

}

func (ua *UserAPI) DeleteUser(w http.ResponseWriter, r *http.Request) {

}