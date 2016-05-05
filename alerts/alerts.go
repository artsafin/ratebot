package alerts

import (
    "sync"
    "fmt"
)

type Alert struct {
    chatId int
    instr string
    op string
    value float32
}

func (a *Alert) isEqual(other *Alert) bool {
    return a.instr == other.instr && a.op == other.op && a.value == other.value
}

func (a *Alert) String() string {
    return fmt.Sprintf("%s %s %v", a.instr, a.op, a.value)
}

func (a *Alert) Matches(value float32) bool {
    return opIsMatch(a.op, value, a.value)
}

var alertsIndex = make(map[string]map[int][]*Alert)

var lock sync.Mutex

func allocIndex(alert *Alert) {
    if _, ok := alertsIndex[alert.instr]; !ok {
        alertsIndex[alert.instr] = make(map[int][]*Alert)
    }
    if _, ok := alertsIndex[alert.instr][alert.chatId]; !ok {
        alertsIndex[alert.instr][alert.chatId] = make([]*Alert, 0, 2)
    }
}

func NewAlert(chatId int, instr string, op string, value float32) (*Alert, error) {
    if !hasInstrument(instr) {
        return nil, fmt.Errorf("Instrument not found")
    }
    op = opNormalize(op)
    if !opIsNormal(op) {
        return nil, fmt.Errorf("Invalid OPERATION")
    }

    return &Alert{chatId, instr, op, value}, nil
}

func AddAlert(alert *Alert) {
    lock.Lock()
    defer lock.Unlock()

    allocIndex(alert)

    for _, v := range alertsIndex[alert.instr][alert.chatId] {
        if v.isEqual(alert) {
            return
        }
    }

    alertsIndex[alert.instr][alert.chatId] = append(alertsIndex[alert.instr][alert.chatId], alert)
}

func RemoveAlert(alert Alert) {
    lock.Lock()
    defer lock.Unlock()

    alerts, ok := alertsIndex[alert.instr][alert.chatId]
    if !ok {
        return
    }

    for i, v := range alerts {
        if v.isEqual(&alert) {
            alertsIndex[alert.instr][alert.chatId] = append(alertsIndex[alert.instr][alert.chatId][:i], alertsIndex[alert.instr][alert.chatId][i+1:]...)
            break
        }
    }
}

func GetRegisteredInstr() []string {
    lock.Lock()
    defer lock.Unlock()

    keys := make([]string, 0, len(alertsIndex))

    for k := range alertsIndex {
        keys = append(keys, k)
    }

    return keys
}

func GetAlertsByInstr(instr string) []*Alert {
    lock.Lock()
    defer lock.Unlock()

    byChat, ok := alertsIndex[instr]
    if !ok {
        return []*Alert{}
    }

    res := make([]*Alert, len(byChat))
    for _, chatAlerts := range byChat {
        res = append(res, chatAlerts...)
    }

    return res
}

func GetAlertsByChatId(chatId int) ([]*Alert, error) {
    lock.Lock()
    defer lock.Unlock()

    res := make([]*Alert, 0)

    for _, byChat := range alertsIndex {
        if alerts, ok := byChat[chatId]; ok {
            res = append(res, alerts...)
        }
    }

    if len(res) > 0 {
        return res, nil
    }
    return res, fmt.Errorf("chat not found")
}

func AlertsArrayToStrings(alerts []*Alert) []string {
    reply := make([]string, len(alerts))
    for i, v := range alerts {
        reply[i] = v.String()
    }
    return reply
}
