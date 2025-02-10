package pgctx_test

import (
	"github.com/a-novel-kit/context"
	pgctx "github.com/a-novel-kit/context/pgbun"
	"github.com/a-novel-kit/context/pgbun/test/migrations"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func countPG(t *testing.T, ctx context.Context) int {
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
	ctx, err := pgctx.NewContext(context.Background(), nil)
	require.NoError(t, err)

	db, err := pgctx.Context(ctx)
	require.NoError(t, err)
	require.NotNil(t, db)

	require.NoError(t, db.QueryRowContext(ctx, "SELECT 1;").Err())
}

func TestPGContextOKMigrations(t *testing.T) {
	ctx, err := pgctx.NewContext(context.Background(), &migrations.Migrations)
	require.NoError(t, err)

	require.Equal(t, 0, countPG(t, ctx))
}

func TestPGContextTransactionRollbackExplicitly(t *testing.T) {
	ctx, err := pgctx.NewContext(context.Background(), &migrations.Migrations)
	require.NoError(t, err)

	tx, cancel, err := pgctx.NewContextTX(ctx, nil)
	require.NoError(t, err)

	db, err := pgctx.Context(tx)
	require.NoError(t, err)

	t.Run("InsertInTX", func(t *testing.T) {
		require.Equal(t, 0, countPG(t, tx))

		_, err = db.ExecContext(tx, "INSERT INTO test (id, content) VALUES (?, ?);", uuid.New(), "foobarqux")
		require.NoError(t, err)

		require.Equal(t, 1, countPG(t, tx))
	})

	t.Run("InvisibleFromParent", func(t *testing.T) {
		require.Equal(t, 0, countPG(t, ctx))
	})

	t.Run("Rollback", func(t *testing.T) {
		require.NoError(t, cancel(false))
		require.Equal(t, 0, countPG(t, ctx))
	})
}

func TestPGContextTransactionRollbackAuto(t *testing.T) {
	ctx, err := pgctx.NewContext(context.Background(), &migrations.Migrations)
	require.NoError(t, err)

	parentCTX, parentCancel := context.WithCancel(ctx)

	tx, _, err := pgctx.NewContextTX(parentCTX, nil)
	require.NoError(t, err)

	db, err := pgctx.Context(tx)
	require.NoError(t, err)

	t.Run("InsertInTX", func(t *testing.T) {
		require.Equal(t, 0, countPG(t, tx))

		_, err = db.ExecContext(tx, "INSERT INTO test (id, content) VALUES (?, ?);", uuid.New(), "foobarqux")
		require.NoError(t, err)

		require.Equal(t, 1, countPG(t, tx))
	})

	t.Run("InvisibleFromParent", func(t *testing.T) {
		require.Equal(t, 0, countPG(t, parentCTX))
	})

	t.Run("Rollback", func(t *testing.T) {
		parentCancel()
		require.Equal(t, 0, countPG(t, ctx))
	})
}

func TestPGContextTransactionCommit(t *testing.T) {
	ctx, err := pgctx.NewContext(context.Background(), &migrations.Migrations)
	require.NoError(t, err)

	tx, cancel, err := pgctx.NewContextTX(ctx, nil)
	require.NoError(t, err)

	db, err := pgctx.Context(tx)
	require.NoError(t, err)

	t.Run("InsertInTX", func(t *testing.T) {
		require.Equal(t, 0, countPG(t, tx))

		_, err = db.ExecContext(tx, "INSERT INTO test (id, content) VALUES (?, ?);", uuid.New(), "foobarqux")
		require.NoError(t, err)

		require.Equal(t, 1, countPG(t, tx))
	})

	t.Run("InvisibleFromParent", func(t *testing.T) {
		require.Equal(t, 0, countPG(t, ctx))
	})

	t.Run("Commit", func(t *testing.T) {
		require.NoError(t, cancel(true))
		require.Equal(t, 1, countPG(t, ctx))

		db, err = pgctx.Context(ctx)
		res, err := db.ExecContext(ctx, "DELETE FROM test;")
		require.NoError(t, err)

		rowsAffected, err := res.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAffected)
	})
}
