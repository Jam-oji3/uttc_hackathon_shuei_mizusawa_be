package main

import (
	"context"
	"hackathon/controller"
	"hackathon/infra/firebase"
	"hackathon/infra/mysql"
	"hackathon/usecase"
	"log"
	"net/http"
)

func main() {
	db, err := mysql.InitDB()
	if err != nil {

		log.Fatalf("fail: InitDB(), %v\n", err)
	}
	defer db.Close()
	mysql.CloseDBWithSysCall(db)

	ctx := context.Background()

	firebaseAuthRepo, err := firebase.NewFirebaseAuthRepository(ctx)
	if err != nil {
		log.Fatalf("fail: InitFirebaseAuthRepo(), %v\n", err)
	}

	userRepo := mysql.NewUserRepository()
	txExecutor := mysql.NewTxExecutor()

	AuthUC := usecase.NewAuthUserUseCase(firebaseAuthRepo, userRepo, db)
	registerUC := usecase.NewUserRegisterUseCase(txExecutor, userRepo, db)

	AuthC := controller.NewAuthUserController(AuthUC)
	UserRegisterC := controller.NewUserRegisterController(registerUC)

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		AuthC.ServeHTTP(w, r)
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		UserRegisterC.ServeHTTP(w, r)
	})

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("fail: ListenAndServe(), %v\n", err)
	}
}
