package architecture

import (
	"errors"
	"strings"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
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
)

const (
	BAD_FORMAT      string = "BAD_FORMAT"
	UNKNOWN_COMMAND        = "UNKNOWN_COMMAND"
	NOT_IGNORED            = "NOT_IGNORED"
	NOT_FOUND              = "NOT_FOUND"
)

type Command struct {
	Name           CommandName
	RawCommand     string
	Params         map[string]string
	WaitingForMore bool
	MoreToSend     bool
	Err            error
	Job            Job
}

func NewCommand() Command {
	return Command{
		MoreToSend: false,
		Params:     map[string]string{},
	}
}

func NewDefaultCommand() Command {
	return Command{
		Name:       USE,
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

	command.Job = *NewJob(
		uuid.NewV1().String(),
		pri,
		delay,
		ttr,
		bytes,
		command.Params["data"],
	)
	log.Println("PROTOCOL new job: ", command.Job)
	return nil
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
				command.Err = errors.New(BAD_FORMAT)
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
				command.Err = errors.New(BAD_FORMAT)
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
			// watch <tube>\r\n
			command.Name = WATCH
			if len(parts) != 2 {
				command.Err = errors.New(BAD_FORMAT)
				return true, command.Err
			}
			command.RawCommand = rawCommand
			command.Params = map[string]string{
				"tube": parts[1],
			}
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		case "ignore":
			// ignore <tube>\r\n
			command.Name = IGNORE
			if len(parts) != 2 {
				command.Err = errors.New(BAD_FORMAT)
				return true, command.Err
			}
			command.RawCommand = rawCommand
			command.Params = map[string]string{
				"tube": parts[1],
			}
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		case "reserve":
			// reserve\r\n
			command.Name = RESERVE
			if len(parts) != 1 {
				command.Err = errors.New(BAD_FORMAT)
				return true, command.Err
			}
			command.RawCommand = rawCommand
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		case "reserve-with-timeout":
			// reserve-with-timeout <seconds>\r\n
			command.Name = RESERVE_WITH_TIMEOUT
			if len(parts) != 2 {
				command.Err = errors.New(BAD_FORMAT)
				return true, command.Err
			}
			command.Params = map[string]string{
				"timeout": parts[1],
			}
			command.RawCommand = rawCommand
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		case "delete":
			// delete <id>\r\n
			command.Name = DELETE
			if len(parts) != 2 {
				command.Err = errors.New(BAD_FORMAT)
				return true, command.Err
			}
			command.Params = map[string]string{
				"id": parts[1],
			}
			command.RawCommand = rawCommand
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		case "release":
			// release <id> <pri> <delay>\r\n
			command.Name = RELEASE
			if len(parts) != 4 {
				command.Err = errors.New(BAD_FORMAT)
				return true, command.Err
			}
			command.Params = map[string]string{
				"id":    parts[1],
				"pri":   parts[2],
				"delay": parts[3],
			}
			command.RawCommand = rawCommand
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		case "bury":
			// bury <id> <pri>\r\n
			command.Name = BURY
			if len(parts) != 3 {
				command.Err = errors.New(BAD_FORMAT)
				return true, command.Err
			}
			command.Params = map[string]string{
				"id":  parts[1],
				"pri": parts[2],
			}
			command.RawCommand = rawCommand
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		case "touch":
			// touch <id>\r\n
			command.Name = TOUCH
			if len(parts) != 2 {
				command.Err = errors.New(BAD_FORMAT)
				return true, command.Err
			}
			command.Params = map[string]string{
				"id": parts[1],
			}
			command.RawCommand = rawCommand
			command.WaitingForMore = false
			return !command.WaitingForMore, nil
		default:
			command.Err = errors.New(UNKNOWN_COMMAND)
			return true, command.Err
		}
	} else {
		// second round
		// log.Println("GOT MORE", command)
		switch command.Name {
		case PUT:
			command.Params["data"] = rawCommand
			err := command.createJobFromParams()
			// log.Println("GOT MORE PUT", c, err)
			command.Err = err
			return true, err
		}
	}
	return true, nil
}

func (command *Command) Reply() (bool, string) {

	//- "BAD_FORMAT\r\n" The client sent a command line that was not well-formed.
	//   This can happen if the line does not end with \r\n, if non-numeric
	//   characters occur where an integer is expected, if the wrong number of
	//   arguments are present, or if the command line is mal-formed in any other
	//   way.
	//
	// - "UNKNOWN_COMMAND\r\n" The client sent a command that the server does not
	//   know.
	// TODO
	// - "OUT_OF_MEMORY\r\n" The server cannot allocate enough memory for the job.
	//   The client should try again later.
	//
	// - "INTERNAL_ERROR\r\n" This indicates a bug in the server. It should never
	//   happen. If it does happen, please report it at
	//   http://groups.google.com/group/beanstalk-talk.
	switch command.Name {
	case USE:
		if command.Err == nil {
			return false, "USING " + command.Params["tube"]
		}
	case PUT:

		// 	 - "INSERTED <id>\r\n" to indicate success.
		//
		//	 - <id> is the integer id of the new job
		// TODO
		//	 - "BURIED <id>\r\n" if the server ran out of memory trying to grow the
		//	   priority queue data structure.
		//
		//	   - <id> is the integer id of the new job
		//
		//	 - "EXPECTED_CRLF\r\n" The job body must be followed by a CR-LF pair, that is,
		//	   "\r\n". These two bytes are not counted in the job size given by the client
		//	   in the put command line.
		//
		//	 - "JOB_TOO_BIG\r\n" The client has requested to put a job with a body larger
		//	   than max-job-size bytes.
		//
		//	 - "DRAINING\r\n" This means that the server has been put into "drain mode" and
		//	   is no longer accepting new jobs. The client should try another server or
		//	   disconnect and try again later. To put the server in drain mode, send the
		//	   SIGUSR1 signal to the process.

		if command.Err == nil {
			return false, "INSERTED " + command.Job.Id()
		}
	case WATCH:
		// WATCHING <count>\r\n
		if command.Err == nil {
			return false, "WATCHING " + command.Params["count"]
		}
	case IGNORE:
		// - "WATCHING <count>\r\n" to indicate success.
		//
		// - <count> is the integer number of tubes currently in the watch list.
		//
		// - "NOT_IGNORED\r\n" if the client attempts to ignore the only tube in its
		// watch list.
		if command.Err == nil {
			return false, "WATCHING " + command.Params["count"]
		}

	// RESERVE will return a newly-reserved job. If no job is available to be reserved,
	// beanstalkd will wait to send a response until one becomes available. Once a
	// job is reserved for the client, the client has limited time to run (TTR) the
	// job before the job times out. When the job times out, the server will put the
	// job back into the ready queue. Both the TTR and the actual time left can be
	// found in response to the stats-job command.
	//
	// If more than one job is ready, beanstalkd will choose the one with the
	// smallest priority value. Within each priority, it will choose the one that
	// was received first.
	//
	// A timeout value of 0 will cause the server to immediately return either a
	// response or TIMED_OUT.  A positive value of timeout will limit the amount of
	// time the client will block on the reserve request until a job becomes
	// available.
	//
	// During the TTR of a reserved job, the last second is kept by the server as a
	// safety margin, during which the client will not be made to wait for another
	// job. If the client issues a reserve command during the safety margin, or if
	// the safety margin arrives while the client is waiting on a reserve command,
	// the server will respond with:
	//
	//    DEADLINE_SOON\r\n
	//
	// This gives the client a chance to delete or release its reserved job before
	// the server automatically releases it.
	//
	//    TIMED_OUT\r\n
	//
	// If a non-negative timeout was specified and the timeout exceeded before a job
	// became available, or if the client's connection is half-closed, the server
	// will respond with TIMED_OUT.
	//
	// Otherwise, the only other response to this command is a successful reservation
	// in the form of a text line followed by the job body:
	//
	//    RESERVED <id> <bytes>\r\n
	//    <data>\r\n
	//
	// - <id> is the job id -- an integer unique to this job in this instance of
	//   beanstalkd.
	//
	// - <bytes> is an integer indicating the size of the job body, not including
	//   the trailing "\r\n".
	//
	// - <data> is the job body -- a sequence of bytes of length <bytes> from the
	//   previous line. This is a verbatim copy of the bytes that were originally
	//   sent to the server in the put command for this job.
	case RESERVE:
		if command.Err == nil {
			if !command.MoreToSend {
				command.MoreToSend = true
				return true, fmt.Sprintf("RESERVED %s %d", command.Job.Id(), command.Job.Bytes)
			} else {
				return false, command.Job.Data
			}
		}
	case RESERVE_WITH_TIMEOUT:
		if command.Err == nil {
			return false, "INSERTED " + command.Job.Id()
		}
	case DELETE:
		// The client then waits for one line of response, which may be:
		//
		// - "DELETED\r\n" to indicate success.
		//
		// - "NOT_FOUND\r\n" if the job does not exist or is not either reserved by the
		// client, ready, or buried. This could happen if the job timed out before the
		// client sent the delete command.
		if command.Err == nil {
			return false, "DELETED"
		}
	case RELEASE:
		// The client expects one line of response, which may be:
		//
		// - "RELEASED\r\n" to indicate success.
		//
		// - "BURIED\r\n" if the server ran out of memory trying to grow the priority
		// queue data structure.
		//
		// - "NOT_FOUND\r\n" if the job does not exist or is not reserved by the client.
		if command.Err == nil {
			return false, "RELEASED"
		}
	case BURY:
		// There are two possible responses:
		//
		// - "BURIED\r\n" to indicate success.
		//
		// - "NOT_FOUND\r\n" if the job does not exist or is not reserved by the client.
		if command.Err == nil {
			return false, "BURIED"
		}
	case TOUCH:
		// There are two possible responses:
		//
		// - "TOUCHED\r\n" to indicate success.
		//
		// - "NOT_FOUND\r\n" if the job does not exist or is not reserved by the client.
		if command.Err == nil {
			return false, "INSERTED " + command.Job.Id()
		}
	}
	return false, command.Err.Error()
}
