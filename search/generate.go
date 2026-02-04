package search

//go:generate buf generate

//go:generate mockgen -destination=searchpb/mock.go -package=searchpb mall/search/searchpb SearchServiceClient,SearchServiceServer

//go:generate mockgen -destination=internal/application/mock.go -package=application mall/search/internal/application Application

//go:generate mockgen -destination=internal/domain/mock.go -package=domain mall/search/internal/domain CustomerCacheRepository,CustomerRepository,OrderRepository,ProductCacheRepository,ProductRepository,StoreCacheRepository,StoreRepository
