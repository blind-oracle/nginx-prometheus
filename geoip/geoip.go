package geoip

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type GeoIP struct {
	dbCountry *geoip2.Reader
}

// New returns an instance of GeoIP
func New(countryDB string) (g *GeoIP, err error) {
	g = &GeoIP{}
	if g.dbCountry, err = geoip2.Open(countryDB); err != nil {
		return
	}

	return
}

// LookupCountry looks up country in a geoip db
func (g *GeoIP) LookupCountry(ip net.IP) (name string, err error) {
	c, err := g.dbCountry.Country(ip)
	if err != nil {
		return
	}

	name = c.Country.IsoCode
	return
}

// Metadata returns maxmind metadata
func (g *GeoIP) Metadata() maxminddb.Metadata {
	return g.dbCountry.Metadata()
}
