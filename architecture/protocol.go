package architecture

import (
	"errors"
	"strings"
	// "fmt"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
)

type CommandName string

const (
	USE     CommandName = "use"
	PUT                 = "put"
	RESERVE             = "reserve"
)

type Command struct {
	Name           CommandName
	RawCommand     string
	Params         map[string]string
	WaitingForMore bool
	Err            error
	Job            *Job
}

func NewDefaultCommand() Command {
	return Command{
		USE,
		"use default",
		map[string]string{
			"tube": "default",
		},
		false,
		nil,
		nil,
	}
}

func (command *Command) createJobFromParams() (*Job, error) {
	pri, e1 := strconv.ParseInt(command.Params["pri"], 10, 0)
	if e1 != nil {
		return nil, e1
	}
	delay, e2 := strconv.ParseInt(command.Params["delay"], 10, 0)
	if e2 != nil {
		return nil, e2
	}
	ttr, e3 := strconv.ParseInt(command.Params["ttr"], 10, 0)
	if e3 != nil {
		return nil, e3
	}
	bytes, e4 := strconv.ParseInt(command.Params["bytes"], 10, 0)
	if e4 != nil {
		return nil, e4
	}

	command.Job = NewJob(
		uuid.NewV1().String(),
		pri,
		delay,
		ttr,
		bytes,
		command.Params["data"],
	)
	log.Println("PROTOCOL new job: ", command.Job)
	return command.Job, nil
}

// Parse keeps track of the state of the command and it will be called multiple times for commands such as  'put'
// where the length of command spans multiple lines
func (command *Command) Parse(rawCommand string) (bool, error) {
	// check if this command has already been waiting for a second round
	if !command.WaitingForMore {
		// first round
		parts := strings.Split(rawCommand, " ")
		switch strings.ToLower(parts[0]) {
		case "use":
			// use <tube>\r\n
			command.Name = USE
			if len(parts) > 2 {
				command.Err = errors.New("BAD_FORMAT")
				return true, command.Err
			}
			command.RawCommand = rawCommand
			command.Params = map[string]string{
				"tube": parts[1],
			}
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		case "put":
			// put <pri> <delay> <ttr> <bytes>\r\n <data>\r\n
			command.Name = PUT
			if len(parts) != 5 {
				command.Err = errors.New("BAD_FORMAT")
				return true, command.Err
			}
			command.RawCommand = rawCommand
			command.Params = map[string]string{
				"pri":   parts[1],
				"delay": parts[2],
				"ttr":   parts[3],
				"bytes": parts[4],
			}
			command.WaitingForMore = true
			return !command.WaitingForMore, nil
		case "watch":
			//
		case "reserve":
			// reserve\r\n
			command.Name = RESERVE
			if len(parts) != 1 {
				command.Err = errors.New("BAD_FORMAT")
				return true, command.Err
			}
			command.RawCommand = rawCommand
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		default:
			command.Err = errors.New("UNKNOWN_COMMAND")
			return true, command.Err
		}
	} else {
		// second round
		// log.Println("GOT MORE", command)
		switch command.Name {
		case PUT:
			command.Params["data"] = rawCommand
			_, err := command.createJobFromParams()
			// log.Println("GOT MORE PUT", c, err)
			command.Err = err
			return true, err
		}
	}
	return true, nil
}

func (command *Command) Reply() string {
	/**
	 	- "BAD_FORMAT\r\n" The client sent a command line that was not well-formed.
		   This can happen if the line does not end with \r\n, if non-numeric
		   characters occur where an integer is expected, if the wrong number of
		   arguments are present, or if the command line is mal-formed in any other
		   way.

		 - "UNKNOWN_COMMAND\r\n" The client sent a command that the server does not
		   know.
		TODO
		- "OUT_OF_MEMORY\r\n" The server cannot allocate enough memory for the job.
		   The client should try again later.

		 - "INTERNAL_ERROR\r\n" This indicates a bug in the server. It should never
		   happen. If it does happen, please report it at
		   http://groups.google.com/group/beanstalk-talk.


		 */
	switch command.Name {
	case USE:
		if command.Err == nil {
			return "USING " + command.Params["tube"]
		} else {
			return command.Err.Error()
		}
	case PUT:
		/**
		 	 - "INSERTED <id>\r\n" to indicate success.

			 - <id> is the integer id of the new job
		TODO
			 - "BURIED <id>\r\n" if the server ran out of memory trying to grow the
			   priority queue data structure.

			   - <id> is the integer id of the new job

			 - "EXPECTED_CRLF\r\n" The job body must be followed by a CR-LF pair, that is,
			   "\r\n". These two bytes are not counted in the job size given by the client
			   in the put command line.

			 - "JOB_TOO_BIG\r\n" The client has requested to put a job with a body larger
			   than max-job-size bytes.

			 - "DRAINING\r\n" This means that the server has been put into "drain mode" and
			   is no longer accepting new jobs. The client should try another server or
			   disconnect and try again later. To put the server in drain mode, send the
			   SIGUSR1 signal to the process.
		 */
		if command.Err == nil {
			return "INSERTED " + command.Job.Id()
		} else {
			return command.Err.Error()
		}
	case RESERVE:
	default:
		return command.Err.Error()

	}
	return ""
}
