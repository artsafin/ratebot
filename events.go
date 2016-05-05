package main

import (
    "github.com/artsafin/ratebot/alerts"
    "time"
    "sync"
)

type Tick struct {
    ts time.Time
    value float32
}

type InstrEvent struct {
    instr string
    event string
}

type AlertEvent struct {
    alert *alerts.Alert
    tick *Tick
}

type SyncInstrColl struct {
    coll map[string]bool
    lock sync.Mutex
}

func (this *SyncInstrColl) remove(instr string) {
    this.lock.Lock()
    defer this.lock.Unlock()

    delete(this.coll, instr)
}
func (this *SyncInstrColl) add(instr string) {
    this.lock.Lock()
    defer this.lock.Unlock()

    this.coll[instr] = true
}

const HIST_LEN = 100
var instrs = SyncInstrColl{coll: make(map[string]bool)}
var history = make(map[string]bool)

func ListenTicks(addr string, instrChan <-chan InstrEvent, alertEventChan chan<- AlertEvent) {
    for ev := range instrChan {
        if ev.event == "remove" {
            instrs.remove(ev.instr)
        } else if ev.event == "add" {
            instrs.add(ev.instr)
        }
    }
}

func checkInstrCriteria(instr string, tick *Tick, alertEventChan chan<- AlertEvent) {
    for _, alert := range alerts.GetAlertsByInstr(instr) {
        if alert.Matches(tick.value) {
            ev := &AlertEvent{alert, tick}
            alertEventChan <- *ev
        }
    }
}
