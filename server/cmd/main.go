// $ protoc -I=shared --go_out=server shared/packets.proto
package main

import (
	"fmt"
	"server/pkg/packets"

	"google.golang.org/protobuf/proto"
)

func main() {
	packet := &packets.Packet{
		SenderId: 69,
		Msg:      packets.NewChat("Hello"),
	}
	fmt.Println(packet)
	data, err := proto.Marshal(packet)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
