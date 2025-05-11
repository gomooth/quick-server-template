package platformfilter

import "server-api/repository/types/platformtypes"

type OpenAPP struct {
	ID    uint
	IDs   []uint
	AppID string
	State *platformtypes.OpenAPPState
}
