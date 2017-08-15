package hosts

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type HostMapping map[string][]string

const START_DELIMITER = "### Detours Start ###"
const END_DELIMITER = "### Detours End ###"

var ipv4_matcher *regexp.Regexp

func init() {
	// Very loose IPv4 matcher
	ipv4_matcher = regexp.MustCompile("^([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})\\s+(.+)")
}

func ExtractHostBlock(src io.Reader) (HostMapping, error) {
	scanner := bufio.NewScanner(src)
	start_found := false
	blocks := HostMapping{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == START_DELIMITER {
			start_found = true
			continue
		}

		if start_found {
			if line == END_DELIMITER {
				return blocks, nil
			}

			ipv4_match := ipv4_matcher.FindStringSubmatch(line)
			if ipv4_match != nil {
				ip := ipv4_match[1] // TODO: Validation?
				hosts := strings.Fields(ipv4_match[2])
				// TODO: Tailing comments
				// TODO: Don't clobber
				blocks[ip] = hosts
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return blocks, err
	}

	return blocks, nil
}

func UpsertHostBock(hosts map[string]string, sink io.Writer) error {
	// TODO
	return nil
}
