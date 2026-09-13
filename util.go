package hermeti

import (
	"strings"
)

// take strings of the form "foo=bar" and return a map
func stringsToMap(kvs []string) map[string]string {
	m := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		x := strings.Split(kv, "=")
		if len(x) == 2 {
			m[x[0]] = x[1]
		}
		if len(x) == 1 {
			m[kv] = ""
		}
		if len(x) > 2 {
			m[x[0]] = strings.Join(x[1:], "=")
		}
	}
	return m
}
