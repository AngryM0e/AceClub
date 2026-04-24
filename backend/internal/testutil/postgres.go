package testutil

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AngryM0e/AceClub/Backend/internal/repository/postgres/pgutils"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

type TestContainer struct {
	Container *postgres.PostgresContainer
	ConnStr   string
}

func SetupTest(t *testing.T, ctx context.Context) (*sql.DB, func(), func()) {
	tc := SetupTestContainer(t, ctx)

	db, err := sql.Open("postgres", tc.ConnStr)
	require.NoError(t, err)

	tc.RunMigrations(t, "../../../migrations")

	cleanupData := func() {
		_, err := db.Exec("TRUNCATE TABLE users CASCADE")
		if err != nil {
			t.Logf("Failed to truncate: %v", err)
		}
	}

	fullCleanup := func() {
		db.Close()
		tc.Cleanup(t)
	}

	return db, cleanupData, fullCleanup
}

func SetupTestContainer(t *testing.T, ctx context.Context) *TestContainer {
	container, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}
	connStr = connStr + "&sslmode=disable"

	return &TestContainer{
		Container: container,
		ConnStr:   connStr,
	}
}

func (tc *TestContainer) RunMigrations(t *testing.T, migrationsPath string) {
	err := pgutils.RunMigrations(migrationsPath, tc.ConnStr)
	if err != nil {
		t.Fatal(err)
	}
}

func (tc *TestContainer) Cleanup(t *testing.T) {
	if err := tc.Container.Terminate(context.Background()); err != nil {
		t.Logf("failed to terminate container: %v", err)
	}
}
