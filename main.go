package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/handler"
	"github.com/bhivam/saangees-backend/middleware"
)

func main() {
	logger := log.New(os.Stdout, "AUTH API :: ", log.LstdFlags)

	tokenStore := data.NewInMemoryTokenStore()
	userStore := data.NewInMemoryUserStore(tokenStore)

	router := http.NewServeMux()

	auth := middleware.GetAuthMiddlewareFunc(userStore, logger)

	userHandler := handler.NewUserHandler(logger, userStore)
	tokenHandler := handler.NewTokenHandler(logger, userStore, tokenStore)

	router.Handle("POST /token", http.HandlerFunc(tokenHandler.CreateToken))

	router.Handle("POST /user/create", http.HandlerFunc(userHandler.CreateUser))
	router.Handle("GET /user/list", http.HandlerFunc(userHandler.ListUsers)) // Admin Auth
	router.Handle("GET /user", http.HandlerFunc(userHandler.GetUser))        // Base Auth

	CORS := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"https://sangees-kitchen.vercel.app",
			"https://www.saangeeskitchen.com/",
		},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler

	server := &http.Server{
		Addr:         ":3000",
		Handler:      middleware.Logging(CORS(auth(router)), logger),
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
