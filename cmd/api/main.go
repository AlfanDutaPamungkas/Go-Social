package main

import (
	"time"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/db"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/env"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/mailer"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@secureDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description
func main() {
	env.LoadEnv()

	cfg := config{
		addr: env.GetEnv("PORT", ":8080"),
		apiURL: env.GetEnv("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetEnv("DB_ADDR", "postgres://user:password@localhost:5432/mydb?sslmode=disable"),
			maxOpenConns: env.GetIntEnv("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetEnv("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetEnv("env", "DEVELOPMENT"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
			smtp: smtpConfig{
				host:     env.GetEnv("SMTP_HOST", "smtp.example.com"),
				port:     env.GetEnv("SMTP_PORT", "587"),
				username: env.GetEnv("SMTP_USERNAME", ""),
				password: env.GetEnv("SMTP_PASSWORD", ""),
			},
		},
		frontendURL: env.GetEnv("FRONTEND_URL", "http://localhost:4000"),
	}

	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("db connection pool established")

	store := store.NewStorage(db)

	mailer := mailer.NewSMTPMailer(
		cfg.mail.smtp.host,
		cfg.mail.smtp.port,
		cfg.mail.smtp.username,
		cfg.mail.smtp.password,
	)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
