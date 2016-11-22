package architecture

import (
	"errors"
	"strings"
	// "fmt"
	"github.com/satori/go.uuid"
	"strconv"
	"log"
)

type CommandName string

const (
	USE CommandName = "use"
	PUT             = "put"
	RESERVE         = "reserve"
)

type Command struct {
	Name           CommandName
	RawCommand     string
	Params         map[string]string
	WaitingForMore bool
	Err            error
	Job				*Job
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

// Parse keeps track of the state of the command and it will be called multiple times for commands such as  'put'
// where the length of command spans multiple lines
func (command *Command) Parse(rawCommand string) (bool, error) {
	// check if this command has already been waiting for a second round
	if !command.WaitingForMore {
		// first round
		parts := strings.Split(rawCommand, " ")
		switch strings.ToLower(parts[0]) {
		case "use":
			// fmt.Println(len(parts))
			if len(parts) > 2 {
				command.Err = errors.New("invalid format")
				return true, command.Err
			}
			command.Name = USE
			command.RawCommand = rawCommand
			command.Params = map[string]string{
				"tube": parts[1],
			}
			return true, nil
		case "put":
			// put <pri> <delay> <ttr> <bytes>\r\n <data>\r\n
			if len(parts) != 5 {
				command.Err = errors.New("invalid format")
				return true, command.Err
			}
			command.Name = PUT
			command.RawCommand = rawCommand
			command.Params = map[string]string{
				"pri":   parts[1],
				"delay": parts[2],
				"ttr":   parts[3],
				"bytes": parts[4],
			}
			command.WaitingForMore = true
			return false, nil
		}
	} else {
		// second round
		log.Println("GOT MORE", command)
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

func (command *Command) Reply() string {
	switch command.Name {
	case USE:
		if command.Err == nil {
			return "USING " + command.Params["tube"]
		}
		return ""
	case PUT:
		if command.Err == nil {
			return "INSERTED " + command.Job.Id()
		}
	}
	return ""
}
