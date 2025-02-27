package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/handler"
	"github.com/bhivam/saangees-backend/middleware"
)

func main() {
	logger := log.New(os.Stdout, "AUTH API :: ", log.LstdFlags)

	err := godotenv.Load()
	if err != nil {
		logger.Println("Did not load .env file :: ", err)
	}

	db_cnx_string := os.Getenv("DB_CONNECTION_STRING")
	if db_cnx_string == "" {
		logger.Println("DB_CONNECTION_STRING is not set")
		return
	}

	db, err := gorm.Open(postgres.Open(db_cnx_string), &gorm.Config{})
	if err != nil {
		logger.Println("Error connecting to database :: ", err)
		return
	}

	tokenStore := data.NewPostgresTokenStore(db)
	userStore := data.NewPostgresUserStore(db)
	itemStore := data.NewPostgresItemStore(db)

	router := http.NewServeMux()

	auth := middleware.GetAuthMiddlewareFunc(userStore, logger)

	userHandler := handler.NewUserHandler(logger, userStore)
	tokenHandler := handler.NewTokenHandler(logger, userStore, tokenStore)
	itemHandler := handler.NewItemHandler(logger, itemStore)

	router.Handle("POST /token", http.HandlerFunc(tokenHandler.CreateToken)) // Base Auth (kind of)
	router.Handle("DELETE /token", http.HandlerFunc(tokenHandler.DeleteToken))

	router.Handle("POST /user/create", http.HandlerFunc(userHandler.CreateUser))
	router.Handle("GET /user/list", http.HandlerFunc(userHandler.ListUsers)) // Admin Auth
	router.Handle("GET /user", http.HandlerFunc(userHandler.GetUser))        // Base Auth

	router.Handle("POST /item/create", http.HandlerFunc(itemHandler.CreateItem))  // Admin Auth
	router.Handle("PUT /item/update", http.HandlerFunc(itemHandler.UpdateItem))   // Admin Auth
	router.Handle("GET /item/week", http.HandlerFunc(itemHandler.GetItemsByWeek)) // Base Auth
	router.Handle("GET /item", http.HandlerFunc(itemHandler.GetItem))             // Base Auth

	CORS := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"https://sangees-kitchen.vercel.app",
			"https://www.saangeeskitchen.com",
		},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
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
