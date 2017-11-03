package hosts

import (
	"io"
	"strings"
)

type HostMapping map[string][]string

func (m HostMapping) Add(ip, hosts string) {
	// No clobbering
	if val, exists := m[ip]; exists {
		m[ip] = append(val, strings.Fields(hosts)...)
	} else {
		m[ip] = strings.Fields(hosts)
	}
}

func (m HostMapping) Remove(host string) {
	for ip, hosts := range m {
		new_hosts := []string{}
		for _, h := range hosts {
			if host != h {
				new_hosts = append(new_hosts, h)
			}
		}
		m[ip] = new_hosts
	}
}

func (m HostMapping) Write(dst io.Writer) error {
	var err error
	for ip, hosts := range m {
		_, err = dst.Write([]byte(ip))
		if err != nil {
			return err
		}
		_, err = dst.Write([]byte("\t"))
		if err != nil {
			return err
		}
		_, err = dst.Write([]byte(strings.Join(hosts, " ")))
		if err != nil {
			return err
		}
		_, err = dst.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
