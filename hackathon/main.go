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
	postRepo := mysql.NewPostsRepository()
	likeRepo := mysql.NewLikesRepository()
	repostRepo := mysql.NewRepostsRepository()
	txExecutor := mysql.NewTxExecutor()

	AuthUC := usecase.NewAuthUserUseCase(firebaseAuthRepo, userRepo, db)
	postCreateUC := usecase.NewPostCreateUseCase(txExecutor, postRepo, db)
	postGetRecentUC := usecase.NewPostGetRecentUseCase(postRepo, db)
	likeCreateUC := usecase.NewLikeCreateUseCase(txExecutor, likeRepo, db)
	likeDeleteUC := usecase.NewLikeDeleteUseCase(txExecutor, likeRepo, db)
	repostCreateUC := usecase.NewRepostCreateUseCase(txExecutor, repostRepo, db)
	repostDeleteUC := usecase.NewRepostDeleteUseCase(txExecutor, repostRepo, db)
	registerUC := usecase.NewUserRegisterUseCase(txExecutor, userRepo, db)

	AuthC := controller.NewAuthUserController(AuthUC)
	postCreateC := controller.NewPostCreateController(postCreateUC)
	postGetRecentC := controller.NewPostGetRecentController(postGetRecentUC)
	likeC := controller.NewLikeController(likeCreateUC, likeDeleteUC)
	repostC := controller.NewRepostController(repostCreateUC, repostDeleteUC)
	userRegisterC := controller.NewUserRegisterController(registerUC)

	mux := http.NewServeMux()
	mux.Handle("/auth", AuthC)
	mux.Handle("/posts/", postCreateC)
	mux.Handle("/posts/recent", postGetRecentC)
	mux.Handle("/likes", likeC)
	mux.Handle("/reposts/", repostC)
	mux.Handle("/users", userRegisterC)

	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("fail: ListenAndServe(), %v\n", err)
	}
}
