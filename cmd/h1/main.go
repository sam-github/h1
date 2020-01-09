package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

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

	var (
		waitInfo     []h1.Report
		waitTriage   []h1.Report
		waitFix      []h1.Report
		waitDisclose []h1.Report
	)

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
			var kind *[]h1.Report
			switch *report.State {
			case "needs-more-info":
				kind = &waitInfo
			case "new":
				kind = &waitTriage
			case "triaged":
				// if report.RawAssignee == nil {\
				//   kind = &waitChampion
				// }
				kind = &waitFix
			case "duplicate", "informative", "not-applicable", "resolved":
				// XXX Do we need a way of marking an issue as "do not disclose"?
				if report.DisclosedAt == nil {
					kind = &waitDisclose
				}
			case "spam":
				continue
			default:
				fmt.Printf("Unhandled: %v (%v)\n", *report.ID, *report.State)
				continue
			}
			if kind == nil {
				continue
			}
			*kind = append(*kind, report)
		}
	}
	fmt.Printf("# Open Reports\n")

	list("Waiting for more info", waitInfo)
	list("Waiting for triage", waitTriage)
	list("Waiting for fix", waitFix)
	list("Waiting for disclosure", waitDisclose)
}

func daysWaiting(report h1.Report) int {
	waitingSince := *report.CreatedAt
	if report.TriagedAt != nil {
		waitingSince = *report.TriagedAt
	}
	if report.ClosedAt != nil {
		waitingSince = *report.ClosedAt
	}
	waiting := time.Since(waitingSince.Time) / (24 * time.Hour)

	return int(waiting)
}

func sortByDaysWaiting(reports []h1.Report) {
	sort.Slice(reports, func(i, j int) bool {
		return daysWaiting(reports[i]) > daysWaiting(reports[j])
	})
}

func assignee(report h1.Report) string {
	if report.RawAssignee != nil {
		switch v := report.Assignee().(type) {
		case *h1.User:
			return *v.Username
		case *h1.Group:
			return *v.Name
		}
	}
	return ""
}

func list(h string, reports []h1.Report) {
	if len(reports) == 0 {
		return
	}

	sortByDaysWaiting(reports)

	fmt.Printf("## %s\n", h)
	for _, report := range reports {
		var champion string
		// Assignment to the Node.js Team means "no champion" :-(.
		if a := assignee(report); a != "" && a != "Node.js Team" {
			champion = fmt.Sprintf(" (%s)", a)
		}
		fmt.Printf("* %d days: https://hackerone.com/reports/%s%s\n",
			daysWaiting(report),
			*report.ID,
			champion,
		)
	}
}
