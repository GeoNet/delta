package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ozym/delta/meta"
)

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var data string
	flag.StringVar(&data, "data", "../data", "base data directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "An example program to examine ...\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	// load network details into a map ...
	netmap := make(map[string]meta.Network)
	{
		var nets meta.Networks
		if err := meta.LoadList(filepath.Join(data, "networks.csv"), &nets); err != nil {
			panic(err)
		}

		for _, n := range nets {
			netmap[n.Code] = n
		}
	}
	//fmt.Println(netmap)

	// load station details
	markmap := make(map[string]meta.Mark)
	{
		var marks meta.Marks
		if err := meta.LoadList(filepath.Join(data, "marks.csv"), &marks); err != nil {
			panic(err)
		}

		for _, m := range marks {
			markmap[m.Code] = m
		}
	}
	//fmt.Println(netmap)

	// sort the keys on output
	var keys []string
	for k, _ := range markmap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// a simple loop and print
	for _, k := range keys {
		v, ok := markmap[k]
		if !ok {
			panic("invalid mark key: " + k)
		}
		{
			n, ok := netmap[v.Network]
			if !ok {
				panic("unable to find network: " + v.Network)
			}
			j, err := json.MarshalIndent(n, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))
		}
		j, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(j))
	}

}
