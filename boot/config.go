package boot

import (
	"fmt"
)

type Param struct {
	ConfigFilename  string
	RegisterServers []ServerType

	CMDParam *CMDBootParam
}

func (m Param) LogCategory() string {
	if len(m.RegisterServers) > 1 {
		return ""
	}
	if m.CMDParam != nil {
		return "command"
	}

	switch m.RegisterServers[0] {
	case InitServerTypeOpenAPI:
		return "openapi"
	case InitServerTypeCronjob:
		return "cronjob"
	case InitServerTypeConsumer:
		return "consumer"
	default:
		return ""
	}
}

type ServerType int

const (
	InitServerTypeWeb ServerType = iota
	InitServerTypeOpenAPI
	InitServerTypeCronjob
	InitServerTypeConsumer
)

type CMDBootParam struct {
	Name    string
	Timeout int
	Args    []string
}

type FlagSlice []string

func (f *FlagSlice) String() string {
	return fmt.Sprintf("%v", []string(*f))
}

func (f *FlagSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}
