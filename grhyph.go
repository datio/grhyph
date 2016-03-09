package grhyph

import (
	"fmt"
	"github.com/datio/grhyph/definitions"
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
		Seperator            string
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
	Seperator:            "/",
	MinHyphenationLength: 2,
	CombineConsonantsFk:  true,
}

func GetDefaultOptions() Options {
	return defaultOptions
}

var speechSoundRe *regexp.Regexp = regexp.MustCompile(grhyph.SpeechSoundRe)

func stringTospeechSounds(s string) ([]SpeechSound, error) {
	submatchGroups := speechSoundRe.SubexpNames()
	submatches := speechSoundRe.FindAllStringSubmatch(s, -1)
	speechSounds := make([]SpeechSound, len(submatches))

	// It's not possible to use PCRE-like lookaheads in RE2.
	// The greeklish equivelant of the consonant combination 'νθ' (/nθ/ -> 'nth')
	// must split as ["n" "th"] instead of ["nt" "h"].
	// todo: Extra rules/tests for (sth: σθ|στη), (nth: νθ|ντη).
	nuTauRe, err := regexp.Compile(grhyph.NuTauRe)
	if err != nil {
		return speechSounds, err
	}

	var (
		eventualVowelsExist  bool
		immediateVowelExists bool
		immediateConsonants  int
	)

	for submatchIndex := len(submatches) - 1; submatchIndex >= 0; submatchIndex-- {
		// Greeklish speechSound seperation fix for 'Nu' 'Tau' 'H' combinations.
		h := submatches[submatchIndex][0]

		// If the current index includes an 'H' vowel, match the previous one to a 'NuTau' combination.
		if (h == "h" || h == "H") && submatchIndex > 0 {
			nt := nuTauRe.FindStringSubmatch(submatches[submatchIndex-1][0])

			if len(nt) > 1 && len(submatches)-1 != submatchIndex &&
				!(submatchIndex > 0 && submatchIndex+1 < len(submatches) && len(submatches[submatchIndex+1][3]) > 0) {

				// Replace [νn][τt] with the [νn] match.
				submatches[submatchIndex-1][0] = nt[1]
				submatches[submatchIndex-1][2] = nt[1]

				// Append the [τt] match before [h].
				// The speechSound 'h' normally belongs to the "vowels" group.
				// Move it to the "consonants" group.
				submatches[submatchIndex][0] = nt[2] + h
				submatches[submatchIndex][2] = nt[2] + h
				submatches[submatchIndex][1] = "" //remove from "vowels".
			}
		}

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
		case "punctuation":
			eventualVowelsExist = false
		}

		if submatchGroups[speechSoundGroupIndex] != "consonants" {
			immediateConsonants = 0
		} else {
			immediateConsonants++
		}
	}

	return speechSounds, nil
}

func (h *Hyphenation) Hyphenate() (string, error) {
	ss, err := stringTospeechSounds(h.Input)
	if err != nil {
		return "", err
	}

	h.SpeechSounds = ss
	h.WSCRe = grhyph.GetWSCRe(h.Options.CombineConsonantsDn, h.Options.CombineConsonantsKv,
		h.Options.CombineConsonantsPf, h.Options.CombineConsonantsFk)

	if !h.Options.UseGrhyphRules {
		return plainHyphenation(ss, h.Options, h.WSCRe)
	} else {
		return h.regexpHyphenation()
	}

	return "", nil
}

var synizesisVowelsRe *regexp.Regexp = regexp.MustCompile(grhyph.SynizesisVowelsRe)

// Hyphenate without using the GrhyphRules set from the definitions.
func plainHyphenation(ss []SpeechSound, o Options, wSCRe *regexp.Regexp) (string, error) {
	if len(ss) <= 1 || len(ss) < o.MinHyphenationLength {
		return speechSoundJoin(ss), nil
	}

	var hyphenedBytes []byte

	for i := 0; i < len(ss); i++ {
		if ss[i].Group == "consonants" && ss[i].ImmediateVowelExists {
			hyphenedBytes = append(hyphenedBytes, ss[i].Match...)
			continue
		} else if ss[i].Group == "vowels" {
			if ss[i].EventualVowelsExist && ss[i].ImmediateConsonants == 1 {
				hyphenedBytes = append(hyphenedBytes, fmt.Sprintf("%s%s", ss[i].Match, o.Seperator)...)
				continue
			} else if ss[i].ImmediateConsonants >= 1 && !ss[i].EventualVowelsExist {
				hyphenedBytes = append(hyphenedBytes, ss[i].Match...)
				continue
			} else if ss[i].EventualVowelsExist && ss[i].ImmediateConsonants > 1 {
				hyphenedConsonants := consonantHyphenation(i+1, ss[i].ImmediateConsonants, ss, o, wSCRe)
				hyphenedBytes = append(hyphenedBytes, fmt.Sprintf("%s%s", ss[i].Match, hyphenedConsonants)...)
				i += ss[i].ImmediateConsonants
				continue
			} else if ss[i].ImmediateVowelExists {
				// Flag useful for quick end-of-the-line splitting/hyphenation.
				if o.QuickSynizesis &&
					synizesisVowelsRe.MatchString(fmt.Sprintf("%s%s", ss[i].Match, ss[i+1].Match)) {
					hyphenedBytes = append(hyphenedBytes, ss[i].Match...)
					continue
				}
				hyphenedBytes = append(hyphenedBytes, fmt.Sprintf("%s%s", ss[i].Match, o.Seperator)...)
				continue
			}
		}

		hyphenedBytes = append(hyphenedBytes, ss[i].Match...)
	}

	return string(hyphenedBytes[:]), nil
}

func consonantHyphenation(startIndex int, consonantsN int,
	ss []SpeechSound, o Options, wSCRe *regexp.Regexp) string {
	var hyphenenedConsonants []byte

	endIndex := startIndex + consonantsN
	for i := startIndex; i < endIndex; i++ {
		if i == endIndex-1 {
			hyphenenedConsonants = append(hyphenenedConsonants, fmt.Sprintf("%s%s", o.Seperator, ss[i].Match)...)
			break
		}

		consonantsPair := fmt.Sprintf("%s%s", ss[i].Match, ss[i+1].Match)
		if wSCRe.MatchString(consonantsPair) {
			hyphenenedConsonants = append(hyphenenedConsonants, o.Seperator...)
			for ; i < endIndex; i++ {
				hyphenenedConsonants = append(hyphenenedConsonants, ss[i].Match...)
			}
			break
		} else {
			hyphenenedConsonants = append(hyphenenedConsonants, ss[i].Match...)
			continue
		}
	}

	return string(hyphenenedConsonants[:])
}

func (h *Hyphenation) regexpHyphenation() (string, error) {
	if len(h.Input) <= 1 || len(h.Input) < h.Options.MinHyphenationLength {
		return h.Input, nil
	}

	var err error
	var hyphenedBytes []byte

	// Seperate words using speechSounds punctuations.
	// A bit less complicated than, say, an additional *Regexp based Split would be.
	lastPunctuationIndex := -1
	for i, ss := range h.SpeechSounds {
		start := lastPunctuationIndex + 1
		end := i - 1
		isLastIteration := (i == len(h.SpeechSounds)-1)

		if ss.Group == "punctuation" {
			if start >= 0 && i-start > 1 {
				if (i - start) >= h.Options.MinHyphenationLength {
					hyphenedBytes = append(hyphenedBytes, regexpReplace(h.SpeechSounds[start:i], h.Options, h.WSCRe)...)
				} else {
					hyphenedBytes = append(hyphenedBytes, speechSoundJoin(h.SpeechSounds[start:i])...)
				}
			} else if start >= 0 && end-start == 0 {
				hyphenedBytes = append(hyphenedBytes, h.SpeechSounds[start].Match...)
			}
			hyphenedBytes = append(hyphenedBytes, h.SpeechSounds[i].Match...)
			lastPunctuationIndex = i
		} else if isLastIteration {
			hyphenedBytes = append(hyphenedBytes, regexpReplace(h.SpeechSounds[start:], h.Options, h.WSCRe)...)
		}
	}

	return string(hyphenedBytes[:]), err
}

func speechSoundJoin(ss []SpeechSound) string {
	var joinedMatchesBytes []byte
	for i := 0; i < len(ss); i++ {
		joinedMatchesBytes = append(joinedMatchesBytes, ss[i].Match...)
	}

	return string(joinedMatchesBytes[:])
}

func regexpReplace(ss []SpeechSound, o Options, wSCRe *regexp.Regexp) string {
	joinedSpeechSounds := speechSoundJoin(ss)

	if CachingEnabled {
		ck := CacheKey{joinedSpeechSounds, o}

		if hyphenatedString, ok := Cache[ck]; ok {
			return hyphenatedString
		}
	}

	for _, rule := range grhyph.GrhyphRules {
		if rule.CompiledCustomRe.MatchString(joinedSpeechSounds) {
			repl := strings.Replace(rule.Replacement, "-", o.Seperator, -1)

			var (
				middleRunes      []byte
				toHyphenateLeft  string
				toHyphenateRight string
			)

			for i, rune := range repl {
				currentCharacter := fmt.Sprintf("%c", rune)
				if currentCharacter == ">" {
					toHyphenateLeft = rule.CompiledCustomRe.ReplaceAllString(joinedSpeechSounds, repl[0:i])
					middleRunes = nil
				} else if currentCharacter == "<" && i < len(repl) {
					toHyphenateRight = rule.CompiledCustomRe.ReplaceAllString(joinedSpeechSounds, repl[i+1:])
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
			fmt.Println(rule)
			// todo: return pre-compiled rule, maybe with a record of the line it exists at.

			return toHyphenateLeft + hyphenedMiddle + toHyphenateRight
		}
	}

	hyphened, _ := plainHyphenation(ss, o, wSCRe)

	if CachingEnabled {
		ck := CacheKey{HyphenationInput: joinedSpeechSounds, HyphenationOptions: o}

		if _, ok := Cache[ck]; !ok {
			Cache[ck] = hyphened
		}
	}

	return hyphened
}
