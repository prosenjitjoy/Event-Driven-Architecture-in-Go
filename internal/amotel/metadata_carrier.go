package amotel

import (
	"fmt"
	"mall/internal/ddd"
	"strconv"

	"go.opentelemetry.io/otel/propagation"
)

type MetadataCarrier ddd.Metadata

var _ propagation.TextMapCarrier = (*MetadataCarrier)(nil)

func (mc MetadataCarrier) Get(key string) string {
	metadata := ddd.Metadata(mc)

	switch v := metadata.Get(key).(type) {
	case nil:
		return ""
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (mc MetadataCarrier) Set(key, value string) {
	metadata := ddd.Metadata(mc)

	metadata.Set(key, value)
}

func (mc MetadataCarrier) Keys() []string {
	metadata := ddd.Metadata(mc)

	return metadata.Keys()
}
