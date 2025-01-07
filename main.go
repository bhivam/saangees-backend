package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/handler"
	"github.com/bhivam/saangees-backend/middleware"
	"github.com/bhivam/saangees-backend/util"
)

func main() {
	logger := log.New(os.Stdout, "AUTH API :: ", log.LstdFlags)

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "01234567890123456789012345678901"
	}

	if len(secretKey) != 32 {
		logger.Fatal("SECRET_KEY must be 32 bytes long")
	}

	userStore := data.NewInMemoryUserStore()
	sessionStore := data.NewInMemorySessionStore()
	tokenMaker := util.NewJWTMaker(secretKey)

	router := mux.NewRouter()
	userHandler := handler.NewUserHandler(logger, userStore)
	sessionHandler := handler.NewSessionHandler(logger, userStore, sessionStore, tokenMaker)

	authMiddleware := middleware.GetAuthMiddlewareFunc(tokenMaker, logger, false)
	adminAuthMiddleware := middleware.GetAuthMiddlewareFunc(tokenMaker, logger, true)

	postRouter := router.Methods("POST").Subrouter()

	// TODO create user should have an admin and basic version
	postRouter.Handle("/user/create", http.HandlerFunc(userHandler.CreateUser))
	postRouter.Handle("/user/login", http.HandlerFunc(sessionHandler.LoginUser))
	postRouter.Handle(
		"/user/logout",
		http.HandlerFunc(sessionHandler.LogoutUser),
	)

	postRouter.Handle(
		"/token/revoke/{id}",
		adminAuthMiddleware(http.HandlerFunc(sessionHandler.RevokeSession)),
	)

	getRouter := router.Methods("GET").Subrouter()

	getRouter.Handle("/user/list", adminAuthMiddleware(http.HandlerFunc(userHandler.ListUsers)))
	getRouter.Handle("/user", authMiddleware(http.HandlerFunc(userHandler.GetUser)))

	getRouter.Handle("/token/refresh", http.HandlerFunc(sessionHandler.RefreshToken))

	CORS := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler

	server := &http.Server{
		Addr:         ":3000",
		Handler:      middleware.Logging(CORS(router), logger),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Println("Starting server on port 3000")
		err := server.ListenAndServe()

		logger.Println("Shutting Down :: ", err)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
