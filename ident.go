package ident

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const DEFAULT_PORT string = "113"

type Resp struct {
	OS      string
	Charset string
	ID      string
}

type Error struct {
	Err        string
	Unexpected string
}

func (e Error) Error() string {
	if e.Unexpected != "" {
		return fmt.Sprintf("Unexpected response from ident server: %s", e.Unexpected)
	}
	return fmt.Sprintf("Ident server error: %s", e.Err)
}

func Query(network, host, srcPort, dstPort string) (Resp, error) {
	conn, err := net.Dial(network, host)
	if err != nil {
		return Resp{}, err
	}

	_, err = fmt.Fprintf(conn, "%s, %s\n", srcPort, dstPort)
	if err != nil {
		return Resp{}, err
	}

	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return Resp{}, err
	}

	fields := strings.SplitN(strings.TrimSpace(resp), ":", 4)
	if len(fields) < 3 {
		return Resp{}, Error{Unexpected: resp}
	}

	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	switch fields[1] {
	case "USERID":
		if len(fields) != 4 {
			return Resp{}, Error{Unexpected: resp}
		}

		os := ""
		charset := "ASCII"

		srcinfo := strings.SplitN(fields[2], ",", 2)
		if len(srcinfo) == 2 {
			os = srcinfo[0]
			charset = srcinfo[1]
		} else {
			os = srcinfo[0]
		}

		return Resp{
			OS:      os,
			Charset: charset,
			ID:      fields[3],
		}, nil
	case "ERROR":
		if len(fields) != 3 {
			return Resp{}, Error{Unexpected: resp}
		}

		return Resp{}, Error{Err: fields[3]}
	}
	return Resp{}, Error{Unexpected: resp}
}
