package main

import (
	"context"
	"hackathon/controller"
	"hackathon/infra/firebase"
	"hackathon/infra/mysql"
	"hackathon/usecase"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
)

func main() {
	db, err := mysql.InitDB()
	if err != nil {
		log.Fatalf("fail: InitDB(), %v\n", err)
	}

	mysql.CloseDBWithSysCall(db)

	ctx := context.Background()
	firebaseAuthRepo, err := firebase.NewFirebaseAuthRepository(ctx)
	if err != nil {
		log.Fatalf("fail: InitFirebaseAuthRepo(), %v\n", err)
	}

	userRepo := mysql.NewUserRepository()
	postsRepo := mysql.NewPostsRepository()
	txExecutor := mysql.NewTxExecutor()

	AuthUC := usecase.NewAuthUserUseCase(firebaseAuthRepo, userRepo, db)
	postCreateUC := usecase.NewPostCreateUseCase(txExecutor, postsRepo, db)
	registerUC := usecase.NewUserRegisterUseCase(txExecutor, userRepo, db)

	AuthC := controller.NewAuthUserController(AuthUC)
	PostCreatC := controller.NewPostCreateController(postCreateUC)
	UserRegisterC := controller.NewUserRegisterController(registerUC)

	mux := http.NewServeMux()
	mux.Handle("/auth", AuthC)
	mux.Handle("/posts", PostCreatC)
	mux.Handle("/users", UserRegisterC)

	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("fail: ListenAndServe(), %v\n", err)
	}
}
