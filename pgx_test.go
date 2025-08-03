package postgres

import (
	"testing"
	"time"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func init() {
	// hack to make sure location it has the same location object as returned by pgx driver.
	time.Local, _ = time.LoadLocation("Asia/Jakarta")
}

func TestAdapterPgx_specs(t *testing.T) {
	adapter := MustOpen(dsn(), WithDriver("pgx"))
	defer adapter.Close()

	repo := rel.New(adapter)
	AdapterSpecs(t, repo)
}

func TestAdapterPgx_Transaction_commitError(t *testing.T) {
	adapter := MustOpen(dsn(), WithDriver("pgx"))
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit(ctx))
}

func TestAdapterPgx_Transaction_rollbackError(t *testing.T) {
	adapter := MustOpen(dsn(), WithDriver("pgx"))
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback(ctx))
}

func TestAdapterPgx_Exec_error(t *testing.T) {
	adapter := MustOpen(dsn(), WithDriver("pgx"))
	defer adapter.Close()

	_, _, err := adapter.Exec(ctx, "error", nil)
	assert.NotNil(t, err)
}

func TestAdapterPgx_InvalidDriverPanic(t *testing.T) {
	assert.Panics(t, func() {
		MustOpen("postgres://test:test@localhost:1111/test?sslmode=disable&timezone=Asia/Jakarta", WithDriver("pgx/v4"))
	})
}
