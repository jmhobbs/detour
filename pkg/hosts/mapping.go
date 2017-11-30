package hosts

import (
	"io"
	"strings"
)

type Hostname string
type IPAddress string
type HostMapping map[Hostname]IPAddress

func (m HostMapping) Add(ip IPAddress, host Hostname) {
	m[host] = ip
}

func (m HostMapping) AddMany(ip IPAddress, hosts string) {
	for _, host := range strings.Fields(hosts) {
		m[Hostname(host)] = ip
	}
}

func (m HostMapping) Remove(host Hostname) {
	delete(m, host)
}

func (m HostMapping) Write(dst io.Writer) error {
	var err error

	inverted := map[IPAddress][]string{}
	for host, ip := range m {
		hostnames, exists := inverted[ip]
		if !exists {
			inverted[ip] = []string{string(host)}
		} else {
			inverted[ip] = append(hostnames, string(host))
		}
	}

	for ip, hosts := range inverted {
		_, err = dst.Write([]byte(ip))
		if err != nil {
			return err
		}
		_, err = dst.Write([]byte("\t"))
		if err != nil {
			return err
		}

		_, err = dst.Write([]byte(strings.Join(hosts, "  ")))
		if err != nil {
			return err
		}

		_, err = dst.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}
	return nil
}
