package architecture

type CommandParseOptions struct {
	ExpectedLength int
	WaitingForMore bool
	Params         []string
	Name           CommandName
}

var cmdParseOptions map[CommandName]CommandParseOptions

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
}
