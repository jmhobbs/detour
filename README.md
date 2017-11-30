[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/jmhobbs/detour/pkg/hosts)

# Detour

Detour is a CLI tool to manage your [hosts file](https://en.wikipedia.org/wiki/Hosts_(file)).

# Commands

The detour binary currently has three commands.

## set <hostname> <ip>

Creates or updates a detour for a domain.

    $ detour set example.com 127.0.0.1
    2017/11/29 22:24:19 Detoured example.com to 127.0.0.1

## list

Lists all the detours active in your hosts file.

    $ detour list
    127.0.0.1          example.com
    
## unset <hostname>

Remove a detour for a given hostname.

    $ detour unset example.com
    2017/11/29 22:24:29 Removed detour to example.com
    $ detour unset example.com
    2017/11/29 22:25:55 No detour for example.com