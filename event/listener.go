package event

import (
    "time"
    "log"
)

type ChangedEvent struct {
    Src *Source
    Ts time.Time
    Symbol string
    New Quote
}

func NewListenChannel() chan ChangedEvent {
    return make(chan ChangedEvent)
}

func ListenSource(duration time.Duration, src *Source, ch chan<- ChangedEvent) {
    tickCh := time.Tick(duration)

    for ts := range tickCh {
        err := src.loadFromUrl()
        if err != nil {
            log.Println("error:", err)
        }

        log.Println("src loaded", len(src.quotes))

        for sym, q := range src.quotes {
            qPrev, ok := src.quotesPrev[sym]
            if !ok || qPrev.Bid != q.Bid || qPrev.Ask != q.Ask {
                ch <- ChangedEvent{src, ts, sym, q}
            }
        }
    }
}