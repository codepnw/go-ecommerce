package config

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config interface {
	App() AppConfig
	Db() DbConfig
	Jwt() JwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

func LoadConfig(path string) Config {
	env, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed: %v", err)
	} 

	return &config{
		app: &app{
			host: env["APP_HOST"],
			port: func() int {
				result, err := strconv.Atoi(env["APP_PORT"])
				if err != nil {
					log.Fatalf("load APP_PORT failed: %v", err)
				}
				return result
			}(),
			name:    env["APP_NAME"],
			version: env["APP_VERSION"],
			readTimeout: func() time.Duration {
				res, err := strconv.Atoi(env["APP_READ_TIMEOUT"])
				if err != nil {
					log.Fatalf("load APP_READ_TIMEOUT failed: %v", err)
				}
				return time.Duration(int64(res) * int64(math.Pow10(9)))
			}(),
			writeTimeout: func() time.Duration {
				res, err := strconv.Atoi(env["APP_WRITE_TIMEOUT"])
				if err != nil {
					log.Fatalf("load APP_WRITE_TIMEOUT failed: %v", err)
				}
				return time.Duration(int64(res) * int64(math.Pow10(9)))
			}(),
			bodyLimit: func() int {
				result, err := strconv.Atoi(env["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatalf("load APP_BODY_LIMIT failed: %v", err)
				}
				return result
			}(),
			fileLimit: func() int {
				result, err := strconv.Atoi(env["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatalf("load APP_FILE_LIMIT failed: %v", err)
				}
				return result
			}(),
		},
		db: &db{
			driver: env["DB_DRIVER"],
			host:   env["DB_HOST"],
			port: func() int {
				result, err := strconv.Atoi(env["DB_PORT"])
				if err != nil {
					log.Fatalf("load DB_PORT failed: %v", err)
				}
				return result
			}(),
			protocal: env["DB_PROTOCAL"],
			username: env["DB_USERNAME"],
			password: env["DB_PASSWORD"],
			database: env["DB_DATABASE"],
			sslMode:  env["DB_SSL_MODE"],
			maxConnections: func() int {
				result, err := strconv.Atoi(env["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("load DB_MAX_CONNECTIONS failed: %v", err)
				}
				return result
			}(),
		},
		jwt: &jwt{
			adminKey:  env["JWT_ADMIN_KEY"],
			secertKey: env["JWT_SECRET_KEY"],
			apiKey:    env["JWT_API_KEY"],
			accessExpiresAt: func() int {
				result, err := strconv.Atoi(env["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("load JWT_ACCESS_EXPIRES failed: %v", err)
				}
				return result
			}(),
			refreshExpiresAt: func() int {
				result, err := strconv.Atoi(env["JWT_REFRESH_EXPIRES"])
				if err != nil {
					log.Fatalf("load JWT_REFRESH_EXPIRES failed: %v", err)
				}
				return result
			}(),
		},
	}
}
