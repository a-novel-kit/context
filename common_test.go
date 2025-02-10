package context_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/context"
)

func TestExtractValue(t *testing.T) {
	t.Parallel()

	value := "bar"
	ctx := context.WithValue(context.Background(), context.CtxKey("foo"), value)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		extracted, err := context.ExtractValue[string](ctx, "foo")
		require.NoError(t, err)
		require.Equal(t, value, extracted)
	})

	t.Run("WrongType", func(t *testing.T) {
		t.Parallel()

		_, err := context.ExtractValue[int](ctx, "foo")
		require.ErrorIs(t, err, context.ErrUnsupportedContext)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		_, err := context.ExtractValue[string](ctx, "bar")
		require.ErrorIs(t, err, context.ErrUnsupportedContext)
	})
}
