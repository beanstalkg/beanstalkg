package architecture

type CommandParseOptions struct {
	ExpectedLength int
	WaitingForMore bool
	Params         []string
	ParamInteger   []bool
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
			ParamInteger:   []bool{ false },
		},
		PUT: {
			Name:           PUT,
			ExpectedLength: 5,
			WaitingForMore: true,
			Params:         []string{"pri", "delay", "ttr", "bytes"},
			ParamInteger:   []bool{ true, true, true, true },
		},
		WATCH: {
			Name:           WATCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
			ParamInteger:   []bool{ false },
		},
		IGNORE: {
			Name:           IGNORE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
			ParamInteger:   []bool{ false },
		},
		RESERVE: {
			Name:           RESERVE,
			ExpectedLength: 1,
			WaitingForMore: false,
			Params:         []string{},
			ParamInteger:   []bool{},
		},
		RESERVE_WITH_TIMEOUT: {
			Name:           RESERVE_WITH_TIMEOUT,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"timeout"},
			ParamInteger:   []bool{ true },
		},
		DELETE: {
			Name:           DELETE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
			ParamInteger:   []bool{ false },
		},
		RELEASE: {
			Name:           RELEASE,
			ExpectedLength: 4,
			WaitingForMore: false,
			Params:         []string{"id", "pri", "delay"},
			ParamInteger:   []bool{ false, true, true},
		},
		BURY: {
			Name:           BURY,
			ExpectedLength: 3,
			WaitingForMore: false,
			Params:         []string{"id", "pri"},
			ParamInteger:   []bool{ false, true },
		},
		TOUCH: {
			Name:           TOUCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
			ParamInteger:   []bool{ false },
		},
		KICK: {
			Name:           KICK,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"bound"},
			ParamInteger:   []bool{ true },
		},
		KICK_JOB: {
			Name:           KICK_JOB,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
			ParamInteger:   []bool{ false },
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
