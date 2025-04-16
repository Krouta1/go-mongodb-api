package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/krouta1/go-mongodb-api/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getSession())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	http.ListenAndServe("localhost:8080", r)
}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		fmt.Println("Error connecting to MongoDB (Dial):", err)
		panic(err)
	}

	err = s.Ping() // Try a simple ping
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		s.Close() // Close the session if ping fails
		panic(err)
	}

	fmt.Println("Successfully established MongoDB session.")
	return s
}
