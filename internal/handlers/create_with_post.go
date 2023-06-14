package handlers

import (
	"fmt"
	"mattermost-bot/internal/models"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
)

func CreateWithPost(w http.ResponseWriter, r *http.Request) {
	c, err := createCallRequest(r)
	if err != nil {
		handleError(w, err)
		return
	}

	log.Debug().
		Interface("call_request", c).
		Send()

	if isAllowedChannel(c.Context.Channel.Id) {
		group, err := getProductGroup(c.Values["product"].(string))
		if err != nil {
			handleError(w, err)
			return
		}

		// Get real username because "c.Context.ActingUser.Username" provides nickname instead
		user, rsp, err := appclient.AsBot(c.Context).GetUser(c.Context.Post.UserId, "")
		if (err != nil) || (rsp.StatusCode != 200) {
			handleError(w, err)
			log.Error().Err(err).Str("user_id", c.Context.Post.UserId).Int("status_code", rsp.StatusCode).Msg("cannot get user information")
			return
		}

		// Calculate internal system id to avoid duplicate issue creation
		systemID := md5SystemID(
			user.Username,
			c.Values["product"].(string),
			c.Values["summary"].(string),
			c.Context.Post.Id,
		)

		// Compose issue description text
		additionalInfo := ""
		if c.Values["additional"] != nil {
			additionalInfo = fmt.Sprintf("Additional info: %s\n\n", c.Values["additional"].(string))
		}
		postMessage := c.Context.Post.Message
		messageLink := fmt.Sprintf("Mattermost message link: %s\n\n",
			c.Context.ExpandedContext.MattermostSiteURL+"/"+c.Context.ExpandedContext.Team.Name+"/"+"pl/"+c.Context.Post.Id)
		issueDescription := additionalInfo + messageLink + postMessage

		// Compose connector issue object
		data := models.NewTask{}.CreateWithDefaults()
		data.CreateTask.Task.Summary = c.Values["summary"].(string)
		data.CreateTask.Task.Description = issueDescription
		data.CreateTask.Task.Reporter = user.Username
		data.CreateTask.Task.ExternalIds.ExternalID = append(data.CreateTask.Task.ExternalIds.ExternalID, models.TaskExternalID{SystemID: systemID})
		data.CreateTask.Task.CustomFields.CustomField = append(data.CreateTask.Task.CustomFields.CustomField, models.TaskCustomField{Name: "Assigned Group+", Value: group})

		resp, err := createIssue(cfg.APIURL, data)
		if err != nil {
			handleError(w, err)
			return
		}

		message := fmt.Sprintf("Thanks for the report @%s!\n\nIssue created: %s\n%s\n%s", c.Context.ActingUser.Username, resp.CreateTaskResponse.Status, resp.CreateTaskResponse.IssueKey, resp.CreateTaskResponse.Error.Message)
		messageRootID := c.Context.Post.Id
		if c.Context.Post.RootId != "" {
			messageRootID = c.Context.Post.RootId
		}
		if err = createPostAsBot(c, c.Context.Channel.Id, messageRootID, message); err != nil {
			handleError(w, err)
			return
		}

		//Add bot reaction to post
		if err = addReactionAsBot(c, c.Context.BotUserID, c.Context.Post.Id); err != nil {
			handleError(w, err)
			return
		}

		sendTooltipResponse(w, "See details above")
	} else {
		sendTooltipResponse(w, "Channel is not allowed")
	}
}
