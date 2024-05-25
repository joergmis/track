package track

import (
	_ "embed"
)

//go:generate sh -c "printf %s $(git describe --tags --dirty) > VERSION.txt"
//go:embed VERSION.txt
var VersionString string
