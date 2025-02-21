package main

import (
	"time"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/auth"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/db"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/env"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/mailer"
	ratelimiter "github.com/AlfanDutaPamungkas/Go-Social/internal/rate_limiter"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
	"github.com/AlfanDutaPamungkas/Go-Social/internal/store/cache"
	"github.com/redis/go-redis/v9"
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

// @BasePath					/v1
//
// @secureDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	env.LoadEnv()

	cfg := config{
		addr:   env.GetEnv("PORT", ":8080"),
		apiURL: env.GetEnv("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetEnv("DB_ADDR", "postgres://user:password@localhost:5432/mydb?sslmode=disable"),
			maxOpenConns: env.GetIntEnv("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetEnv("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:   env.GetEnv("REDIS_ADDR", "localhost:6379"),
			pw:     env.GetEnv("REDIS_PW", ""),
			db:     env.GetIntEnv("REDIS_DB", 0),
			enable: env.GetBoolEnv("REDIS_ENABLED", false),
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
		frontendURL: env.GetEnv("FRONTEND_URL", "http://localhost:5173"),
		auth: authConfig{
			basic: basicConfig{
				user: env.GetEnv("AUTH_USERNAME", ""),
				pass: env.GetEnv("AUTH_PASS", ""),
			},
			token: tokenConfig{
				secret: env.GetEnv("AUTH_TOKEN_SECRET", ""),
				exp:    time.Hour * 24 * 3,
				iss:    "gophersocial",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetIntEnv("RATE_LIMITER_REQUESTS_COUNT", 20),
			TimeFrame:           time.Second * 5,
			Enabled:             env.GetBoolEnv("RATE_LIMITER_ENABLED", true),
		},
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

	var rdb *redis.Client
	if cfg.redisCfg.enable {
		rdb = cache.NewRedisClient(
			cfg.redisCfg.addr,
			cfg.redisCfg.pw,
			cfg.redisCfg.db,
		)
		logger.Info("redis cache connection established")
	}

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	mailer := mailer.NewSMTPMailer(
		cfg.mail.smtp.host,
		cfg.mail.smtp.port,
		cfg.mail.smtp.username,
		cfg.mail.smtp.password,
	)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		rateLimiter: rateLimiter,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
