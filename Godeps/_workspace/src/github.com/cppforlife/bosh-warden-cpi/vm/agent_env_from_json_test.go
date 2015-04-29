package vm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/bosh-warden-cpi/vm"
)

var _ = Describe("NewAgentEnvFromJSON", func() {
	Context("when json is valid", func() {
		It("returns agent env parsed from JSON string", func() {
			agentEnvJSON := `{
        "agent_id": "fake-agent-id",

        "vm": {
          "name": "fake-vm-name",
          "id": "fake-vm-id"
        },

        "networks": {
          "fake-net-name": {
            "ip":      "fake-ip",
            "netmask": "fake-netmask",
            "gateway": "fake-gateway",

            "dns":     ["fake-dns"],
            "default": ["fake-default"],

            "mac": "fake-mac",

            "cloud_properties": {"fake-cp-key": "fake-cp-value"}
          }
        },

        "disks": {
          "persistent": {
            "fake-persistent-id": "fake-persistent-path"
          }
        },

        "env": {"fake-env-key": "fake-env-value"}
      }`

			expectedAgentEnv := AgentEnv{
				AgentID: "fake-agent-id",

				VM: VMSpec{
					Name: "fake-vm-name",
					ID:   "fake-vm-id",
				},

				Networks: NetworksSpec{
					"fake-net-name": NetworkSpec{
						IP:      "fake-ip",
						Netmask: "fake-netmask",
						Gateway: "fake-gateway",

						DNS:     []string{"fake-dns"},
						Default: []string{"fake-default"},

						MAC: "fake-mac",

						CloudProperties: map[string]interface{}{
							"fake-cp-key": "fake-cp-value",
						},
					},
				},

				Disks: DisksSpec{
					Persistent: PersistentSpec{
						"fake-persistent-id": "fake-persistent-path",
					},
				},

				Env: map[string]interface{}{
					"fake-env-key": "fake-env-value",
				},
			}

			agentEnv, err := NewAgentEnvFromJSON([]byte(agentEnvJSON))
			Expect(err).ToNot(HaveOccurred())
			Expect(agentEnv).To(Equal(expectedAgentEnv))
		})
	})

	Context("when json is not valid", func() {
		It("returns error", func() {
			agentEnv, err := NewAgentEnvFromJSON([]byte(`-`))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid character"))
			Expect(agentEnv).To(Equal(AgentEnv{}))
		})
	})
})
