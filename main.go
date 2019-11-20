package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"

	"github.com/apparentlymart/go-cidr/cidr"
)

const (
	usage = `Usage: cidrviz A=1.2.3.4/8 B=1.1.1.1/10 C=2.2.2.2/20`
)

var (
	argRegexp = regexp.MustCompile("^([a-zA-Z0-9]+)=(.+)$")
)

func main() {
	namedSubnets := parseArgsOrExit()

	ipMap := make(map[string]bool)
	var ips []net.IP
	putIP := func(ip net.IP) {
		if _, ok := ipMap[ip.String()]; !ok {
			ipMap[ip.String()] = true
			ips = append(ips, ip)
		}
	}

	var names []string
	for name, c := range namedSubnets {
		names = append(names, name)
		a, b := cidr.AddressRange(c)
		putIP(a)
		putIP(b)
	}

	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})

	sort.Strings(names)
	for _, ip := range ips {
		fmt.Printf("%-15s ", ip)
		for _, name := range names {
			c := namedSubnets[name]
			if c.Contains(ip) {
				fmt.Printf("%s", name)
			} else {
				fmt.Printf("-")
			}
		}
		fmt.Println("")
	}
}

func parseArgsOrExit() map[string]*net.IPNet {
	flag.Parse()

	errAndExit := func(message string) {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n%s\n", message, usage)
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		errAndExit("No arguments given")
	}

	subnets := make(map[string]*net.IPNet)
	for _, arg := range flag.Args() {
		matches := argRegexp.FindAllStringSubmatch(arg, -1)
		if len(matches) == 0 {
			errAndExit(fmt.Sprintf("Invalid argument: %s", arg))
		}

		name, rawCIDR := matches[0][1], matches[0][2]
		if len(name) > 1 {
			errAndExit(fmt.Sprintf("Name can only be a single alphanumeric character: %s", name))
		}
		_, cidr, err := net.ParseCIDR(rawCIDR)
		if err != nil {
			errAndExit(fmt.Sprintf("Invalid CIDR format: %s", rawCIDR))
		}

		if _, ok := subnets[name]; ok {
			errAndExit(fmt.Sprintf("Multiple arguments with same name: %s", name))
		}

		subnets[name] = cidr
	}

	return subnets
}
