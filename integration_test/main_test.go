//go:build integration
// +build integration

package integration_test

import (
	"database/sql"
	"math/rand"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"hoang/m/handlers"
	"hoang/m/logPkg"
	"hoang/m/models"
	"hoang/m/utils"

	"github.com/uptrace/bun/driver/pgdriver"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	logPkg.Init()
	log := logPkg.GetLog()
	log.Infoln("Create DB Connection")
	// Open a PostgreSQL database.
	dsn := utils.GetENV("POSTGRES_URL", "postgres://postgres:12345@localhost:55432/postgres?sslmode=disable")
	pgDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	pgDB.SetMaxIdleConns(500)
	models.CreateDBConnection(pgDB)

	handler := handlers.SetupRoutes()

	testServer = httptest.NewServer(handler)
	defer testServer.Close()
	code := m.Run()
	os.Exit(code)
}
