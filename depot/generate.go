package depot

//go:generate buf generate

//go:generate mockgen -destination=depotpb/mock.go -package=depotpb mall/depot/depotpb DepotServiceClient,DepotServiceServer

//go:generate mockgen -destination=internal/application/mock.go -package=application mall/depot/internal/application App

//go:generate mockgen -destination=internal/domain/mock.go -package=domain mall/depot/internal/domain OrderRepository,ProductCacheRepository,ProductRepository,ShoppingListRepository,StoreCacheRepository,StoreRepository
