package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/krouta1/go-mongodb-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
	}

	oid := bson.ObjectIdHex(id)
	user_model := models.User{}

	if err := uc.session.DB("go-mongodb-api").C("users").FindId(oid).One(&user_model); err != nil {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	user_json, err := json.Marshal(user_model)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader((http.StatusOK)) // 200
	fmt.Fprintf(w, "%s \n", user_json)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user_model := models.User{}

	json.NewDecoder(r.Body).Decode(&user_model)
	user_model.Id = bson.NewObjectId()

	uc.session.DB("go-mongodb-api").C("users").Insert(user_model)

	user_json, err := json.Marshal(user_model)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader((http.StatusCreated)) // 201
	fmt.Fprintf(w, "%s \n", user_json)

}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
	}

	oid := bson.ObjectIdHex(id)

	if err := uc.session.DB("go-mongodb-api").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	w.WriteHeader((http.StatusOK)) // 200
	fmt.Fprint(w, "Deleted user", oid, "\n")
}
