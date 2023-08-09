package test

import (
	"database/sql"
	"fmt"
	"log"
	"lucassantoss1701/bank/internal/infra/web/webserver"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetConnection() *sql.DB {
	return db
}

func UpSQLite3(t *testing.T) error {
	var err error
	db, err = sql.Open("sqlite3", "./test.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("./migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Migrations applied successfully")
	return nil
}

func DownSQLite3(t *testing.T) {

	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove("./test.sqlite")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Migrations down successfully and database file removed")

}

func MockWebServer() webserver.WebServer {
	return webserver.WebServer{}

}
