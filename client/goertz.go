package client

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/shopspring/decimal"
)

func (g GoertzCard) String() string {
	return fmt.Sprintf("GoertzCard{ean=%s, balance=%s)", g.Ean, g.Balance.StringFixed(2))
}

var jar *cookiejar.Jar

func getCollector() (*colly.Collector, chan error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		// colly.Debugger(&debug.LogDebugger{}),
	)
	c.SetRequestTimeout(5 * time.Second)
	c.Async = true

	c.SetCookieJar(jar)

	c.OnRequest(func(r *colly.Request) {
		log.Printf("%4s %s", r.Method, r.URL)
	})

	errChan := make(chan error, 1)
	c.OnError(func(_ *colly.Response, err error) {
		errChan <- err
	})

	return c, errChan
}

func NewClient(username, password string) GoertzClient {
	jar, _ = cookiejar.New(nil)
	return GoertzClient{
		username: username,
		password: password,
	}
}

func (g GoertzClient) Login() error {
	c, errChan := getCollector()
	defer close(errChan)

	var response *colly.Response
	c.OnResponse(func(r *colly.Response) {
		response = r
	})

	// authenticate
	err := c.PostMultipart("https://vorteilskarte.baeckergoertz.de/", map[string][]byte{
		"login_logout":         []byte("login"),
		"loginform_abgesendet": []byte("ja"),
		"login[ean]":           []byte(g.username),
		"login[pin_jetzt]":     []byte(g.password),
	})

	if err != nil {
		return err
	}

	c.Wait()

	if response == nil {
		return fmt.Errorf("did not receive a response")
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}

func (g GoertzClient) GetCards() ([]GoertzCard, error) {
	c, errChan := getCollector()
	defer close(errChan)

	cards := make([]GoertzCard, 0)

	var response *colly.Response
	c.OnResponse(func(r *colly.Response) {
		response = r
	})

	c.OnHTML("table.kartenliste", func(e *colly.HTMLElement) {
		goquerySelection := e.DOM
		goquerySelection.Find("tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
			re := regexp.MustCompile("[^\\d]+")
			ean := re.ReplaceAllString(strings.TrimSpace(s.Find("td.ean").Text()), "")

			// balance
			strValue := strings.TrimSpace(s.Find("td.kontostand").Text())
			strValue = strings.Replace(strValue, " €", "", 1)
			strValue = strings.Replace(strValue, ",", ".", 1)
			value, err := decimal.NewFromString(strValue)
			if err != nil {
				g.errChan <- err
				return
			}

			strValue, _ = s.Find("td.bezeichnung input").Attr("value")
			label := strings.TrimSpace(strValue)
			log.Printf("label is %s", label)

			// internal card id
			re = regexp.MustCompile("karten_verwalten\\[(\\d+)\\]\\[kartenbezeichnung\\]")
			strValue, exists := s.Find("td.bezeichnung input").Attr("name")
			if !exists {
				g.errChan <- errors.New("could not find internal ID of card")
				return
			}
			sm := re.FindStringSubmatch(strValue)
			if len(sm) == 0 {
				g.errChan <- errors.New("could not find internal ID of card")
				return
			}

			card := GoertzCard{
				client:      &g,
				kn:          sm[1],
				Ean:         ean,
				Description: label,
				Balance:     value,
			}

			cards = append(cards, card)
		})
	})

	err := c.Visit("https://vorteilskarte.baeckergoertz.de/Karten_verwalten")
	if err != nil {
		return nil, err
	}

	c.Wait()

	// XXX: server always returns 200
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// check for errors during processing
	select {
	case err := <-g.errChan:
		return nil, err
	default:
	}

	if len(cards) == 0 {
		return nil, errors.New("invalid login or no cards found")
	}

	return cards, nil
}

func (g GoertzCard) Store(m string) error {
	c, errChan := getCollector()
	defer close(errChan)

	var response *colly.Response
	c.OnResponse(func(r *colly.Response) {
		response = r
	})

	requestData := map[string]string{
		"karten_verwalten_abgesendet": "ja",
		"karten_verwalten_send":       "speichern",
	}
	requestData[fmt.Sprintf("karten_verwalten[%s][kartenbezeichnung]", g.kn)] = m

	err := c.Post("https://vorteilskarte.baeckergoertz.de/Karten_verwalten", requestData)
	if err != nil {
		return err
	}

	c.Wait()

	// XXX: server always returns 200
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}
