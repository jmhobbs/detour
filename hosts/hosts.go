package hosts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

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

				mapping.Add(ip, hosts)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return mapping, err
	}

	return mapping, nil
}

// TODO: Switch to a generic interface over os.File
// TODO: Mac vs Windows line endings
func UpsertHostBock(mapping HostMapping, sink *os.File) error {
	// This could be optimized by reading from sink, writing to a temp file, then doing a rename
	scanner := bufio.NewScanner(sink)
	var buffer bytes.Buffer
	start_found := false

	// Extract lines that aren't ours and put in a buffer
	for scanner.Scan() {
		line := scanner.Text()
		if line == START_DELIMITER {
			start_found = true
			continue
		}

		if start_found {
			if line == END_DELIMITER {
				start_found = false
			}
			continue
		}

		buffer.WriteString(line)
		buffer.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Now write our block
	buffer.WriteByte('\n')
	buffer.WriteString(START_DELIMITER)
	buffer.WriteByte('\n')
	for ip, hosts := range mapping {
		buffer.WriteString(fmt.Sprintf("%-15s\t%s\n", ip, strings.Join(hosts, " ")))
	}
	buffer.WriteString(END_DELIMITER)

	// This is our os.File specific call
	sink.Seek(0, 0)
	sink.Truncate(0)
	_, err := buffer.WriteTo(sink)

	return err
}
