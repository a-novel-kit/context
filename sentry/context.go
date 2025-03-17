package sentryctx

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

func LastEventID(ctx context.Context) sentry.EventID {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.LastEventID()
	}

	return sentry.CurrentHub().LastEventID()
}

func Clone(ctx context.Context) *sentry.Hub {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.Clone()
	}

	return sentry.CurrentHub().Clone()
}

func Client(ctx context.Context) *sentry.Client {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.Client()
	}

	return sentry.CurrentHub().Client()
}

func Scope(ctx context.Context) *sentry.Scope {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.Scope()
	}

	return sentry.CurrentHub().Scope()
}

func PushScope(ctx context.Context) *sentry.Scope {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.PushScope()
	}

	return sentry.CurrentHub().PushScope()
}

func PopScope(ctx context.Context) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.PopScope()

		return
	}

	sentry.CurrentHub().PopScope()
}

func BindClient(ctx context.Context, client *sentry.Client) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.BindClient(client)

		return
	}

	sentry.CurrentHub().BindClient(client)
}

func WithScope(ctx context.Context, f func(scope *sentry.Scope)) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(f)

		return
	}

	sentry.CurrentHub().WithScope(f)
}

func ConfigureScope(ctx context.Context, f func(scope *sentry.Scope)) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.ConfigureScope(f)

		return
	}

	sentry.CurrentHub().ConfigureScope(f)
}

func CaptureEvent(ctx context.Context, event *sentry.Event) *sentry.EventID {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.CaptureEvent(event)
	}

	return sentry.CurrentHub().CaptureEvent(event)
}

func CaptureMessage(ctx context.Context, message string) *sentry.EventID {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.CaptureMessage(message)
	}

	return sentry.CurrentHub().CaptureMessage(message)
}

func CaptureException(ctx context.Context, exception error) *sentry.EventID {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.CaptureException(exception)
	}

	return sentry.CurrentHub().CaptureException(exception)
}

func CaptureCheckIn(ctx context.Context, checkIn *sentry.CheckIn, monitorConfig *sentry.MonitorConfig) *sentry.EventID {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.CaptureCheckIn(checkIn, monitorConfig)
	}

	return sentry.CurrentHub().CaptureCheckIn(checkIn, monitorConfig)
}

func AddBreadcrumb(ctx context.Context, breadcrumb *sentry.Breadcrumb, hint *sentry.BreadcrumbHint) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.AddBreadcrumb(breadcrumb, hint)

		return
	}

	sentry.CurrentHub().AddBreadcrumb(breadcrumb, hint)
}

func Recover(ctx context.Context, err interface{}) *sentry.EventID {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.Recover(err)
	}

	return sentry.CurrentHub().Recover(err)
}

func RecoverWithContext(ctx context.Context, err interface{}) *sentry.EventID {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.RecoverWithContext(ctx, err)
	}

	return sentry.CurrentHub().RecoverWithContext(ctx, err)
}

func Flush(ctx context.Context, timeout time.Duration) bool {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.Flush(timeout)
	}

	return sentry.CurrentHub().Flush(timeout)
}

func GetTraceparent(ctx context.Context) string {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.GetTraceparent()
	}

	return sentry.CurrentHub().GetTraceparent()
}

func GetBaggage(ctx context.Context) string {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.GetBaggage()
	}

	return sentry.CurrentHub().GetBaggage()
}
