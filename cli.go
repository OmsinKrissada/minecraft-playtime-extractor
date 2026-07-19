package main

import (
	"fmt"
	"math"
)

func printInfo() {
	pt, err := getAllPlaytime()
	if err != nil {
		fmt.Printf("unable to read playtime data: %v\n", err)
		return
	}
	ls, err := getAllLastSeen()
	if err != nil {
		fmt.Printf("unable to read lastseen data: %v\n", err)
		return
	}
	transformed := transformResponse(pt, ls)
	worldRunTime := float64(getWorldRunTime()) / TICKS_IN_AN_HOUR
	worldRunTime = math.Round(worldRunTime*100) / 100

	maxLen := len("World run time:")
	for _, u := range transformed {
		if len(u.Name) > maxLen {
			maxLen = len(u.Name)
		}
	}

	fmt.Printf("%-*v %8.2f hrs\n", maxLen+2, "World run time:", worldRunTime)
	fmt.Println("Playtime:")
	for _, u := range transformed {
		name := u.Name
		if name == "" {
			name = "<unknown>"
		}
		fmt.Printf("  %-*v %8.2f hrs\n", maxLen, name, u.PlaytimeHr)
	}
}
