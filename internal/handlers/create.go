package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mattermost-bot/internal/models"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
)

func createIssue(url string, requestData models.NewTask) (*models.NewTaskResponse, error) {
	var responseData *models.NewTaskResponse

	data, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	resp, err := Client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("internal server error: %v", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}

func Create(w http.ResponseWriter, r *http.Request) {
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
		user, rsp, err := appclient.AsBot(c.Context).GetUser(c.Context.ActingUser.Id, "")
		if (err != nil) || (rsp.StatusCode != 200) {
			handleError(w, err)
			log.Error().Err(err).Str("user_id", c.Context.ActingUser.Id).Int("status_code", rsp.StatusCode).Msg("cannot get user information")
			return
		}

		// Calculate internal system id to avoid duplicate issue creation
		systemID := md5SystemID(
			user.Username,
			c.Values["product"].(string),
			c.Values["summary"].(string),
			c.Values["description"].(string))

		data := models.NewTask{}.CreateWithDefaults()
		data.CreateTask.Task.Summary = c.Values["summary"].(string)
		data.CreateTask.Task.Description = c.Values["description"].(string)
		data.CreateTask.Task.Reporter = user.Username
		data.CreateTask.Task.ExternalIds.ExternalID = append(data.CreateTask.Task.ExternalIds.ExternalID, models.TaskExternalID{SystemID: systemID})
		data.CreateTask.Task.CustomFields.CustomField = append(data.CreateTask.Task.CustomFields.CustomField, models.TaskCustomField{Name: "Assigned Group+", Value: group})

		resp, err := createIssue(cfg.APIURL, data)
		if err != nil {
			handleError(w, err)
			return
		}

		message := fmt.Sprintf("Issue created: %s\n%s\n%s", resp.CreateTaskResponse.Status, resp.CreateTaskResponse.IssueKey, resp.CreateTaskResponse.Error.Message)
		sendMessageAsBot(w, c, message)
		sendTooltipResponse(w, "Created a post in your DM channel")
	} else {
		sendTooltipResponse(w, "Channel is not allowed")
	}
}
