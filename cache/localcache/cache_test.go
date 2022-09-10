package localcache

import (
	"context"
	"encoding/binary"
	"math"
	"testing"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/bytedance/sonic"
	pool "github.com/libp2p/go-buffer-pool"
)

func TestTime(t *testing.T) {
	now := time.Now().Unix()
	t.Log(now, now%0xffff)
	b := binary.BigEndian.AppendUint16(nil, uint16(now%0xffff))
	t.Log(b, len(b))
}

func BenchmarkLoad(b *testing.B) {
	mock, _ := sonic.Marshal(&foo{Id: math.MaxInt64, Name: "你好, 世界"})
	fc := WithFastCache[int64, *foo](1, 10, &mockLoader{
		itemTemplate: mock,
	})
	ids := make([]int64, 0, 20)
	for i := 0; i < 20; i++ {
		ids = append(ids, fastrand.Int63())
	}
	for i := 0; i < b.N; i++ {
		_, err := fc.MGet(context.Background(), ids)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

type mockLoader struct {
	itemTemplate []byte
}

func (m *mockLoader) Load(ctx context.Context, ids []int64) (map[int64][]byte, error) {
	results := make(map[int64][]byte, len(ids))
	for _, id := range ids {
		v, _ := sonic.Marshal(&foo{
			Id:     id,
			Name:   "你好，世界",
			Value1: 1,
			Value2: 2,
			Value3: []*Nest1{{Key: "abc"}},
		})
		results[id] = v
	}
	return results, nil
}

func (m *mockLoader) KeyFactory(id int64) []byte {
	buff := pool.Get(16)
	binary.BigEndian.AppendUint64(buff[:0], uint64(id))
	return buff
}

func (m *mockLoader) KeyRelease(bytes []byte) {
	pool.Put(bytes)
}

func (m *mockLoader) Unmarshal(bytes []byte) (*foo, error) {
	var v = &foo{}
	err := sonic.Unmarshal(bytes, v)
	return v, err
}

type Nest1 struct {
	Key    string
	Height int64
	Width  int64
}

type foo struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Value1 int64
	Value2 int64
	Value3 []*Nest1
}
