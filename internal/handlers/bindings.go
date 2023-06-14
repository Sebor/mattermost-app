package handlers

import (
	"mattermost-bot/internal/meta"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
)

func Bindings(w http.ResponseWriter, r *http.Request) {
	var bindings = meta.CommandBindings

	c, err := createCallRequest(r)
	if err != nil {
		handleError(w, err)
		return
	}

	log.Debug().
		Interface("call_request", c).
		Send()

	if c.Context.Channel != nil {
		if isAllowedChannel(c.Context.Channel.Id) {
			bindings = meta.FullBindings
		}
	}

	httputils.WriteJSON(w, apps.NewDataResponse(bindings))
}
