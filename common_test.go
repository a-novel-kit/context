package context_test

import (
	"github.com/a-novel-kit/context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractValue(t *testing.T) {
	value := "bar"
	ctx := context.WithValue(context.Background(), context.CtxKey("foo"), value)

	t.Run("OK", func(t *testing.T) {
		extracted, err := context.ExtractValue[string](ctx, "foo")
		require.NoError(t, err)
		require.Equal(t, value, extracted)
	})

	t.Run("WrongType", func(t *testing.T) {
		_, err := context.ExtractValue[int](ctx, "foo")
		require.ErrorIs(t, err, context.ErrUnsupportedContext)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := context.ExtractValue[string](ctx, "bar")
		require.ErrorIs(t, err, context.ErrUnsupportedContext)
	})
}
