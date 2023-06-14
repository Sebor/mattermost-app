package meta

import (
	"embed"
	"mattermost-bot/internal/config"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

//go:embed static
var StaticFiles embed.FS

var cfg = config.GetConfig()

// Manifest declares the app's metadata. It must be provided for the app to be
// installable. In this example, the following permissions are requested:
//   - Create posts as a bot.
//   - Add icons to the channel header that will call back into your app when
//     clicked.
//   - Add a /-command with a callback.
var Manifest = apps.Manifest{
	// App ID must be unique across all Mattermost Apps.
	AppID: "mattermost-bot",

	// App's release/version.
	Version: "v1.0.0",

	// A (long) display name for the app.
	DisplayName: cfg.AppDisplayName,

	// The icon for the app's bot account, same icon is also used for bindings
	// and forms.
	Icon: "gojira.png",

	// HomepageURL is required for an app to be installable.
	HomepageURL: cfg.AppHomePage,

	// Need ActAsBot to post back to the user.
	RequestedPermissions: []apps.Permission{
		apps.PermissionActAsBot,
		apps.PermissionActAsUser,
	},

	// Add UI elements: a /-command, and a channel header button.
	RequestedLocations: []apps.Location{
		apps.LocationChannelHeader,
		apps.LocationCommand,
		apps.LocationPostMenu,
	},

	// Expand bindings requests
	Bindings: apps.NewCall("/bindings").WithExpand(apps.Expand{
		Channel: apps.ExpandAll,
	}),

	// Running the app as an HTTP service is the only deployment option
	// supported.
	Deploy: apps.Deploy{
		HTTP: &apps.HTTP{
			RootURL: cfg.AppRootURL,
		},
	},

	OnInstall: apps.NewCall("/install").WithExpand(apps.Expand{
		ActingUserAccessToken: apps.ExpandAll,
		Team:                  apps.ExpandID,
		TeamMember:            apps.ExpandSummary,
	}),
}

// CommandBindings describes the details for the command UI
var CommandBindings = []apps.Binding{
	{
		Location: "/command",
		Bindings: []apps.Binding{
			{
				// For commands, Location is not necessary, it will be defaulted to the label.
				Icon:        "bottle.png",
				Label:       "jira",
				Description: "Create jira issue", // appears in autocomplete.
				Hint:        "[ create ]",        // appears in autocomplete, usually indicates as to what comes after choosing the option.
				Bindings: []apps.Binding{
					{
						Label: "create",
						Form:  &CreateForm,
					},
				},
			},
		},
	},
}

// FullBindings describes the details for the command and graphical UI
var FullBindings = []apps.Binding{
	{
		Location: apps.LocationCommand,
		Bindings: []apps.Binding{
			{
				// For commands, Location is not necessary, it will be defaulted to the label.
				Icon:        "icon.png",
				Label:       "jira",
				Description: "Create jira issue", // appears in autocomplete.
				Hint:        "[ create ]",        // appears in autocomplete, usually indicates as to what comes after choosing the option.
				Bindings: []apps.Binding{
					{
						Label: "create",
						Form:  &CreateForm,
					},
				},
			},
		},
	},
	{
		Location: apps.LocationChannelHeader,
		Bindings: []apps.Binding{
			{
				// Location: "send-button",      // an app-chosen string.
				Icon:  "jira.png",          // reuse the App icon for the channel header.
				Label: "Create Jira issue", // appearance in the "more..." menu.
				Form:  &CreateForm,         // the form to display.
			},
		},
	},
	{
		Location: apps.LocationPostMenu,
		Bindings: []apps.Binding{
			{
				Icon:  "mattermost.png",
				Label: "CREATE JIRA ISSUE",
				Form:  &CreateWithPostForm,
			},
		},
	},
}

// CreateForm is used to display the modal after clicking on the channel header
// button. It is also used for `/<bot_name>` sub-command's autocomplete.
var CreateForm = apps.Form{
	Title:  "Create Jira Issue",
	Header: "Create Jira Issue in TEST project",
	Icon:   "jira.png",
	Fields: []apps.Field{
		{
			Type:             "text",
			Name:             "product",
			AutocompleteHint: "pr1 | pr2 | pr3 | pr4 | pr5",
			IsRequired:       true,
		},
		{
			Type:             "text",
			Name:             "summary",
			AutocompleteHint: "Issue Summary",
			IsRequired:       true,
		},
		{
			Type:             "text",
			Name:             "description",
			AutocompleteHint: "Issue Description",
			IsRequired:       true,
		},
	},
	Submit: apps.NewCall("/create").WithExpand(apps.Expand{
		ActingUserAccessToken: apps.ExpandAll,
		ActingUser:            apps.ExpandAll,
		Channel:               apps.ExpandAll,
	}),
}

// CreateWithPostForm is used to display the modal after clicking on the post menu
// button.
var CreateWithPostForm = apps.Form{
	Title:  "Create Jira Issue With Post Message",
	Header: "Create Jira Issue in TEST project with post message",
	Icon:   "icon.png",
	Fields: []apps.Field{
		{
			Type:             "text",
			Name:             "product",
			AutocompleteHint: "pr1 | pr2 | pr3 | pr4 | pr5",
			IsRequired:       true,
		},
		{
			Type:             "text",
			Name:             "summary",
			AutocompleteHint: "Issue Summary",
			IsRequired:       true,
		},
		{
			Type:             "text",
			Name:             "additional",
			AutocompleteHint: "Additional information",
			IsRequired:       false,
		},
	},
	Submit: apps.NewCall("/create/withpost").WithExpand(apps.Expand{
		ActingUserAccessToken: apps.ExpandAll,
		ActingUser:            apps.ExpandAll,
		Channel:               apps.ExpandAll,
		Post:                  apps.ExpandAll,
		Team:                  apps.ExpandAll,
	}),
}
