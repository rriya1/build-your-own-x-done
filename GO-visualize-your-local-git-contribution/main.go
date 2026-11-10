package main

import (
	"flag" //flag is a library to parse command line arguments
)

func stats(email string) {
	print("stats")
}

func main() {

	var folder string
	var email string

	flag.StringVar(&folder, "add", "", "Add a folder to add to the scan")
	flag.StringVar(&email, "email", "rriya1_be20@thapar.edu", "Enter the email to scan the stats for")
	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	stats(email)
}
