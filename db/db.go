package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/marcioaso/consult/config"
)

type DB struct {
	Instance *sql.DB
}

var Db = DB{}

func (d *DB) Close() {
	d.Instance.Close()
}

func (d *DB) Connect() {
	conf := config.AppConfig

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/consult", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort) // Replace with your details
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	d.Instance = db
}
func (d *DB) CreateTables() {
	if d.Instance == nil {
		d.Connect()
	}

	// Create table if not exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS bybit_responses (
		id INT AUTO_INCREMENT PRIMARY KEY,
		t BIGINT NOT NULL,
		v VARCHAR(255),
		o VARCHAR(255),
		c VARCHAR(255),
		h VARCHAR(255),
		l VARCHAR(255),
		s VARCHAR(255),
		sn VARCHAR(255),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := d.Instance.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
