package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"leetkob/backend/internal/api"
)

func main() {
	databaseURL, email, password, role := os.Getenv("DATABASE_URL"), strings.ToLower(strings.TrimSpace(os.Getenv("USER_EMAIL"))), os.Getenv("USER_PASSWORD"), os.Getenv("USER_ROLE")
	if databaseURL == "" || email == "" || len(password) < 12 || (role != "student" && role != "ta" && role != "admin") {
		log.Fatal("Set DATABASE_URL, USER_EMAIL, USER_PASSWORD (at least 12 bytes), and USER_ROLE (student, ta, or admin)")
	}
	hash, err := api.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	id := uuid.NewString()
	if _, err := pool.Exec(context.Background(), "INSERT INTO users (id,email,password_hash,role) VALUES ($1,$2,$3,$4)", id, email, hash, role); err != nil {
		log.Fatal("Could not create user; check the database and email uniqueness")
	}
	log.Printf("Created %s user %s", role, id)
}
