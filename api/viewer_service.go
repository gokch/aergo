package api

import (
	"context"
	"math/big"
	"reflect"
	"time"

	"github.com/aergoio/aergo-lib/log"
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/pkg/component"
	"github.com/aergoio/aergo/types"
	"github.com/aergoio/aergo/types/typesconnect"
	connect_go "github.com/bufbuild/connect-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	logger              = log.NewLogger("api")
	defaultActorTimeout = time.Second * 3
)

var _ typesconnect.ViewerServiceClient = (*ViewerService)(nil)

type ViewerService struct {
	hub         *component.ComponentHub
	actorHelper p2pcommon.ActorService
	msgHelper   message.Helper
}

func (vs *ViewerService) GetServerInfo(ctx context.Context, in *connect_go.Request[types.Empty]) (*connect_go.Response[types.ChainInfo], error) {
	return nil, nil
}

func (vs *ViewerService) GetMempoolInfo(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.SingleBytes], error) {
	return nil, nil
}

func (vs *ViewerService) GetBlock(ctx context.Context, blockHash *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Block], error) {
	ca := vs.actorHelper.GetChainAccessor()
	block, err := ca.GetBlock(blockHash.Msg.GetValue())
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(block), nil
}

func (vs *ViewerService) GetBlockByNum(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Block], error) {
	ca := vs.actorHelper.GetChainAccessor()
	blockNo := big.NewInt(0).SetBytes(in.Msg.GetValue()).Uint64()
	blockHash, err := ca.GetHashByNo(blockNo)
	if err != nil {
		return nil, err
	}
	block, err := ca.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(block), nil
}

func (vs *ViewerService) GetBlockMetadata(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.BlockMetadata], error) {
	block, err := vs.GetBlock(ctx, in)
	if err != nil {
		return nil, err
	}
	meta := block.Msg.GetMetadata()
	return connect_go.NewResponse(meta), nil
}

func (vs *ViewerService) GetBlockList(ctx context.Context, in *connect_go.Request[types.ListParams]) (*connect_go.Response[types.BlockHeaderList], error) {
	return nil, nil
}

func (vs *ViewerService) GetTx(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Tx], error) {
	result, err := vs.actorHelper.CallRequestDefaultTimeout(message.MemPoolSvc, &message.MemPoolExist{Hash: in.Msg.Value})
	if err != nil {
		return nil, err
	}
	tx, err := vs.msgHelper.ExtractTxFromResponse(result)
	if err != nil {
		return nil, err
	} else if tx == nil {
		// TODO try find tx in blockchain, but chainservice doesn't have method yet.
		return nil, status.Errorf(codes.NotFound, "not found")
	}
	return connect_go.NewResponse(tx), nil
}

func (vs *ViewerService) GetTxInBlock(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.TxInBlock], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.GetTx{TxHash: in.Msg.Value}, defaultActorTimeout, "rpc.(*AergoRPCService).GetBlockTX").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(message.GetTxRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(&types.TxInBlock{Tx: rsp.Tx, TxIdx: rsp.TxIds}), rsp.Err
}

// receipt
func (vs *ViewerService) GetReceipt(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Receipt], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.GetReceipt{TxHash: in.Msg.Value}, defaultActorTimeout, "rpc.(*AergoRPCService).GetReceipt").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(message.GetReceiptRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(rsp.Receipt), rsp.Err
}

// event
func (vs *ViewerService) GetEventList(ctx context.Context, in *connect_go.Request[types.FilterInfo]) (*connect_go.Response[types.EventList], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.ListEvents{Filter: in.Msg}, defaultActorTimeout, "rpc.(*AergoRPCService).ListEvents").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(*message.ListEventsRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(&types.EventList{Events: rsp.Events}), rsp.Err
}

// contract
func (vs *ViewerService) QueryContract(ctx context.Context, in *connect_go.Request[types.Query]) (*connect_go.Response[types.SingleBytes], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.GetQuery{Contract: in.Msg.ContractAddress, Queryinfo: in.Msg.Queryinfo}, defaultActorTimeout, "rpc.(*AergoRPCService).QueryContract").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(message.GetQueryRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(&types.SingleBytes{Value: rsp.Result}), rsp.Err
}

func (vs *ViewerService) QueryContractState(ctx context.Context, in *connect_go.Request[types.StateQuery]) (*connect_go.Response[types.StateQueryProof], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.GetStateQuery{ContractAddress: in.Msg.ContractAddress, StorageKeys: in.Msg.StorageKeys, Root: in.Msg.Root, Compressed: in.Msg.Compressed}, defaultActorTimeout, "rpc.(*AergoRPCService).GetStateQuery").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(message.GetStateQueryRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(rsp.Result), rsp.Err
}

func (vs *ViewerService) GetABI(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.ABI], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.GetABI{Contract: in.Msg.Value}, defaultActorTimeout, "rpc.(*AergoRPCService).GetABI").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(message.GetABIRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(rsp.ABI), rsp.Err
}

// account
func (vs *ViewerService) GetAccountState(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.State], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.GetState{Account: in.Msg.Value}, defaultActorTimeout, "rpc.(*AergoRPCService).GetState").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(message.GetStateRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(rsp.State), rsp.Err
}

func (vs *ViewerService) GetAccountStateAndProof(ctx context.Context, in *connect_go.Request[types.AccountAndRoot]) (*connect_go.Response[types.AccountProof], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.GetStateAndProof{Account: in.Msg.Account, Root: in.Msg.Root, Compressed: in.Msg.Compressed}, defaultActorTimeout, "rpc.(*AergoRPCService).GetStateAndProof").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(message.GetStateAndProofRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(rsp.StateProof), rsp.Err
}

func (vs *ViewerService) GetName(ctx context.Context, in *connect_go.Request[types.Name]) (*connect_go.Response[types.NameInfo], error) {
	result, err := vs.hub.RequestFuture(message.ChainSvc,
		&message.GetNameInfo{Name: in.Msg.Name, BlockNo: in.Msg.BlockNo}, defaultActorTimeout, "rpc.(*AergoRPCService).GetName").Result()
	if err != nil {
		return nil, err
	}
	rsp, ok := result.(*message.GetNameInfoRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	if rsp.Err == types.ErrNameNotFound {
		return connect_go.NewResponse(rsp.Owner), status.Errorf(codes.NotFound, rsp.Err.Error())
	}
	return connect_go.NewResponse(rsp.Owner), rsp.Err
}
