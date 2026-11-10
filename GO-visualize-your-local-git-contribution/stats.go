package main

import (
	"sort"
	"time"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func stats(email string){
	commits := processRepositories(email)
	printCommitStats(commits)
}

func processRepositories(email string) map[int]int{
	dotFilePath := fileDotPath()
	repos := parseFilePathToSlice(dotFilePath)
	daysInMap := daysInLastSixMonths

	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i> 0; i--{
		commits[i] = 0
	}

	for _,path := range repos{
		commits = fillCommits(email, path, commits)
	}

	return commits
}

func fillCommits(email string, repoPath string, commits map[int]int) map[int]int{
	repo, err := git.PlainOpen(repoPath)
	if err != nil{
		fmt.Fprintf(os.Stderr, "skipping %s: %v\n", repoPath, err)
		return commits
	}

	//head returns a git referecne to where repo is pointing to
	ref, err := repo.Head()
	if err != nil{
		fmt.Fprintf(os.Stderr, "skipping %s: %v\n", repoPath, err)
		return commits
	}

	//a github log action being performed - from gets the commits
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil{
		fmt.Fprintf(os.Stderr, "skipping %s: %v\n", repoPath, err)
		return commits
	}

	offset := calcOffset()

	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When)

		if(c.Author.Email != email){
			return nil
		}

		if daysAgo != outOfRange{
			commits[daysAgo + offset]++
		}

		return nil
	})

	// shallow clones end with "object not found" when the walk reaches
	// the missing parents: keep the commits counted so far
	if err != nil{
		fmt.Fprintf(os.Stderr, "partial history for %s: %v\n", repoPath, err)
	}

	return commits
}

func getBegginningOfDay(t time.Time) time.Time{
	year, month, date := t.Date()
	startOfDay := time.Date(year, month, date, 0, 0, 0, 0, t.Location())
	return startOfDay
}

func countDaysSinceDate(date time.Time) int{
	days := 0
	now := getBegginningOfDay(time.Now())

	for date.Before(now){
		date = date.Add(time.Hour * 24)
		days ++
		if days > daysInLastSixMonths{
			return outOfRange
		}
	}
	return days
}

func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday{
	case time.Sunday:
		offset = 7
	case time.Saturday:
		offset = 6
	case time.Friday:
		offset = 5
	case time.Thursday:
		offset = 4
	case time.Wednesday:
		offset = 3
	case time.Tuesday:
		offset = 2
	case time.Monday:
		offset = 1
	}

	return offset
}

func printCommitStats(commits map[int]int){
	keys := sortMapIntoSlice(commits)
	cols := buildCols(keys, commits)
	printCells(cols)
}

func sortMapIntoSlice(commits map[int]int) []int{
	var keys[] int
	for k := range commits{
		keys = append(keys, k)
	}

	sort.Ints(keys)
	return keys
}

func buildCols(keys []int, commits map[int]int) map[int]column{
	cols := make(map[int]column)
	col := column{}

	for _,k := range keys{
		week := int(k/7)
		dayInWeek := k%7

		if dayInWeek == 0{
			col = column{}
		}

		col = append(col, commits[k])

		if dayInWeek == 6{
			cols[week] = col
		}
	}
	return cols
}

func printCells(cols map[int]column){
	printMonths()

	for j:=6; j>=0; j--{
		for i := weeksInLastSixMonths + 1; i>=0; i--{
			if i == weeksInLastSixMonths + 1{
				printDayCol(j)
			}

			if col, ok := cols[i]; ok{
				if i == 0 && j == calcOffset()-1{
					printCell(col[j], true)
					continue
				} else{
					if len(col) > j{
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

func printMonths(){
	week := getBegginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()

	fmt.Printf("		")
	for{
		if week.Month()!= month{
			fmt.Printf("%s ", week.Month().String()[:3])
			month =week.Month()
		} else{
			fmt.Printf("		")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

func printDayCol(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	// case 2:
	// 	out = " Tues "
	case 3:
		out = " Wed "
	// case 4:
	// 	out = " Thurs"
	case 5:
		out = " Fri "
	// case 6:
	// 	out = " Sat "
	}

	fmt.Print(out)
}

// printCell given a cell value prints it with a different format
// based on the value amount, and on the `today` flag.
func printCell(val int, today bool) {
    escape := "\033[0;37;30m"
    switch {
    case val > 0 && val < 5:
        escape = "\033[1;30;47m"
    case val >= 5 && val < 10:
        escape = "\033[1;30;43m"
    case val >= 10:
        escape = "\033[1;30;42m"
    }

    if today {
        escape = "\033[1;37;45m"
    }

    if val == 0 {
        fmt.Print(escape + "  - " + "\033[0m")
        return
    }

    str := "  %d "
    switch {
    case val >= 10:
        str = " %d "
    case val >= 100:
        str = "%d "
    }

    fmt.Printf(escape+str+"\033[0m", val)
}
