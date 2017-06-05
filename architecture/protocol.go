package architecture

import (
	"errors"
	"fmt"
	"strings"

	"github.com/satori/go.uuid"
	// "log"
)

// TODO extract protocol error messages into a error helper

type CommandName string

const (
	USE                  CommandName = "use"
	PUT                              = "put"
	WATCH                            = "watch"
	IGNORE                           = "ignore"
	RESERVE                          = "reserve"
	RESERVE_WITH_TIMEOUT             = "reserve-with-timeout"
	DELETE                           = "delete"
	RELEASE                          = "release"
	BURY                             = "bury"
	TOUCH                            = "touch"
	QUIT                             = "quit"
	KICK                             = "kick"
	KICK_JOB                         = "kick-job"
)

const (
	BAD_FORMAT      string = "BAD_FORMAT"
	UNKNOWN_COMMAND        = "UNKNOWN_COMMAND"
	NOT_IGNORED            = "NOT_IGNORED"
	NOT_FOUND              = "NOT_FOUND"
	EXPECTED_CRLF          = "EXPECTED_CRLF"
	JOB_TOO_BIG            = "JOB_TOO_BIG"
	TIMED_OUT              = "TIMED_OUT"
)

const MAX_JOB_SIZE int64 = 65536 // 2^16

type Command struct {
	Name           CommandName
	RawCommand     string
	WaitingForMore bool
	MoreToSend     bool
	Err            error
	Job            Job

	Tube     string
	ID       string
	Priority int64
	Delay    int64
	TTR      int64
	Bytes    int64
	Timeout  int64
	Bound    int64
	Count    int64
}

func NewCommand() Command {
	return Command{
		MoreToSend: false,
	}
}

func NewDefaultCommand() Command {
	return Command{
		Name:           USE,
		RawCommand:     "use default",
		WaitingForMore: false,
	}
}

func (command *Command) Copy() Command {
	return Command{
		Name:           command.Name,
		RawCommand:     command.RawCommand,
		WaitingForMore: command.WaitingForMore,
		MoreToSend:     command.MoreToSend,
		Err:            command.Err,
		Job:            command.Job,
		Tube:           command.Tube,
		ID:             command.ID,
		Priority:       command.Priority,
		Delay:          command.Delay,
		TTR:            command.TTR,
		Bytes:          command.Bytes,
		Timeout:        command.Timeout,
		Bound:          command.Bound,
	}
}

func (command *Command) createJobFromParams(payload string) error {
	bytes := command.Bytes
	if bytes > MAX_JOB_SIZE {
		return errors.New(JOB_TOO_BIG)
	} else if bytes != int64(len(payload)) {
		return errors.New(EXPECTED_CRLF)
	}

	command.Job = *NewJob(
		uuid.NewV1().String(),
		command.Priority,
		command.Delay,
		command.TTR,
		command.Bytes,
		payload,
	)
	// log.Println("PROTOCOL new job: ", command.Job)
	return nil
}

// Parse keeps track of the state of the command and it will be called multiple times for commands such as  'put'
// where the length of command spans multiple lines
func (command *Command) Parse(rawCommand string) (bool, error) {
	// check if this command has already been waiting for a second round
	if !command.WaitingForMore {
		// first round
		parts := strings.Split(rawCommand, " ")
		name := CommandName(strings.ToLower(parts[0]))

		// Check for valid command.
		opts, ok := cmdParseOptions[name]
		if !ok {
			// Unknown command.
			command.Err = errors.New(UNKNOWN_COMMAND)
			return true, command.Err
		}

		command.Name = name

		// Were we given the proper number of parameters?
		if len(parts) != opts.ExpectedLength {
			command.Err = errors.New(BAD_FORMAT)
			return true, command.Err
		}

		// Store command info.  For future logging, maybe?
		command.RawCommand = rawCommand
		command.WaitingForMore = opts.WaitingForMore
		for i, setterFunc := range opts.Params {
			if err := setterFunc(command, parts[i+1]); err != nil {
				return !command.WaitingForMore, err
			}
		}
		log.Debug("PROTOCOL command after parsing ", command)

		return !command.WaitingForMore, nil
	}

	// second round; PUT is the only valid command when WaitingForMore.
	// log.Println("GOT MORE", command)
	if command.Name == PUT {
		command.RawCommand += ("\r\n" + rawCommand)
		err := command.createJobFromParams(rawCommand)
		// log.Println("GOT MORE PUT", c, err)
		command.Err = err
		return true, err
	}

	// Unknown command, secondround.
	return true, nil
}

func (command *Command) Reply() (bool, string) {
	if err := command.Err; err != nil {
		return false, err.Error()
	}

	// cmdReplyOptions = cmdParseOptions - RESERVE, RESERVE_WITH_TIMEOUT
	if opts, ok := cmdReplyOptions[command.Name]; ok {
		// Take care of everything except PUT, RESERVE, RESERVE_WITH_TIMEOUT, and TOUCH
		if !opts.UseJobID {
			// DELETE, RELEASE, BURY
			if len(opts.Param) == 0 {
				return false, opts.Message
			}

			// USE, WATCH, IGNORE
			if getter, ok := cmdGetterFuncs[opts.Param]; ok {
				return false, strings.Join([]string{opts.Message, getter(command)}, " ")
			}

			return false, opts.Message
		}

		// PUT, TOUCH
		return false, strings.Join([]string{opts.Message, command.ID}, " ")
	}

	// RESERVE, RESERVE_WITH_TIMEOUT
	if !command.MoreToSend {
		command.MoreToSend = true
		return true, fmt.Sprintf("RESERVED %s %d", command.ID, command.Job.Bytes)
	}

	return false, command.Job.Data
}
