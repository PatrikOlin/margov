package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("margov: ")
	input := flag.String("in", "horoskop", "input file")
	n := flag.Int("n", 2, "number of words to use as prefix")
	runs := flag.Int("runs", 1, "number of runs to generate")
	wordsPerRun := flag.Int("words", 150, "number of words per run")
	startOnCapital := flag.Bool("capital", true, "start output with capitalized prefix")
	stopAtSentence := flag.Bool("sentence", true, "end output at a sentence ending punctuation mark")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	m, err := NewMargovFromFile(*input, *n)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < *runs; i++ {
		err = m.Output(os.Stdout, *wordsPerRun, *startOnCapital, *stopAtSentence)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	}

}
