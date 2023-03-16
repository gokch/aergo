package api

import (
	"context"
	"reflect"

	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/pkg/component"
	"github.com/aergoio/aergo/types"
	"github.com/aergoio/aergo/types/typesconnect"
	"github.com/bufbuild/connect-go"
	connect_go "github.com/bufbuild/connect-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ typesconnect.WalletServiceClient = (*WalletApi)(nil)

type WalletApi struct {
	hub         *component.ComponentHub
	actorHelper p2pcommon.ActorService
	msgHelper   message.Helper
}

func (w *WalletApi) Create(context.Context, *connect_go.Request[types.Personal]) (*connect_go.Response[types.Account], error) {
	return nil, nil
}

func (w *WalletApi) Import(ctx context.Context, in *connect_go.Request[types.ImportFormat]) (*connect_go.Response[types.Account], error) {
	msg := &message.ImportAccount{OldPass: in.Msg.Oldpass, NewPass: in.Msg.Newpass}
	if in.Msg.Wif != nil {
		msg.Wif = in.Msg.Wif.Value
	} else if in.Msg.Keystore != nil {
		msg.Keystore = in.Msg.Keystore.Value
	} else {
		return nil, status.Errorf(codes.Internal, "require either wif or keystore contents")
	}
	result, err := w.hub.RequestFutureResult(message.AccountsSvc,
		msg,
		defaultActorTimeout, "rpc.(*AergoRPCService).ImportAccount")
	if err != nil {
		if err == component.ErrHubUnregistered {
			return nil, status.Errorf(codes.Unavailable, "Unavailable personal feature")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rsp, ok := result.(*message.ImportAccountRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(rsp.Account), rsp.Err
}

func (w *WalletApi) Export(ctx context.Context, in *connect_go.Request[types.Personal]) (*connect_go.Response[types.SingleBytes], error) {
	result, err := w.hub.RequestFutureResult(message.AccountsSvc,
		&message.ExportAccount{Account: in.Msg.Account, Pass: in.Msg.Passphrase, AsKeystore: true},
		defaultActorTimeout, "rpc.(*AergoRPCService).ExportAccount")
	if err != nil {
		if err == component.ErrHubUnregistered {
			return nil, status.Errorf(codes.Unavailable, "Unavailable personal feature")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rsp, ok := result.(*message.ExportAccountRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(&types.SingleBytes{Value: rsp.Wif}), rsp.Err
}

func (w *WalletApi) Accounts(context.Context, *connect_go.Request[types.Empty]) (*connect_go.Response[types.AccountList], error) {
	result, err := w.hub.RequestFutureResult(message.AccountsSvc,
		&message.GetAccounts{}, defaultActorTimeout, "rpc.(*AergoRPCService).GetAccounts")
	if err != nil {
		if err == component.ErrHubUnregistered {
			return nil, status.Errorf(codes.Unavailable, "Unavailable personal feature")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rsp, ok := result.(*message.GetAccountsRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(rsp.Accounts), nil
}

func (w *WalletApi) Make(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.SingleBytes], error) {
	return nil, nil
}

func (w *WalletApi) Sign(ctx context.Context, in *connect_go.Request[types.Tx]) (*connect_go.Response[types.Tx], error) {
	result, err := w.hub.RequestFutureResult(message.AccountsSvc,
		&message.SignTx{Tx: in.Msg}, defaultActorTimeout, "rpc.(*AergoRPCService).SignTX")
	if err != nil {
		if err == component.ErrHubUnregistered {
			return nil, status.Errorf(codes.Unavailable, "Unavailable personal feature")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rsp, ok := result.(*message.SignTxRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	return connect_go.NewResponse(rsp.Tx), rsp.Err
}

func (w *WalletApi) Verify(ctx context.Context, in *connect_go.Request[types.Tx]) (*connect_go.Response[types.VerifyResult], error) {
	result, err := w.hub.RequestFutureResult(message.AccountsSvc,
		&message.VerifyTx{Tx: in.Msg}, defaultActorTimeout, "rpc.(*AergoRPCService).VerifyTX")
	if err != nil {
		if err == component.ErrHubUnregistered {
			return nil, status.Errorf(codes.Unavailable, "Unavailable personal feature")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rsp, ok := result.(*message.VerifyTxRsp)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal type (%v) error", reflect.TypeOf(result))
	}
	ret := &types.VerifyResult{Tx: rsp.Tx}
	if rsp.Err == types.ErrSignNotMatch {
		ret.Error = types.VerifyStatus_VERIFY_STATUS_SIGN_NOT_MATCH
	} else {
		ret.Error = types.VerifyStatus_VERIFY_STATUS_OK
	}
	return connect.NewResponse(ret), nil
}

func (w *WalletApi) Send(ctx context.Context, in *connect_go.Request[types.TxList]) (*connect_go.Response[types.CommitResultList], error) {
	if in.Msg.Txs == nil {
		return nil, status.Errorf(codes.InvalidArgument, "input tx is empty")
	}
	w.hub.Get(message.MemPoolSvc)
	// p := newPutter(ctx, in.Txs, rpc.hub, defaultActorTimeout<<2)
	// err := p.Commit()
	// if err == nil {
	// 	results := &types.CommitResultList{Results: p.rs}
	// 	return results, nil
	// } else {
	// 	return nil, err
	// }
	return nil, nil
}

func (w *WalletApi) Deposit(context.Context, *connect_go.Request[types.SingleBytes]) (*connect_go.Response[types.SingleBytes], error) {
	return nil, nil
}
