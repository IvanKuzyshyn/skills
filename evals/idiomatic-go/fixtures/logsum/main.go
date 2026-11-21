package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"example.com/logsum/parser"
)

func main() {
	path := flag.String("f", "", "log file")
	skip := flag.Bool("skip-bad", false, "skip malformed lines")
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	entries, err := parser.ParseAll(parser.ParserConfig{TimeLayout: "2006-01-02T15:04:05Z07:00"}, f, *skip, false)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	counts := map[string]int{}
	perRequest := map[string]int{}
	for _, e := range entries {
		counts[e.Level]++
		if e.Level == "ERROR" {
			perRequest[e.RequestId]++
		}
	}
	var levels []string
	for l := range counts { levels = append(levels, l) }
	sort.Strings(levels)
	for _, l := range levels {
		fmt.Printf("%-6s %d\n", l, counts[l])
	}
	type kv struct {
		k string
		v int
	}
	var top []kv
	for k, v := range perRequest {
		top = append(top, kv{k, v})
	}
	sort.Slice(top, func(i, j int) bool { return top[i].v > top[j].v })
	for i, x := range top {
		if i == 5 {
			break
		}
		fmt.Printf("%s %d\n", x.k, x.v)
	}
}
