package alerts

func GetInstr() []string {
    return []string{"EURRUB", "EURUSD"}
}

func hasInstrument(instr string) bool {
    for _, v := range GetInstr() {
        if v == instr {
            return true
        }
    }

    return false
}