package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/bosh-warden-cpi/action"
	fakestem "github.com/cppforlife/bosh-warden-cpi/stemcell/fakes"
	bwcvm "github.com/cppforlife/bosh-warden-cpi/vm"
	fakevm "github.com/cppforlife/bosh-warden-cpi/vm/fakes"
)

var _ = Describe("CreateVM", func() {
	var (
		stemcellFinder *fakestem.FakeFinder
		vmCreator      *fakevm.FakeCreator
		action         CreateVM
	)

	BeforeEach(func() {
		stemcellFinder = &fakestem.FakeFinder{}
		vmCreator = &fakevm.FakeCreator{}
		action = NewCreateVM(stemcellFinder, vmCreator)
	})

	Describe("Run", func() {
		var (
			stemcellCID  StemcellCID
			resourcePool VMCloudProperties
			networks     Networks
			diskLocality []DiskCID
			env          Environment
		)

		BeforeEach(func() {
			stemcellCID = StemcellCID("fake-stemcell-id")
			resourcePool = VMCloudProperties{}
			networks = Networks{"fake-net-name": Network{IP: "fake-ip"}}
			diskLocality = []DiskCID{"fake-disk-id"}
			env = Environment{"fake-env-key": "fake-env-value"}
		})

		It("tries to find stemcell with given stemcell cid", func() {
			stemcellFinder.FindFound = true
			vmCreator.CreateVM = fakevm.NewFakeVM("fake-vm-id")

			_, err := action.Run("fake-agent-id", stemcellCID, resourcePool, networks, diskLocality, env)
			Expect(err).ToNot(HaveOccurred())

			Expect(stemcellFinder.FindID).To(Equal("fake-stemcell-id"))
		})

		Context("when stemcell is found with given stemcell cid", func() {
			var (
				stemcell *fakestem.FakeStemcell
			)

			BeforeEach(func() {
				stemcell = fakestem.NewFakeStemcell("fake-stemcell-id")
				stemcellFinder.FindStemcell = stemcell
				stemcellFinder.FindFound = true
			})

			It("returns id for created VM", func() {
				vmCreator.CreateVM = fakevm.NewFakeVM("fake-vm-id")

				id, err := action.Run("fake-agent-id", stemcellCID, resourcePool, networks, diskLocality, env)
				Expect(err).ToNot(HaveOccurred())
				Expect(id).To(Equal(VMCID("fake-vm-id")))
			})

			It("creates VM with requested agent ID, stemcell, and networks", func() {
				vmCreator.CreateVM = fakevm.NewFakeVM("fake-vm-id")

				_, err := action.Run("fake-agent-id", stemcellCID, resourcePool, networks, diskLocality, env)
				Expect(err).ToNot(HaveOccurred())

				Expect(vmCreator.CreateAgentID).To(Equal("fake-agent-id"))
				Expect(vmCreator.CreateStemcell).To(Equal(stemcell))
				Expect(vmCreator.CreateNetworks).To(Equal(networks.AsVMNetworks()))
				Expect(vmCreator.CreateEnvironment).To(Equal(
					bwcvm.Environment{"fake-env-key": "fake-env-value"},
				))
			})

			It("returns error if creating VM fails", func() {
				vmCreator.CreateErr = errors.New("fake-create-err")

				id, err := action.Run("fake-agent-id", stemcellCID, resourcePool, networks, diskLocality, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-create-err"))
				Expect(id).To(Equal(VMCID("")))
			})
		})

		Context("when stemcell is not found with given cid", func() {
			It("returns error because VM cannot be created without a stemcell", func() {
				stemcellFinder.FindFound = false

				id, err := action.Run("fake-agent-id", stemcellCID, resourcePool, networks, diskLocality, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Expected to find stemcell"))
				Expect(id).To(Equal(VMCID("")))
			})
		})

		Context("when stemcell finding fails", func() {
			It("returns error because VM cannot be created without a stemcell", func() {
				stemcellFinder.FindErr = errors.New("fake-find-err")

				id, err := action.Run("fake-agent-id", stemcellCID, resourcePool, networks, diskLocality, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-find-err"))
				Expect(id).To(Equal(VMCID("")))
			})
		})
	})
})
