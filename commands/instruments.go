package commands

import (
    "strings"
)

func instruments(known []string) string {
    return "Supported instruments:\n\n" + strings.Join(known, "\n")
}
