package handlers

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
)

func Install(w http.ResponseWriter, r *http.Request) {
	c, err := createCallRequest(r)
	if err != nil {
		handleError(w, err)
		return
	}

	log.Debug().
		Interface("call_request", c).
		Send()

	team, rsp, err := appclient.AsActingUser(c.Context).GetTeamByName(cfg.TeamName, "")
	if (err != nil) || (rsp.StatusCode != 200) {
		handleError(w, err)
		log.Error().Err(err).Str("team_name", cfg.TeamName).Int("status_code", rsp.StatusCode).Msg("cannot get team information")
		return
	}

	_, rsp, err = appclient.AsActingUser(c.Context).AddTeamMember(team.Id, c.Context.BotUserID)
	if (err != nil) || (rsp.StatusCode != 201) {
		handleError(w, err)
		log.Error().Err(err).Str("team_id", team.Id).Int("status_code", rsp.StatusCode).Msg("cannot add bot to team")
		return
	}
	log.Info().Str("team_name", team.Name).Str("team_id", team.Id).Str("bot_id", c.Context.BotUserID).Msg("bot added to team")

	for ch := range cfg.ChannelAllowedIDs.GetChannels() {
		_, rsp, err := appclient.AsActingUser(c.Context).AddChannelMember(ch, c.Context.BotUserID)
		if (err != nil) || (rsp.StatusCode != 201) {
			log.Error().Err(err).Str("channel", ch).Int("status_code", rsp.StatusCode).Msg("cannot add bot to channel")
		} else {
			log.Info().Str("channel_id", ch).Str("bot_id", c.Context.BotUserID).Msg("bot added to channel")
		}
	}

	httputils.WriteJSON(w, apps.NewDataResponse(nil))
}
