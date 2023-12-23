package config

import (
	"os"

	"github.com/joho/godotenv"
)

/*
Env list
ACCESS_TOKEN_SECRET
REFRESH_TOKEN_SECRET
MONGODB_USERNAME
MONGODB_PASSWORD

Below will be consumed automatically
IMAGEKIT_PUBLIC_KEY: public_+3rAkPsHz8APem/ZFrHbJspD3VI=
IMAGEKIT_PRIVATE_KEY: private_9tdXxngXRqrplX1K667VoRD7R1I=
IMAGEKIT_ENDPOINT_URL: https://ik.imagekit.io/tipbk
*/

type EnvConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	MongodbUsername    string
	MongodbPassword    string
	DatabaseName       string
	MetadataEndpoint   string
}

func GetEnvConfig() *EnvConfig {
	godotenv.Load(".env")

	return &EnvConfig{
		AccessTokenSecret:  os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
		MongodbUsername:    os.Getenv("MONGODB_USERNAME"),
		MongodbPassword:    os.Getenv("MONGODB_PASSWORD"),
		DatabaseName:       os.Getenv("DATABASE_NAME"),
		MetadataEndpoint:   os.Getenv("METADATA_SERVICE_ENDPOINT_URL"),
	}
}
