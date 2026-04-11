package pfilter

import "server-api/repository/platform/pattr"

type OpenAPP struct {
	ID    uint
	IDs   []uint
	AppID string
	State *pattr.OpenAPPState
}
