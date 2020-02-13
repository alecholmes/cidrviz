package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/big"
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

type cidrBoundary struct {
	ip    net.IP
	start bool
}

func main() {
	printGaps, namedSubnets := parseArgsOrExit()

	var names []string
	var boundaries []cidrBoundary
	for name, subnets := range namedSubnets {
		names = append(names, name)
		for _, subnet := range subnets {
			from, to := cidr.AddressRange(subnet)
			boundaries = append(boundaries, cidrBoundary{ip: from, start: true})
			boundaries = append(boundaries, cidrBoundary{ip: to, start: false})
		}
	}

	sort.Strings(names)
	sort.Slice(boundaries, func(i, j int) bool {
		if cmp := bytes.Compare(boundaries[i].ip, boundaries[j].ip); cmp != 0 {
			return cmp < 0
		}
		return boundaries[i].start
	})

	count := 0
	for i, boundary := range boundaries {
		if boundary.start {
			count++
		} else {
			count--
		}

		if i > 0 && bytes.Compare(boundaries[i-1].ip, boundary.ip) == 0 {
			continue
		}

		printIPRow := func(ip net.IP) {
			for _, name := range names {
				containsIP := false
				for _, subnet := range namedSubnets[name] {
					if subnet.Contains(ip) {
						containsIP = true
						break
					}
				}
				if containsIP {
					fmt.Printf("%s", name)
				} else {
					fmt.Printf("-")
				}
			}
		}

		if printGaps && i > 0 && boundary.start {
			lastBoundary := big.NewInt(0).SetBytes(boundaries[i-1].ip)
			currBoundary := big.NewInt(0).SetBytes(boundary.ip)
			distance := big.NewInt(0).Sub(currBoundary, lastBoundary).Int64()
			if distance > 1 {
				betweenIP := net.IP(big.NewInt(0).Add(lastBoundary, big.NewInt(1)).Bytes())
				fmt.Printf("                ")
				printIPRow(betweenIP)
				fmt.Printf(" (%d IPs)", distance-1)
				fmt.Println("")
			}
		}

		fmt.Printf("%-15s ", boundary.ip)
		printIPRow(boundary.ip)
		fmt.Println("")
	}
}

func parseArgsOrExit() (bool, map[string][]*net.IPNet) {
	noGaps := flag.Bool("no-gaps", false, "if true, skips printing gaps between non-adjacent subnets")

	flag.Parse()

	errAndExit := func(message string) {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n%s\n", message, usage)
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		errAndExit("No arguments given")
	}

	subnets := make(map[string][]*net.IPNet)
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

		subnets[name] = append(subnets[name], cidr)
	}

	return !*noGaps, subnets
}
