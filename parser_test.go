package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var (
	msgsGood = []string{
		`10.84.2.10|http|static.site.com|GET|HTTP/1.1|/banners/etc/sirtaki/july_2018/1/fonts/Candara-Bold.ttf|200|0.142|123456|229673`,
		`10.84.2.10|https|static.site.com|GET|HTTP/1.1|/api/foo?a=b|200|0.142|123456|229673`,
		`10.84.2.10|https|static.site.com|GET|HTTP/1.1|/api/bar?a=b|200|0.142|123456|229673`,
		`10.84.2.10|https|static.site.com|GET|HTTP/1.1|/api/bar?a=b|200|0.142|123456|229673`,
	}

	msgsGoodParsed = []*logEntry{
		{
			server:    "1.1.1.1",
			scheme:    "http",
			method:    "GET",
			hostname:  "static.site.com",
			status:    "200",
			protocol:  "HTTP/1.1",
			clientIP:  net.ParseIP("10.84.2.10"),
			uri:       `/banners/etc/sirtaki/july_2018/1/fonts/Candara-Bold.ttf`,
			duration:  0.142,
			bytesSent: 229673,
			bytesRcvd: 123456,
		}, {
			server:    "1.1.1.1",
			scheme:    "https",
			method:    "GET",
			hostname:  "static.site.com",
			status:    "200",
			protocol:  "HTTP/1.1",
			clientIP:  net.ParseIP("10.84.2.10"),
			uri:       `/api/foo`,
			duration:  0.142,
			bytesSent: 229673,
			bytesRcvd: 123456,
		}, {
			server:    "1.1.1.1",
			scheme:    "https",
			method:    "GET",
			hostname:  "static.site.com",
			status:    "200",
			protocol:  "HTTP/1.1",
			clientIP:  net.ParseIP("10.84.2.10"),
			uri:       `/api/bar`,
			duration:  0.142,
			bytesSent: 229673,
			bytesRcvd: 123456,
		}, {
			server:    "1.1.1.1",
			scheme:    "https",
			method:    "GET",
			hostname:  "static.site.com",
			status:    "200",
			protocol:  "HTTP/1.1",
			clientIP:  net.ParseIP("10.84.2.10"),
			uri:       `/api/bar`,
			duration:  0.142,
			bytesSent: 229673,
			bytesRcvd: 123456,
		},
	}

	msgsBad = []string{
		`10.84.2.10|https|static.site.com|GET|HTTP/1.1|/api/bar?a=b|200|0.142a|123456|229673|someshit`,
		`10.84.2.10|https|static.site.com|GET|HTTP/1.1|/api/bar?a=b|200|0.142a|123456|229673`,
		`10.84.2.10|https|static.site.com|GET|HTTP/1.1|/api/bar?a=b|200|0.142|123456a|229673`,
		`10.84.2.10|https|static.site.com|GET|HTTP/1.1|/api/bar?a=b|200|0.142|123456|229673a`,
		`10.84.2.10a|https|static.site.com|GET|HTTP/1.1|/api/bar?a=b|200|0.142|123456|229673`,
	}
)

func Test_parseSyslogMessage(t *testing.T) {
	for k, v := range msgsGood {
		m := format.LogParts{
			"hostname": "1.1.1.1",
			"content":  v,
		}

		l, err := parseSyslogMessage(m)
		assert.Nil(t, err)
		assert.Equal(t, msgsGoodParsed[k], l)
	}

	for _, v := range msgsBad {
		m := format.LogParts{
			"hostname": "1.1.1.1",
			"content":  v,
		}

		_, err := parseSyslogMessage(m)
		assert.NotNil(t, err)
	}
}

func Benchmark_parseSyslogMessage(b *testing.B) {
	msg1 := format.LogParts{
		"hostname": "1.1.1.1",
		"content":  msgsGood[0],
	}

	for i := 0; i < b.N; i++ {
		parseSyslogMessage(msg1)
	}
}
