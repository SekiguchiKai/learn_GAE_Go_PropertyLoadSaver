package util

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	ContentType = "Content-Type"
	ApplicationJson = "application/json"
)


func ResponseHttp(w http.ResponseWriter, status int, obj interface{}) error{
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(ContentType, ApplicationJson)
	w.Write(b)

	return nil
}

func UnmarshalJsonFromRequest(r *http.Request, dst interface{}) error{
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	if err := json.Unmarshal(b, &dst); err != nil {
		return err
	}

	return nil

}