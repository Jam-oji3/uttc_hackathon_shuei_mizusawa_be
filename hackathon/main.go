package main

import (
	"context"
	"hackathon/controller"
	"hackathon/infra/firebase"
	"hackathon/infra/mysql"
	// "hackathon/middleware" // This is no longer needed
	"hackathon/usecase"
	"log"
	"net/http"

	"github.com/rs/cors" // Import the library
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

	// 1. Create a new ServeMux (router) to register your handlers
	mux := http.NewServeMux()
	mux.Handle("/auth", AuthC)
	mux.Handle("/users", UserRegisterC)

	// 2. Configure CORS options using the library
	c := cors.New(cors.Options{
		// Your allowed origins
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://uttc-hackathon-shuei-mizusawa-fe.vercel.app",
		},
		// Methods your frontend is allowed to use
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		// Headers your frontend is allowed to send (important for Authorization)
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		// Allow cookies and credentials to be sent
		AllowCredentials: true,
	})

	// 3. Wrap your router with the CORS middleware
	handler := c.Handler(mux)

	log.Println("Listening on :8080")
	// 4. Start the server with the new CORS-wrapped handler
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("fail: ListenAndServe(), %v\n", err)
	}
}
