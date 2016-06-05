package alerts

import (
    "fmt"
)

type OpMatchFn func(float64, float64) bool

type Op struct {
    aliases []string
    test OpMatchFn
}

func (me *Op) getDisplayName() string {
    if len(me.aliases) > 0 {
        return me.aliases[0]
    } else {
        return "(unknown operation)"
    }
}

func (me *Op) getShortName() string {
    if len(me.aliases) >= 2 {
        return me.aliases[1]
    } else {
        return me.getDisplayName()
    }
}

var ops []*Op

var opsIndex map[string]*Op

var opsFlat []string

func init() {
    ops = []*Op{
        &Op{[]string{"greater than or equals", ">=", "gte"}, func(l float64, r float64) bool {return l >= r}},
        &Op{[]string{"less than or equals", "<=", "lte"}, func(l float64, r float64) bool {return l <= r}},
        &Op{[]string{"greater than", ">", "gt"}, func(l float64, r float64) bool {return l > r}},
        &Op{[]string{"less than", "<", "lt"}, func(l float64, r float64) bool {return l < r}},
        &Op{[]string{"equals", "=", "equal", "eq"}, func(l float64, r float64) bool {return l == r}},
    }

    var numAliases int
    for _, op := range ops {
        numAliases += len(op.aliases)
    }

    opsIndex = make(map[string]*Op, numAliases)
    opsFlat = make([]string, numAliases)

    flatNum := 0
    for _, op := range ops {
        for _, alias := range op.aliases {
            opsIndex[alias] = op
            opsFlat[flatNum] = alias
            flatNum++
        }
    }
}

func FindOpByString(op string) (*Op, error) {
    // fmt.Println(opsIndex)
    if instance, ok := opsIndex[op]; ok {
        return instance, nil
    }
    return nil, fmt.Errorf("Operation not found")
}

func GetSupportedOps() []string {
    return opsFlat
}
