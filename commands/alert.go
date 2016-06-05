package commands

import (
    "strings"
    "strconv"
    "regexp"
    "fmt"
    "github.com/artsafin/ratebot/botservice"
    "github.com/artsafin/ratebot/alerts"
)

func getSymbolMatch(knownSymbols []string, symbol string) *string {
    for _, v := range knownSymbols {
        if strings.ToLower(v) == strings.ToLower(symbol) {
            return &v
        }
    }
    return nil
}

func parseAlertInput(text string, knownSymbols []string) (*[4]string, error) {
    opsExpr := strings.Join(alerts.GetSupportedOps(), "|") //"=|<=|>=|<|>"
    reStr := "^(?:/?alert)?\\s*([[:alpha:]]{3,6})\\s*(" + opsExpr + ")\\s*([0-9.]+)\\s*(?:\\s+for\\s+(\\d+))?\\s*$"

    re, err := regexp.Compile(reStr)
    if err != nil {
        return nil, err
    }

    matches := re.FindStringSubmatch(text)
    if len(matches) != 5 {
        return nil, fmt.Errorf("Invalid format; text=%s; re=%s", text, reStr)
    }

    symbolMatch := getSymbolMatch(knownSymbols, matches[1])
    if symbolMatch == nil {
        return nil, fmt.Errorf("Invalid instrument")
    }

    return &[4]string{*symbolMatch, matches[2], matches[3], matches[4]}, nil
}

func parseAlertInputOld(text string, knownSymbols []string) (*[4]string, error) {
    text = strings.TrimSpace(strings.Replace(strings.Replace(text, "/alert ", "", 1), "alert ", "", 1))

    result := new([4]string)

    for _, v := range knownSymbols {
        if len(text) >= len(v) && strings.ToLower(text[:len(v)]) == strings.ToLower(v) {
            result[0] = v
            text = strings.TrimSpace(text[len(v):])
            break
        }
    }

    if result[0] == "" {
        return result, fmt.Errorf("Instrument not found")
    }

    for _, v := range alerts.GetSupportedOps() {
        if len(text) >= len(v) && strings.ToLower(text[:len(v)]) == strings.ToLower(v) {
            result[1] = v
            text = strings.TrimSpace(text[len(v):])
            break
        }
    }

    if result[1] == "" {
        return result, fmt.Errorf("Operation not found")
    }

    if text == "" {
        return result, fmt.Errorf("Value not found")
    }

    result[2] = text

    // fmt.Println(result)

    return result, nil
}

func alert(chatId int64, userId int, fields *[4]string) (*botservice.HandlerReply, error) {
    opValue, err := strconv.ParseFloat(fields[2], 64)

    if err != nil {
        return nil, botservice.BotPublicErr("Unable to parse <code>VALUE</code>\n\n" + help(2))
    }

    ttl, err := strconv.Atoi(fields[3])
    if err != nil {
        ttl = 1
    }

    alert, err := alerts.NewAlert(chatId, userId, fields[0], fields[1], float64(opValue), ttl)
    if err != nil {
        return nil, botservice.BotPublicErr(err.Error())
    }

    err = alerts.Add(alert)
    if err != nil {
        return nil, botservice.BotPublicErr("Cannot subscribe to " + alert.StringShort() + ": " + err.Error())
    }

    return &botservice.HandlerReply{Text: "Subscribed to " + alert.String()}, nil
}
