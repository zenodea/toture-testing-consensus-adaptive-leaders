package common

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
	"toture-test/util"
)

type Network struct {
	Id                  int
	ListenAddress       string
	RemoteAddresses     map[int]string
	IncomingConnections map[int]*bufio.Reader
	OutgoingConnections map[int]*bufio.Writer
	OutMutex            map[int]*sync.Mutex
	OutChan             chan *RPCPairPeer
	logger              *util.Logger
	rpcTable            map[uint8]*RPCPair // map each RPC type (message type) to its unique number
}

type NetworkConfig struct {
	ListenAddress   string
	RemoteAddresses map[int]string // to connect to
}

func NewNetwork(Id int, config *NetworkConfig, outChan chan *RPCPairPeer, logger *util.Logger) *Network {
	n := &Network{
		Id:                  Id,
		ListenAddress:       config.ListenAddress,
		RemoteAddresses:     config.RemoteAddresses,
		IncomingConnections: make(map[int]*bufio.Reader),
		OutgoingConnections: make(map[int]*bufio.Writer),
		OutMutex:            make(map[int]*sync.Mutex),
		OutChan:             outChan,
		logger:              logger,
		rpcTable:            make(map[uint8]*RPCPair),
	}

	n.logger.Debug(fmt.Sprintf("network created %v\n", n), 3)

	return n
}

func (n *Network) RegisterRPC(msgObj Serializable, code uint8) {
	n.rpcTable[code] = &RPCPair{Code: code, Obj: msgObj}
	n.logger.Debug("Registered RPC code "+strconv.Itoa(int(code)), 3)
}

// connect to all remote nodes

func (n *Network) ConnectRemotes() error {
	n.logger.Debug(fmt.Sprintf("connecting to remotes %v\n", n.RemoteAddresses), 3)
	for id, address := range n.RemoteAddresses {
		var b [4]byte
		bs := b[:4]

		for true {
			conn, err := net.Dial("tcp", address)
			if err == nil {
				n.OutgoingConnections[id] = bufio.NewWriter(conn)
				n.OutMutex[id] = &sync.Mutex{}
				binary.LittleEndian.PutUint16(bs, uint16(n.Id))
				_, err := conn.Write(bs)
				if err != nil {
					panic("Error while connecting to replica " + strconv.Itoa(id))
				}
				n.logger.Debug("Outgoing TCP Connected to "+strconv.Itoa(id), 3)
				break
			} else {
				n.logger.Debug("Error while connecting to "+strconv.Itoa(id)+" "+err.Error()+" retrying....", 3)
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}

// listen to self.ListenAddress until all expected peers are connected

func (n *Network) Listen() error {
	n.logger.Debug("Listening on "+n.ListenAddress, 3)
	counter := 0
	var b [4]byte
	bs := b[:4]
	Listener, err_ := net.Listen("tcp", n.ListenAddress)
	if err_ != nil {
		panic("Error while listening to incoming connections" + err_.Error())
	}
	for counter < len(n.RemoteAddresses) {
		conn, err := Listener.Accept()
		if err != nil {
			panic(err.Error() + fmt.Sprintf("%v", err.Error()))
		}
		if _, err := io.ReadFull(conn, bs); err != nil {
			panic(err.Error() + fmt.Sprintf("%v", err.Error()))
		}
		id := int(binary.LittleEndian.Uint16(bs))

		n.IncomingConnections[id] = bufio.NewReader(conn)
		go n.HandleReadStream(n.IncomingConnections[id], id)
		n.logger.Debug("Incoming TCP Connected from "+strconv.Itoa(id), 3)
		counter++
	}

	n.logger.Debug("All incoming connections are established", 3)
	return nil
}

// read from reader

func (n *Network) HandleReadStream(reader *bufio.Reader, id int) error {
	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			n.logger.Debug("Error while reading message code: connection broken from "+strconv.Itoa(id)+fmt.Sprintf(" %v", err.Error()), 3)
			return err
		}
		if rpair, present := n.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				n.logger.Debug("Error while unmarshalling from "+strconv.Itoa(id)+fmt.Sprintf(" %v", err.Error()), 3)
				return err
			}
			n.OutChan <- &RPCPairPeer{
				RpcPair: &RPCPair{
					Code: msgType,
					Obj:  obj,
				},
				Peer: id,
			}
			n.logger.Debug("Pushed a message from "+strconv.Itoa(id), -1)
		} else {
			n.logger.Debug("Error received unknown message type from "+strconv.Itoa(id), 3)
			return nil
		}
	}
	return nil
}

// send to peer

func (n *Network) Send(rpc *RPCPairPeer) error {
	w := n.OutgoingConnections[rpc.Peer]
	if w == nil {
		panic("remote machine not found" + strconv.Itoa(rpc.Peer))
	}
	n.OutMutex[rpc.Peer].Lock()
	err := w.WriteByte(rpc.RpcPair.Code)
	if err != nil {
		n.logger.Debug("Error writing message code byte:"+err.Error(), 0)
		n.OutMutex[rpc.Peer].Unlock()
		return nil
	}
	err = rpc.RpcPair.Obj.Marshal(w)
	if err != nil {
		n.logger.Debug("Error while marshalling:"+err.Error(), 0)
		n.OutMutex[rpc.Peer].Unlock()
		return nil
	}
	err = w.Flush()
	if err != nil {
		n.logger.Debug("Error while flushing:"+err.Error(), 0)
		n.OutMutex[rpc.Peer].Unlock()
		return nil
	}
	n.OutMutex[rpc.Peer].Unlock()
	n.logger.Debug("Sent message to "+strconv.Itoa(rpc.Peer), 0)
	return nil
}

// broadcast to all peers

func (n *Network) Broadcast(rpc *RPCPair) error {
	for id, _ := range n.OutgoingConnections {
		n.Send(&RPCPairPeer{
			RpcPair: rpc,
			Peer:    id,
		})
	}
	n.logger.Debug("Broadcast message", 0)
	return nil
}

// get remote addresses from nodes

func GetRemoteAddresses(nodes []*Node) map[int]string {
	remoteAddresses := make(map[int]string)
	for _, node := range nodes {
		remoteAddresses[node.Id] = node.Ip + ":10080"
	}
	return remoteAddresses
}
