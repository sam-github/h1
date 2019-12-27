package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// when needed: https://github.com/spf13/cobra

func main() {
	log.SetOutput(os.Stderr)

	url := "https://api.hackerone.com/v1/reports?filter[program][]=nodejs"
	user := "at-sam-github"
	token, err := ioutil.ReadFile(".token")
	if err != nil {
		log.Fatalf("read file: %s", err)
	}
	pass := strings.TrimSpace(string(token)) // Strip trailing newline, etc.

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("new req: %s", err)
	}

	req.SetBasicAuth(user, pass)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("http.get: %s", err)
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatalf("read body: %s", err)
	}

	if rsp.StatusCode != http.StatusOK {
		log.Fatalf("%s %s", rsp.Status, body)
	}

	fmt.Printf("%s\n", body)
}
