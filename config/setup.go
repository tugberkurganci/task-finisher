package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDB() *sql.DB {

	dbURI := EnvPostgresURI()

	conn, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata oluştu: %v\n", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("Veritabanına ping atılırken hata oluştu: %v\n", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        payload TEXT,
        deadline TIMESTAMP,
        retry INTEGER,
        max_retries INTEGER,
        priority INTEGER,
        task_interval  TEXT, 
        completed BOOLEAN,
		error BOOLEAN
    )`)
	if err != nil {
		log.Fatalf("Failed to create tasks table: %v", err)
	}

	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL
	)
`
	_, err = conn.Exec(createUserTableSQL)
	if err != nil {
		log.Fatalf("users tablosu oluşturulurken hata oluştu: %v\n", err)
	}
	return conn

}
func EnvPostgresURI() string {

	host := "postgres"
	port := "5432"
	user := "postgres"
	password := "test"
	dbname := "konzek"

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func GetDB() *sql.DB {
	if db == nil {
		db = ConnectDB()
	}
	return db
}
