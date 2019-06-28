package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
	"github.com/robfig/cron"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	// a job for pre-calculating channel recommendations for users.
	preCalcJob *cron.Cron
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

// OnActivate will be run on plugin activation.
func (p *Plugin) OnActivate() error {
	p.API.RegisterCommand(getCommand())
	err := p.initStore()
	if err != nil {
		return err
	}

	c := cron.New()

	if err := c.AddFunc("@daily", func() { // Run once a day
		p.preCalculateRecommendations()
	}); err != nil {
		return err
	}

	c.Start()

	p.preCalcJob = c
	return nil
}

// OnDeactivate will be run on plugin deactivation.
func (p *Plugin) OnDeactivate() error {
	p.preCalcJob.Stop()
	return nil
}
