package main

import (
	"db/controller"
	"db/dao"
	"db/infra"
	"db/usecase"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := infra.InitDB()
	if err != nil {
		println("fail: InitDB(), %v\n", err)
	}
	defer db.Close()

	infra.CloseDBWithSysCall(db)

	userDAO := &dao.UserDAO{DB: db}

	searchUC := &usecase.SearchUserUseCase{UserDAO: userDAO}
	registerUC := &usecase.RegisterUserUseCase{UserDAO: userDAO}

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller := &controller.SearchUserController{UseCase: searchUC}
			controller.ServeHTTP(w, r)
		case http.MethodPost:
			controller := &controller.RegisterUserController{UseCase: registerUC}
			controller.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("fail: ListenAndServe(), %v\n", err)
	}
}
