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
		{"μυστηριώδης", "μυ-στη-ρι-ώ-δης"},
		{"mustiriwdis", "mu-sti-ri-w-dis"},
		{"musthriwdhs", "musthri-w-dhs"},
		{"στήριγμα", "στή-ριγ-μα"},
		{"stirigma", "sti-rig-ma"},
		{"sthrigma", "sthrig-ma"},
		{"παρίσθμιος", "πα-ρί-σθμι-ος"},
		{"paris8mios", "pa-ri-s8mi-os"},
		{"paristhmios", "pa-risthmi-os"},
		{"άσθμα", "ά-σθμα"},
		{"as8ma", "a-s8ma"},
		{"asthma", "asthma"},
		{"διάστημα", "δι-ά-στη-μα"},
		{"diastima", "di-a-sti-ma"},
		{"diasthma", "di-asthma"},
		{"διαισθητικός", "δι-αι-σθη-τι-κός"},
		{"diais8htikos", "di-ai-s8h-ti-kos"},
		{"επιπρόσθετα", "ε-πι-πρό-σθε-τα"},
		{"epiprostheta", "e-pi-prosthe-ta"},
		{"ευπρόσδεκτος", "ευ-πρόσ-δε-κτος"},
		{"euprosdektos", "eu-pros-de-ktos"},
		{"euprosthektos", "eu-prosthe-ktos"},
		{"τρισδιάστατος", "τρισ-δι-ά-στα-τος"},
		{"trisdiastatos", "tris-di-a-sta-tos"},
		{"tristhiastatos", "tristhi-a-sta-tos"},
		{"αντηλιακό", "α-ντη-λι-α-κό"},
		{"antiliako", "a-nti-li-a-ko"},
		{"anthliako", "anthli-a-ko"},
		{"συνθλίβω", "συν-θλί-βω"},
		{"sun8livw", "sun-8li-vw"},
		{"sunthlivw", "sunthli-vw"},
		{"αντήχηση", "α-ντή-χη-ση"},
		{"anthhisi", "anthhi-si"},
		{"anthhhsh", "anthhhsh"},
		{"άνθρωπος", "άν-θρω-πος"},
		{"an8rwpos", "an-8rw-pos"},
		{"anthrwpos", "anthrw-pos"},
		{"ενθαρρύνω", "εν-θαρ-ρύ-νω"},
		{"en8arrunw", "en-8ar-ru-nw"},
		{"entharrunw", "enthar-ru-nw"},
		{"πανδαιμόνιο", "παν-δαι-μό-νι-ο"},
		{"pandaimonio", "pan-dai-mo-ni-o"},
		{"panthaimonio", "panthai-mo-ni-o"},
		{"ταυτότητα", "ταυ-τό-τη-τα"},
		{"tautotita", "tau-to-ti-ta"},
		{"tautothta", "tau-toth-ta"},
		{"συντηρητικός", "συ-ντη-ρη-τι-κός"},
		{"sintiritikos", "si-nti-ri-ti-kos"},
		{"sunthritikos", "sunthri-ti-kos"},
		{"αλκμιόνη", "αλκ-μι-ό-νη"},
		{"aλkμiοnη", "aλk-μi-ο-nη"},
		{"φαστφουντάδικο", "φαστ-φου-ντά-δι-κο"},
		{"grafete Ellhnika", "gra-fe-te El-lh-ni-ka"},
		{`δυο
      γραμμές.`, `δυ-ο
      γραμ-μές.`},
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
		// {"drasthriothta", "dra-sth-ri-othta"},
		// {"lupanthhrio", "lu-pa-nth-ri-o"},
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
		// {"drasthriothta", "dra-sth-ri-othta"},
		// {"lupanthhrio", "lu-pa-nth-rio"},
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
			b.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", "παλιόπαλιοpalio",
				"πα-λι-ό-πα-λι-ο-pa-li-o", hyphenedText)
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
			b.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", "παλιόπαλιοpalio", "πα-λιό-πα-λιο-pa-lio",
				hyphenedText)
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
			b.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", "παλιόπαλιοpalio", "πα-λιό-πα-λιο-pa-lio",
				hyphenedText)
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
			b.Errorf("(%s) Hyphenated value does not match: expected %s, got %s", "παλιόπαλιοpalio", "πα-λιό-πα-λιο-pa-lio",
				hyphenedText)
		}
	}
}
