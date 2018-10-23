package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type logEntry struct {
	// Tags
	server        string
	scheme        string
	method        string
	hostname      string
	status        string
	protocol      string
	uri           string
	clientIP      net.IP
	clientCountry string

	// Fields
	duration  float64
	bytesSent uint64
	bytesRcvd uint64
}

func parseSyslogMessage(msg format.LogParts) (l *logEntry, err error) {
	content := msg["content"].(string)

	a := strings.Split(content, "|")
	if len(a) != 10 {
		return nil, fmt.Errorf("Wrong number of fields in message: %s", content)
	}

	l = &logEntry{
		server:   msg["hostname"].(string),
		scheme:   a[1],
		hostname: a[2],
		method:   a[3],
		protocol: a[4],
		uri:      strings.Split(a[5], "?")[0],
		status:   a[6],
	}

	if l.clientIP = net.ParseIP(a[0]); l.clientIP == nil {
		return nil, fmt.Errorf("Unable to parse clientIP")
	}

	if l.duration, err = strconv.ParseFloat(a[7], 64); err != nil {
		return nil, fmt.Errorf("Unable to parse duration as float: %s", err)
	}

	if l.bytesRcvd, err = strconv.ParseUint(a[8], 10, 64); err != nil {
		return nil, fmt.Errorf("Unable to parse bytesRcvd as uint: %s", err)
	}

	if l.bytesSent, err = strconv.ParseUint(a[9], 10, 64); err != nil {
		return nil, fmt.Errorf("Unable to parse bytesSent as uint: %s", err)
	}

	return
}
