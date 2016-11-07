package main

import (
	"flag"
	"fmt"
	"github.com/datio/grhyph"
)

func main() {
	// minHyphenationLength := flag.Int("min-length", 2, "Minimum length of a word to hyphenate.")

	disableCaching := flag.Bool("no-cache", false, "Disable caching when rule-hyphenating words.")

	quickSynizesis := flag.Bool("quick-synizesis", false, `Combine nearby vowels, whenever synizesis is prone to occur.`)

	// The separator defaults to a single soft hyphen (U+00AD SOFT HYPHEN: "­").
	separator := flag.String("separator", "­", `The separator to append between syllables.`)

	useGrhyphRules := flag.Bool("use-rules", false, `Match and replace using rules, based on regular expressions,
	 as defined in the definitions.go file.`)

	flag.Parse()

	grhyph.CachingEnabled = !*disableCaching

	hyphenationOptions := grhyph.GetDefaultOptions()

	// if *minHyphenationLength > 1 {
	// hyphenationOptions.MinHyphenationLength = *minHyphenationLength
	// }

	hyphenationOptions.QuickSynizesis = *quickSynizesis
	hyphenationOptions.Separator = *separator
	hyphenationOptions.UseGrhyphRules = *useGrhyphRules

	h := grhyph.Hyphenation{
		Options: hyphenationOptions,
	}

	inputs := flag.Args()
	for _, input := range inputs {
		h.Input = input

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			fmt.Println(fmt.Errorf("grhyph err:\n%v", err))
			return
		}

		fmt.Println(hyphenedText)
	}
}
