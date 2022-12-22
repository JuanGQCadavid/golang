//Basename remove directory components and a .suffix
// e.g., a => a, a.go => a, a/b/c/d/a.go => a, a/b.c.d => b.c
package main

import "fmt"

func main() {
	fmt.Println(basename("a/b/c/d.go"))
	fmt.Println(basename("c.d.go"))
	fmt.Println(basename("abc"))
}

func basename(s string) string {
	// Discard last '/' and everiting before.

	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}

	// Preserve everything before last '.'
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}

	return s

}
