package server

import (
	"strings"
	"errors"
	// "fmt"
)

type Command struct {
	Name string
	RawCommand string
	Params map[string]string
	WaitingForMore bool
}

func (c *Command) Parse(rawCommand string) (bool, error) {
	// check if this command has already been waiting for a second round
	if c.WaitingForMore {
		switch c.Name {
		case "put":
			c.Params["data"] = rawCommand
			return true, nil
		}
	} else {
		parts := strings.Split(rawCommand, " ")
		switch strings.ToLower(parts[0]) {
		case "use":
			// fmt.Println(len(parts))
			if len(parts) > 2 {
				return true, errors.New("invalid format")
			}
			c.Name = "use"
			c.RawCommand = rawCommand
			c.Params = map[string]string{
				"tube": parts[1],
			}
			return true, nil
		case "put":
			// put <pri> <delay> <ttr> <bytes>\r\n <data>\r\n
			if len(parts) != 5 {
				return true, errors.New("invalid format")
			}
			c.Name = "put"
			c.RawCommand = rawCommand
			c.Params = map[string]string{
				"pri": parts[1],
				"delay": parts[2],
				"ttr": parts[3],
				"bytes": parts[4],
			}
			c.WaitingForMore = true
			return false, nil
		}
	}
	return true, nil
}
