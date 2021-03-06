package raft

import (
	"bytes"
	"context"
	"github.com/processout/grpc-go-pool"
	"github.com/readystock/golog"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

type grpcPipeline struct {
	conn         *grpcpool.ClientConn
	trans        *GrpcTransport
	doneCh       chan AppendFuture
	inProgressCh chan *appendFuture
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

func newGrpcPipeline(transport *GrpcTransport, conn *grpcpool.ClientConn) *grpcPipeline {
	p := &grpcPipeline{
		conn:         conn,
		trans:        transport,
		doneCh:       make(chan AppendFuture, rpcMaxPipeline),
		inProgressCh: make(chan *appendFuture, rpcMaxPipeline),
		shutdownCh:   make(chan struct{}),
	}
	return p
}

func (pipe *grpcPipeline) AppendEntries(args *AppendEntriesRequest, resp *AppendEntriesResponse) (AppendFuture, error) {
	ctx := context.Background()
	// Create a new future
	future := &appendFuture{
		start: time.Now(),
		args:  args,
		resp:  resp,
	}
	future.init()
	go func() {
		// golog.Debugf("[%s] sending append entries on pipeline", pipe.trans.LocalAddr())
		client := NewRaftServiceClient(pipe.conn.ClientConn)
		result, err := client.AppendEntries(ctx, args)
		if err != nil {
			future.respond(err)
		}
		future.resp = result
		pipe.doneCh <- future
	}()
	return future, nil
}

// Consumer returns a channel that can be used to consume complete futures.
func (pipe *grpcPipeline) Consumer() <-chan AppendFuture {
	return pipe.doneCh
}

// Closed is used to shutdown the pipeline connection.
func (pipe *grpcPipeline) Close() error {
	pipe.shutdownLock.Lock()
	defer pipe.shutdownLock.Unlock()
	if pipe.shutdown {
		return nil
	}

	// Release the connection
	pipe.conn.Close()

	pipe.shutdown = true
	close(pipe.shutdownCh)
	return nil
}

func NewGrpcTransport(server *grpc.Server, localAddr string) (*GrpcTransport, error) {
	return newGrpcTransport(server, localAddr)
}

type WithProtocolVersion interface {
	GetProtocolVersion() ProtocolVersion
}

type GrpcTransport struct {
	connPool        lake
	consumeChan     chan RPC
	heartbeatFn     func(RPC)
	heartbeatFnLock sync.Mutex

	shutdown     bool
	shutdownChan chan struct{}
	shutdownLock sync.Mutex

	serverAddressProvider ServerAddressProvider

	gRPC      *grpc.Server
	svc       service
	localAddr string
}

func newGrpcTransport(server *grpc.Server, localAddr string) (*GrpcTransport, error) {
	transport := &GrpcTransport{
		connPool:     newLake(),
		consumeChan:  make(chan RPC),
		shutdownChan: make(chan struct{}),
		gRPC:         server,
		localAddr:    localAddr,
	}
	transport.svc = service{
		appendEntries:   transport.handleAppendEntriesCommand,
		requestVote:     transport.handleRequestVoteCommand,
		installSnapshot: transport.handleInstallSnapshotCommand,
	}
	RegisterRaftServiceServer(server, transport.svc)
	golog.Infof("[%s] finished starting grpc transport", transport.LocalAddr())
	return transport, nil
}

func (transport *GrpcTransport) Consumer() <-chan RPC {
	return transport.consumeChan
}

func (transport *GrpcTransport) LocalAddr() ServerAddress {
	return ServerAddress(transport.localAddr)
}

func (transport *GrpcTransport) AppendEntriesPipeline(id ServerID, target ServerAddress) (AppendPipeline, error) {
	// Get a connection
	ctx := context.Background()
	conn, err := transport.connPool.GetConnection(ctx, target)
	if err != nil {
		return nil, err
	}
	// Create the pipeline
	return newGrpcPipeline(transport, conn), nil
}

func (transport *GrpcTransport) AppendEntries(id ServerID, target ServerAddress, args *AppendEntriesRequest, resp *AppendEntriesResponse) error {
	result, err := transport.executeTransportClient(context.Background(), id, target, func(ctx context.Context, client RaftServiceClient) (result interface{}, err error) {
		return client.AppendEntries(ctx, args)
	})
	if err == nil {
		*resp = *result.(*AppendEntriesResponse)
	}
	return err
}

func (transport *GrpcTransport) RequestVote(id ServerID, target ServerAddress, args *RequestVoteRequest, resp *RequestVoteResponse) error {
	result, err := transport.executeTransportClient(context.Background(), id, target, func(ctx context.Context, client RaftServiceClient) (result interface{}, err error) {
		return client.RequestVote(ctx, args)
	})
	if err == nil {
		*resp = *result.(*RequestVoteResponse)
	}
	return err
}

func (transport *GrpcTransport) InstallSnapshot(id ServerID, target ServerAddress, args *InstallSnapshotRequest, resp *InstallSnapshotResponse, data io.Reader) error {
	ctx := context.Background()
	golog.Debugf("[%s] sending install snapshot to %s", transport.LocalAddr(), target)
	conn, err := transport.connPool.GetConnection(ctx, target)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := NewRaftServiceClient(conn.ClientConn)
	request := *args
	bytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}
	wrapper := &InstallSnapshotRequestWrapper{
		Request:  &request,
		Snapshot: bytes,
	}
	result, err := client.InstallSnapshot(ctx, wrapper)
	if err == nil {
		*resp = *result
	}
	return err
}

func (transport *GrpcTransport) EncodePeer(id ServerID, addr ServerAddress) []byte {
	address := transport.getProviderAddressOrFallback(id, addr)
	return []byte(address)
}

func (transport *GrpcTransport) DecodePeer(data []byte) ServerAddress {
	return ServerAddress(data)
}

func (transport *GrpcTransport) SetHeartbeatHandler(cb func(rpc RPC)) {
	transport.heartbeatFnLock.Lock()
	defer transport.heartbeatFnLock.Unlock()
	transport.heartbeatFn = cb
}

// Implementing the WithClose interface.
func (transport *GrpcTransport) Close() error {
	transport.shutdownLock.Lock()
	defer transport.shutdownLock.Unlock()

	if !transport.shutdown {
		close(transport.shutdownChan)
		transport.shutdown = true
	}
	return nil
}

func (transport *GrpcTransport) getProviderAddressOrFallback(id ServerID, target ServerAddress) ServerAddress {
	if transport.serverAddressProvider != nil {
		serverAddressOverride, err := transport.serverAddressProvider.ServerAddr(id)
		if err != nil {
			golog.Warnf("unable to get address for server id %v, using fallback address %v: %v", id, target, err)
		} else {
			return serverAddressOverride
		}
	}
	return target
}

func (transport *GrpcTransport) executeTransportClient(
	ctx context.Context,
	id ServerID,
	target ServerAddress,
	call func(ctx context.Context, client RaftServiceClient) (interface{}, error)) (result interface{}, err error) {
	conn, err := transport.connPool.GetConnection(ctx, target)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := NewRaftServiceClient(conn.ClientConn)
	return call(ctx, client)
}

func (transport *GrpcTransport) listen() {
	RegisterRaftServiceServer(transport.gRPC, transport.svc)
}

func (transport *GrpcTransport) handleAppendEntriesCommand(ctx context.Context, request *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	// Create the RPC object
	respCh := make(chan RPCResponse, 1)
	rpc := RPC{
		RespChan: respCh,
	}

	// Decode the command
	isHeartbeat := false

	rpc.Command = request

	// Check if this is a heartbeat
	if request.Term != 0 && request.Leader != nil &&
		request.PrevLogEntry == 0 && request.PrevLogTerm == 0 &&
		len(request.Entries) == 0 && request.LeaderCommitIndex == 0 {
		isHeartbeat = true
	}

	// Check for heartbeat fast-path
	if isHeartbeat {
		transport.heartbeatFnLock.Lock()
		fn := transport.heartbeatFn
		transport.heartbeatFnLock.Unlock()
		if fn != nil {
			fn(rpc)
			// golog.Debugf("[%s] append entries command is heartbeat", transport.LocalAddr())
			goto RESP
		}
	}

	// Dispatch the RPC
	select {
	case transport.consumeChan <- rpc:
		// golog.Debugf("[%s] dispatching append entries request to consumer", transport.LocalAddr())
	case <-transport.shutdownChan:
		return nil, ErrTransportShutdown
	}

	// Wait for response
RESP:
	select {
	case resp := <-respCh:
		// Send the error first
		if resp.Error != nil {
			return nil, resp.Error
		}
		rsp := (resp.Response).(*AppendEntriesResponse)
		return rsp, resp.Error

	case <-transport.shutdownChan:
		return nil, ErrTransportShutdown
	}
}

func (transport *GrpcTransport) handleRequestVoteCommand(ctx context.Context, request *RequestVoteRequest) (*RequestVoteResponse, error) {
	golog.Infof("[%s] received request vote command", transport.LocalAddr())
	// Create the RPC object
	respCh := make(chan RPCResponse, 1)
	rpc := RPC{
		RespChan: respCh,
	}

	// Decode the command
	rpc.Command = *request
	// Dispatch the RPC
	select {
	case transport.consumeChan <- rpc:
		golog.Debugf("[%s] dispatching request vote to consumer", transport.LocalAddr())
	case <-transport.shutdownChan:
		return nil, ErrTransportShutdown
	}

	select {
	case resp := <-respCh:
		golog.Debugf("[%s] received vote response from consumer", transport.LocalAddr())
		// Send the error first
		if resp.Error != nil {
			return nil, resp.Error
		}
		rsp := (resp.Response).(*RequestVoteResponse)
		return rsp, resp.Error

	case <-transport.shutdownChan:
		return nil, ErrTransportShutdown
	}
}

func (transport *GrpcTransport) handleInstallSnapshotCommand(ctx context.Context, request *InstallSnapshotRequestWrapper) (*InstallSnapshotResponse, error) {
	golog.Infof("[%s] received install snapshot command", transport.LocalAddr())
	// Create the RPC object
	respCh := make(chan RPCResponse, 1)
	rpc := RPC{
		RespChan: respCh,
	}

	rpc.Command = request.Request

	rpc.Reader = bytes.NewReader(request.Snapshot)

	// Dispatch the RPC
	select {
	case transport.consumeChan <- rpc:
		golog.Debugf("[%s] dispatching append entries request to consumer", transport.LocalAddr())
	case <-transport.shutdownChan:
		return nil, ErrTransportShutdown
	}

	select {
	case resp := <-respCh:
		golog.Debugf("[%s] received append entries response from consumer", transport.LocalAddr())
		// Send the error first
		if resp.Error != nil {
			return nil, resp.Error
		}
		rsp := (resp.Response).(*InstallSnapshotResponse)
		return rsp, resp.Error

	case <-transport.shutdownChan:
		return nil, ErrTransportShutdown
	}
}
