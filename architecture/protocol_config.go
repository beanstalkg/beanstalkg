package architecture

type CommandParseOptions struct {
	ExpectedLength int
	WaitingForMore bool
	Params         []string
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
		},
		PUT: {
			Name:           PUT,
			ExpectedLength: 5,
			WaitingForMore: true,
			Params:         []string{"pri", "delay", "ttr", "bytes"},
		},
		WATCH: {
			Name:           WATCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
		},
		IGNORE: {
			Name:           IGNORE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
		},
		RESERVE: {
			Name:           RESERVE,
			ExpectedLength: 1,
			WaitingForMore: false,
			Params:         []string{},
		},
		RESERVE_WITH_TIMEOUT: {
			Name:           RESERVE_WITH_TIMEOUT,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"timeout"},
		},
		DELETE: {
			Name:           DELETE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
		},
		RELEASE: {
			Name:           RELEASE,
			ExpectedLength: 4,
			WaitingForMore: false,
			Params:         []string{"id", "pri", "delay"},
		},
		BURY: {
			Name:           BURY,
			ExpectedLength: 3,
			WaitingForMore: false,
			Params:         []string{"id", "pri"},
		},
		TOUCH: {
			Name:           TOUCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
		},
		QUIT: {
			Name:           QUIT,
			ExpectedLength: 1,
			WaitingForMore: false,
			Params:         []string{},
		KICK: {
			Name:           KICK,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"bound"},
		},
		KICK_JOB: {
			Name:           KICK_JOB,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
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
