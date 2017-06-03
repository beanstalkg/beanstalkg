package architecture

import (
	"strconv"
)

type CommandParseOptions struct {
	ExpectedLength int
	WaitingForMore bool
	Params         []string
	ParamValidators []interface{}
	Name           CommandName
}

type CommandReplyOptions struct {
	Result   bool
	Message  string
	Param    string
	UseJobID bool
}

var cmdParseOptions map[CommandName]CommandParseOptions
var cmdReplyOptions map[CommandName]CommandReplyOptions

func init() {
	cmdParseOptions = map[CommandName]CommandParseOptions{
		USE: {
			Name:           USE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
			ParamValidators: []interface{}{ strconv.Atoi },
		},
		PUT: {
			Name:           PUT,
			ExpectedLength: 5,
			WaitingForMore: true,
			Params:         []string{"pri", "delay", "ttr", "bytes"},
			ParamValidators: []interface{}{ strconv.Atoi, strconv.Atoi, strconv.Atoi, strconv.Atoi },
		},
		WATCH: {
			Name:           WATCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
			ParamValidators: []interface{}{ nil },
		},
		IGNORE: {
			Name:           IGNORE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
			ParamValidators: []interface{}{ nil },
		},
		RESERVE: {
			Name:           RESERVE,
			ExpectedLength: 1,
			WaitingForMore: false,
			Params:         []string{},
			ParamValidators: []interface{}{},
		},
		RESERVE_WITH_TIMEOUT: {
			Name:           RESERVE_WITH_TIMEOUT,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"timeout"},
			ParamValidators: []interface{}{ strconv.Atoi },
		},
		DELETE: {
			Name:           DELETE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
			ParamValidators: []interface{}{ nil },
		},
		RELEASE: {
			Name:           RELEASE,
			ExpectedLength: 4,
			WaitingForMore: false,
			Params:         []string{"id", "pri", "delay"},
			ParamValidators: []interface{}{ nil, strconv.Atoi, strconv.Atoi },
		},
		BURY: {
			Name:           BURY,
			ExpectedLength: 3,
			WaitingForMore: false,
			Params:         []string{"id", "pri"},
			ParamValidators: []interface{}{ nil, strconv.Atoi },
		},
		TOUCH: {
			Name:           TOUCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
			ParamValidators: []interface{}{ nil },
		},
		KICK: {
			Name:           KICK,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"bound"},
			ParamValidators: []interface{}{ strconv.Atoi },
		},
		KICK_JOB: {
			Name:           KICK_JOB,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
			ParamValidators: []interface{}{ nil },
		},
	}

	cmdReplyOptions = map[CommandName]CommandReplyOptions{
		USE: {
			Result:   false,
			Message:  "USING",
			Param:    "tube",
			UseJobID: false,
		},
		PUT: {
			Result:   false,
			Message:  "INSERTED",
			Param:    "",
			UseJobID: true,
		},
		WATCH: {
			Result:   false,
			Message:  "WATCHING",
			Param:    "count",
			UseJobID: false,
		},
		IGNORE: {
			Result:   false,
			Message:  "WATCHING",
			Param:    "count",
			UseJobID: false,
		},
		DELETE: {
			Result:   false,
			Message:  "DELETED",
			Param:    "",
			UseJobID: false,
		},
		RELEASE: {
			Result:   false,
			Message:  "RELEASED",
			Param:    "",
			UseJobID: false,
		},
		BURY: {
			Result:   false,
			Message:  "BURIED",
			Param:    "",
			UseJobID: false,
		},
		TOUCH: {
			Result:   false,
			Message:  "INSERTED",
			Param:    "",
			UseJobID: true,
		},
		KICK: {
			Result:   false,
			Message:  "KICKED",
			Param:    "",
			UseJobID: false,
		},
		KICK_JOB: {
			Result:   false,
			Message:  "KICKED",
			Param:    "",
			UseJobID: true,
		},
	}
}
