package grhyph

import (
	"testing"
)

type hyphenationTest struct {
	input      string
	hyphenated string
}

func TestPlain(t *testing.T) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "-"

	h := Hyphenation{
		Options: hyphenationOptions,
	}

	tests := []hyphenationTest{
		{"αλκμιόνη", "αλκ-μι-ό-νη"},
		{"aλkμiοnη", "aλk-μi-ο-nη"},
		{"grafete Ellhnika", "gra-fe-te El-lh-ni-ka"},
		{`δυο
      γραμμές.`, `δυ-ο
      γραμ-μές.`},
		{"Mode Plagal - Εμένα μου το παν τα πουλιά", "Mo-de Pla-gal - Ε-μέ-να μου το παν τα που-λι-ά"},
		{"Μια και δυο. Μία και δύο", "Μι-α και δυ-ο. Μί-α και δύ-ο"},
		{"χελιδόνια", "χε-λι-δό-νι-α"},
		{"αηδόνια", "α-η-δό-νι-α"},
		{"English words are hyphenated using Greek/Greeklish grammar. Their detection and exclusion has to happen outside of this package.", "En-gli-sh w-ords a-re h-y-ph-e-na-ted u-sing Gre-ek/Gre-e-kli-sh gram-mar. Their de-tec-ti-on and exc-lu-si-on h-as to h-ap-pen ou-tsi-de of this pac-ka-ge."},
	}

	for _, test := range tests {
		h.Input = test.input

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			panic(err)
		}

		if hyphenedText != test.hyphenated {
			t.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", test.input, test.hyphenated, hyphenedText)
		}
	}
}

func TestQuickSynizesis(t *testing.T) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "/"
	hyphenationOptions.QuickSynizesis = true

	h := Hyphenation{
		Options: hyphenationOptions,
	}

	tests := []hyphenationTest{
		{"αλκμιόνη", "αλκ/μιό/νη"},
		{"Μια και δυο. Μία και δύο", "Μια και δυο. Μί/α και δύ/ο"},
		{"δυο μαύρα χελιδόνια", "δυο μαύ/ρα χε/λι/δό/νια"},
		{"αηδόνια", "αη/δό/νια"},
		{"αραχνοΰφαντος", "α/ρα/χνο/ΰ/φα/ντος"},
		{"αραχνούφαντος", "α/ρα/χνού/φα/ντος"},
	}

	for _, test := range tests {
		h.Input = test.input

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			panic(err)
		}

		if hyphenedText != test.hyphenated {
			t.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", test.input, test.hyphenated, hyphenedText)
		}
	}
}

// Use no cache.
func TestRules(t *testing.T) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "-"
	hyphenationOptions.UseGrhyphRules = true

	h := Hyphenation{
		Options: hyphenationOptions,
	}

	tests := []hyphenationTest{
		{"αλκμιόνη", "αλκ-μι-ό-νη"},
		// {"Μια και δυο. Μία και δύο", "Μια και δυο. Μί-α και δύ-ο"}, // todo: μ, δ first.
		// {"δυο μαύρα χελιδόνια", "δυο μαύ/ρα χε/λι/δό/νια"}, // todo: χ first.
		{"αηδόνια", "αη-δό-νια"},
		{"αουτσαιντερ", "α-ου-τσα-ι-ντερ"},
		{"αραχνοΰφαντος", "α-ρα-χνο-ΰ-φα-ντος"},
		{"αραχνούφαντος", "α-ρα-χνο-ύ-φα-ντος"},
		{"aidiniwtis", "a-i-di-niw-tis"},
	}

	for _, test := range tests {
		h.Input = test.input

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			panic(err)
		}

		if hyphenedText != test.hyphenated {
			t.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", test.input, test.hyphenated, hyphenedText)
		}
	}
}

func TestQuickSynizesisRules(t *testing.T) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "-"
	hyphenationOptions.QuickSynizesis = true
	hyphenationOptions.UseGrhyphRules = true

	h := Hyphenation{
		Options: hyphenationOptions,
	}

	tests := []hyphenationTest{
		{"αλκμιόνη", "αλκ-μιό-νη"},
		{"aidiniwtis", "a-i-di-niw-tis"},
	}

	for _, test := range tests {
		h.Input = test.input

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			panic(err)
		}

		if hyphenedText != test.hyphenated {
			t.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", test.input, test.hyphenated, hyphenedText)
		}
	}
}

func TestConsonantCombinations(t *testing.T) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "-"
	hyphenationOptions.CombineConsonantsDn = true
	hyphenationOptions.CombineConsonantsKv = true
	hyphenationOptions.CombineConsonantsPf = true
	hyphenationOptions.CombineConsonantsFk = true

	h := Hyphenation{
		Input:   "adnakvapfafka",
		Options: hyphenationOptions,
	}

	hyphenedText, err := h.Hyphenate()
	if err != nil {
		panic(err)
	}

	if hyphenedText != "a-dna-kva-pfa-fka" {
		t.Error("Incorrect hyphenation when every combination option is set to true.")
	}

	h.Options.CombineConsonantsDn = false
	h.Options.CombineConsonantsKv = false
	h.Options.CombineConsonantsPf = false
	h.Options.CombineConsonantsFk = false

	hyphenedText, _ = h.Hyphenate()

	if hyphenedText != "ad-nak-vap-faf-ka" {
		t.Error("Incorrect hyphenation when every combination option is set to false.")
	}
}

func BenchmarkPlain(b *testing.B) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "-"

	h := Hyphenation{
		Options: hyphenationOptions,
	}

	for i := 0; i < b.N; i++ {
		h.Input = "παλιόπαλιοpalio"

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			panic(err)
		}

		if hyphenedText != "πα-λι-ό-πα-λι-ο-pa-li-o" {
			b.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", "παλιόπαλιοpalio", "πα-λι-ό-πα-λι-ο-pa-li-o", hyphenedText)
		}
	}
}

func BenchmarkQuickSynizesis(b *testing.B) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "-"
	hyphenationOptions.QuickSynizesis = true

	h := Hyphenation{
		Options: hyphenationOptions,
	}

	for i := 0; i < b.N; i++ {
		h.Input = "παλιόπαλιοpalio"

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			panic(err)
		}

		if hyphenedText != "πα-λιό-πα-λιο-pa-lio" {
			b.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", "παλιόπαλιοpalio", "πα-λιό-πα-λιο-pa-lio", hyphenedText)
		}
	}
}

func BenchmarkRules(b *testing.B) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "-"
	hyphenationOptions.UseGrhyphRules = true

	h := Hyphenation{
		Options: hyphenationOptions,
	}

	CachingEnabled = true

	for i := 0; i < b.N; i++ {
		h.Input = "παλιόπαλιοpalio"

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			panic(err)
		}

		if hyphenedText != "πα-λιό-πα-λιο-pa-lio" {
			b.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", "παλιόπαλιοpalio", "πα-λιό-πα-λιο-pa-lio", hyphenedText)
		}
	}
}

func BenchmarkRulesNoCache(b *testing.B) {
	hyphenationOptions := GetDefaultOptions()

	hyphenationOptions.Seperator = "-"
	hyphenationOptions.UseGrhyphRules = true

	h := Hyphenation{
		Options: hyphenationOptions,
	}

	CachingEnabled = false

	for i := 0; i < b.N; i++ {
		h.Input = "παλιόπαλιοpalio"

		hyphenedText, err := h.Hyphenate()
		if err != nil {
			panic(err)
		}

		if hyphenedText != "πα-λιό-πα-λιο-pa-lio" {
			b.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", "παλιόπαλιοpalio", "πα-λιό-πα-λιο-pa-lio", hyphenedText)
		}
	}
}
