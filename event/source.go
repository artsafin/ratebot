package event

import (
    // "log"
    // "fmt"
    "net/http"
    // "io/ioutil"
    "encoding/json"
    "strconv"
)

func loadToVar(url string, v interface{}) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    dec := json.NewDecoder(resp.Body)

    err = dec.Decode(v)
    if err != nil {
        return err
    }

    return nil
}

type Quote struct {
    Bid float64
    Ask float64
}

type exactTypedVal struct {
    Symbol string
    Bid string
    Ask string
}

type Source struct {
    url string
    KnownSymbols []string
    quotesPrev map[string]Quote
    quotes map[string]Quote
}

func (me *Source) copyToHistory() {
    if len(me.quotes) > 0 && len(me.quotes) > len(me.quotesPrev) {
        me.quotesPrev = make(map[string]Quote, len(me.quotes))
    }

    for k, v := range me.quotes {
        me.quotesPrev[k] = v
    }
}

func (me *Source) loadFromUrl() error {
    var arr []exactTypedVal
    err := loadToVar(me.url, &arr)
    if err != nil {
        return err
    }

    me.copyToHistory()

    // Check if resize is needed
    if me.quotes == nil || len(arr) > 0 && len(me.quotes) != len(arr) {
        me.quotes = make(map[string]Quote, len(arr))
    }

    if len(arr) > 0 {
        me.KnownSymbols = make([]string, len(arr))
    }   

    for k, v := range arr {
        bid, _ := strconv.ParseFloat(v.Bid, 64)
        ask, _ := strconv.ParseFloat(v.Ask, 64)
        // fmt.Printf("symbol:%v bid=%v %v (%v)", v.Symbol, v.Bid, bid, errBid)

        me.quotes[v.Symbol] = Quote{Bid: bid, Ask: ask}
        me.KnownSymbols[k] = v.Symbol
    }

    return nil
}

func NewSource(url string) (*Source, error) {
    src := Source{url: url, KnownSymbols: make([]string, 0)}

    err := src.loadFromUrl()
    if err != nil {
        return nil, err
    }

    return &src, nil
}
