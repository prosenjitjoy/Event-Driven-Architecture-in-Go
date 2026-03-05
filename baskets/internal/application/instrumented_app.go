package application

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type instrumentedApp struct {
	App
	basketsStarted    metric.Int64Counter
	basketsCheckedOut metric.Int64Counter
	basketsCanceled   metric.Int64Counter
}

var _ App = (*instrumentedApp)(nil)

func NewInstrumentedApp(app App, basketsStarted, basketsCheckedOut, basketsCanceled metric.Int64Counter) App {
	return instrumentedApp{
		App:               app,
		basketsStarted:    basketsStarted,
		basketsCheckedOut: basketsCheckedOut,
		basketsCanceled:   basketsCanceled,
	}
}

func (a instrumentedApp) StartBasket(ctx context.Context, start StartBasket) error {
	if err := a.App.StartBasket(ctx, start); err != nil {
		return err
	}

	a.basketsStarted.Add(ctx, 1, metric.WithAttributes(attribute.String("type", "StartBasket")))

	return nil
}

func (a instrumentedApp) CancelBasket(ctx context.Context, cancel CancelBasket) error {
	if err := a.App.CancelBasket(ctx, cancel); err != nil {
		return err
	}

	a.basketsCanceled.Add(ctx, 1, metric.WithAttributes(attribute.String("type", "CancelBasket")))

	return nil
}

func (a instrumentedApp) CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error {
	if err := a.App.CheckoutBasket(ctx, checkout); err != nil {
		return err
	}

	a.basketsCheckedOut.Add(ctx, 1, metric.WithAttributes(attribute.String("type", "CheckoutBasket")))

	return nil
}
