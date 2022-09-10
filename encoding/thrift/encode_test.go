package thrift

import (
	"math/rand"
	"testing"
	"time"
)
import "github.com/cloudwego/frugal"

type frugalStruct struct {
	Id    int64   `frugal:"1,default,i64" json:"id,string"`
	Name  string  `frugal:"2,default,string"`
	Name2 *string `frugal:"3,optional,string"`
	Name3 string  `frugal:"4,default,string"`
}

func TestEncode(t *testing.T) {
	rand.Seed(int64(time.Now().Nanosecond()))
	val := &frugalStruct{Id: rand.Int63(), Name: "2"}
	size := frugal.EncodedSize(val)
	buff := make([]byte, size)
	n, err := frugal.EncodeObject(buff, nil, val)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("size", size, "output", []byte(buff), "written", n)
}
