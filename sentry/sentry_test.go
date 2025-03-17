package sentryctx_test

import (
	"testing"
	"time"

	"github.com/getsentry/sentry-go"

	sentryctx "github.com/a-novel-kit/context/sentry"
)

func TestSentryUnconfigured(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	// Make sure we can call every function without panicking.
	sentryctx.LastEventID(ctx)
	sentryctx.Clone(ctx)
	sentryctx.Client(ctx)
	sentryctx.Scope(ctx).SetExtra("key", "value")
	sentryctx.PushScope(ctx)
	sentryctx.PopScope(ctx)
	sentryctx.BindClient(ctx, nil)
	sentryctx.WithScope(ctx, func(scope *sentry.Scope) {
		scope.SetExtra("key", "value")
	})
	sentryctx.ConfigureScope(ctx, func(scope *sentry.Scope) {
		scope.SetExtra("key", "value")
	})
	sentryctx.CaptureException(ctx, nil)
	sentryctx.CaptureMessage(ctx, "message")
	sentryctx.CaptureEvent(ctx, &sentry.Event{})
	sentryctx.CaptureCheckIn(ctx, nil, nil)
	sentryctx.AddBreadcrumb(ctx, &sentry.Breadcrumb{}, nil)
	sentryctx.Recover(ctx, nil)
	sentryctx.RecoverWithContext(ctx, nil)
	sentryctx.Flush(ctx, time.Second)
	sentryctx.GetTraceparent(ctx)
	sentryctx.GetBaggage(ctx)
}

func TestSentryConfigured(t *testing.T) {
	t.Parallel()

	cloned := sentry.CurrentHub().Clone()
	ctx := sentry.SetHubOnContext(t.Context(), cloned)

	// Make sure we can call every function without panicking.
	sentryctx.LastEventID(ctx)
	sentryctx.Clone(ctx)
	sentryctx.Client(ctx)
	sentryctx.Scope(ctx).SetExtra("key", "value")
	sentryctx.PushScope(ctx)
	sentryctx.PopScope(ctx)
	sentryctx.BindClient(ctx, nil)
	sentryctx.WithScope(ctx, func(scope *sentry.Scope) {
		scope.SetExtra("key", "value")
	})
	sentryctx.ConfigureScope(ctx, func(scope *sentry.Scope) {
		scope.SetExtra("key", "value")
	})
	sentryctx.CaptureException(ctx, nil)
	sentryctx.CaptureMessage(ctx, "message")
	sentryctx.CaptureEvent(ctx, &sentry.Event{})
	sentryctx.CaptureCheckIn(ctx, nil, nil)
	sentryctx.AddBreadcrumb(ctx, &sentry.Breadcrumb{}, nil)
	sentryctx.Recover(ctx, nil)
	sentryctx.RecoverWithContext(ctx, nil)
	sentryctx.Flush(ctx, time.Second)
	sentryctx.GetTraceparent(ctx)
	sentryctx.GetBaggage(ctx)
}
