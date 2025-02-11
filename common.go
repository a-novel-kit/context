package context

import (
	"context"
	"errors"
	"fmt"
)

var ErrUnsupportedContext = errors.New("unsupported context")

// ExtractValue extracts a value from a context with the correct typing, It fails if the value is missing or has
// the wrong type.
func ExtractValue[T any](ctx context.Context, key any) (T, error) {
	var zero T

	raw := ctx.Value(key)
	if raw == nil {
		return zero, fmt.Errorf("(ExtractContext) %w: missing key %s", ErrUnsupportedContext, key)
	}

	value, ok := raw.(T)
	if !ok {
		return zero, fmt.Errorf("(ExtractContext) %w: invalid type for key %s", ErrUnsupportedContext, key)
	}

	return value, nil
}
