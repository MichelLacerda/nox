package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)
	source := string(input)

	lines := strings.Split(source, "\n")
	var formatted []string
	indent := 0
	indentStr := "    "

	usingFixRegex := regexp.MustCompile(`using\s+package\s+\.(\w+)`)
	plusRegex := regexp.MustCompile(`\s*\+\s*`)
	commaRegex := regexp.MustCompile(`\s*,\s*`)
	assignRegex := regexp.MustCompile(`\s*=\s*`)
	openBlockRegex := regexp.MustCompile(`(?i)(\b(func|if|else if|else|for|class|with)\b[^\{]*)\{`)
	letCleanup := regexp.MustCompile(`^(export\s+)?let\s+([a-zA-Z0-9_]+)\s*=\s*(.*)`)
	elseRegex := regexp.MustCompile(`^(else\b|else if\b)`)
	oneLineBlock := regexp.MustCompile(`^\s*(class|func|if|for|while|with|else.*)\s*.*\{\s*\}$`)
	stringRegex := regexp.MustCompile(`"([^"\\]*(\\.[^"\\]*)*)"`)
	exportRegex := regexp.MustCompile(`(?m)^\s*export\s+(let|func|class)\b`)

	inMultilineComment := false
	var multilineCommentBlock []string

	for _, raw := range lines {
		raw = strings.TrimRight(raw, "\r\n")
		line := strings.TrimSpace(raw)

		if inMultilineComment {
			multilineCommentBlock = append(multilineCommentBlock, raw)
			if strings.Contains(line, "*/") {
				inMultilineComment = false
				formatted = append(formatted, multilineCommentBlock...)
				multilineCommentBlock = nil
			}
			continue
		}

		if strings.HasPrefix(line, "/*") && !strings.Contains(line, "*/") {
			inMultilineComment = true
			multilineCommentBlock = append(multilineCommentBlock, raw)
			continue
		}

		if line == "" {
			formatted = append(formatted, "")
			continue
		}

		if strings.HasPrefix(line, "//") {
			formatted = append(formatted, raw)
			continue
		}

		// split comment off the end
		var comment string
		if idx := strings.Index(line, "//"); idx != -1 {
			comment = line[idx:]
			line = strings.TrimSpace(line[:idx])
		}

		originalStrings := []string{}
		line = stringRegex.ReplaceAllStringFunc(line, func(s string) string {
			originalStrings = append(originalStrings, s)
			return fmt.Sprintf("__STR%d__", len(originalStrings)-1)
		})

		line = strings.ReplaceAll(line, "< =", "<=")
		line = strings.ReplaceAll(line, "> =", ">=")
		line = strings.ReplaceAll(line, "= =", "==")
		line = strings.ReplaceAll(line, "! =", "!=")

		line = strings.ReplaceAll(line, "<=", "__LE__")
		line = strings.ReplaceAll(line, ">=", "__GE__")
		line = strings.ReplaceAll(line, "==", "__EQ__")
		line = strings.ReplaceAll(line, "!=", "__NE__")

		line = usingFixRegex.ReplaceAllString(line, "using package.$1")
		line = assignRegex.ReplaceAllString(line, " = ")
		line = plusRegex.ReplaceAllString(line, " + ")
		line = commaRegex.ReplaceAllString(line, ", ")

		line = strings.ReplaceAll(line, "__LE__", "<=")
		line = strings.ReplaceAll(line, "__GE__", ">=")
		line = strings.ReplaceAll(line, "__EQ__", "==")
		line = strings.ReplaceAll(line, "__NE__", "!=")

		for i, str := range originalStrings {
			placeholder := fmt.Sprintf("__STR%d__", i)
			line = strings.Replace(line, placeholder, str, 1)
		}

		// Ensure export keyword is preserved once
		if exportRegex.MatchString(line) {
			line = exportRegex.ReplaceAllString(line, "export $1")
		}

		if letCleanup.MatchString(line) {
			matches := letCleanup.FindStringSubmatch(line)
			exportPrefix := strings.TrimSpace(matches[1]) // pode ser "export" ou ""
			line = fmt.Sprintf("%s let %s = %s", exportPrefix, matches[2], matches[3])
			line = strings.TrimSpace(line)
		}

		if oneLineBlock.MatchString(line) {
			line = strings.Repeat(indentStr, indent) + line
			if comment != "" {
				line += " " + comment
			}
			formatted = append(formatted, line)
			continue
		}

		if strings.HasSuffix(line, "= {") {
			line = strings.Repeat(indentStr, indent) + line
			if comment != "" {
				line += " " + comment
			}
			formatted = append(formatted, line)
			indent++
			continue
		}

		if strings.HasPrefix(line, "} else") {
			if indent > 0 {
				indent--
			}
			line = strings.Repeat(indentStr, indent) + line
			if comment != "" {
				line += " " + comment
			}
			formatted = append(formatted, line)
			indent++
			continue
		}

		if line == "{" || strings.HasSuffix(line, "{") {
			line = strings.Repeat(indentStr, indent) + line
			if comment != "" {
				line += " " + comment
			}
			formatted = append(formatted, line)
			indent++
			continue
		}

		if openBlockRegex.MatchString(line) {
			line = openBlockRegex.ReplaceAllStringFunc(line, func(m string) string {
				i := strings.LastIndex(m, "{")
				if i > 0 {
					before := strings.TrimRight(m[:i], " ")
					return before + " {"
				}
				return m
			})
			line = strings.Repeat(indentStr, indent) + line
			if comment != "" {
				line += " " + comment
			}
			formatted = append(formatted, line)
			indent++
			continue
		}

		if line == "}" || strings.HasPrefix(line, "}") {
			if indent > 0 {
				indent--
			}
			line = strings.Repeat(indentStr, indent) + line
			if comment != "" {
				line += " " + comment
			}
			formatted = append(formatted, line)
			continue
		}

		if elseRegex.MatchString(line) {
			if indent > 0 {
				indent--
				line = strings.Repeat(indentStr, indent) + line
				indent++
			} else {
				line = strings.Repeat(indentStr, indent) + line
			}
			if comment != "" {
				line += " " + comment
			}
			formatted = append(formatted, line)
			continue
		}

		line = strings.Repeat(indentStr, indent) + strings.TrimSpace(line)
		if comment != "" {
			line += " " + comment
		}
		formatted = append(formatted, line)
	}

	writer := bufio.NewWriter(os.Stdout)
	for i, f := range formatted {
		if i < len(formatted)-1 {
			fmt.Fprint(writer, f+"\n")
		} else {
			fmt.Fprint(writer, f)
		}
	}
	writer.Flush()
}
