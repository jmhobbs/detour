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
	mapping := HostMapping{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == START_DELIMITER {
			start_found = true
			continue
		}

		if start_found {
			if line == END_DELIMITER {
				return mapping, nil
			}

			ipv4_match := ipv4_matcher.FindStringSubmatch(line)
			if ipv4_match != nil {
				ip := ipv4_match[1] // TODO: Validation?
				hosts := ipv4_match[2]

				// Strip trailing comments
				if idx := strings.Index(hosts, "#"); idx != -1 {
					hosts = hosts[0:idx]
				}

				// No clobbering
				if val, exists := mapping[ip]; exists {
					mapping[ip] = append(val, strings.Fields(hosts)...)
				} else {
					mapping[ip] = strings.Fields(hosts)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return mapping, err
	}

	return mapping, nil
}

func UpsertHostBock(hosts map[string]string, sink io.Writer) error {
	// TODO
	return nil
}
