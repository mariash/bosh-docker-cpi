package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/bosh-warden-cpi/action"
	fakedisk "github.com/cppforlife/bosh-warden-cpi/disk/fakes"
)

var _ = Describe("DeleteDisk", func() {
	var (
		diskFinder *fakedisk.FakeFinder
		action     DeleteDisk
	)

	BeforeEach(func() {
		diskFinder = &fakedisk.FakeFinder{}
		action = NewDeleteDisk(diskFinder)
	})

	Describe("Run", func() {
		It("tries to find disk with given disk cid", func() {
			_, err := action.Run("fake-disk-id")
			Expect(err).ToNot(HaveOccurred())

			Expect(diskFinder.FindID).To(Equal("fake-disk-id"))
		})

		Context("when disk is found with given disk cid", func() {
			var (
				disk *fakedisk.FakeDisk
			)

			BeforeEach(func() {
				disk = fakedisk.NewFakeDisk("fake-disk-id")
				diskFinder.FindDisk = disk
				diskFinder.FindFound = true
			})

			It("deletes disk", func() {
				_, err := action.Run("fake-disk-id")
				Expect(err).ToNot(HaveOccurred())

				Expect(disk.DeleteCalled).To(BeTrue())
			})

			It("returns error if deleting disk fails", func() {
				disk.DeleteErr = errors.New("fake-delete-err")

				_, err := action.Run("fake-disk-id")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-delete-err"))
			})
		})

		Context("when disk is not found with given cid", func() {
			It("does not return error", func() {
				diskFinder.FindFound = false

				_, err := action.Run("fake-disk-id")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when disk finding fails", func() {
			It("does not return error", func() {
				diskFinder.FindErr = errors.New("fake-find-err")

				_, err := action.Run("fake-disk-id")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-find-err"))
			})
		})
	})
})
