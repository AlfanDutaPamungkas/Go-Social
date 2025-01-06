package main

import (
	"log"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/db"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/env"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
)

const version = "0.0.1"

func main() {
	env.LoadEnv()

	cfg := config{
		addr: env.GetEnv("PORT", ":8080"),
		db: dbConfig{
			addr:         env.GetEnv("DB_ADDR", "postgres://user:password@localhost:5432/mydb?sslmode=disable"),
			maxOpenConns: env.GetIntEnv("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetEnv("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetEnv("env", "DEVELOPMENT"),
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Print("db connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
