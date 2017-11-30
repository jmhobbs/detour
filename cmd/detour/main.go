package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmhobbs/detour/pkg/hosts"
)

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "list":
		list()
	case "set":
		set()
	case "unset":
		unset()
	default:
		usage()
	}

}

func usage() {
	fmt.Println(`usage: detour <command> ...

Commands
  list                  - Show all active detours
  set <hostname> <ip>   - Set a detour
  unset <hostname>      - Remove a detour
`)
	os.Exit(1)
}

func list() {
	file, err := os.Open("/etc/hosts")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mapping, err := hosts.ExtractHostBlock(file)
	if err != nil {
		log.Fatal(err)
	}

	for host, ip := range mapping {
		fmt.Printf("%-15s    %s\n", ip, host)
	}
}

func set() {
	if len(os.Args) != 4 {
		usage()
	}
	// TODO: Test if args match IP and hostname formats

	file, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	mapping, err := hosts.ExtractHostBlock(file)
	if err != nil {
		log.Fatal(err)
	}

	mapping.Add(hosts.IPAddress(os.Args[3]), hosts.Hostname(os.Args[2]))

	file.Seek(0, 0)
	err = hosts.UpsertHostBock(mapping, file)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Detoured %s to %s\n", os.Args[2], os.Args[3])
}

func unset() {
	if len(os.Args) != 3 {
		usage()
	}

	file, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	mapping, err := hosts.ExtractHostBlock(file)
	if err != nil {
		log.Fatal(err)
	}

	if mapping[hosts.Hostname(os.Args[2])] == "" {
		log.Printf("No detour for %s\n", os.Args[2])
		return
	}

	mapping.Remove(hosts.Hostname(os.Args[2]))

	file.Seek(0, 0)
	err = hosts.UpsertHostBock(mapping, file)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Removed detour to %s\n", os.Args[2])
}
