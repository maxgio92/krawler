package rpm

import log "github.com/sirupsen/logrus"

const (
	metadataPath      = "repodata/repomd.xml"
	metadataDataXPath = "//repomd/data"
	dataPackageXPath  = "//package"
	primary           = "primary"
)

var (
	logger = log.New()
)
