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
		USE: CommandParseOptions{
			Name:           USE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
		},
		PUT: CommandParseOptions{
			Name:           PUT,
			ExpectedLength: 5,
			WaitingForMore: true,
			Params:         []string{"pri", "delay", "ttr", "bytes"},
		},
		WATCH: CommandParseOptions{
			Name:           WATCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
		},
		IGNORE: CommandParseOptions{
			Name:           IGNORE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"tube"},
		},
		RESERVE: CommandParseOptions{
			Name:           RESERVE,
			ExpectedLength: 1,
			WaitingForMore: false,
			Params:         []string{},
		},
		RESERVE_WITH_TIMEOUT: CommandParseOptions{
			Name:           RESERVE_WITH_TIMEOUT,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"timeout"},
		},
		DELETE: CommandParseOptions{
			Name:           DELETE,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
		},
		RELEASE: CommandParseOptions{
			Name:           RELEASE,
			ExpectedLength: 4,
			WaitingForMore: false,
			Params:         []string{"id", "pri", "delay"},
		},
		BURY: CommandParseOptions{
			Name:           BURY,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id", "pri"},
		},
		TOUCH: CommandParseOptions{
			Name:           TOUCH,
			ExpectedLength: 2,
			WaitingForMore: false,
			Params:         []string{"id"},
		},
	}

	cmdReplyOptions = map[CommandName]CommandReplyOptions{
		USE: CommandReplyOptions{
			Result:   false,
			Message:  "USING",
			Param:    "tube",
			UseJobID: false,
		},
		PUT: CommandReplyOptions{
			Result:   false,
			Message:  "INSERTED",
			Param:    "",
			UseJobID: true,
		},
		WATCH: CommandReplyOptions{
			Result:   false,
			Message:  "WATCHING",
			Param:    "count",
			UseJobID: false,
		},
		IGNORE: CommandReplyOptions{
			Result:   false,
			Message:  "WATCHING",
			Param:    "count",
			UseJobID: false,
		},
		DELETE: CommandReplyOptions{
			Result:   false,
			Message:  "DELETED",
			Param:    "",
			UseJobID: false,
		},
		RELEASE: CommandReplyOptions{
			Result:   false,
			Message:  "RELEASED",
			Param:    "",
			UseJobID: false,
		},
		BURY: CommandReplyOptions{
			Result:   false,
			Message:  "BURIED",
			Param:    "",
			UseJobID: false,
		},
		TOUCH: CommandReplyOptions{
			Result:   false,
			Message:  "INSERTED",
			Param:    "",
			UseJobID: true,
		},
	}
}
