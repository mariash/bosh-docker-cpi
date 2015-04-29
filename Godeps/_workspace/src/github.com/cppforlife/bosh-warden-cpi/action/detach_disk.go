package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	bwcdisk "github.com/cppforlife/bosh-warden-cpi/disk"
	bwcvm "github.com/cppforlife/bosh-warden-cpi/vm"
)

type DetachDisk struct {
	vmFinder   bwcvm.Finder
	diskFinder bwcdisk.Finder
}

func NewDetachDisk(vmFinder bwcvm.Finder, diskFinder bwcdisk.Finder) DetachDisk {
	return DetachDisk{
		vmFinder:   vmFinder,
		diskFinder: diskFinder,
	}
}

func (a DetachDisk) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {
	vm, found, err := a.vmFinder.Find(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding VM '%s'", vmCID)
	}

	if !found {
		return nil, bosherr.New("Expected to find VM '%s'", vmCID)
	}

	disk, found, err := a.diskFinder.Find(string(diskCID))
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding disk '%s'", diskCID)
	}

	if !found {
		return nil, bosherr.New("Expected to find disk '%s'", diskCID)
	}

	err = vm.DetachDisk(disk)
	if err != nil {
		return nil, bosherr.WrapError(err, "Detaching disk '%s' to VM '%s'", diskCID, vmCID)
	}

	return nil, nil
}
