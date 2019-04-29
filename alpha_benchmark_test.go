// +build bench

package goradix

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/karrick/gobls"
)

const dictionaryPathname = "/usr/share/dict/words"

var urls []string
var indexes []int

func init() {
	fh, err := os.Open(dictionaryPathname)
	if err != nil {
		log.Fatal(err)
	}

	words := make([]string, 0, 26)
	initials := make(map[rune]struct{})

	scanner := gobls.NewScanner(fh)
	for scanner.Scan() {
		if len(words) > 5 {
			continue
		}
		word := scanner.Text()
		for _, rune := range word {
			if _, ok := initials[rune]; ok {
				break
			}
			if unicode.IsUpper(rune) && len(word) >= 5 && len(word) <= 10 {
				initials[rune] = struct{}{}
				words = append(words, word)
			}
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	urls = make([]string, 0, 1000) // make a few urls
	for i := 0; i < cap(urls); i++ {
		parts := make([]string, 0, 10) // each url has a few parts
		for j := 0; j < cap(parts); j++ {
			wi := rand.Intn(len(words))
			parts = append(parts, words[wi])
		}
		urls = append(urls, "/"+strings.Join(parts, "/"))
	}

	// Create a randomly generated list of indexes for the urls so that
	// different benchmarks for this run access the same sequence of urls.
	indexes = rand.Perm(len(urls))
}

func BenchmarkUrlMap(b *testing.B) {
	b.Skip("not looking at url now")
	m := make(map[string]int)
	for i, path := range urls {
		m[path] = i
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, ok := m[urls[indexes[i%len(indexes)]]]
		if got, want := ok, true; got != want {
			b.Fatalf("GOT: %v; WANT: %v", got, want)
		}
	}
}

func BenchmarkUrlTrie(b *testing.B) {
	b.Skip("not looking at url now")
	root := new(Alpha)
	for i, path := range urls {
		root.Store(path, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, ok := root.Load(urls[indexes[i%len(indexes)]])
		if got, want := ok, true; got != want {
			b.Fatalf("GOT: %v; WANT: %v", got, want)
		}
	}
}
