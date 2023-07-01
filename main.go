package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type ARB map[string]any

func (a ARB) JSON() []byte {
	// write keys sorted
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetIndent(strings.Repeat(" ", 20), "\t")
	io.WriteString(b, "{\n")
	keys := []string{}
	for k := range a {
		keys = append(keys, k)
	}
	// put @ version below the one without
	sort.Slice(keys, func(i, j int) bool {
		k1 := keys[i]
		k2 := keys[j]
		if strings.HasPrefix(k1, "@") {
			k1 = k1[1:] + "@"
		}
		if strings.HasPrefix(k2, "@") {
			k2 = k2[1:] + "@"
		}
		return k1 < k2
	})
	for i, k := range keys {
		if i > 0 {
			fmt.Fprintf(b, ",\n")
		}
		v := a[k]
		fmt.Fprintf(b, "%q:%s", k, padding(k))
		if s, ok := v.(string); ok {
			fmt.Fprintf(b, "%q", s)
		} else {
			enc.Encode(v)
		}
	}
	io.WriteString(b, "\n}")
	return b.Bytes()
}

func padding(key string) string {
	return strings.Repeat(" ", 20-len(key))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("arb <source>.arb <target1>.arb <target2>.arb")
		return
	}
	sourceFile := os.Args[1]
	source := readARB(sourceFile)
	writeARB(source, sourceFile)
	for t := 2; t < len(os.Args); t++ {
		targetFile := os.Args[t]
		target := readARB(targetFile)
		sync(source, target)
		writeARB(target, targetFile)
	}
}

func sync(source, target ARB) {
	for ik, iv := range source {
		if _, ok := target[ik]; !ok {
			// no key in target
			// value is string ?
			if ivs, ok := iv.(string); ok {
				target[ik] = fmt.Sprintf("?%s?", ivs)
			} else {
				target[ik] = iv
			}
		}
	}
}

func writeARB(arb ARB, to string) {
	os.WriteFile(to, arb.JSON(), os.ModePerm)
}

func readARB(from string) (arb ARB) {
	data, err := os.ReadFile(from)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &arb)
	if err != nil {
		fmt.Println("error reading:", from)
		panic(err)
	}
	return
}
