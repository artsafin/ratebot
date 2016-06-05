package commands

func help(helpSection byte) string {
    head := `
Bot is able to notify you on currency rate changes.
All commands are available both with or without leading slash.

/help - Show this message
/symbols - Show supported instruments
/alerts - Show currently set alerts
/alert - Subscribe to a new alert
/forget - Unsubscribe from alert
/cancel - Cancel current operation

`
    newa := `<b>Alert subscription syntax</b>
The following commands are identical:
    <code>/alert SYMBOL OPERATION VALUE</code>
    <code>alert SYMBOL OPERATION VALUE</code>
    <code>SYMBOL OPERATION VALUE</code>

where:
<code>SYMBOL</code> is one of the supported instruments (see /symbols), case insensitive
<code>OPERATION</code> is one of:
    = (or eq, equals)
    &lt; (or lt, less than)
    &lt;= (or lte, less than or equals)
    &gt; (or gt, greater than)
    &gt;= (or gte, greater than or equals)
<code>VALUE</code> is a decimal number with dot as a decimal part separator

`
    if helpSection == 0 {
        return head + newa
    }

    var res string
    if helpSection & 1 == 1 {
        res += head
    }
    if helpSection & 2 == 2 {
        res += newa
    }
    return res
}