package architecture

import (
	"errors"
	"fmt"
	"strings"

	"github.com/satori/go.uuid"
	// "log"
	"strconv"
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
	ClientId       string
	RawCommand     string
	Params         map[string]string
	WaitingForMore bool
	MoreToSend     bool
	Err            error
	Job            Job
}

func NewCommand(clientId string) Command {
	return Command{
		ClientId: clientId,
		MoreToSend: false,
		Params:     map[string]string{},
	}
}

func NewDefaultCommand(clientId string) Command {
	return Command{
		Name:       USE,
		ClientId: clientId,
		RawCommand: "use default",
		Params: map[string]string{
			"tube": "default",
		},
		WaitingForMore: false,
	}
}

func (command *Command) Copy() Command {
	paramsCopy := map[string]string{}
	for k, v := range command.Params {
		paramsCopy[k] = v
	}
	return Command{
		Name:           command.Name,
		ClientId: command.ClientId,
		RawCommand:     command.RawCommand,
		Params:         paramsCopy,
		WaitingForMore: command.WaitingForMore,
		MoreToSend:     command.MoreToSend,
		Err:            command.Err,
		Job:            command.Job,
	}
}

func (command *Command) createJobFromParams() error {
	pri, e1 := strconv.ParseInt(command.Params["pri"], 10, 0)
	if e1 != nil {
		return errors.New(BAD_FORMAT)
	}
	delay, e2 := strconv.ParseInt(command.Params["delay"], 10, 0)
	if e2 != nil {
		return errors.New(BAD_FORMAT)
	}
	ttr, e3 := strconv.ParseInt(command.Params["ttr"], 10, 0)
	if e3 != nil {
		return errors.New(BAD_FORMAT)
	}
	bytes, e4 := strconv.ParseInt(command.Params["bytes"], 10, 0)
	if e4 != nil {
		return errors.New(BAD_FORMAT)
	}

	if bytes > MAX_JOB_SIZE {
		return errors.New(JOB_TOO_BIG)
	} else if bytes != int64(len(command.Params["data"])) {
		return errors.New(EXPECTED_CRLF)
	}

	command.Job = *NewJob(
		uuid.NewV1().String(),
		pri,
		delay,
		ttr,
		bytes,
		command.Params["data"],
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
		command.Params = map[string]string{}
		command.RawCommand = rawCommand
		for i, paramName := range opts.Params {
			command.Params[paramName] = parts[i+1]
		}
		command.WaitingForMore = opts.WaitingForMore
		log.Debug("PROTOCOL command after parsing ", command)

		return !command.WaitingForMore, nil
	}

	// second round; PUT is the only valid command when WaitingForMore.
	// log.Println("GOT MORE", command)
	if command.Name == PUT {
		command.Params["data"] = rawCommand
		command.RawCommand += ("\r\n" + rawCommand)
		err := command.createJobFromParams()
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
			return false, strings.Join([]string{opts.Message, command.Params[opts.Param]}, " ")
		}

		// PUT, TOUCH
		return false, strings.Join([]string{opts.Message, command.Job.Id()}, " ")
	}

	// RESERVE, RESERVE_WITH_TIMEOUT
	if !command.MoreToSend {
		command.MoreToSend = true
		return true, fmt.Sprintf("RESERVED %s %d", command.Job.Id(), command.Job.Bytes)
	}

	return false, command.Job.Data
}
