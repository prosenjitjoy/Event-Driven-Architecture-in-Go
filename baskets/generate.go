package baskets

//go:generate buf generate

//go:generate mockgen -destination=basketspb/mock.go -package=basketspb mall/baskets/basketspb BasketServiceClient,BasketServiceServer

//go:generate mockgen -destination=internal/application/mock.go -package=application mall/baskets/internal/application App

//go:generate mockgen -destination=internal/domain/mock.go -package=domain mall/baskets/internal/domain BasketRepository,ProductCacheRepository,ProductRepository,StoreCacheRepository,StoreRepository
