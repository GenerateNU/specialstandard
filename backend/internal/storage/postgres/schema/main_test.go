package schema_test

import (
	"log"
	"os"
	"specialstandard/internal/storage/postgres/testutil"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	// Setup shared container
	db, err := testutil.GetSharedTestDB()
	if err != nil {
		log.Fatal(err)
	}
	testDB = db.Pool

	// Run tests
	code := m.Run()

	// Cleanup
	testutil.Shutdown()
	os.Exit(code)
}
