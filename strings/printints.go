//IntsToStrings is like fmt.Spint but adds commas
package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(IntsToStrings([]int{1, 2, 3, 4, 5, 6}))
}

func IntsToStrings(values []int) string {
	var buf bytes.Buffer

	buf.WriteByte('[')

	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}

	buf.WriteByte(']')

	return buf.String()

}
