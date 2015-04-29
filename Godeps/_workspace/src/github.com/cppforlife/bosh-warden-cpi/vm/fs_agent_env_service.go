package vm

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

type fsAgentEnvService struct {
	wardenFileService WardenFileService
	settingsPath      string
	logger            boshlog.Logger
	logTag            string
}

func NewFSAgentEnvService(
	wardenFileService WardenFileService,
	logger boshlog.Logger,
) AgentEnvService {
	return fsAgentEnvService{
		wardenFileService: wardenFileService,
		settingsPath:      "/var/vcap/bosh/warden-cpi-agent-env.json",
		logger:            logger,
		logTag:            "FSAgentEnvService",
	}
}

func (s fsAgentEnvService) Fetch() (AgentEnv, error) {
	var agentEnv AgentEnv

	contents, err := s.wardenFileService.Download(s.settingsPath)
	if err != nil {
		return AgentEnv{}, bosherr.WrapError(err, "Downloading agent env from container")
	}

	err = json.Unmarshal(contents, &agentEnv)
	if err != nil {
		return AgentEnv{}, bosherr.WrapError(err, "Unmarshalling agent env")
	}

	s.logger.Debug(s.logTag, "Fetched agent env: %#v", agentEnv)

	return agentEnv, nil
}

func (s fsAgentEnvService) Update(agentEnv AgentEnv) error {
	s.logger.Debug(s.logTag, "Updating agent env: %#v", agentEnv)

	jsonBytes, err := json.Marshal(agentEnv)
	if err != nil {
		return bosherr.WrapError(err, "Marshalling agent env")
	}

	return s.wardenFileService.Upload(s.settingsPath, jsonBytes)
}
