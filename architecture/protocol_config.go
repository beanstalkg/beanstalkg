package architecture

import (
	"errors"
	"strconv"
)

type CommandParseOptions struct {
	ExpectedLength int
	WaitingForMore bool
	Params         []cmdPropertySetter
	Name           CommandName
}

type CommandReplyOptions struct {
	Result   bool
	Message  string
	Param    string
	UseJobID bool
}

type cmdPropertySetter func(*Command, string) error
type cmdPropertyGetter func(*Command) string

var cmdParseOptions map[CommandName]CommandParseOptions
var cmdReplyOptions map[CommandName]CommandReplyOptions
var cmdGetterFuncs map[string]cmdPropertyGetter

const (
	COUNT string = "count"
	TUBE         = "tube"
)

func init() {
	cmdParseOptions = map[CommandName]CommandParseOptions{
		USE: {
			Name:           USE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetTube},
		},
		PUT: {
			Name:           PUT,
			ExpectedLength: 5,
			WaitingForMore: true,
			Params:         []cmdPropertySetter{cmdSetPriority, cmdSetDelay, cmdSetTTR, cmdSetBytes},
		},
		WATCH: {
			Name:           WATCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetTube},
		},
		IGNORE: {
			Name:           IGNORE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetTube},
		},
		RESERVE: {
			Name:           RESERVE,
			ExpectedLength: 1,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{},
		},
		RESERVE_WITH_TIMEOUT: {
			Name:           RESERVE_WITH_TIMEOUT,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetTimeout},
		},
		DELETE: {
			Name:           DELETE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetID},
		},
		RELEASE: {
			Name:           RELEASE,
			ExpectedLength: 4,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetID, cmdSetPriority, cmdSetDelay},
		},
		BURY: {
			Name:           BURY,
			ExpectedLength: 3,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetID, cmdSetPriority},
		},
		TOUCH: {
			Name:           TOUCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetID},
		},
		QUIT: {
			Name:           QUIT,
			ExpectedLength: 1,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{},
		},
		KICK: {
			Name:           KICK,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetBound},
		},
		KICK_JOB: {
			Name:           KICK_JOB,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []cmdPropertySetter{cmdSetID},
		},
	}

	cmdReplyOptions = map[CommandName]CommandReplyOptions{
		USE: {
			Result:   false,
			Message:  "USING",
			Param:    TUBE,
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
			Param:    COUNT,
			UseJobID: false,
		},
		IGNORE: {
			Result:   false,
			Message:  "WATCHING",
			Param:    COUNT,
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

	cmdGetterFuncs = map[string]cmdPropertyGetter{
		COUNT: cmdGetCount,
		TUBE:  cmdGetTube,
	}
}

// Auxiliary setters.
func cmdSetTube(command *Command, s string) error {
	command.Tube = s

	return nil
}

func cmdSetID(command *Command, s string) error {
	command.ID = s

	return nil
}

func cmdSetPriority(command *Command, s string) (err error) {
	command.Priority, err = cmdParseInt64(s)

	return
}

func cmdSetDelay(command *Command, s string) (err error) {
	command.Delay, err = cmdParseInt64(s)

	return
}

func cmdSetTTR(command *Command, s string) (err error) {
	command.TTR, err = cmdParseInt64(s)

	return
}

func cmdSetBytes(command *Command, s string) (err error) {
	command.Bytes, err = cmdParseInt64(s)

	return
}

func cmdSetTimeout(command *Command, s string) (err error) {
	command.Timeout, err = cmdParseInt64(s)

	return
}

func cmdSetBound(command *Command, s string) (err error) {
	command.Bound, err = cmdParseInt64(s)

	return
}

func cmdParseInt64(s string) (i int64, err error) {
	if i, err = strconv.ParseInt(s, 10, 0); err != nil {
		return 0, errors.New(BAD_FORMAT)
	}

	return i, nil
}

// Auxiliary getters.
func cmdGetCount(command *Command) string {
	return strconv.FormatInt(command.Count, 10)
}

func cmdGetTube(command *Command) string {
	return command.Tube
}
