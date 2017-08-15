package hosts

import (
	"reflect"
	"strings"
	"testing"
)

func TestIPv4Matcher(t *testing.T) {
	positive := []string{
		"196.128.1.1  some.domain.com",
		"127.0.0.1\ta.com b.com",
		"8.8.8.8    dns.thegoogle.com",
	}

	negative := []string{
		"no.com ok.com",
		"# this is a comment",
		"   # so is this",
		"127.0.0.0.1  too.many.dots.com",
	}

	for _, str := range positive {
		if !ipv4_matcher.MatchString(str) {
			t.Errorf("Failed to match: '%s'\n", str)
		}
	}

	for _, str := range negative {
		if ipv4_matcher.MatchString(str) {
			t.Errorf("Should not have matched: '%s'\n", str)
		}
	}
}

func TestExtractHostBlock(t *testing.T) {
	src := strings.NewReader(`
127.0.0.1  this.is.a.domain
### Detours Start ###
127.0.0.1  capture.me
0.0.0.0			and.capture.me but.also.me
### Detours End ###
`)

	expected := map[string][]string{
		"127.0.0.1": []string{"capture.me"},
		"0.0.0.0":   []string{"and.capture.me", "but.also.me"},
	}

	mapping, err := ExtractHostBlock(src)
	if err != nil {
		t.Fatal(err)
	}

	for ip, hosts := range expected {
		if !reflect.DeepEqual(mapping[ip], hosts) {
			t.Errorf("Error mapping '%s', expected '%v' got '%v'", ip, hosts, mapping[ip])
		}
	}
}
