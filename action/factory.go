package action

import (
	"github.com/fsouza/go-dockerclient"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	bwcaction "github.com/cppforlife/bosh-warden-cpi/action"

	"github.com/mariash/bosh-docker-cpi/action/container"
	cfg "github.com/mariash/bosh-docker-cpi/config"
)

type factory struct {
	actions map[string]bwcaction.Action
}

func NewFactory(client *docker.Client, config cfg.Config, fs boshsys.FileSystem, logger boshlog.Logger) bwcaction.Factory {
	containerCreator := container.NewCreator(client, config, logger)
	containerFinder := container.NewFinder(client)
	settingsUpdaterFactory := container.NewSettingsUpdaterFactory(client, config)

	return &factory{
		actions: map[string]bwcaction.Action{
			"create_stemcell": NewCreateStemcell(client),
			"create_vm":       NewCreateVM(containerCreator, settingsUpdaterFactory),
			"has_vm":          NewHasVM(containerFinder),
			"delete_vm":       NewDeleteVM(client),
		},
	}
}

func (f *factory) Create(method string) (bwcaction.Action, error) {
	action, found := f.actions[method]
	if !found {
		return nil, bosherr.New("Could not create action with method %s", method)
	}

	return action, nil
}
