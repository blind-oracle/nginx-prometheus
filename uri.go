//go:generate sh ./generate.sh

package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/armon/go-radix"
)

func uriLoad(filename string) (uris *radix.Tree, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	uris = radix.New()

	lines := strings.Split(string(data), "\n")
	// Sort prefixes ascending
	sort.Strings(lines)

	for _, l := range lines {
		if l = strings.TrimSpace(l); len(l) == 0 {
			continue
		}

		if pfx, _, ok := uris.LongestPrefix(l); ok {
			return nil, fmt.Errorf("'%s' is a subprefix of '%s'", pfx, l)
		}

		uris.Insert(l, 1)
	}

	return
}
