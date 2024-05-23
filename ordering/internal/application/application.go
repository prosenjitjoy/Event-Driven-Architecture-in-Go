package application

import (
	"context"
	"fmt"

	"mall/internal/ddd"
	"mall/ordering/internal/domain"
)

type CreateOrder struct {
	ID         string
	CustomerID string
	PaymentID  string
	Items      []*domain.Item
}

type CancelOrder struct {
	ID string
}

type ReadyOrder struct {
	ID string
}

type CompleteOrder struct {
	ID        string
	InvoiceID string
}

type GetOrder struct {
	ID string
}

type App interface {
	CreateOrder(ctx context.Context, cmd CreateOrder) error
	CancelOrder(ctx context.Context, cmd CancelOrder) error
	ReadyOrder(ctx context.Context, cmd ReadyOrder) error
	CompleteOrder(ctx context.Context, cmd CompleteOrder) error
	GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error)
}

type Application struct {
	orders          domain.OrderRepository
	domainPublisher ddd.EventPublisher
}

var _ App = (*Application)(nil)

func New(orders domain.OrderRepository, domainPublisher ddd.EventPublisher) *Application {
	return &Application{
		orders:          orders,
		domainPublisher: domainPublisher,
	}
}

func (a Application) CreateOrder(ctx context.Context, cmd CreateOrder) error {
	order, err := domain.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return fmt.Errorf("create order command: %w", err)
	}

	// publish domain events
	if err = a.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return err
	}

	err = a.orders.Save(ctx, order)
	if err != nil {
		return fmt.Errorf("create order command: %w", err)
	}

	return nil
}

func (a Application) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := a.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Cancel(); err != nil {
		return err
	}

	// publish domain events
	if err = a.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return err
	}

	if err = a.orders.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (a Application) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	order, err := a.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Ready(); err != nil {
		return nil
	}

	// publish domain events
	if err = a.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return err
	}

	if err = a.orders.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (a Application) CompleteOrder(ctx context.Context, cmd CompleteOrder) error {
	order, err := a.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = order.Complete(cmd.InvoiceID)
	if err != nil {
		return nil
	}

	return a.orders.Update(ctx, order)
}

func (a Application) GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error) {
	order, err := a.orders.Find(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("get order query: %w", err)
	}

	return order, nil
}
