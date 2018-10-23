package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_URI(t *testing.T) {
	uri, err := uriLoad("uriPrefixes.txt")
	assert.Nil(t, err)

	_, _, ok := uri.LongestPrefix("/v10/account/myitems_singletab.json?a=b")
	assert.True(t, ok)

	_, _, ok = uri.LongestPrefix("/v10/account/myitems_singletab.jso?a=b")
	assert.False(t, ok)
}

func Benchmark_URI(b *testing.B) {
	uri, _ := uriLoad("uriPrefixes.txt")

	for i := 0; i < b.N; i++ {
		uri.LongestPrefix("/v10/account/myitems_singletab.json?a=b")
	}
}
