package clistock

import (
    "container/list"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
    "time"
)

const (
    apiUrl         = "http://query.yahooapis.com/v1/public/yql?q="
	queryString    = "select * from yahoo.finance.historicaldata where symbol in (%s) and " +
                     "startDate = '%s' and endDate = '%s'"
	configSettings = "&format=json&env=http://datatables.org/alltables.env"
)

type Quote struct {
    Name          string `json:"Name"`
    Date          string `json:"Date"`
	Symbol        string `json:"Symbol"`
	Open          float64 `json:"Open,string"`
    High          float64 `json:"High,string"`
    Low           float64 `json:"Low,string"`
    Close         float64 `json:"Close,string"`
    Volume        float64 `json:"Volume,string"`
}

type Quotes []Quote

type YahooQuoteResponse struct {
	Query struct {
		Results struct {
			Quotes Quotes `json:"quote"`
		}
	}
}

func CreateSymbolQuery(symbols []string) string {
    for i, symbol := range symbols {
        symbols[i] = `"` + symbol + `"`
    }
    return strings.Join(symbols, ",")
}

func GetHistoricalData(symbols []string, start *time.Time, end *time.Time) (quotes map[string]*list.List, err error) {

    if len(symbols) == 0 {
		err = fmt.Errorf("Must have at least one symbol.")
		return
	}

    var symbolsString = CreateSymbolQuery(symbols)
    var urlString = apiUrl + url.QueryEscape(fmt.Sprintf(queryString, symbolsString,
                                             start.Format("2006-01-02"), end.Format("2006-01-02"))) + configSettings

    response, err := http.Get(urlString)
    if err != nil {
        return
    }

    defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return
    }

    var yqr YahooQuoteResponse
    err = json.Unmarshal(contents, &yqr)
    if err != nil {
    	return
    }

    var results = yqr.Query.Results.Quotes
    quotes = make(map[string]*list.List)
    for _, value := range results {
        var symbol = value.Symbol
        if _, ok := quotes[symbol]; ! ok {
            quotes[symbol] = list.New()
        }
        quotes[symbol].PushFront(value)
    }

    return
}
