package common

import (
	"sync"
)

var (
	// RequestIdMap utils.Uuid.GetId()
	RequestIdMap = new(sync.Map)
)
