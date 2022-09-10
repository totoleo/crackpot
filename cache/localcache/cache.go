package localcache

import "context"

type Cache[K comparable, T any] interface {
	MGet(ctx context.Context, ids []K) ([]T, error)
}

type Storage[K comparable, T any] interface {
	Load(ctx context.Context, ids []K) (map[K][]byte, error)
	KeyFactory(id K) []byte

	// KeyRelease 如果希望使用 sync.Pool 可以在该方法进行归还操作
	KeyRelease([]byte)

	// Unmarshal  bytes 在多次调用中可能会被复用，请避免使用 unsafe 方法
	Unmarshal(bytes []byte) (T, error)
}
