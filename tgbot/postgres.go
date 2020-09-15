package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func GetPostgres(config PostgresConfig) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panicf("Error on creating Postgres client: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Panicf("Error on checking Postgres connection: %v", err)
	}
	return db
}
