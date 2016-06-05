package alerts

import (
    "sync"
    "fmt"
)

type Alert struct {
    ChatId int64
    userId int
    instr string
    op *Op
    value float64
    ttl int
}

func (a *Alert) isEqual(other *Alert) bool {
    return a.instr == other.instr && a.op == other.op && a.value == other.value
}

func (a *Alert) String() string {
    return fmt.Sprintf("%s %s %v (repeat %v)", a.instr, a.op.getShortName(), a.value, a.ttl)
}

func (a *Alert) StringShort() string {
    return fmt.Sprintf("%s %s %v", a.instr, a.op.getShortName(), a.value)
}

func (a *Alert) Test(value float64) bool {
    return a.op.test(value, a.value)
}
func (a *Alert) TryExpire() {
    a.ttl--
    if a.ttl <= 0 {
        RemoveAlert(a)
    }
}

var alertsIndex = make(map[string]map[int][]*Alert)

var lock sync.Mutex

func allocIndex(alert *Alert) {
    if _, ok := alertsIndex[alert.instr]; !ok {
        alertsIndex[alert.instr] = make(map[int][]*Alert)
    }
    if _, ok := alertsIndex[alert.instr][alert.userId]; !ok {
        alertsIndex[alert.instr][alert.userId] = make([]*Alert, 0, 2)
    }
}

func NewAlert(chatId int64, userId int, instr string, opStr string, value float64, ttl int) (*Alert, error) {
    if len(instr) < 3 {
        return nil, fmt.Errorf("Invalid instrument")
    }

    op, err := FindOpByString(opStr)
    if err != nil {
        return nil, err
    }

    return &Alert{chatId, userId, instr, op, value, ttl}, nil
}

func Add(alert *Alert) error {
    lock.Lock()
    defer lock.Unlock()

    allocIndex(alert)

    for _, v := range alertsIndex[alert.instr][alert.userId] {
        if v.isEqual(alert) {
            return fmt.Errorf("Already subscribed")
        }
    }

    alertsIndex[alert.instr][alert.userId] = append(alertsIndex[alert.instr][alert.userId], alert)

    return nil
}

func RemoveAlert(alert *Alert) {
    lock.Lock()
    defer lock.Unlock()

    alerts, ok := alertsIndex[alert.instr][alert.userId]
    if !ok {
        return
    }

    for i, v := range alerts {
        if v.isEqual(alert) {
            alertsIndex[alert.instr][alert.userId] = append(alertsIndex[alert.instr][alert.userId][:i], alertsIndex[alert.instr][alert.userId][i+1:]...)
            break
        }
    }
}

func GetAlertsBySymbol(instr string) []*Alert {
    lock.Lock()
    defer lock.Unlock()

    byUser, ok := alertsIndex[instr]
    if !ok {
        return []*Alert{}
    }

    res := make([]*Alert, 0)
    for _, userAlerts := range byUser {
        res = append(res, userAlerts...)
    }

    return res
}

func GetAlertsByUserId(userId int) []*Alert {
    lock.Lock()
    defer lock.Unlock()

    res := make([]*Alert, 0)

    for _, byUser := range alertsIndex {
        if alerts, ok := byUser[userId]; ok {
            res = append(res, alerts...)
        }
    }

    return res
}

func FilterUserAlerts(userId int, filterCb func(*Alert) bool) []*Alert {
    lock.Lock()
    defer lock.Unlock()

    res := make([]*Alert, 0)

    for _, byUser := range alertsIndex {
        if alerts, ok := byUser[userId]; ok {
            for _, alert := range alerts {
                if filterCb(alert) {
                    res = append(res, alert)
                }
            }
        }
    }

    return res
}

func FindUserAlert(userId int, filterCb func(*Alert) bool) *Alert {
    all := FilterUserAlerts(userId, filterCb)
    if len(all) > 0 {
        return all[0]
    }
    return nil
}

func AlertsArrayToStrings(alerts []*Alert) []string {
    reply := make([]string, len(alerts))
    for i, v := range alerts {
        reply[i] = v.String()
    }
    return reply
}
