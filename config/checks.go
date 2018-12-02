package config

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func GenerateCheckEntry(balance decimal.Decimal, notified bool) string {
	return fmt.Sprintf("%d:%s:%t", time.Now().UTC().Unix(), balance.StringFixed(2), notified)
}

func NeedsNotification(checkEntry string, balance decimal.Decimal, threshold decimal.Decimal) bool {
	thresholdMet := balance.LessThanOrEqual(threshold)

	split := strings.Split(checkEntry, ":")
	if len(split) != 3 {
		log.Printf("error splitting: len %d", len(split))
		return thresholdMet
	}

	tStr, balanceStr, notifiedStr := split[0], split[1], split[2]

	tInt, err := strconv.ParseInt(tStr, 10, 64)
	if err != nil {
		log.Printf("error parsing time: %s", err)
		return thresholdMet
	}
	t := time.Unix(tInt, 0)

	balance, err = decimal.NewFromString(balanceStr)
	if err != nil {
		log.Printf("error parsing decimal: %s", err)
		return thresholdMet
	}

	notified, err := strconv.ParseBool(notifiedStr)
	if err != nil {
		log.Printf("error parsing bool: %s", err)
		return thresholdMet
	}

	log.Printf("time: %s, value: %s, notified: %t", t.String(), balance.StringFixed(2), notified)

	// notify every n hours
	_ = time.Now().Sub(t).Hours() < 72

	return !notified && thresholdMet
}
