package common

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/network/protocol"
)

var (
	rnLoggerFormat = "remote %s"
)

type RemoteNode struct {
	Node
	conn          net.Conn
	sessionKey    []byte
	privateIP     net.IP
	publicAddress string
	logger        log.Logger
	lastHeartbeat time.Time
}

func NewRemoteNode(conn net.Conn, sessionKey []byte, privateIP net.IP) *RemoteNode {
	return &RemoteNode{
		conn:          conn,
		sessionKey:    sessionKey,
		privateIP:     privateIP,
		publicAddress: conn.RemoteAddr().String(),
		logger:        log.L(fmt.Sprintf(rnLoggerFormat, privateIP.String())),
		lastHeartbeat: time.Now(),
	}
}

func (rn *RemoteNode) SendToInterface(payload []byte) error {
	return protocol.WriteEncodeTransfer(rn.conn, payload)
}

func (rn *RemoteNode) SendPack(pack *protocol.Packet) (err error) {
	if err = protocol.EncodeAndWrite(rn.conn, pack); err != nil {
		err = fmt.Errorf("error on write transfer message, %v", err)
	}
	return
}

func (rn *RemoteNode) Close() {
	defer rn.conn.Close()
	rn.logger.Debug("closing...")
}

func (rn *RemoteNode) listen(ln *LocalNode) {
	defer rn.logger.Debug("listener stopped...")
	defer func() {
		ln.NetTable().RemoveRemoteNode(rn.privateIP)
	}()

	rn.logger.Debug("listening...")

	for {
		pack, err := protocol.Decode(rn.conn)
		if err != nil {
			rn.logger.Error("decode error, %v", err)
			if err == io.EOF {
				break
			}
			continue
		}
		rn.logger.Debug("received, %+v", pack)

		switch pack.Data.Type {
		case protocol.TypeTransfer:
			rn.logger.Debug("Writing to interface...")

		case protocol.TypeHeartbeat:
			rn.logger.Debug("heardbeat received, %v", pack.Data.Msg)
			rn.lastHeartbeat = time.Now()
		}
	}
}
