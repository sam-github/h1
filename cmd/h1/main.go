package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/uber-go/hackeroni/h1"
)

// when needed: https://github.com/spf13/cobra

var debug bool

func main() {
	file := flag.String("token", ".token", "`file` contains 'ID[@PROGRAM]:TOKEN' lines")
	prog := flag.String("program", "nodejs", "Program name")
	priv := flag.Bool("private", false, "Include private information")
	flag.BoolVar(&debug, "debug", false, "Include debug information")
	flag.Parse()

	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	var id string
	var token string
	config, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Failed to %s", err)
	}
	comment := regexp.MustCompile(`#.*$`)
	clear := string(comment.ReplaceAll(config, []byte("")))
	// rx: `ID [@ PROGRAM] : TOKEN`
	rx := `(?m)^\s*([^\s:@]+)\s*(?:@\s*([^\s:]*)\s*)?:\s*([^\s]+)\s*$`
	for _, ipt := range regexp.MustCompile(rx).FindAllStringSubmatch(clear, -1) {
		i := ipt[1]
		p := ipt[2]
		t := ipt[3]
		// If `@ PROGRAM` is included, it must match the program.
		if p != "" && p != *prog {
			continue
		}
		id = i
		token = t
		break
	}

	if id == "" || token == "" {
		log.Fatalf("Failed to find `ID[@PROGRAM]:TOKEN` in %s", *file)
	}

	auth := h1.APIAuthTransport{
		APIIdentifier: strings.TrimSpace(id),
		APIToken:      strings.TrimSpace(token),
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
			case "informative", "not-applicable", "resolved":
				if report.DisclosedAt == nil {
					// If a closed issue needs disclosure, assign it. Otherwise, treat
					// it as "nothing left to do".
					if report.RawAssignee != nil {
						kind = &waitDisclose
					}
				}
			case "duplicate":
				continue
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

	list(*priv, "Waiting for more info", waitInfo, false)
	list(*priv, "Waiting for triage", waitTriage, false)
	list(*priv, "Waiting for fix", waitFix, false)
	list(*priv, "Waiting for disclosure", waitDisclose, true)
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

func list(priv bool, h string, reports []h1.Report, withState bool) {
	if len(reports) == 0 {
		return
	}

	sortByDaysWaiting(reports)

	fmt.Printf("\n## %s\n", h)
	for _, report := range reports {
		var champion string
		// Assignment to the Node.js Team means "no champion" :-(.
		if a := assignee(report); a != "" && a != "Node.js Team" {
			champion = fmt.Sprintf(" (%s)", a)
		}

		var state string
		if withState {
			state = fmt.Sprintf(" (%s)", *report.State)
		}
		fmt.Printf("\n* %d days:%s https://hackerone.com/reports/%s%s\n",
			daysWaiting(report),
			state,
			*report.ID,
			champion,
		)
		if !priv {
			continue
		}
		fmt.Printf("> %s\n", *report.Title)

		if debug {
			fmt.Printf("report.BountyAwardedAt %+v\n", report.BountyAwardedAt)
			fmt.Printf("report.Bounties %+v\n", report.Bounties)
			fmt.Printf("report.Activities %+v\n", report.Activities)
		}
	}
}
