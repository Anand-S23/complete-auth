package config

import (
	"fmt"
	"log"
	"os"
)

type EnvVars struct {
    PRODUCTION           bool
    PORT                 string
    FE_URI               string
    DB_URI               string
    JWT_SECRET           []byte
    COOKIE_HASH_KEY      []byte
    COOKIE_BLOCK_KEY     []byte
    GOOGLE_CLIENT_SECRET string
    GOOGLE_CLIENT_ID     string
    GOOGLE_CALLBACK_URI  string
}

func LoadEnv() (*EnvVars) {
    envMode       := MustGetEnv("MODE")
    port          := MustGetEnv("PORT")
    feURI         := MustGetEnv("FE_URI")
    dbURI         := MustGetEnv("DB_URI")
    secret        := MustGetEnv("JWT_SECRET")
    hashKey       := MustGetEnv("COOKIE_HASH_KEY")
    blockKey      := MustGetEnv("COOKIE_BLOCK_KEY")
    gClientId     := MustGetEnv("GOOGLE_CLIENT_ID")
    gClientSecret := MustGetEnv("GOOGLE_CLIENT_SECRET")
    gCallbackUri  := MustGetEnv("GOOGLE_CALLBACK_URI")

    return &EnvVars {
        PRODUCTION: (envMode == "production"),
        FE_URI: feURI,
        DB_URI: dbURI,
        JWT_SECRET: []byte(secret),
        COOKIE_HASH_KEY: []byte(hashKey),
        COOKIE_BLOCK_KEY: []byte(blockKey),
        GOOGLE_CLIENT_ID: gClientId,
        GOOGLE_CLIENT_SECRET: gClientSecret,
        GOOGLE_CALLBACK_URI: gCallbackUri,
        PORT: port,
    }
}

func MustGetEnv(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
        message := fmt.Sprintf("Must provide %s variable in .env file", env)
        log.Fatal(message)
	}

	return variable
} 

