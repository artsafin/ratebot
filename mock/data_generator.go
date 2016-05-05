package main

import (
    "math/rand"
    "strings"
)

func generateData(pairs []string) map[string]float32 {
    res := make(map[string]float32)

    for k, v := range pairs {
        res[strings.ToUpper(v)] = 10 * float32(k + 1)  + rand.Float32()
    }

    return res
}