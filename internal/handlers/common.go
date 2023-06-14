package handlers

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"mattermost-bot/internal/config"
	"mattermost-bot/internal/logger"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
)

var (
	Client *http.Client
	cfg    = config.GetConfig()
	log    = logger.GetLogger()
)

func init() {
	if cfg.InsecureSkipVerify {
		log.Warn().Msg("InsecureSkipVerify is enabled")
	}
	Client = &http.Client{
		Timeout:   cfg.HTTPClientTimeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify}},
	}
}

func handleError(w http.ResponseWriter, err error) {
	log.Error().Err(err).Send()
	httputils.WriteErrorIfNeeded(w, err)
}

func sendTooltipResponse(w http.ResponseWriter, msg string) error {
	if err := httputils.WriteJSON(w, apps.NewTextResponse(msg)); err != nil {
		return err
	}
	return nil
}

func createCallRequest(req *http.Request) (apps.CallRequest, error) {
	c := apps.CallRequest{}
	if err := json.NewDecoder(req.Body).Decode(&c); err != nil {
		return c, err
	}
	return c, nil
}

func sendMessageAsBot(w http.ResponseWriter, c apps.CallRequest, msg string) {
	_, err := appclient.AsBot(c.Context).DM(c.Context.ActingUser.Id, msg)
	if err != nil {
		handleError(w, err)
	}
}

func isAllowedChannel(id string) bool {
	if _, ok := cfg.ChannelAllowedIDs.GetChannels()[id]; ok {
		return true
	}
	return false
}

func getProductGroup(product string) (string, error) {
	if group, ok := cfg.ProductGroup[product]; ok {
		return group, nil
	}
	return "", fmt.Errorf("incorrect product name")
}

func md5SystemID(strings ...string) string {
	str := ""
	for _, s := range strings {
		str += s
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("remote_addr", r.RemoteAddr).
			Str("http_method", r.Method).
			Str("uri", r.RequestURI).
			Str("user_agent", r.UserAgent()).
			Send()
		handler.ServeHTTP(w, r)
	})
}

func createPostAsBot(c apps.CallRequest, chID, msgRootID, msg string) error {
	_, err := appclient.AsBot(c.Context).CreatePost(&model.Post{
		ChannelId: chID,
		RootId:    msgRootID,
		Message:   msg,
	})
	if err != nil {
		return err
	}
	return nil
}

func addReactionAsBot(c apps.CallRequest, userID, postID string) error {
	_, _, err := appclient.AsBot(c.Context).SaveReaction(&model.Reaction{
		UserId:    userID,
		PostId:    postID,
		EmojiName: cfg.PostEmojiName,
	})
	if err != nil {
		return err
	}

	return nil
}
