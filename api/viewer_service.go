package api

import (
	"context"

	"github.com/aergoio/aergo/api/typesconnect"
	"github.com/aergoio/aergo/consensus"
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/pkg/component"
	"github.com/aergoio/aergo/types"
	connect_go "github.com/bufbuild/connect-go"
)

var _ typesconnect.ViewerServiceClient = (*ViewerService)(nil)

type ViewerService struct {
	hub               *component.ComponentHub
	actorHelper       p2pcommon.ActorService
	consensusAccessor consensus.ConsensusAccessor //TODO refactor with actorHelper
	msgHelper         message.Helper
}

func (vs *ViewerService) Ping(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}

func (vs *ViewerService) GetChainInfo(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.ChainInfo], error) {
	return nil, nil
}

func (vs *ViewerService) GetNodeState(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}

func (vs *ViewerService) GetMetric(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}

func (vs *ViewerService) GetBestBlock(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	ca := vs.actorHelper.GetChainAccessor()
	bestBlock, err := ca.GetBestBlock()
	if err != nil {
		return nil, err
	}
	_ = bestBlock

	return nil, nil
}

func (vs *ViewerService) GetBlock(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Block], error) {
	ca := vs.actorHelper.GetChainAccessor()
	bestBlock, err := ca.GetBlock()
	if err != nil {
		return nil, err
	}
	_ = bestBlock

	return nil, nil
}

func (vs *ViewerService) GetBlockByNum(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Block], error) {

}

func (vs *ViewerService) GetBlockList(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.BlockHeaderList], error) {
	return nil, nil
}

func (vs *ViewerService) GetBlockStream(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}

func (vs *ViewerService) GetTx(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Tx], error) {
	return nil, nil
}

func (vs *ViewerService) GetTxInBlock(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.TxInBlock], error) {
	return nil, nil
}

// receipt
func (vs *ViewerService) GetReceipt(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Receipt], error) {
	return nil, nil
}

func (vs *ViewerService) GetReceiptInBlock(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}

// event
func (vs *ViewerService) GetEventList(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.EventList], error) {
	return nil, nil
}

func (vs *ViewerService) GetEventStream(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}

// contract
func (vs *ViewerService) GetContract(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}
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

// nft
func (vs *ViewerService) GetNFT(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.Empty], error) {
	return nil, nil
}
