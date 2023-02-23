package api

import (
	"net/http"

	"github.com/aergoio/aergo/api/typesconnect"
)

func RegisterMux() {
	viewerService := &ViewerService{}
	path, handler := typesconnect.NewViewerServiceHandler(viewerService)
	http.DefaultServeMux.Handle(path, handler)
}
