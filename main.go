package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const WordsPath = "/usr/share/dict/words"

func Exit(c int, v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(c)
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		Exit(2, os.Args[0], "input [length]")
	}
	var err error
	input := []byte(os.Args[1])
	if !Alphagram(input) {
		Exit(2, "input argument must contain only letters")
	}
	maxlen := len(input)
	if len(os.Args) == 3 {
		maxlen, err = strconv.Atoi(os.Args[2])
		if err != nil {
			Exit(2, "Invalid length argument:", err.Error())
		} else if maxlen > len(input) {
			Exit(2, "length argument must not be greater than len(input)")
		}
	}
	f, err := os.Open(WordsPath)
	if err != nil {
		Exit(1, "Error opening words file:", err)
	}
	defer f.Close()
	w := bufio.NewWriter(os.Stdout)
	r := bufio.NewScanner(f)
	var buf []byte
	for r.Scan() {
		word := r.Bytes()
		buf = append(buf[:0], word...)
		// ignore words that are too long or contain non-letters
		if len(buf) <= maxlen && Alphagram(buf) && Match(buf, input) {
			w.Write(word)
			w.WriteByte('\n')
		}
	}
	//w.Flush()
	if err := r.Err(); err != nil {
		Exit(1, "Error reading words file:", err)
	}
}

// es
// aaccderrss

// Match returns true if s is in t, which must each already be Alphagram processed.
func Match(s, t []byte) bool {
	//fmt.Printf("%q %q\n", t, s)
	i := 0
outer:
	for _, b := range s {
		for _, c := range t[i:] {
			i++
			if c == b {
				continue outer
			} else if c > b {
				return false
			}
		}
		return false
	}
	return true
}

// Alphagram sorts and down-cases s, returning false if s contains non-letters.
func Alphagram(s []byte) bool {
	for i, b := range s {
		if 'A' <= b && b <= 'Z' {
			s[i] = b + 'a' - 'A'
		} else if b < 'a' || 'z' < b {
			return false
		}
	}
	sort.Sort(alpha(s))
	return true
}

type alpha []byte

func (s alpha) Len() int           { return len(s) }
func (s alpha) Less(i, j int) bool { return s[i] < s[j] }
func (s alpha) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
