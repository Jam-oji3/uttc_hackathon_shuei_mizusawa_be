package main

import (
	"github.com/gorilla/mux"
	"hackathon/controller"
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
	//ctx := context.Background()
	//firebaseAuthRepo, err := firebase.NewFirebaseAuthRepository(ctx)
	if err != nil {
		log.Fatalf("fail: InitFirebaseAuthRepo(), %v\n", err)
	}

	userRepo := mysql.NewUsersRepository()
	postRepo := mysql.NewPostsRepository()
	likeRepo := mysql.NewLikesRepository()
	repostRepo := mysql.NewRepostsRepository()
	txExecutor := mysql.NewTxExecutor()

	//AuthUC := usecase.NewAuthUserUseCase(firebaseAuthRepo, userRepo, db)
	postCreateUC := usecase.NewPostCreateUseCase(txExecutor, postRepo, db)
	postGetRecentUC := usecase.NewPostGetRecentUseCase(postRepo, db)
	postGetRepliesUC := usecase.NewPostGetRepliesUseCase(postRepo, db)
	postFindByIdUC := usecase.NewPostFindByIdUseCase(postRepo, db)
	likeCreateUC := usecase.NewLikeCreateUseCase(txExecutor, likeRepo, db)
	likeDeleteUC := usecase.NewLikeDeleteUseCase(txExecutor, likeRepo, db)
	repostCreateUC := usecase.NewRepostCreateUseCase(txExecutor, repostRepo, db)
	repostDeleteUC := usecase.NewRepostDeleteUseCase(txExecutor, repostRepo, db)
	userRegisterUC := usecase.NewUserRegisterUseCase(txExecutor, userRepo, db)
	userFindProfileUC := usecase.NewUserFindProfileUseCase(userRepo, db)

	//AuthC := controller.NewAuthUserController(AuthUC)
	postCreateC := controller.NewPostCreateController(postCreateUC)
	postGetRecentC := controller.NewPostGetRecentController(postGetRecentUC)
	postGetRepliesC := controller.NewPostGetRepliesController(postGetRepliesUC)
	postFindByIdC := controller.NewPostFindByIdController(postFindByIdUC)
	likeC := controller.NewLikeController(likeCreateUC, likeDeleteUC)
	repostC := controller.NewRepostController(repostCreateUC, repostDeleteUC)
	userRegisterC := controller.NewUserRegisterController(userRegisterUC)
	userFindProfileC := controller.NewUserFindProfileController(userFindProfileUC)

	r := mux.NewRouter()

	// RESTfulエンドポイント
	r.Handle("/posts", postCreateC).Methods("POST")
	r.Handle("/posts/recent", postGetRecentC).Methods("GET")
	r.Handle("/posts/{postId}", postFindByIdC).Methods("GET")
	r.Handle("/posts/{postId}/replies", postGetRepliesC).Methods("GET")
	r.Handle("/likes", likeC).Methods("POST")
	r.Handle("/likes", likeC).Methods("DELETE")
	r.Handle("/reposts", repostC).Methods("POST")
	r.Handle("/reposts", repostC).Methods("DELETE")
	r.Handle("/users", userRegisterC).Methods("POST")
	r.Handle("/users/{username}", userFindProfileC).Methods("GET")

	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("fail: ListenAndServe(), %v\n", err)
	}
}
