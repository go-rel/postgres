package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/go-rel/rel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
)

func init() {
	// hack to make sure location it has the same location object as returned by pq driver.
	time.Local, _ = time.LoadLocation("Asia/Jakarta")
}

func pgxOpen(dsn string) (rel.Adapter, error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	database, err := sql.Open("pgx", stdlib.RegisterConnConfig(config))
	if err != nil {
		return nil, err
	}
	return New(database), err
}

func TestAdapterPgx_specs(t *testing.T) {
	adapter, err := pgxOpen(dsn())
	assert.Nil(t, err)
	defer adapter.Close()

	repo := rel.New(adapter)
	AdapterSpecs(t, repo)
}

func TestAdapterPgx_Transaction_commitError(t *testing.T) {
	adapter, err := pgxOpen(dsn())
	assert.Nil(t, err)
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit(ctx))
}

func TestAdapterPgx_Transaction_rollbackError(t *testing.T) {
	adapter, err := pgxOpen(dsn())
	assert.Nil(t, err)
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback(ctx))
}

func TestAdapterPgx_Exec_error(t *testing.T) {
	adapter, err := pgxOpen(dsn())
	assert.Nil(t, err)
	defer adapter.Close()

	_, _, err = adapter.Exec(ctx, "error", nil)
	assert.NotNil(t, err)
}
