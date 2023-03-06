package api

import (
	"context"

	"github.com/aergoio/aergo/consensus"
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/pkg/component"
	"github.com/aergoio/aergo/types"
	"github.com/aergoio/aergo/types/typesconnect"
	connect_go "github.com/bufbuild/connect-go"
)

var _ typesconnect.ViewerServiceClient = (*ViewerService)(nil)

type ViewerService struct {
	hub               *component.ComponentHub
	actorHelper       p2pcommon.ActorService
	consensusAccessor consensus.ConsensusAccessor
	msgHelper         message.Helper
}

func (vs *ViewerService) Ping(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.SingleBytes], error) {
	return nil, nil
}

func (vs *ViewerService) GetChainInfo(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.ChainInfo], error) {
	return nil, nil
}

func (vs *ViewerService) GetNodeState(context.Context, *connect_go.Request[types.NodeReq]) (*connect_go.Response[types.SingleBytes], error) {
	return nil, nil
}

func (vs *ViewerService) GetMetric(context.Context, *connect_go.Request[types.MetricsRequest]) (*connect_go.Response[types.Metrics], error) {
	return nil, nil
}

func (vs *ViewerService) GetBestBlock(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Block], error) {
	ca := vs.actorHelper.GetChainAccessor()
	bestBlock, err := ca.GetBestBlock()
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(bestBlock), nil
}

func (vs *ViewerService) GetBlock(ctx context.Context, blockHash *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Block], error) {
	ca := vs.actorHelper.GetChainAccessor()
	block, err := ca.GetBlock(blockHash.Msg.GetValue())
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(block), nil
}

func (vs *ViewerService) GetBlockMetadata(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.BlockMetadata], error) {
	return nil, nil
}

func (vs *ViewerService) GetBlockByNum(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Block], error) {
	return nil, nil
}

func (vs *ViewerService) GetBlockList(context.Context, *connect_go.Request[types.ListParams]) (*connect_go.Response[types.BlockHeaderList], error) {
	return nil, nil
}

func (vs *ViewerService) GetTx(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Tx], error) {
	return nil, nil
}

func (vs *ViewerService) GetTxInBlock(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.TxInBlock], error) {
	return nil, nil
}

// receipt
func (vs *ViewerService) GetReceipt(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Receipt], error) {
	return nil, nil
}

func (vs *ViewerService) GetReceiptInBlock(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.Receipt], error) {
	return nil, nil
}

// event
func (vs *ViewerService) GetEvent(context.Context, *connect_go.Request[types.FilterInfo]) (*connect_go.Response[types.Event], error) {
	return nil, nil
}

func (vs *ViewerService) GetEventList(context.Context, *connect_go.Request[types.FilterInfo]) (*connect_go.Response[types.EventList], error) {
	return nil, nil
}

// contract
func (vs *ViewerService) QueryContract(context.Context, *connect_go.Request[types.Query]) (*connect_go.Response[types.SingleBytes], error) {
	return nil, nil
}

// account
func (vs *ViewerService) GetAccount(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Account], error) {
	return nil, nil
}

func (vs *ViewerService) GetName(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}

// token
func (vs *ViewerService) GetToken(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}

func (vs *ViewerService) GetABI(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.ABI], error) {
	return nil, nil
}

func (vs *ViewerService) GetAccountState(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.State], error) {
	return nil, nil
}

func (vs *ViewerService) GetNameInfo(context.Context, *connect_go.Request[types.Name]) (*connect_go.Response[types.NameInfo], error) {
	return nil, nil
}

// nft
func (vs *ViewerService) GetNFT(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}
