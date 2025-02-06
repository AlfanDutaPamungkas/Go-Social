package main

import (
	"log"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/db"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/env"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
)

func main() {
	env.LoadEnv()

	addr := env.GetEnv("DB_ADDR", "postgres://user:password@localhost:5432/mydb?sslmode=disable")
	conn, err := db.New(addr, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store)
}
