package main

import (
	"flag"
	"log"

	radix "github.com/armon/go-radix"
	"github.com/blind-oracle/nginx-prometheus/geoip"
	syslog "gopkg.in/mcuadros/go-syslog.v2"
)

var (
	geoipDB *geoip.GeoIP
	uriTree *radix.Tree
	debug   bool
)

func receiveSyslog(ch syslog.LogPartsChannel) {
	var (
		l       *logEntry
		country string
		err     error
		ok      bool
	)

	for msg := range ch {
		if debug {
			log.Printf("%s: %s", msg["hostname"], msg["content"])
		}

		if l, err = parseSyslogMessage(msg); err != nil {
			log.Printf("Unable to parse message: %s", err)
			continue
		}

		// If we have a prefix list - try to match URI against it
		if uriTree != nil {
			_, _, ok = uriTree.LongestPrefix(l.uri)
		} else {
			ok = false
		}

		// If the prefix list is disabled or URI was not found - treat it as unknown.
		// Otherwise we can have potential for unbounded memory growth in Prometheus module.
		if !ok {
			l.uri = "Unknown"
		}

		if geoipDB != nil {
			if country, err = geoipDB.LookupCountry(l.clientIP); err != nil {
				log.Printf("Unable to lookup country for IP %s: %s", l.clientIP, err)
			} else if country != "" {
				l.clientCountry = country
			}
		}

		prometheusObserve(l)
	}
}

func main() {
	var (
		listenSyslog   string
		listenHTTP     string
		uriPrefixFile  string
		geoipCountryDB string
	)

	flag.StringVar(&listenSyslog, "listenSyslog", "0.0.0.0:1514", "ip:port to listen for syslog messages")
	flag.StringVar(&listenHTTP, "listenHTTP", "0.0.0.0:11080", "ip:port to listen for http requests")
	flag.StringVar(&uriPrefixFile, "uriPrefixFile", "", "file with allowed URI prefixes - one per line. If not specified - no per-URI stats gathered")
	flag.StringVar(&geoipCountryDB, "geoipCountryDB", "", "path to MaxMind GeoIP country DB. If not specified - no lookups performed")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
	flag.Parse()

	if uriPrefixFile != "" {
		uriTree, err := uriLoad(uriPrefixFile)
		if err != nil {
			log.Fatalf("Unable to load URIs: %s", err)
		}

		log.Printf("URI prefixes loaded: %d", uriTree.Len())
	}

	if geoipCountryDB != "" {
		geoipDB, err := geoip.New(geoipCountryDB)
		if err != nil {
			log.Fatalf("Unable to load GeoIP database: %s", err)
		}

		md := geoipDB.Metadata()
		log.Printf("GeoIP database loaded (%d nodes, %d build epoch)", md.NodeCount, md.BuildEpoch)
	}

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	log.Printf("Starting Syslog UDP listener on %s", listenSyslog)
	srv := syslog.NewServer()
	srv.SetFormat(syslog.RFC3164)
	srv.SetHandler(handler)
	srv.ListenUDP(listenSyslog)
	srv.Boot()
	go receiveSyslog(channel)

	log.Printf("Starting HTTP listener on %s", listenHTTP)
	go httpInit(listenHTTP)
	srv.Wait()
}
