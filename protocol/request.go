package protocol

import (
	"bytes"
	"errors"
	// "io"
	"log"
	"net"
	"regexp"
	// "strings"
	"strconv"
)

var (
	commandPattern    *regexp.Regexp
	setCommandPattern *regexp.Regexp
	getCommandPattern *regexp.Regexp
)

func init() {
	commandPattern = regexp.MustCompile(`\A\S+`)
	setCommandPattern = regexp.MustCompile(`\A(\S+) (\S+) (\S+) (\S+) (\S+)( (\S+))?\r\n(.+)\r\n`)
	getCommandPattern = regexp.MustCompile(`\A(\S+) (\S+)\r\n`)
}

type RequestReader struct {
	conn   net.Conn
	buf    []byte
	buffer bytes.Buffer
}

type Request interface {
}

type SetRequest struct {
	Command string
	Key     string
	Flags   int
	Exptime int
	Bytes   int
	Noreply int
	Data    string
}

type GetRequest struct {
	Command string
	Key     string
}

//
//  Create a request reader.
//
func NewRequestReader(conn net.Conn) *RequestReader {
	reader := &RequestReader{
		conn: conn,
		buf:  make([]byte, 1024),
	}
	return reader
}

//
//  Read next request and returns a new request
//
func (r *RequestReader) Read() (Request, error) {
	for i := 0; i < 3; i++ {
		n, err := r.conn.Read(r.buf)
		if err != nil {
			log.Fatal(err)
		}
		r.buffer.WriteString(string(r.buf))
		for i := 0; i < n; i++ {
			r.buf[i] = 0
		}
		if n == 2 {
			return r.MakeRequest()
		}
	}
	return nil, errors.New("Failed to parse input")
}

func (r *RequestReader) MakeRequest() (Request, error) {
	s := r.buffer.String()
	r.buffer.Reset()
	command := commandPattern.FindString(s)
	if command == "set" {
		return newSetRequest(s)
	} else if command == "get" {
		return newGetRequest(s)
	} else {

	}
	return nil, nil
}

func newSetRequest(s string) (Request, error) {
	elms := setCommandPattern.FindStringSubmatch(s)
	if elms == nil {
		return nil, errors.New("Bad format request")
	}
	flags, err := strconv.Atoi(elms[3])
	if err != nil {
		return nil, err
	}
	exptime, err := strconv.Atoi(elms[4])
	if err != nil {
		return nil, err
	}
	bytes, err := strconv.Atoi(elms[5])
	if err != nil {
		return nil, err
	}
	var noreply int
	if len(elms[7]) > 0 {
		noreply, err = strconv.Atoi(elms[7])
		if err != nil {
			return nil, err
		}
	}
	req := &SetRequest{
		Command: "set",
		Key:     elms[2],
		Flags:   flags,
		Exptime: exptime,
		Bytes:   bytes,
		Noreply: noreply,
		Data:    elms[8],
	}
	return req, nil
}

func newGetRequest(s string) (Request, error) {
	elms := getCommandPattern.FindStringSubmatch(s)
	if elms == nil {
		return nil, errors.New("Bad format request")
	}
	req := &GetRequest{
		Command: "get",
		Key:     elms[2],
	}
	return req, nil
}
