package pgctx

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

type postgresKey struct{}

var ErrInvalidPostgresContext = errors.New("invalid postgres context")

const (
	// PostgresDSNEnv is the env variable that should contain the Postgres connection string (Data Source Name or DSN).
	//
	// The DSN string should be in the format:
	//
	//	postgres://[user]:[password]@[host]:[port]/[database]?[...pgoptions]
	PostgresDSNEnv = "DSN"

	PingTimeout = 10 * time.Second
)

var ErrNoDSN = errors.New("missing DSN environment variable")

func NewContext(ctx context.Context, migrations *embed.FS) (context.Context, error) {
	// Get the URL connection string from the environment.
	dsn := os.Getenv(PostgresDSNEnv)
	if dsn == "" {
		return nil, ErrNoDSN
	}

	return NewContextWithOptions(ctx, migrations, pgdriver.WithDSN(dsn))
}

// NewContextWithOptions returns a new context with a shared PG connection.
func NewContextWithOptions(
	ctx context.Context, migrations *embed.FS, options ...pgdriver.Option,
) (context.Context, error) {
	// Open a connection to the database.
	sqldb := sql.OpenDB(pgdriver.NewConnector(options...))

	// Make a temporary assignation. If something goes wrong, it is unnecessary and misleading to assign a value
	// to the global variable.
	client := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	// Wait for connection to be established.
	start := time.Now()
	for err := client.PingContext(context.Background()); err != nil; err = client.PingContext(context.Background()) {
		if time.Since(start) > PingTimeout {
			return nil, fmt.Errorf("ping database: %w", err)
		}
	}

	if migrations != nil {
		// Apply migrations.
		mig := migrate.NewMigrations()

		if err := mig.Discover(migrations); err != nil {
			return nil, fmt.Errorf("discover mig: %w", err)
		}

		migrator := migrate.NewMigrator(client, mig)
		if err := migrator.Init(ctx); err != nil {
			return nil, fmt.Errorf("create migrator: %w", err)
		}

		if _, err := migrator.Migrate(ctx); err != nil {
			return nil, fmt.Errorf("apply mig: %w", err)
		}
	}

	ctxPG := context.WithValue(ctx, postgresKey{}, bun.IDB(client))
	// Close clients on context termination.
	context.AfterFunc(ctxPG, func() {
		_ = client.Close()
		_ = sqldb.Close()
	})

	// Use type wrap so the assertion works properly when we try to extract the value.
	return ctxPG, nil
}

// NewContextTX creates a new Postgres context where the virtual database is replaced by a transaction.
//
// The parent context MUST contain a virtual Postgres instance, either created with NewContext or NewContextTX.
//
// The returned cancel function MUST be called with the argument set to true for the transaction to be committed.
// Omitting this call will automatically roll back the whole transaction.
func NewContextTX(ctx context.Context, opts *sql.TxOptions) (context.Context, func(commit bool) error, error) {
	pg, err := Context(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("extract pg: %w", err)
	}

	tx, err := pg.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("begin tx: %w", err)
	}

	var done bool

	ctxTx, cancelFn := context.WithCancel(context.WithValue(ctx, postgresKey{}, bun.IDB(&tx)))
	context.AfterFunc(ctxTx, func() {
		if !done {
			// If context is canceled without calling the cancel function, abort.
			// If the cancel function was already called, this will return an error,
			// so we ignore it.
			_ = tx.Rollback()
		}
	})

	cancelFnAugmented := func(commit bool) error {
		defer cancelFn()

		if commit {
			done = true

			return tx.Commit()
		}

		return nil
	}

	return ctxTx, cancelFnAugmented, nil
}

// Context extracts the Postgres virtual database from the context.
func Context(ctx context.Context) (bun.IDB, error) {
	db, ok := ctx.Value(postgresKey{}).(bun.IDB)
	if !ok {
		return nil, fmt.Errorf(
			"(pgctx) extract pg: %w: got type %T, expected %T",
			ErrInvalidPostgresContext,
			ctx.Value(postgresKey{}), bun.IDB(nil),
		)
	}

	return db, nil
}
