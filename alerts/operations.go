package alerts

func GetOps() []string {
    return []string {
        ">=", "gte", "greater than or equals",
        "<=", "lte", "less than or equals",
        "<", "lt", "less than",
        "=", "equals", "equal", "eq",
        ">", "gt", "greater than",
    }
}

func opNormalize(op string) string {
    if op == "equals" || op == "eq" || op == "equal" {
        return "="
    }

    if op == "gt" || op == "greater than" {
        return ">"
    }

    if op == "gte" || op == "greater than or equals" {
        return ">="
    }

    if op == "lt" || op == "less than" {
        return "<"
    }

    if op == "lte" || op == "less than or equals" {
        return "<="
    }

    return op
}

func opIsNormal(op string) bool {
    return op == "=" || op == ">" || op == ">=" || op == "<" || op == "<=";
}

func opIsMatch(op string, left float32, right float32) bool {
    return op == "=" && left == right ||
        op == "<" && left < right ||
        op == "<=" && left <= right ||
        op == ">" && left > right ||
        op == ">=" && left >= right
}