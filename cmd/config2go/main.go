package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	debug     = flag.Bool("debug", false, "Set debug mode")
	omitempty = flag.Bool("omitempty", false, "Set omitempty mode")
	short     = flag.Bool("short", false, "Set short struct name mode")
	local     = flag.Bool("local", false, "Use local struct mode")
	example   = flag.Bool("example", false, "Use example tag mode")
	prefix    = flag.String("prefix", "", "Set struct name prefix")
	suffix    = flag.String("suffix", "", "Set struct name suffix")
	name      = flag.String("name", DefaultStructName, "Set struct name")
)

func main() {
	flag.Parse()
	SetDebug(*debug)

	opt := Options{
		UseOmitempty:   *omitempty,
		UseShortStruct: *short,
		UseLocal:       *local,
		UseExample:     *example,
		Prefix:         *prefix,
		Suffix:         *suffix,
		Name:           "config",
	}
	file, err := os.Open("config.json")
	if err != nil {
		os.Exit(0)
	}
	parsed, err := Parse(file, opt)
	if err != nil {
		panic(err)
	}
	title := "package main\n"

	fmt.Println(title + parsed)
}
