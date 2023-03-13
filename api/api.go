package api

import (
	"net/http"

	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/types/typesconnect"
)

func RegisterAPI(actorHelper p2pcommon.ActorService, msgHelper message.Helper) {
	// Register the viewer service
	viewerService := &ViewerService{
		actorHelper: actorHelper,
		msgHelper:   msgHelper,
	}
	path, handler := typesconnect.NewViewerServiceHandler(viewerService)
	http.DefaultServeMux.Handle(path, handler)
}
