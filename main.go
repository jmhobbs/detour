package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmhobbs/detour/hosts"
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
  list                  - Show all current detours
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

	for ip, hosts := range mapping {
		for _, host := range hosts {
			fmt.Printf("%-15s    %s\n", ip, host)
		}
	}
}

func set() {
}

func unset() {
}
