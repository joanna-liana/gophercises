package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // needed by migrate lib
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbString := "postgres://" + os.Getenv("DB_USER") +
	":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_HOST") +
	":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") +
	"?sslmode=disable"

	m, err := migrate.New("file://migrations/scripts", dbString)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
