package main

import (
	"fmt"
	"math"
)

func printInfo() {
	pt := getAllPlaytime()
	ls := getAllLastSeen()
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
		fmt.Printf("  %-*v %8.2f hrs\n", maxLen, u.Name, u.PlaytimeHr)
	}
}
