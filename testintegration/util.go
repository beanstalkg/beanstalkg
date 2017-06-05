package testintegration

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

// A Conn represents a connection to a beanstalkd server.
type Conn struct {
	c *textproto.Conn
}

var (
	space      = []byte{' '}
	crnl       = []byte{'\r', '\n'}
	yamlHead   = []byte{'-', '-', '-', '\n'}
	nl         = []byte{'\n'}
	colonSpace = []byte{':', ' '}
	minusSpace = []byte{'-', ' '}
)

// NewConn returns a new Conn using conn for I/O.
func NewConn(conn io.ReadWriteCloser) *Conn {
	c := new(Conn)
	c.c = textproto.NewConn(conn)
	return c
}

// Dial connects to the given address on the given network using net.Dial
// and then returns a new Conn for the connection.
func dial(network, addr string) (*Conn, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return NewConn(c), nil
}

// Close closes the underlying network connection.
func (c *Conn) Close() error {
	return c.c.Close()
}

func (c *Conn) cmd(body []byte, op string, args ...interface{}) (req, error) {
	r := req{c.c.Next(), op}
	c.c.StartRequest(r.id)
	if body != nil {
		args = append(args, len(body))
	}
	c.printLine(string(op), args...)
	if body != nil {
		c.c.W.Write(body)
		c.c.W.Write(crnl)
	}
	err := c.c.W.Flush()
	if err != nil {
		return req{}, ConnError{c, op, err}
	}
	c.c.EndRequest(r.id)
	return r, nil
}

// does not flush
func (c *Conn) printLine(cmd string, args ...interface{}) {
	io.WriteString(c.c.W, cmd)
	for _, a := range args {
		c.c.W.Write(space)
		fmt.Fprint(c.c.W, a)
	}
	c.c.W.Write(crnl)
}

func (c *Conn) readResp(r req, readBody bool, f string, a ...interface{}) (body []byte, err error) {
	c.c.StartResponse(r.id)
	defer c.c.EndResponse(r.id)
	line, err := c.c.ReadLine()
	for strings.HasPrefix(line, "WATCHING ") || strings.HasPrefix(line, "USING ") {
		line, err = c.c.ReadLine()
	}
	if err != nil {
		return nil, ConnError{c, r.op, err}
	}
	toScan := line
	if readBody {
		var size int
		toScan, size, err = parseSize(toScan)
		if err != nil {
			return nil, ConnError{c, r.op, err}
		}
		body = make([]byte, size+2) // include trailing CR NL
		_, err = io.ReadFull(c.c.R, body)
		if err != nil {
			return nil, ConnError{c, r.op, err}
		}
		body = body[:size] // exclude trailing CR NL
	}

	err = scan(toScan, f, a...)
	if err != nil {
		return nil, ConnError{c, r.op, err}
	}
	return body, nil
}

func scan(input, format string, a ...interface{}) error {
	_, err := fmt.Sscanf(input, format, a...)
	if err != nil {
		return findRespError(input)
	}
	return nil
}

type req struct {
	id uint
	op string
}

// ConnError records an error message from the server and the operation
// and connection that caused it.
type ConnError struct {
	Conn *Conn
	Op   string
	Err  error
}

func (e ConnError) Error() string {
	return e.Op + ": " + e.Err.Error()
}

// Error messages returned by the server.
var (
	ErrBadFormat  = errors.New("bad command format")
	ErrBuried     = errors.New("buried")
	ErrDeadline   = errors.New("deadline soon")
	ErrDraining   = errors.New("draining")
	ErrInternal   = errors.New("internal error")
	ErrJobTooBig  = errors.New("job too big")
	ErrNoCRLF     = errors.New("expected CR LF")
	ErrNotFound   = errors.New("not found")
	ErrNotIgnored = errors.New("not ignored")
	ErrOOM        = errors.New("server is out of memory")
	ErrTimeout    = errors.New("timeout")
	ErrUnknown    = errors.New("unknown command")
)

var respError = map[string]error{
	"BAD_FORMAT":      ErrBadFormat,
	"BURIED":          ErrBuried,
	"DEADLINE_SOON":   ErrDeadline,
	"DRAINING":        ErrDraining,
	"EXPECTED_CRLF":   ErrNoCRLF,
	"INTERNAL_ERROR":  ErrInternal,
	"JOB_TOO_BIG":     ErrJobTooBig,
	"NOT_FOUND":       ErrNotFound,
	"NOT_IGNORED":     ErrNotIgnored,
	"OUT_OF_MEMORY":   ErrOOM,
	"TIMED_OUT":       ErrTimeout,
	"UNKNOWN_COMMAND": ErrUnknown,
}

type unknownRespError string

func (e unknownRespError) Error() string {
	return "unknown response: " + string(e)
}

func findRespError(s string) error {
	if err := respError[s]; err != nil {
		return err
	}
	return unknownRespError(s)
}

func parseSize(s string) (string, int, error) {
	i := strings.LastIndex(s, " ")
	if i == -1 {
		return "", 0, findRespError(s)
	}
	n, err := strconv.Atoi(s[i+1:])
	if err != nil {
		return "", 0, err
	}
	return s[:i], n, nil
}

func parseDict(dat []byte) map[string]string {
	if dat == nil {
		return nil
	}
	d := make(map[string]string)
	if bytes.HasPrefix(dat, yamlHead) {
		dat = dat[4:]
	}
	for _, s := range bytes.Split(dat, nl) {
		kv := bytes.SplitN(s, colonSpace, 2)
		if len(kv) != 2 {
			continue
		}
		d[string(kv[0])] = string(kv[1])
	}
	return d
}

func parseList(dat []byte) []string {
	if dat == nil {
		return nil
	}
	l := []string{}
	if bytes.HasPrefix(dat, yamlHead) {
		dat = dat[4:]
	}
	for _, s := range bytes.Split(dat, nl) {
		if !bytes.HasPrefix(s, minusSpace) {
			continue
		}
		l = append(l, string(s[2:]))
	}
	return l
}

type dur time.Duration

func (d dur) String() string {
	return strconv.FormatInt(int64(time.Duration(d)/time.Second), 10)
}
