package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-rel/primaryreplica"
	"github.com/go-rel/rel"
	"github.com/go-rel/sql/specs"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

func init() {
	// hack to make sure location it has the same location object as returned by pq driver.
	time.Local, _ = time.LoadLocation("Asia/Jakarta")
}

func dsn() string {
	if os.Getenv("POSTGRESQL_DATABASE") != "" {
		return os.Getenv("POSTGRESQL_DATABASE") + "?sslmode=disable&timezone=Asia/Jakarta"
	}

	return "postgres://rel:rel@localhost:25432/rel_test?sslmode=disable&timezone=Asia/Jakarta"
}

func TestAdapter_Name(t *testing.T) {
	adapter := MustOpen(dsn())
	defer adapter.Close()

	assert.Equal(t, Name, adapter.Name())
}

func AdapterSpecs(t *testing.T, repo rel.Repository) {
	// Prepare tables
	teardown := specs.Setup(repo)
	defer teardown()

	// Migration Specs
	specs.Migrate()

	// Query Specs
	specs.Query(t, repo)
	specs.QueryJoin(t, repo)
	specs.QueryJoinAssoc(t, repo)
	specs.QueryNotFound(t, repo)
	specs.QueryWhereSubQuery(t, repo)

	// Preload specs
	specs.PreloadHasMany(t, repo)
	specs.PreloadHasManyWithQuery(t, repo)
	specs.PreloadHasManySlice(t, repo)
	specs.PreloadHasOne(t, repo)
	specs.PreloadHasOneWithQuery(t, repo)
	specs.PreloadHasOneSlice(t, repo)
	specs.PreloadBelongsTo(t, repo)
	specs.PreloadBelongsToWithQuery(t, repo)
	specs.PreloadBelongsToSlice(t, repo)

	// Aggregate Specs
	specs.Aggregate(t, repo)

	// Insert Specs
	specs.Insert(t, repo)
	specs.InsertHasMany(t, repo)
	specs.InsertHasOne(t, repo)
	specs.InsertBelongsTo(t, repo)
	specs.Inserts(t, repo)
	specs.InsertAll(t, repo)
	specs.InsertOnConflictIgnore(t, repo)
	specs.InsertOnConflictReplace(t, repo)
	specs.InsertAllOnConflictIgnore(t, repo)
	specs.InsertAllOnConflictReplace(t, repo)
	specs.InsertAllPartialCustomPrimary(t, repo)

	// Update Specs
	specs.Update(t, repo)
	specs.UpdateNotFound(t, repo)
	specs.UpdateHasManyInsert(t, repo)
	specs.UpdateHasManyUpdate(t, repo)
	specs.UpdateHasManyReplace(t, repo)
	specs.UpdateHasOneInsert(t, repo)
	specs.UpdateHasOneUpdate(t, repo)
	specs.UpdateBelongsToInsert(t, repo)
	specs.UpdateBelongsToUpdate(t, repo)
	specs.UpdateAtomic(t, repo)
	specs.Updates(t, repo)
	specs.UpdateAny(t, repo)

	// Delete specs
	specs.Delete(t, repo)
	specs.DeleteBelongsTo(t, repo)
	specs.DeleteHasOne(t, repo)
	specs.DeleteHasMany(t, repo)
	specs.DeleteAll(t, repo)
	specs.DeleteAny(t, repo)

	// Constraint specs
	specs.UniqueConstraintOnInsert(t, repo)
	specs.UniqueConstraintOnUpdate(t, repo)
	specs.ForeignKeyConstraintOnInsert(t, repo)
	specs.ForeignKeyConstraintOnUpdate(t, repo)
	specs.CheckConstraintOnInsert(t, repo)
	specs.CheckConstraintOnUpdate(t, repo)
}

func TestAdapter_specs(t *testing.T) {
	if os.Getenv("TEST_PRIMARY_REPLICA") == "true" {
		t.Log("Skipping single node specs")
		return
	}

	adapter := MustOpen(dsn())
	defer adapter.Close()

	repo := rel.New(adapter)
	AdapterSpecs(t, repo)
}

func TestAdapter_PrimaryReplica_specs(t *testing.T) {
	if os.Getenv("TEST_PRIMARY_REPLICA") != "true" {
		t.Log("Skipping primary replica specs")
		return
	}

	adapter := primaryreplica.New(
		MustOpen("postgres://rel:rel@localhost:25432/rel_test?sslmode=disable&timezone=Asia/Jakarta"),
		MustOpen("postgres://rel:rel@localhost:25433/rel_test?sslmode=disable&timezone=Asia/Jakarta"),
	)

	defer adapter.Close()

	repo := rel.New(adapter)
	AdapterSpecs(t, repo)
}

func TestAdapter_Transaction_commitError(t *testing.T) {
	adapter := MustOpen(dsn())
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit(ctx))
}

func TestAdapter_Transaction_rollbackError(t *testing.T) {
	adapter := MustOpen(dsn())
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback(ctx))
}

func TestAdapter_Exec_error(t *testing.T) {
	adapter, err := Open(dsn())
	assert.Nil(t, err)
	defer adapter.Close()

	_, _, err = adapter.Exec(ctx, "error", nil)
	assert.NotNil(t, err)
}

func TestAdapter_TableBuilder(t *testing.T) {
	adapter := MustOpen(dsn())
	defer adapter.Close()

	tests := []struct {
		result string
		table  rel.Table
	}{
		{
			result: `ALTER TABLE "table" DROP CONSTRAINT "key";`,
			table: rel.Table{
				Op:   rel.SchemaAlter,
				Name: "table",
				Definitions: []rel.TableDefinition{
					rel.Key{Op: rel.SchemaDrop, Name: "key", Type: rel.ForeignKey},
				},
			},
		},
		{
			result: `ALTER TABLE "table" DROP CONSTRAINT "key";`,
			table: rel.Table{
				Op:   rel.SchemaAlter,
				Name: "table",
				Definitions: []rel.TableDefinition{
					rel.Key{Op: rel.SchemaDrop, Name: "key", Type: rel.UniqueKey},
				},
			},
		},
		{
			result: `ALTER TABLE "table" DROP CONSTRAINT "key";`,
			table: rel.Table{
				Op:   rel.SchemaAlter,
				Name: "table",
				Definitions: []rel.TableDefinition{
					rel.Key{Op: rel.SchemaDrop, Name: "key", Type: rel.PrimaryKey},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.result, func(t *testing.T) {
			assert.Equal(t, test.result, adapter.(*Postgres).TableBuilder.Build(test.table))
		})
	}
}
