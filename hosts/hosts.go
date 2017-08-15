package hosts

import (
	"bufio"
	"io"
	"regexp"
)

const START_DELIMITER = "### Detours Start ###"
const END_DELIMITER = "### Detours End ###"

var hosts_line_matcher *regexp.Regexp

func init() {
	// Very loose IPv4 matcher
	hosts_line_matcher = regexp.MustCompile("^([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})\\s+(.+)")
}

func ExtractHostBlock(src io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(src)
	start_found := false
	blocks := []string{}

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

			if hosts_line_matcher.MatchString(line) {
				blocks = append(blocks, line)
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
