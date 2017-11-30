package hosts

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/aryann/difflib"
)

func TestIPv4Matcher(t *testing.T) {
	positive := []string{
		"196.128.1.1  some.domain.com",
		"127.0.0.1\ta.com  b.com",
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
127.0.0.1  no.clobber.please.com
192.168.1.1  this.com # but not this
### Detours End ###
`)

	expected := HostMapping{
		"capture.me":            "127.0.0.1",
		"no.clobber.please.com": "127.0.0.1",
		"and.capture.me":        "0.0.0.0",
		"but.also.me":           "0.0.0.0",
		"this.com":              "192.168.1.1",
	}

	mapping, err := ExtractHostBlock(src)
	if err != nil {
		t.Fatal(err)
	}

	for host, ip := range expected {
		if mapping[host] != ip {
			t.Errorf("Error mapping '%s', expected '%v' got '%v'", host, ip, mapping[host])
		}
	}
}

func TestUpsertHostBlockInsert(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "test-detours")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	hosts_content := "127.0.0.1  ignore.this.com"
	mapping := HostMapping{"google.com": "192.168.1.1", "velvetcache.org": "192.168.1.1"}
	expected := []byte(`127.0.0.1  ignore.this.com

### Detours Start ###
192.168.1.1	google.com  velvetcache.org
### Detours End ###`)

	if _, err := tmpfile.WriteString(hosts_content); err != nil {
		t.Fatal(err)
	}

	tmpfile.Seek(0, 0)
	if err = UpsertHostBock(mapping, tmpfile); err != nil {
		t.Error(err)
	}

	tmpfile.Seek(0, 0)
	b, err := ioutil.ReadAll(tmpfile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, b) {
		diff := difflib.Diff(strings.Split(string(expected), "\n"), strings.Split(string(b), "\n"))
		t.Error("Block doesn't match expectations")
		for _, d := range diff {
			t.Log(d)
		}
	}
}

func TestUpsertHostBlockUpdate(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "test-detours")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	hosts_content := `
127.0.0.1  ignore.this.com

192.168.1.1 and.this.com  too.com
### Detours Start ###
127.0.0.1  spacing.doesnt.matter.com cool.com
### Detours End ###
127.0.0.1  this-came-after.org`

	mapping := HostMapping{"google.com": "192.168.1.1", "velvetcache.org": "192.168.1.1", "example.net": "127.0.0.1"}
	expected := []byte(`
127.0.0.1  ignore.this.com

192.168.1.1 and.this.com  too.com
127.0.0.1  this-came-after.org

### Detours Start ###
192.168.1.1	google.com  velvetcache.org
127.0.0.1	example.net
### Detours End ###`)

	if _, err := tmpfile.WriteString(hosts_content); err != nil {
		t.Fatal(err)
	}

	tmpfile.Seek(0, 0)
	if err = UpsertHostBock(mapping, tmpfile); err != nil {
		t.Error(err)
	}

	tmpfile.Seek(0, 0)
	b, err := ioutil.ReadAll(tmpfile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, b) {
		diff := difflib.Diff(strings.Split(string(expected), "\n"), strings.Split(string(b), "\n"))
		t.Error("Block doesn't match expectations")
		for _, d := range diff {
			t.Log(d)
		}
	}
}
