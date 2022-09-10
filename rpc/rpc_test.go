package main

import (
	"testing"

	at "github.com/apache/thrift/lib/go/thrift"

	"github.com/cloudwego/kitex/pkg/remote/codec/thrift"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/invoke"

	"github.com/totoleo/crackpot/rpc/kitex_gen/a/b/c"
	"github.com/totoleo/crackpot/rpc/kitex_gen/a/b/c/echo"
)

func TestCall(t *testing.T) {
	code := thrift.NewThriftCodecWithConfig(thrift.FastRead | thrift.FastWrite|thrift.FrugalRead)
	v := echo.NewInvoker(new(EchoImpl), server.WithPayloadCodec(code))

	codec:=utils.NewThriftMessageCodec()
	codec.Encode("xxx",at.CALL,1,&c.EchoHelloArgs{})
	message := invoke.NewMessage(nil, nil)
	err := code.Marshal()
	message.SetRequestBytes(err)
	v.Call(message)
}
