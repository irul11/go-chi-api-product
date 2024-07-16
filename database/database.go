package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func LoadDatabase() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
	// ================= Database =====================
	var (
		db_host     string = os.Getenv("DATABASE_HOST")
		db_name     string = os.Getenv("DATABASE_NAME")
		db_port     string = os.Getenv("DATABASE_PORT")
		db_username string = os.Getenv("DATABASE_USERNAME")
		db_password string = os.Getenv("DATABASE_PASSWORD")
	)
	db_url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_username, db_password, db_host, db_port, db_name)

	db, err := sql.Open("postgres", db_url)

	if err != nil {
		panic(err)
	}

	// defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	DB = db
	log.Println("Successfully connected to PostgreSQL")
}
