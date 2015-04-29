package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	bwcstem "github.com/cppforlife/bosh-warden-cpi/stemcell"
)

type DeleteStemcell struct {
	stemcellFinder bwcstem.Finder
}

func NewDeleteStemcell(stemcellFinder bwcstem.Finder) DeleteStemcell {
	return DeleteStemcell{stemcellFinder: stemcellFinder}
}

func (a DeleteStemcell) Run(stemcellCID StemcellCID) (interface{}, error) {
	stemcell, found, err := a.stemcellFinder.Find(string(stemcellCID))
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding stemcell '%s'", stemcellCID)
	}

	if found {
		err := stemcell.Delete()
		if err != nil {
			return nil, bosherr.WrapError(err, "Deleting stemcell '%s'", stemcellCID)
		}
	}

	return nil, nil
}
