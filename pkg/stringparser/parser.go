package stringparser

import (
	"strings"
)

// ParseArgs mem-parsing input command string menjadi map flag ke nilai,
// khususnya mengambil seluruh sisa string setelah '=' untuk flag --flag=nilai tanpa memecahnya.
func ParseArgs(input string) map[string]string {
	flags := make(map[string]string)

	args := Tokenize(input)

	startIndex := 0
	if len(args) > 0 && strings.HasPrefix(args[0], "/") {
		startIndex = 1
	}

	i := startIndex
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			// Jika ada '='
			if strings.Contains(arg, "=") {
				splitIndex := strings.Index(arg, "=")
				key := arg[2:splitIndex]
				prefix := "--" + key + "="
				pos := strings.Index(input, prefix)
				if pos != -1 {
					// Mulai dari setelah '='
					valStart := pos + len(prefix)
					// Cari flag berikutnya dengan pola " --"
					nextFlagPos := strings.Index(input[valStart:], " --")
					if nextFlagPos != -1 {
						value := input[valStart : valStart+nextFlagPos]
						flags[key] = strings.TrimSpace(value)
					} else {
						// Ambil sampai akhir string
						value := input[valStart:]
						flags[key] = strings.TrimSpace(value)
					}
				} else {
					val := arg[splitIndex+1:]
					flags[key] = val
				}
				i++
			} else {
				key := arg[2:]
				valueTokens := []string{}
				i++
				for i < len(args) && !strings.HasPrefix(args[i], "--") {
					valueTokens = append(valueTokens, args[i])
					i++
				}
				flags[key] = strings.Join(valueTokens, " ")
			}
		} else {
			i++
		}
	}

	return flags
}

// Tokenize memecah string input dengan aturan nilai dalam tanda kutip dianggap satu token.
func Tokenize(input string) []string {
	var tokens []string
	inQuotes := false
	current := strings.Builder{}

	for i := 0; i < len(input); i++ {
		c := input[i]

		if c == '"' {
			inQuotes = !inQuotes
			continue
		}

		if c == ' ' && !inQuotes {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(c)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
