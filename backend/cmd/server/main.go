package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/api"
)

func main() {
	dbURL, origin := os.Getenv("DATABASE_URL"), os.Getenv("APP_ORIGIN")
	if dbURL == "" || origin == "" {
		log.Fatal("DATABASE_URL and APP_ORIGIN are required")
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	store := pgxstore.New(pool)
	defer store.StopCleanup()
	sessions := scs.New()
	sessions.Store = store
	sessions.Lifetime = 24 * time.Hour
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.SameSite = http.SameSiteLaxMode
	sessions.Cookie.Secure = os.Getenv("COOKIE_SECURE") != "false"
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{Addr: ":" + port, Handler: (&api.Server{DB: pool, Sessions: sessions, Origin: origin}).Routes(), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("Listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}
