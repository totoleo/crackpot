package localcache

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/libp2p/go-buffer-pool"
)

const timeMaxLen = 2
const timeFilter = 0xffff

type fastcacheVariant[K comparable, T any] struct {
	local *fastcache.Cache

	loader Storage[K, T]

	Now      func() int64
	lifetime uint16
}

func WithFastCache[K comparable, T any](cap int, lifetime uint64, loader Storage[K, T]) Cache[K, T] {
	return &fastcacheVariant[K, T]{
		local:    fastcache.New(cap),
		lifetime: uint16(lifetime),
		loader:   loader,
		Now: func() int64 {
			return time.Now().Unix()
		},
	}
}

func (f *fastcacheVariant[K, T]) MGet(ctx context.Context, ids []K) ([]T, error) {
	var items = make([]T, 0, len(ids))
	var (
		missing []K
		err     error
	)
	items, missing, err = f.fromLocal(ctx, items, ids)
	if err != nil {
		return nil, err
	}
	if len(missing) == 0 {
		return items, nil
	}
	items, err = f.loadMore(ctx, items, missing)
	if err != nil {
		return items, err
	}
	return items, nil
}

func (f *fastcacheVariant[K, T]) fromLocal(ctx context.Context, items []T, ids []K) ([]T, []K, error) {
	var missing = make([]K, 0, len(ids))
	now := uint16(f.Now() % timeFilter)
	buf := f.getBuff(4096)
	for _, id := range ids {
		key := f.loader.KeyFactory(id)
		buff := f.local.Get(buf[:0], key)
		f.loader.KeyRelease(key)
		if len(buff) == 0 {
			missing = append(missing, id)
			continue
		}
		duration := now - binary.BigEndian.Uint16(buff[:timeMaxLen])
		if duration > f.lifetime {
			missing = append(missing, id)
			continue
		}
		val, err := f.loader.Unmarshal(buff[timeMaxLen:])
		if err != nil {
			return nil, nil, fmt.Errorf("cache err:%w", err)
		}
		items = append(items, val)
	}
	return items, missing, nil
}

func (f *fastcacheVariant[K, T]) loadMore(ctx context.Context, items []T, missing []K) ([]T, error) {

	fromLoader, err := f.loader.Load(ctx, missing)
	if err != nil {
		return nil, fmt.Errorf("cache] error from loader:%w", err)
	}
	now := uint64(f.Now())

	timeBuff := f.getBuff(timeMaxLen)
	binary.BigEndian.AppendUint16(timeBuff[:0], uint16(now%timeFilter))
	for _, id := range missing {
		content, ok := fromLoader[id]
		if !ok {
			//TODO:: not found?
			continue
		}
		buff := f.getBuff(len(content) + timeMaxLen)

		copy(buff, timeBuff)             //cache time
		copy(buff[timeMaxLen:], content) //cache content

		key := f.loader.KeyFactory(id)
		f.local.Set(key, buff)
		f.loader.KeyRelease(key)
		item, err := f.loader.Unmarshal(content)
		if err != nil {
			//TODO:: return?
			return nil, err
		}
		items = append(items, item)
		f.releaseBuff(buff)
	}
	return items, nil
}

func (f *fastcacheVariant[K, T]) getBuff(l int) []byte {
	return pool.Get(l)
}

func (f *fastcacheVariant[K, T]) releaseBuff(slice []byte) {
	pool.Put(slice)
}
