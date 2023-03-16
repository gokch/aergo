package api

import (
	"net/http"
	"time"

	"github.com/aergoio/aergo-lib/log"
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/pkg/component"
	"github.com/aergoio/aergo/types/typesconnect"
)

var (
	logger              = log.NewLogger("api")
	defaultActorTimeout = time.Second * 3
)

func RegisterAPI(hub *component.ComponentHub, actorHelper p2pcommon.ActorService, msgHelper message.Helper) {
	// Register the viewer service
	viewerPath, viewerHandler := typesconnect.NewViewerServiceHandler(&ViewerApi{
		hub:         hub,
		actorHelper: actorHelper,
		msgHelper:   msgHelper,
	})

	// Register the Wallet service
	walletPath, walletHandler := typesconnect.NewWalletServiceHandler(&WalletApi{
		hub:         hub,
		actorHelper: actorHelper,
		msgHelper:   msgHelper,
	})

	http.DefaultServeMux.Handle(viewerPath, viewerHandler)
	http.DefaultServeMux.Handle(walletPath, walletHandler)
}
