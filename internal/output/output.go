package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

var Format = "table"

var preferred = []string{
	"id", "name", "status", "primary_domain", "domains", "expires", "token_type", "secret_type", "risk_level",
	"action", "outcome", "severity", "category", "occurred_at", "user_email", "message",
}

func Write(raw []byte) error {
	return WriteTo(os.Stdout, raw)
}

func WriteTo(w io.Writer, raw []byte) error {
	if Format == "json" {
		var buf bytes.Buffer
		if err := json.Indent(&buf, raw, "", "  "); err != nil {
			_, err = w.Write(append(bytes.TrimSpace(raw), '\n'))
			return err
		}
		buf.WriteByte('\n')
		_, err := w.Write(buf.Bytes())
		return err
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		_, err = fmt.Fprintln(w, string(raw))
		return err
	}
	return writeHuman(w, value)
}

func writeHuman(w io.Writer, value any) error {
	switch typed := value.(type) {
	case map[string]any:
		if content, ok := typed["content"].([]any); ok {
			if err := writeRows(w, content); err != nil {
				return err
			}
			if total, ok := typed["total"]; ok {
				_, err := fmt.Fprintf(w, "\ntotal: %v\n", total)
				return err
			}
			return nil
		}
		return writeObject(w, typed)
	case []any:
		if len(typed) == 0 {
			_, err := fmt.Fprintln(w, "No results.")
			return err
		}
		if _, ok := typed[0].(map[string]any); ok {
			return writeRows(w, typed)
		}
		for _, item := range typed {
			if _, err := fmt.Fprintln(w, formatScalar(item)); err != nil {
				return err
			}
		}
		return nil
	default:
		_, err := fmt.Fprintln(w, formatScalar(value))
		return err
	}
}

func writeRows(w io.Writer, rows []any) error {
	if len(rows) == 0 {
		_, err := fmt.Fprintln(w, "No results.")
		return err
	}
	maps := make([]map[string]any, 0, len(rows))
	seen := map[string]bool{}
	for _, row := range rows {
		object, ok := row.(map[string]any)
		if !ok {
			_, err := fmt.Fprintln(w, formatScalar(row))
			return err
		}
		maps = append(maps, object)
		for key := range object {
			seen[key] = true
		}
	}
	columns := orderedKeys(seen)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, strings.Join(columns, "\t")); err != nil {
		return err
	}
	for _, row := range maps {
		cells := make([]string, len(columns))
		for i, column := range columns {
			cells[i] = oneLine(formatScalar(row[column]))
		}
		if _, err := fmt.Fprintln(tw, strings.Join(cells, "\t")); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func writeObject(w io.Writer, object map[string]any) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, key := range orderedKeys(object) {
		text := formatScalar(object[key])
		if strings.Contains(text, "\n") {
			if err := tw.Flush(); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, "%s:\n%s\n", key, text); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", key, text); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func orderedKeys(keys any) []string {
	seen := map[string]bool{}
	switch typed := keys.(type) {
	case map[string]bool:
		seen = typed
	case map[string]any:
		for key := range typed {
			seen[key] = true
		}
	}
	var columns []string
	for _, key := range preferred {
		if seen[key] {
			columns = append(columns, key)
			delete(seen, key)
		}
	}
	rest := make([]string, 0, len(seen))
	for key := range seen {
		rest = append(rest, key)
	}
	sort.Strings(rest)
	return append(columns, rest...)
}

func formatScalar(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case map[string]any, []any:
		raw, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprint(typed)
		}
		return string(raw)
	default:
		return fmt.Sprint(typed)
	}
}

func oneLine(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	if len(value) > 80 {
		return value[:77] + "..."
	}
	return value
}
