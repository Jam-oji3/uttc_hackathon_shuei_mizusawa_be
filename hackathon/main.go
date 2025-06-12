package main

import (
	"fmt"
	"hackathon/controller"
	sqlDB "hackathon/infra/db"
	"hackathon/infra/mysql"
	"hackathon/usecase"
	"log"
	"net/http"
)

func main() {
	db, err := sqlDB.InitDB()
	if err != nil {
		fmt.Printf("fail: InitDB(), %v\n", err)
	}
	defer db.Close()
	sqlDB.CloseDBWithSysCall(db)

	userRepo := mysql.NewUserRepository(db)
	registerUC := usecase.NewUserRegisterUseCase(userRepo, db)

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
