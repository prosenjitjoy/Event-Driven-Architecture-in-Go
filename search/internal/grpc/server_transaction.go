package grpc

import (
	"context"
	"database/sql"
	"mall/internal/di"
	"mall/search/internal/application"
	"mall/search/internal/constants"
	"mall/search/searchpb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	searchpb.UnimplementedSearchServiceServer
}

var _ searchpb.SearchServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	searchpb.RegisterSearchServiceServer(registrar, serverTx{c: container})

	return nil
}

func (s serverTx) SearchOrders(ctx context.Context, request *searchpb.SearchOrdersRequest) (resp *searchpb.SearchOrdersResponse, err error) {
	ctx = s.c.Scoped(ctx)

	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, constants.DatabaseTransactionKey).(application.Application)}

	return next.SearchOrders(ctx, request)
}

func (s serverTx) GetOrder(ctx context.Context, request *searchpb.GetOrderRequest) (resp *searchpb.GetOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)

	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, constants.DatabaseTransactionKey).(application.Application)}

	return next.GetOrder(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if err != nil {
		tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
