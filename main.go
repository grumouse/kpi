package main

import (
	"fmt"
	"time"

	"github.com/grumouse/kpi/kpi"
)

func strToTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)

	return t
}

func main() {
	fmt.Println("Start!")

	client := kpi.NewClient()

	req := kpi.KPIRequest{
		PeriodStart:       strToTime("2024-05-01"),
		PeriodEnd:         strToTime("2024-05-31"),
		PeriodKey:         "month",
		IndicatorToMoID:   227373,
		IndicatorToFactID: 0,
		Value:             1,
		FactTime:          strToTime("2024-05-31"),
		IsPlan:            0,
		AuthUserID:        40,
		Comment:           "grushin",
	}

	for i := 0; i < 10; i++ {
		client.Do(&req)
		fmt.Printf("%v sended...\n", i)
	}

	client.Wait()

	fmt.Println("End!")
}
