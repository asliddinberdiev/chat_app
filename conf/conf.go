package conf

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Conf struct {
	App      App
	Postgres Postgres
}

type App struct {
	Host        string
	Port        string
	Version     string
	TokenKey    string
	AccessTime  string
	RefreshTime string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

var Cfg *Conf

func Load(path string) {
	err := godotenv.Load(path + "/.env")
	if err != nil {
		log.Fatal("loading .env file error: ", err)
	}

	conf := viper.New()
	conf.AutomaticEnv()

	Cfg = &Conf{
		App: App{
			Host:        conf.GetString("APP_HOST"),
			Port:        conf.GetString("APP_PORT"),
			Version:     conf.GetString("APP_VERSION"),
			TokenKey:    conf.GetString("APP_TOKEN_KEY"),
			AccessTime:  conf.GetString("APP_ACCESS_TIME"),
			RefreshTime: conf.GetString("APP_REFRESH_TIME"),
		},
		Postgres: Postgres{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DB"),
			SSLMode:  conf.GetString("POSTGRES_SSLMODE"),
		},
	}
}
