package hosts

import "testing"

func TestHostsMatcher(t *testing.T) {
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
		if !hosts_line_matcher.MatchString(str) {
			t.Errorf("Failed to match: '%s'\n", str)
		}
	}

	for _, str := range negative {
		if hosts_line_matcher.MatchString(str) {
			t.Errorf("Should not have matched: '%s'\n", str)
		}
	}
}
