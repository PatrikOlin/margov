package main

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Margov struct {
	n           int
	capitalized int
	suffix      map[string][]string
}

func NewMargovFromFile(filename string, n int) (*Margov, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewMargov(f, n)
}

func NewMargov(r io.Reader, n int) (*Margov, error) {
	m := &Margov{
		n:      n,
		suffix: make(map[string][]string),
	}
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords)
	window := make([]string, 0, n)
	for sc.Scan() {
		word := sc.Text()
		if len(window) > 0 {
			prefix := strings.Join(window, " ")
			m.suffix[prefix] = append(m.suffix[prefix], word)
			if isCapitalized(prefix) {
				m.capitalized++
			}
		}
		window = appendMax(n, window, word)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Margov) Output(w io.Writer, n int, startCapital, stopSentence bool) error {
	bw := bufio.NewWriter(w)

	var i int
	if startCapital {
		i = rand.Intn(m.capitalized)
	} else {
		i = rand.Intn(len(m.suffix))
	}

	var prefix string
	for prefix = range m.suffix {
		if startCapital && !isCapitalized(prefix) {
			continue
		}
		if i == 0 {
			break
		}
		i--
	}

	bw.WriteString(prefix)
	prefixWords := strings.Fields(prefix)
	n -= len(prefixWords)

	for {
		suffixChoices := m.suffix[prefix]
		if len(suffixChoices) == 0 {
			break
		}
		i = rand.Intn(len(suffixChoices))
		suffix := suffixChoices[i]
		bw.WriteByte(' ')
		if _, err := bw.WriteString(suffix); err != nil {
			break
		}
		n--
		if n < 0 && (!stopSentence || isSentenceEnd(suffix)) {
			break
		}

		prefixWords = appendMax(m.n, prefixWords, suffix)
		prefix = strings.Join(prefixWords, " ")
	}

	return bw.Flush()
}

func isCapitalized(s string) bool {
	r, _ := utf8.DecodeLastRuneInString(s)
	return unicode.IsUpper(r)
}

func isSentenceEnd(s string) bool {
	r, _ := utf8.DecodeLastRuneInString(s)
	return r == '.' || r == '?' || r == '!'
}

func appendMax(max int, slice []string, value string) []string {
	if len(slice)+1 > max {
		n := copy(slice, slice[1:])
		slice = slice[:n]
	}
	return append(slice, value)
}
