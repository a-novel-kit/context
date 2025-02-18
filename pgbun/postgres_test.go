package pgctx_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pgctx "github.com/a-novel-kit/context/pgbun"
	"github.com/a-novel-kit/context/pgbun/test/migrations"
)

func countPG(ctx context.Context, t *testing.T) int {
	t.Helper()

	db, err := pgctx.Context(ctx)
	require.NoError(t, err)
	require.NotNil(t, db)

	row := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test;")

	var res int

	require.NoError(t, row.Scan(&res))
	require.NoError(t, row.Err())

	return res
}

func TestPGContextOK(t *testing.T) {
	ctx, err := pgctx.NewContext(t.Context(), nil)
	require.NoError(t, err)

	db, err := pgctx.Context(ctx)
	require.NoError(t, err)
	require.NotNil(t, db)

	require.NoError(t, db.QueryRowContext(ctx, "SELECT 1;").Err())
}

func TestPGContextOKMigrations(t *testing.T) {
	ctx, err := pgctx.NewContext(t.Context(), &migrations.Migrations)
	require.NoError(t, err)

	require.Equal(t, 0, countPG(ctx, t))
}

func TestPGContextTransactionRollbackExplicitly(t *testing.T) {
	ctx, err := pgctx.NewContext(t.Context(), &migrations.Migrations)
	require.NoError(t, err)

	tx, cancel, err := pgctx.NewContextTX(ctx, nil)
	require.NoError(t, err)

	db, err := pgctx.Context(tx)
	require.NoError(t, err)

	t.Log("InsertInTX")
	require.Equal(t, 0, countPG(tx, t))

	_, err = db.ExecContext(tx, "INSERT INTO test (id, content) VALUES (?, ?);", uuid.New(), "foobarqux")
	require.NoError(t, err)

	require.Equal(t, 1, countPG(tx, t))

	t.Log("InvisibleFromParent")
	require.Equal(t, 0, countPG(ctx, t))

	t.Log("Rollback")
	require.NoError(t, cancel(false))
	require.Equal(t, 0, countPG(ctx, t))
}

func TestPGContextTransactionRollbackAuto(t *testing.T) {
	ctx, err := pgctx.NewContext(t.Context(), &migrations.Migrations)
	require.NoError(t, err)

	parentCTX, parentCancel := context.WithCancel(ctx)

	tx, _, err := pgctx.NewContextTX(parentCTX, nil)
	require.NoError(t, err)

	db, err := pgctx.Context(tx)
	require.NoError(t, err)

	t.Log("InsertInTX")
	require.Equal(t, 0, countPG(tx, t))

	_, err = db.ExecContext(tx, "INSERT INTO test (id, content) VALUES (?, ?);", uuid.New(), "foobarqux")
	require.NoError(t, err)

	require.Equal(t, 1, countPG(tx, t))

	t.Log("InvisibleFromParent")
	require.Equal(t, 0, countPG(parentCTX, t))

	t.Log("Rollback")
	parentCancel()
	require.Equal(t, 0, countPG(ctx, t))
}

func TestPGContextTransactionCommit(t *testing.T) {
	ctx, err := pgctx.NewContext(t.Context(), &migrations.Migrations)
	require.NoError(t, err)

	tx, cancel, err := pgctx.NewContextTX(ctx, nil)
	require.NoError(t, err)

	db, err := pgctx.Context(tx)
	require.NoError(t, err)

	t.Log("InsertInTX")
	require.Equal(t, 0, countPG(tx, t))

	_, err = db.ExecContext(tx, "INSERT INTO test (id, content) VALUES (?, ?);", uuid.New(), "foobarqux")
	require.NoError(t, err)

	require.Equal(t, 1, countPG(tx, t))

	t.Log("InvisibleFromParent")
	require.Equal(t, 0, countPG(ctx, t))

	t.Log("Commit")
	require.NoError(t, cancel(true))
	require.Equal(t, 1, countPG(ctx, t))

	db, err = pgctx.Context(ctx)
	require.NoError(t, err)
	res, err := db.ExecContext(ctx, "DELETE FROM test;")
	require.NoError(t, err)

	rowsAffected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)
}

func TestNoDSN(t *testing.T) {
	t.Setenv(pgctx.PostgresDSNEnv, "")
	_, err := pgctx.NewContext(t.Context(), nil)
	require.ErrorIs(t, err, pgctx.ErrNoDSN)
}

func TestBadDSN(t *testing.T) {
	t.Setenv(pgctx.PostgresDSNEnv, "postgres://test:test@localhost:1111/test?sslmode=disable")
	_, err := pgctx.NewContext(t.Context(), nil)
	require.Error(t, err)
}
