package sctp

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/pion/logging"
	"github.com/stretchr/testify/assert"
)

func makePaths() []*Association{
	paths := make([]*Association,5)
	for i:= 0; i < 5; i ++{
		addr := net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 9000+i,
		}

		conn, err := net.ListenUDP("udp", &addr)
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()
	
		config := Config{
			NetConn:       conn,
			LoggerFactory: logging.NewDefaultLoggerFactory(),
		}
		paths[i] = createAssociation(config)
	}
	return paths
}

func TestSenderBufferCreation(t *testing.T){
	t.Run("Creating buffer", func (t *testing.T)  {
		paths := makePaths()
		senderBuff := newSenderBuffer(paths)
		fmt.Print(senderBuff)
		assert.NotNil(t,senderBuff,"should not be nil")
	})
}


func TestSenderBufferPushPop(t *testing.T){
	t.Run("Push and Pop",func (t *testing.T)  {
		paths := makePaths()
		senderBuff := newSenderBuffer(paths)
		senderBuff.bufferSize = 30000
		senderBuff.push(makeDataChunk(0, false, noFragment))
		senderBuff.push(makeDataChunk(1, false, noFragment))
		senderBuff.push(makeDataChunk(2, false, noFragment))
		fmt.Print(senderBuff)
		fmt.Print(senderBuff.pop(0))
		fmt.Print(senderBuff)
	})
}
