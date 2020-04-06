package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// QuickchartPlugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type QuickchartPlugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	// BotId of the created bot account.
	// botID string
}

// MessageWillBePosted will check if there is any qc code in the body of the message to generate a chart
// returns (*model.Post, string)
func (p *QuickchartPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	p.API.LogInfo("starting the quickchart code")

	// configuration := p.getConfiguration()
	if len(post.Message) < 3 {
		return post, ""
	}

	postPrefix := post.Message[0:2]
	p.API.LogInfo("PREFIX", "prefix_content", postPrefix)

	if postPrefix != "qc" {
		return post, ""
	}

	p.API.LogInfo("POST STARTS WITH QC")

	if len(post.Message) > 3 {
		chartJSON := strings.TrimSpace(post.Message[2:])
		p.API.LogInfo("chart content", "formatted json", chartJSON)

		//jsonData := "{\"chart\": {\"type\": \"bar\", \"data\": {\"labels\": [\"Hello\", \"World\"], \"datasets\": [{\"label\": \"Foo\", \"data\": [1, 2]}]}}}"

		// response, err := http.Get("https://quickchart.io/chart?c={type:'bar',data:{labels:[2012,2013,2014,2015,2016],datasets:[{label:'Users',data:[120,60,50,180,120]}]}}")
		response, err := http.Post("https://quickchart.io/chart", "application/json", strings.NewReader(chartJSON))

		if err != nil {
			p.API.LogError("QuickChart could not get a proper response from Quickchart.io", "error", err)

			return nil, "Quickchart error"
		}

		p.API.LogInfo("Reading response from QC POST")
		data, _ := ioutil.ReadAll(response.Body)

		// ioutil.WriteFile("./test.png", data, 777)
		fileInfo, _ := p.API.UploadFile(data, post.ChannelId, "chart.png")
		if post.FileIds != nil {
			post.FileIds = append(post.FileIds, fileInfo.Id)
		} else {
			post.FileIds = []string{fileInfo.Id}
		}
	}

	p.API.LogInfo("finished the quickchart plugin code")
	// post.Message = configuration.TestConfigItem

	p.API.LogError(post.Message, "fileIds", post.FileIds)

	return post, ""
}

// MessageHasBeenPosted is just a text to see if the hooks of thisp[lugin get triggered
// func (p *QuickchartPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
// 	p.API.LogInfo("HasBeenposted triggered in quickchart")
// }
