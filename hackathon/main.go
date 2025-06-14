package main

import (
	"context"
	"hackathon/controller"
	"hackathon/infra/firebase"
	"hackathon/infra/mysql"
	"hackathon/middleware"
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

	allowedOrigins := []string{
		"http://localhost:3000",
		"https://uttc-hackathon-shuei-mizusawa-fe.vercel.app",
	}

	http.Handle("/auth", middleware.CORS(allowedOrigins, AuthC))
	http.Handle("/users", middleware.CORS(allowedOrigins, UserRegisterC))

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("fail: ListenAndServe(), %v\n", err)
	}
}
