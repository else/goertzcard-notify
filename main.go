package main

import (
	"fmt"
	"log"

	"github.com/else/goertzcard-notify/client"
	"github.com/else/goertzcard-notify/config"
	"github.com/else/goertzcard-notify/notify"
	"github.com/shopspring/decimal"
)

const CREDIT_LOW_MSG = "Das Guthaben auf deiner Vorteilskarte beträgt aktuell %s€ und somit weniger als %s€. Du solltest es bald aufladen."

func main() {
	conf, err := config.Load("config.yaml")
	if err != nil {
		panic(err)
	}
	err = conf.Validate()
	if err != nil {
		panic(err)
	}
	for _, acc := range conf.Accounts {
		g := client.NewClient(acc.Credentials.Username, acc.Credentials.Password)
		err := g.Login()
		if err != nil {
			log.Fatal(err)
		}
		cards, err := g.GetCards()
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range cards {
			for _, card := range acc.Cards {
				if c.Ean != card.Ean {
					continue
				}
				log.Println(c)
				var notified bool

				d, err := decimal.NewFromString(card.MinimumAmount)
				if err != nil {
					panic(err)
				}

				if config.NeedsNotification(c.Description, c.Balance, d) {
					log.Println("needs notification")
					n := notify.PushoverNotifier{
						User:  acc.Notifier.Pushover.User,
						Token: acc.Notifier.Pushover.Token,
					}
					err := n.Notify("Goertz Vorteilskarte Guthaben", fmt.Sprintf(CREDIT_LOW_MSG, c.Balance, card.MinimumAmount))
					if err != nil {
						log.Printf("error while notifying: %s", err)
					}
					notified = true
				}

				c.Store(config.GenerateCheckEntry(c.Balance, notified))
			}
		}
	}
}
