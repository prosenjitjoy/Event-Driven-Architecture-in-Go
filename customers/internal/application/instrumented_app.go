package application

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type instrumentedApp struct {
	App
	customersRegistered metric.Int64Counter
}

var _ App = (*instrumentedApp)(nil)

func NewInstrumentedApp(app App, customersRegistered metric.Int64Counter) App {
	return instrumentedApp{
		App:                 app,
		customersRegistered: customersRegistered,
	}
}

func (a instrumentedApp) RegisterCustomer(ctx context.Context, register RegisterCustomer) error {
	if err := a.App.RegisterCustomer(ctx, register); err != nil {
		return err
	}

	a.customersRegistered.Add(ctx, 1, metric.WithAttributes(attribute.String("type", "RegisterCustomer")))

	return nil
}
