package api

import (
	"net/http"
	"google.golang.org/appengine"
)

const (
	Get = "GET"
		Post = "POST"
			Put = "PUT"
				Delete = "DELETE"
)

type UserAPI struct{}

// Battle APIをまとめる
func (ua *UserAPI) HandleBattleAPI(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case Get:
		ua.getUser(w, r)
	case Post:
		ua.PostUser(w, r)
	case Put:
		ua.PutUser(w, r)
	case Delete:
		ua.DeleteUser(w, r)
	}

}
func (ua *UserAPI) getUser(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)



}

func (ua *UserAPI) PostUser(w http.ResponseWriter, r *http.Request) {

}

func (ua *UserAPI) PutUser(w http.ResponseWriter, r *http.Request) {

}

func (ua *UserAPI) DeleteUser(w http.ResponseWriter, r *http.Request) {

}