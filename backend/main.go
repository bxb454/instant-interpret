package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bxb454/instant-interpret/package/config"
	"github.com/bxb454/instant-interpret/package/database"
	"github.com/bxb454/instant-interpret/package/randgen"
	"github.com/bxb454/instant-interpret/package/translationservice"
	"github.com/bxb454/instant-interpret/package/websocket"
	"github.com/rs/cors"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request, ts *translationservice.TranslationService, db *database.Database) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := websocket.NewClient(conn, pool, ts, db, randgen.GenerateRandomUsername(), "en", "en")

	pool.Register <- client
	client.Read()
}

func serveMessageHistory(w http.ResponseWriter, r *http.Request, db *database.Database) {
	messages, err := db.GetAllMessages()
	if err != nil {
		http.Error(w, "Could not fetch messages", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, "Could not encode messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func setupRoutes(ts *translationservice.TranslationService, db *database.Database) {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r, ts, db)
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		serveMessageHistory(w, r, db)
	})
}

func main() {
	fmt.Println("Instant Interpret v0.1.0")
	fmt.Println(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	ctx := context.Background()

	ts, err := translationservice.NewTranslationService(ctx)
	if err != nil {
		log.Fatalf("Failed to create translation service: %v", err)
	}

	fmt.Println(fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DBUser, config.DBPassword, config.DBName))
	db, err := database.NewDatabase(fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DBUser, config.DBPassword, config.DBName))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	defer ts.Close()

	setupRoutes(ts, db)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // Change this to the actual origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	handler := c.Handler(http.DefaultServeMux)

	http.ListenAndServe(":8080", handler)
}
