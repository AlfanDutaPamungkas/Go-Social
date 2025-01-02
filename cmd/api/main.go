package main

import (
	"log"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/env"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
)

func main() {
	env.LoadEnv()

	cfg := config{
		addr: env.GetEnv("PORT", ":8080"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store: store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}