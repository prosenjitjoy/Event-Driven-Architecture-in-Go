package di

import (
	"context"
	"log"
)

func Get(ctx context.Context, key string) any {
	ctn, ok := ctx.Value(containerKey).(*container)
	if !ok {
		log.Panic("container does not exist on context")
	}

	return ctn.Get(key)
}
