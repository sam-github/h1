package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/uber-go/hackeroni/h1"
)

// when needed: https://github.com/spf13/cobra

func main() {
	file := flag.String("token", ".token", "Token file containing `user:apikey`")
	prog := flag.String("program", "nodejs", "Program name")
	flag.Parse()

	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	token, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("read token: %s", err)
	}
	parts := strings.SplitN(string(token), ":", 2)
	if len(parts) != 2 {
		log.Fatalf("failed to find `user:apikey` in %s", *file)
	}

	auth := h1.APIAuthTransport{
		APIIdentifier: strings.TrimSpace(parts[0]),
		APIToken:      strings.TrimSpace(parts[1]),
	}

	client := h1.NewClient(auth.Client())

	filter := h1.ReportListFilter{
		Program: []string{strings.TrimSpace(*prog)},
	}

	var listOpts h1.ListOptions

	fmt.Print("Listing all reports:\n")
	for {
		reports, resp, err := client.Report.List(filter, &listOpts)
		if err != nil {
			panic(err)
		}
		if resp.Links.Next == "" {
			break
		}
		listOpts.Page = resp.Links.NextPageNumber()
		for _, report := range reports {
			fmt.Printf("%s %v\n", *report.ID, *report.Title)
		}
	}
}
