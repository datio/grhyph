package grhyph

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	SpeechSound struct {
		Match                string
		Group                string
		EventualVowelsExist  bool
		ImmediateVowelExists bool
		ImmediateConsonants  int
	}

	Options struct {
		Separator            string
		MinHyphenationLength int
		UseGrhyphRules       bool
		CombineConsonantsDn  bool
		CombineConsonantsKv  bool
		CombineConsonantsPf  bool
		CombineConsonantsFk  bool
		QuickSynizesis       bool
	}

	Hyphenation struct {
		Input        string
		Options      Options
		SpeechSounds []SpeechSound
		WSCRe        *regexp.Regexp // WordStartConsonantsRe may vary between Hyphenation instances.
	}

	CacheKey struct {
		HyphenationInput   string
		HyphenationOptions Options // Have separate caches when options differ.
	}
)

var Cache = map[CacheKey]string{}
var CachingEnabled bool = true

var defaultOptions = Options{
	Separator:            "/",
	MinHyphenationLength: 2,
	CombineConsonantsFk:  true,
}

func GetDefaultOptions() Options {
	return defaultOptions
}

var speechSoundRe *regexp.Regexp = regexp.MustCompile(SpeechSoundRe)

func stringTospeechSounds(s string) ([]SpeechSound, error) {
	submatchGroups := speechSoundRe.SubexpNames()
	submatches := speechSoundRe.FindAllStringSubmatch(s, -1)
	speechSounds := make([]SpeechSound, len(submatches))

	var (
		eventualVowelsExist  bool
		immediateVowelExists bool
		immediateConsonants  int
	)

	for submatchIndex := len(submatches) - 1; submatchIndex >= 0; submatchIndex-- {
		// ["" "vowels" "consonants" "punctuation" "other"]
		speechSoundGroupIndex := 0

		for groupIndex := 1; groupIndex <= 4; groupIndex++ {
			if submatches[submatchIndex][groupIndex] == "" {
				continue
			}

			speechSoundGroupIndex = groupIndex
			break
		}

		if submatchIndex+1 < len(speechSounds) {
			switch speechSounds[submatchIndex+1].Group {
			case "vowels":
				immediateVowelExists = true
			default:
				immediateVowelExists = false
			}
		}

		speechSounds[submatchIndex].Match = submatches[submatchIndex][0]
		speechSounds[submatchIndex].Group = submatchGroups[speechSoundGroupIndex]
		speechSounds[submatchIndex].EventualVowelsExist = eventualVowelsExist
		speechSounds[submatchIndex].ImmediateVowelExists = immediateVowelExists
		speechSounds[submatchIndex].ImmediateConsonants = immediateConsonants

		switch submatchGroups[speechSoundGroupIndex] {
		case "vowels":
			eventualVowelsExist = true
		case "punctuation", "other":
			eventualVowelsExist = false
		}

		if submatchGroups[speechSoundGroupIndex] == "consonants" {
			immediateConsonants++
		} else {
			immediateConsonants = 0
		}
	}

	// Debug:
	// fmt.Println(speechSounds)

	return speechSounds, nil
}

func (h *Hyphenation) Hyphenate() (string, error) {
	speechSounds, err := stringTospeechSounds(h.Input)
	if err != nil {
		return "", err
	}

	h.SpeechSounds = speechSounds
	h.WSCRe = GetWSCRe(h.Options.CombineConsonantsDn, h.Options.CombineConsonantsKv,
		h.Options.CombineConsonantsPf, h.Options.CombineConsonantsFk)

	if !h.Options.UseGrhyphRules {
		return plainHyphenation(speechSounds, h.Options, h.WSCRe)
	} else {
		return h.regexpHyphenation()
	}

	return "", nil
}

var synizesisVowelsRe *regexp.Regexp = regexp.MustCompile(SynizesisVowelsRe)

// Hyphenate without using the GrhyphRules (RegExp exceptions) definitions.
func plainHyphenation(ss []SpeechSound, o Options, wSCRe *regexp.Regexp) (string, error) {
	if len(ss) <= 1 || len(ss) < o.MinHyphenationLength {
		return speechSoundJoin(ss), nil
	}

	var hyphenated []byte

	for i := 0; i < len(ss); i++ {
		if ss[i].Group == "consonants" && ss[i].ImmediateVowelExists {
			hyphenated = append(hyphenated, ss[i].Match...)
			continue
		} else if ss[i].Group == "vowels" {
			if ss[i].EventualVowelsExist && ss[i].ImmediateConsonants == 1 {
				hyphenated = append(hyphenated, fmt.Sprintf("%s%s", ss[i].Match, o.Separator)...)
				continue
			} else if ss[i].ImmediateConsonants >= 1 && !ss[i].EventualVowelsExist {
				hyphenated = append(hyphenated, ss[i].Match...)
				continue
			} else if ss[i].EventualVowelsExist && ss[i].ImmediateConsonants > 1 {
				hyphenedConsonants := consonantHyphenation(i+1, ss[i].ImmediateConsonants, ss, o, wSCRe)
				hyphenated = append(hyphenated, fmt.Sprintf("%s%s", ss[i].Match, hyphenedConsonants)...)
				i += ss[i].ImmediateConsonants
				continue
			} else if ss[i].ImmediateVowelExists {
				// Flag for quick-synizesis / end-of-the-line hyphenation.
				if o.QuickSynizesis &&
					synizesisVowelsRe.MatchString(fmt.Sprintf("%s%s", ss[i].Match, ss[i+1].Match)) {
					hyphenated = append(hyphenated, ss[i].Match...)
					continue
				}
				hyphenated = append(hyphenated, fmt.Sprintf("%s%s", ss[i].Match, o.Separator)...)
				continue
			}
		}

		hyphenated = append(hyphenated, ss[i].Match...)
	}

	return string(hyphenated[:]), nil
}

func consonantHyphenation(startIndex int, consonantsN int,
	ss []SpeechSound, o Options, wSCRe *regexp.Regexp) string {
	var hyphenatedConsonants []byte

	endIndex := startIndex + consonantsN
	for i := startIndex; i < endIndex; i++ {
		if i == endIndex-1 {
			hyphenatedConsonants = append(hyphenatedConsonants, fmt.Sprintf("%s%s", o.Separator, ss[i].Match)...)
			break
		}

		consonantsPair := fmt.Sprintf("%s%s", ss[i].Match, ss[i+1].Match)
		if wSCRe.MatchString(consonantsPair) {
			hyphenatedConsonants = append(hyphenatedConsonants, o.Separator...)
			for ; i < endIndex; i++ {
				hyphenatedConsonants = append(hyphenatedConsonants, ss[i].Match...)
			}
			break
		} else {
			hyphenatedConsonants = append(hyphenatedConsonants, ss[i].Match...)
			continue
		}
	}

	return string(hyphenatedConsonants[:])
}

func (h *Hyphenation) regexpHyphenation() (string, error) {
	if len(h.Input) <= 1 || len(h.Input) < h.Options.MinHyphenationLength {
		return h.Input, nil
	}

	var err error
	var hyphenated []byte

	// Separate multiple input words by using the speechSound punctuation points.
	lastPunctuationIndex := -1
	for i, speechSounds := range h.SpeechSounds {
		start := lastPunctuationIndex + 1
		end := i - 1
		isLastIteration := (i == len(h.SpeechSounds)-1)

		if speechSounds.Group == "punctuation" {
			if start >= 0 && i-start > 1 {
				if (i - start) >= h.Options.MinHyphenationLength {
					hyphenated = append(hyphenated, regexpReplace(h.SpeechSounds[start:i], h.Options, h.WSCRe)...)
				} else {
					hyphenated = append(hyphenated, speechSoundJoin(h.SpeechSounds[start:i])...)
				}
			} else if start >= 0 && end-start == 0 {
				hyphenated = append(hyphenated, h.SpeechSounds[start].Match...)
			}
			hyphenated = append(hyphenated, h.SpeechSounds[i].Match...)
			lastPunctuationIndex = i
		} else if isLastIteration {
			hyphenated = append(hyphenated, regexpReplace(h.SpeechSounds[start:], h.Options, h.WSCRe)...)
		}
	}

	return string(hyphenated[:]), err
}

func speechSoundJoin(speechSounds []SpeechSound) string {
	var joinedMatchesBytes []byte
	for i := 0; i < len(speechSounds); i++ {
		joinedMatchesBytes = append(joinedMatchesBytes, speechSounds[i].Match...)
	}

	return string(joinedMatchesBytes[:])
}

func regexpReplace(speechSounds []SpeechSound, o Options, wSCRe *regexp.Regexp) string {
	joinedSpeechSounds := speechSoundJoin(speechSounds)

	if CachingEnabled {
		cacheKey := CacheKey{joinedSpeechSounds, o}

		if hyphenatedString, ok := Cache[cacheKey]; ok {
			return hyphenatedString
		}
	}

	for _, rule := range GrhyphRules {
		if rule.CompiledCustomRe.MatchString(joinedSpeechSounds) {
			replacement := strings.Replace(rule.Replacement, "-", o.Separator, -1)

			var (
				middleRunes      []byte
				toHyphenateLeft  string
				toHyphenateRight string
			)

			for i, rune := range replacement {
				currentCharacter := fmt.Sprintf("%c", rune)
				if currentCharacter == ">" {
					toHyphenateLeft = rule.CompiledCustomRe.ReplaceAllString(joinedSpeechSounds, replacement[0:i])
					middleRunes = nil
				} else if currentCharacter == "<" && i < len(replacement) {
					toHyphenateRight = rule.CompiledCustomRe.ReplaceAllString(joinedSpeechSounds, replacement[i+1:])
					break
				} else {
					middleRunes = append(middleRunes, currentCharacter...)
				}
			}

			leftSpeechSounds, _ := stringTospeechSounds(toHyphenateLeft)
			toHyphenateLeft = regexpReplace(leftSpeechSounds, o, wSCRe)

			rightSpeechSounds, _ := stringTospeechSounds(toHyphenateRight)
			toHyphenateRight = regexpReplace(rightSpeechSounds, o, wSCRe)

			hyphenedMiddle := rule.CompiledCustomRe.ReplaceAllString(joinedSpeechSounds, string(middleRunes[:]))

			// Debug:
			// fmt.Println(rule)
			// todo: Return the pre-compiled rules, if possible with a record of the line in which they were defined.

			return toHyphenateLeft + hyphenedMiddle + toHyphenateRight
		}
	}

	hyphened, _ := plainHyphenation(speechSounds, o, wSCRe)

	if CachingEnabled {
		cacheKey := CacheKey{HyphenationInput: joinedSpeechSounds, HyphenationOptions: o}

		if _, ok := Cache[cacheKey]; !ok {
			Cache[cacheKey] = hyphened
		}
	}

	return hyphened
}
