package main

import (
    "fmt"
    "os"
    "time"
    "github.com/gizak/termui"
    "github.com/paccorsi/clistock/yahoo"
)

func main() {

    if len(os.Args) < 2 || len(os.Args) > 3 {
        fmt.Println("Provide one or two valid symbol.")
        return
    }

    symbols := os.Args[1:]

    err := termui.Init()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer termui.Close()

    var end = time.Now()
    var start = end.AddDate(0, -6, 0)
    quotes, err := clistock.GetHistoricalData(symbols, &start, &end)
	if err != nil {
		fmt.Println(err)
		return
	}

    for symbol, quote := range quotes {
        var timeseries = make([]float64, 0, len(quotes))
        for value := quote.Front(); value != nil; value = value.Next() {
            timeseries = append(timeseries, value.Value.(clistock.Quote).Close)
        }

        lc := termui.NewLineChart()
    	lc.BorderLabel = symbol
    	lc.Data = timeseries
    	lc.AxesColor = termui.ColorWhite
    	lc.LineColor = termui.ColorGreen | termui.AttrBold
        lc.Mode = "dot"
        lc.Height = 15
        var row = termui.NewRow(termui.NewCol(12, 0, lc))
        termui.Body.AddRows(row)
    }

    termui.Body.Align()
    termui.Render(termui.Body)

    termui.Handle("/sys/kbd/q", func(termui.Event) {
        termui.StopLoop()
    })

    termui.Loop()
}
