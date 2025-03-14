package config

import (
	"flag"
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	Server   Server
	Database Database
	Secret   Secret
	File     File
}

type Server struct {
	Host string
	Port string
}

type Database struct {
	Host string
	Port string
	User string
	Pass string
	Name string
	Tz   string
}

type Secret struct {
	Jwt string
}

type File struct {
	MaxUploadSize string
	UploadPath    string
	LinkCover     string
}

func Get() *Config {
	fileFlag := flag.String("env", "", "file .env location path absolute")
	flag.Parse()

	var err error
	if *fileFlag != "" {
		err = godotenv.Load(*fileFlag)
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatal("error when load .env: ", err.Error())
	}

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
			Tz:   os.Getenv("DB_TZ"),
		},
		Secret: Secret{
			Jwt: os.Getenv("SECRET_JWT"),
		},
		File: File{
			MaxUploadSize: os.Getenv("MAX_UPLOAD_SIZE"),
			LinkCover:     os.Getenv("LINK_COVER"),
			UploadPath:    os.Getenv("UPLOAD_PATH"),
		},
	}
}
