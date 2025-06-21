package main

import (
	"context"
	"github.com/gorilla/mux"
	"hackathon/controller"
	"hackathon/infra/firebase"
	"hackathon/infra/gemini"

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

	geminiClient, err := gemini.NewGeminiGateway(ctx)
	if err != nil {
		log.Fatalf("fail: NewGeminiGateway(), %v\n", err)
	}
	defer geminiClient.Close()

	userRepo := mysql.NewUsersRepository()
	postRepo := mysql.NewPostsRepository()
	likeRepo := mysql.NewLikesRepository()
	repostRepo := mysql.NewRepostsRepository()
	followRepo := mysql.NewFollowsRepository()
	trendRepo := mysql.NewTrendsRepository()
	txExecutor := mysql.NewTxExecutor()
	notificationRepo := mysql.NewNotificationsRepository()
	spoilerRepo := mysql.NewSpoilersRepository()

	authUC := usecase.NewAuthUserUseCase(firebaseAuthRepo, userRepo, db)
	postCreateUC := usecase.NewPostCreateUseCase(txExecutor, geminiClient, postRepo, trendRepo, spoilerRepo, db)
	postGetRecentUC := usecase.NewPostGetRecentUseCase(postRepo, db)
	postGetRepliesUC := usecase.NewPostGetRepliesUseCase(postRepo, db)
	postFindByIdUC := usecase.NewPostFindByIdUseCase(postRepo, db)
	postGetByUserUC := usecase.NewPostGetByUserUseCase(postRepo, db)
	likeCreateUC := usecase.NewLikeCreateUseCase(txExecutor, likeRepo, db)
	likeDeleteUC := usecase.NewLikeDeleteUseCase(txExecutor, likeRepo, db)
	repostCreateUC := usecase.NewRepostCreateUseCase(txExecutor, repostRepo, db)
	repostDeleteUC := usecase.NewRepostDeleteUseCase(txExecutor, repostRepo, db)
	userRegisterUC := usecase.NewUserRegisterUseCase(txExecutor, userRepo, db)
	userFindProfileUC := usecase.NewUserFindProfileUseCase(userRepo, db)
	followCreateUC := usecase.NewFollowCreateUseCase(txExecutor, followRepo, db)
	followDeleteUC := usecase.NewFollowDeleteUseCase(txExecutor, followRepo, db)
	//trendExtractNounsUC := usecase.NewTrendExtractNounsUseCase(txExecutor, trendRepo, db)
	trendGetTopUC := usecase.NewTrendGetTopUseCase(trendRepo, db)
	notificationFetchUC := usecase.NewNotificationFetchUseCase(notificationRepo, db)
	postSearchUC := usecase.NewPostSearchUseCase(postRepo, db)

	authC := controller.NewAuthUserController(authUC)
	postCreateC := controller.NewPostCreateController(authUC, postCreateUC)
	postGetRecentC := controller.NewPostGetRecentController(authUC, postGetRecentUC)
	postGetRepliesC := controller.NewPostGetRepliesController(authUC, postGetRepliesUC)
	postFindByIdC := controller.NewPostFindByIdController(authUC, postFindByIdUC)
	postGetByUserC := controller.NewPostGetByUserController(authUC, postGetByUserUC)
	postSearchC := controller.NewPostSearchController(authUC, postSearchUC)
	likeC := controller.NewLikeController(authUC, likeCreateUC, likeDeleteUC)
	repostC := controller.NewRepostController(authUC, repostCreateUC, repostDeleteUC)
	userRegisterC := controller.NewUserRegisterController(authUC, userRegisterUC)
	userFindProfileC := controller.NewUserFindProfileController(authUC, userFindProfileUC)
	followC := controller.NewFollowController(authUC, followCreateUC, followDeleteUC)
	trendGetTopC := controller.NewTrendGetTopController(trendGetTopUC)
	notificationFetchC := controller.NewNotificationGetController(authUC, notificationFetchUC)

	r := mux.NewRouter()

	// RESTfulエンドポイント
	r.Handle("/auth", authC).Methods("GET")
	r.Handle("/posts", postCreateC).Methods("POST")
	r.Handle("/posts/recent", postGetRecentC).Methods("GET")
	r.Handle("/posts/search", postSearchC).Methods("GET")
	r.Handle("/posts/{postId}", postFindByIdC).Methods("GET")
	r.Handle("/posts/{postId}/replies", postGetRepliesC).Methods("GET")
	r.Handle("/likes", likeC).Methods("POST")
	r.Handle("/likes", likeC).Methods("DELETE")
	r.Handle("/reposts", repostC).Methods("POST")
	r.Handle("/reposts", repostC).Methods("DELETE")
	r.Handle("/users", userRegisterC).Methods("POST")
	r.Handle("/users/{target}/posts", postGetByUserC).Methods("GET")
	r.Handle("/users/{username}", userFindProfileC).Methods("GET")
	r.Handle("/users/{followed}/follow", followC).Methods("POST")
	r.Handle("/users/{followed}/follow", followC).Methods("DELETE")
	r.Handle("/trends/top", trendGetTopC).Methods("GET")
	r.Handle("/notifications", notificationFetchC).Methods("GET")

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
