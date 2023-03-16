package api

import (
	"context"
	"errors"
	"math/big"
	"reflect"

	"github.com/aergoio/aergo-actor/actor"
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/pkg/component"
	"github.com/aergoio/aergo/types"
	"github.com/aergoio/aergo/types/typesconnect"
	connect_go "github.com/bufbuild/connect-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ typesconnect.ViewerServiceClient = (*ViewerApi)(nil)

type ViewerApi struct {
	hub         *component.ComponentHub
	actorHelper p2pcommon.ActorService
	msgHelper   message.Helper
}

func (v *ViewerApi) GetServerInfo(ctx context.Context, in *connect_go.Request[types.Empty]) (*connect_go.Response[types.ChainInfo], error) {
	return nil, nil
}

func (v *ViewerApi) GetMempoolInfo(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.SingleBytes], error) {
	return nil, nil
}

func (v *ViewerApi) GetBlock(ctx context.Context, blockHash *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Block], error) {
	ca := v.actorHelper.GetChainAccessor()
	block, err := ca.GetBlock(blockHash.Msg.GetValue())
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(block), nil
}

func (v *ViewerApi) GetBlockByNum(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Block], error) {
	ca := v.actorHelper.GetChainAccessor()
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

func (v *ViewerApi) GetBlockMetadata(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.BlockMetadata], error) {
	block, err := v.GetBlock(ctx, in)
	if err != nil {
		return nil, err
	}
	meta := block.Msg.GetMetadata()
	return connect_go.NewResponse(meta), nil
}

func (v *ViewerApi) GetBlockList(ctx context.Context, in *connect_go.Request[types.ListParams]) (*connect_go.Response[types.BlockHeaderList], error) {
	if in.Msg.Size > uint32(1000) {
		return nil, status.Errorf(codes.InvalidArgument, "size too big")
	}

	maxFetchSize := in.Msg.Size
	idx := uint32(0)
	hashes := make([][]byte, 0, maxFetchSize)
	blocks := make([]*types.Block, 0, maxFetchSize)
	var err error
	if len(in.Msg.Hash) > 0 {
		hash := in.Msg.Hash
		for idx < maxFetchSize {
			foundBlock, futureErr := extractBlockFromFuture(v.hub.RequestFuture(message.ChainSvc,
				&message.GetBlock{BlockHash: hash}, defaultActorTimeout, "rpc.(*AergoRPCService).ListBlockHeaders#1"))
			if nil != futureErr {
				if idx == 0 {
					err = futureErr
				}
				break
			}
			hashes = append(hashes, foundBlock.BlockHash())
			blocks = append(blocks, foundBlock)
			idx++
			hash = foundBlock.Header.PrevBlockHash
			if len(hash) == 0 {
				break
			}
		}
		if in.Msg.Asc || in.Msg.Offset != 0 {
			err = errors.New("Has unsupported param")
		}
	} else {
		end := types.BlockNo(0)
		start := types.BlockNo(in.Msg.Height) - types.BlockNo(in.Msg.Offset)
		if start >= types.BlockNo(maxFetchSize) {
			end = start - types.BlockNo(maxFetchSize-1)
		}
		if in.Msg.Asc {
			for i := end; i <= start; i++ {
				foundBlock, futureErr := extractBlockFromFuture(v.hub.RequestFuture(message.ChainSvc,
					&message.GetBlockByNo{BlockNo: i}, defaultActorTimeout, "rpc.(*AergoRPCService).ListBlockHeaders#2"))
				if nil != futureErr {
					if i == end {
						err = futureErr
					}
					break
				}
				hashes = append(hashes, foundBlock.BlockHash())
				blocks = append(blocks, foundBlock)
				idx++
			}
		} else {
			for i := start; i >= end; i-- {
				foundBlock, futureErr := extractBlockFromFuture(v.hub.RequestFuture(message.ChainSvc,
					&message.GetBlockByNo{BlockNo: i}, defaultActorTimeout, "rpc.(*AergoRPCService).ListBlockHeaders#2"))
				if nil != futureErr {
					if i == start {
						err = futureErr
					}
					break
				}
				hashes = append(hashes, foundBlock.BlockHash())
				blocks = append(blocks, foundBlock)
				idx++
			}
		}
	}
	return connect_go.NewResponse(&types.BlockHeaderList{Blocks: blocks}), err
}

func extractBlockFromFuture(future *actor.Future) (*types.Block, error) {
	rawResponse, err := future.Result()
	if err != nil {
		return nil, err
	}
	var blockRsp *message.GetBlockRsp
	switch v := rawResponse.(type) {
	case message.GetBlockRsp:
		blockRsp = &v
	case message.GetBestBlockRsp:
		blockRsp = (*message.GetBlockRsp)(&v)
	case message.GetBlockByNoRsp:
		blockRsp = (*message.GetBlockRsp)(&v)
	default:
		return nil, errors.New("Unsupported message type")
	}
	return extractBlock(blockRsp)
}

func extractBlock(from *message.GetBlockRsp) (*types.Block, error) {
	if nil != from.Err {
		return nil, from.Err
	}
	return from.Block, nil

}

func (v *ViewerApi) GetTx(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Tx], error) {
	result, err := v.actorHelper.CallRequestDefaultTimeout(message.MemPoolSvc, &message.MemPoolExist{Hash: in.Msg.Value})
	if err != nil {
		return nil, err
	}
	tx, err := v.msgHelper.ExtractTxFromResponse(result)
	if err != nil {
		return nil, err
	} else if tx == nil {
		// TODO try find tx in blockchain, but chainservice doesn't have method yet.
		return nil, status.Errorf(codes.NotFound, "not found")
	}
	return connect_go.NewResponse(tx), nil
}

func (v *ViewerApi) GetTxInBlock(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.TxInBlock], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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
func (v *ViewerApi) GetReceipt(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Receipt], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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
func (v *ViewerApi) GetEventList(ctx context.Context, in *connect_go.Request[types.FilterInfo]) (*connect_go.Response[types.EventList], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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
func (v *ViewerApi) QueryContract(ctx context.Context, in *connect_go.Request[types.Query]) (*connect_go.Response[types.SingleBytes], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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

func (v *ViewerApi) QueryContractState(ctx context.Context, in *connect_go.Request[types.StateQuery]) (*connect_go.Response[types.StateQueryProof], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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

func (v *ViewerApi) GetABI(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.ABI], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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
func (v *ViewerApi) GetAccountState(ctx context.Context, in *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.State], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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

func (v *ViewerApi) GetAccountStateAndProof(ctx context.Context, in *connect_go.Request[types.AccountAndRoot]) (*connect_go.Response[types.AccountProof], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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

func (v *ViewerApi) GetName(ctx context.Context, in *connect_go.Request[types.Name]) (*connect_go.Response[types.NameInfo], error) {
	result, err := v.hub.RequestFuture(message.ChainSvc,
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
