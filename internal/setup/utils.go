package setup

import (
	"strings"

	teaconst "github.com/dashenmiren/EdgeAPI/internal/const"
)

func ComposeSQLVersion() string {
	var version = teaconst.Version
	if len(teaconst.SQLVersion) == 0 {
		return version
	}

	if strings.Count(version, ".") <= 2 {
		return version + "." + teaconst.SQLVersion
	}
	return version
}
