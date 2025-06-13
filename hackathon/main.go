package main

import (
	"fmt"
	"hackathon/controller"
	"hackathon/infra/mysql"
	"hackathon/usecase"
	"log"
	"net/http"
)

func main() {
	db, err := mysql.InitDB()
	if err != nil {
		fmt.Printf("fail: InitDB(), %v\n", err)
	}
	defer db.Close()
	mysql.CloseDBWithSysCall(db)

	userRepo := mysql.NewUserRepository()
	txExecutor := mysql.NewTxExecutor()
	registerUC := usecase.NewUserRegisterUseCase(txExecutor, userRepo, db)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller := &controller.UserRegisterController{UseCase: registerUC}
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
