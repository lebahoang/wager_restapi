package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"hoang/m/handlers"
	"hoang/m/logPkg"
	"hoang/m/models"
	"hoang/m/utils"

	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	logPkg.Init()
	log := logPkg.GetLog()
	log.Infoln("Create DB Connection")
	// Open a PostgreSQL database.
	dsn := utils.GetENV("POSTGRES_URL", "postgres://postgres:12345@db:5432/postgres?sslmode=disable")
	pgDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	pgDB.SetMaxOpenConns(10)
	models.CreateDBConnection(pgDB)

	handler := handlers.SetupRoutes()

	port := utils.GetENV("API_PORT", "8080")
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.ListenAndServe()
}
