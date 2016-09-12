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
		{"άκαμπτος", "ά-κα-μπτος"},
		{"άλμπατρος", "άλ-μπα-τρος"},
		{"έκθλιψη", "έκ-θλι-ψη"},
		{"έκπληκτος", "έκ-πλη-κτος"},
		{"έμπνευση", "έ-μπνευ-ση"},
		{"ένσφαιρος", "έν-σφαι-ρος"},
		{"ίντσα", "ί-ντσα"},
		{"αεροελεγκτής", "α-ε-ρο-ε-λε-γκτής"},
		{"αισχρολόγος", "αι-σχρο-λό-γος"},
		{"αλτρουισμός", "αλ-τρου-ι-σμός"},
		{"αμφιβληστροειδής", "αμ-φι-βλη-στρο-ει-δής"},
		{"ανεξάντλητος", "α-νε-ξά-ντλη-τος"},
		{"ανυπέρβλητος", "α-νυ-πέρ-βλη-τος"},
		{"απαρτχάιντ", "α-παρτ-χά-ιντ"},
		{"αργκό", "αρ-γκό"},
		{"αρθρογραφία", "αρ-θρο-γρα-φί-α"},
		{"ασύγγνωστος", "α-σύγ-γνω-στος"},
		{"βολφράμιο", "βολ-φρά-μι-ο"},
		{"βούρτσα", "βούρ-τσα"},
		{"γκολτζής", "γκολ-τζής"},
		{"γλεντζές", "γλε-ντζές"},
		{"Δεκέμβριος", "Δε-κέμ-βρι-ος"},
		{"διόπτρα", "δι-ό-πτρα"},
		{"εγγλέζικος", "εγ-γλέ-ζι-κος"},
		{"εισπλέω", "ει-σπλέ-ω"},
		{"εισπνοή", "ει-σπνο-ή"},
		{"εισπράκτορας", "ει-σπρά-κτο-ρας"},
		{"εκδρομέας", "εκ-δρο-μέ-ας"},
		{"εκδρομή", "εκ-δρο-μή"},
		{"εκθρόνιση", "εκ-θρό-νι-ση"},
		{"εκκρεμότητα", "εκ-κρε-μό-τη-τα"},
		{"εκπνοή", "εκ-πνο-ή"},
		{"εκπρόσωπος", "εκ-πρό-σω-πος"},
		{"εκπτωτικός", "εκ-πτω-τι-κός"},
		{"εκστομίζω", "εκ-στο-μί-ζω"},
		{"εκσφενδονισμός", "εκ-σφεν-δο-νι-σμός"},
		{"εκτρέφω", "ε-κτρέ-φω"},
		{"εκφραστικός", "εκ-φρα-στι-κός"},
		{"ελκτικός", "ελ-κτι-κός"},
		{"εμβληματικός", "εμ-βλη-μα-τι-κός"},
		{"ενθρόνιση", "εν-θρό-νι-ση"},
		{"ευστροφία", "ευ-στρο-φί-α"},
		{"εχθροπραξία", "ε-χθρο-πρα-ξί-α"},
		{"ινκόγκνιτο", "ιν-κό-γκνι-το"},
		{"ινστιτούτο", "ιν-στι-τού-το"},
		{"ισχνότητα", "ι-σχνό-τη-τα"},
		{"καλντερίμι", "καλ-ντε-ρί-μι"},
		{"καμτσίκι", "καμ-τσί-κι"},
		{"καρτποστάλ", "καρτ-πο-στάλ"},
		{"κομπλιμέντο", "κο-μπλι-μέ-ντο"},
		{"κύλινδρος", "κύ-λιν-δρος"},
		{"μάνατζμεντ", "μά-να-τζμεντ"},
		{"μαρσπιέ", "μαρ-σπι-έ"}, // todo: Exception rule (synizesis).
		{"μετεγγραφή", "με-τεγ-γρα-φή"},
		{"μπέιζμπολ", "μπέ-ιζ-μπολ"},
		{"μπασκετμπολίστας", "μπα-σκετ-μπο-λί-στας"},
		{"μπαχτσές", "μπα-χτσές"},
		{"νομενκλατούρα", "νο-μεν-κλα-τού-ρα"},
		{"νταρντάνα", "νταρ-ντά-να"},
		{"ντόμπρος", "ντό-μπρος"},
		{"πάμφθηνα", "πάμ-φθη-να"},
		{"πανσπερμία", "παν-σπερ-μί-α"},
		{"παρεκκλήσι", "πα-ρεκ-κλή-σι"},
		{"πορθμός", "πορθ-μός"},
		{"πορτμαντό", "πορ-τμα-ντό"}, // todo: Exception rules for 'πορτ-μα-ντό' (etymological hyphenation) and related.
		{"πορτμπαγκάζ", "πορτ-μπα-γκάζ"},
		{"προσβλέπω", "προ-σβλέ-πω"},
		{"πρόσκληση", "πρό-σκλη-ση"},
		{"πρόσκρουση", "πρό-σκρου-ση"},
		{"πρόσκτηση", "πρό-σκτη-ση"},
		{"πρόσπτωση", "πρό-σπτω-ση"},
		{"ράφτρα", "ρά-φτρα"},
		{"ροσμπίφ", "ρο-σμπίφ"},
		{"σάλτσα", "σάλ-τσα"},
		{"σεντράρισμα", "σε-ντρά-ρι-σμα"},
		{"στιλπνός", "στιλ-πνός"},
		{"συγκλονιστικός", "συ-γκλο-νι-στι-κός"},
		{"σφυρίχτρα", "σφυ-ρί-χτρα"},
		{"σύμπτωση", "σύ-μπτω-ση"},
		{"σύντμηση", "σύ-ντμη-ση"},
		{"τερπνότητα", "τερ-πνό-τη-τα"},
		{"τζαμτζής", "τζαμ-τζής"},
		{"Τουρκμενιστάν", "Τουρκ-με-νι-στάν"},
		{"τουρμπίνα", "τουρ-μπί-να"},
		{"τροτσκισμός", "τρο-τσκι-σμός"},
		{"τσουγκράνα", "τσου-γκρά-να"},
		{"υπαρκτός", "υ-παρ-κτός"},
		{"υπερδραστήριος", "υ-περ-δρα-στή-ρι-ος"},
		{"υπερκράτος", "υ-περ-κρά-τος"},
		{"υπερπλήρης", "υ-περ-πλή-ρης"},
		{"υπερπροσπάθεια", "υ-περ-προ-σπά-θει-α"},
		{"υπερσκελίζω", "υ-περ-σκε-λί-ζω"},
		{"υπερσταθμός", "υ-περ-σταθ-μός"},
		{"υπερσύγχρονος", "υ-περ-σύγ-χρο-νος"},
		{"υπερτραφής", "υ-περ-τρα-φής"},
		{"υπερχρονίζω", "υ-περ-χρο-νί-ζω"},
		{"φλαμίνγκο", "φλα-μίν-γκο"},
		{"φολκλορισμός", "φολ-κλο-ρι-σμός"},
		{"χάντμπολ", "χά-ντμπολ"},
		{"χαρτζιλίκι", "χαρ-τζι-λί-κι"},
		{"μυστηριώδης", "μυ-στη-ρι-ώ-δης"},
		{"mustiriwdis", "mu-sti-ri-w-dis"},
		{"musthriwdhs", "musthri-wdhs"},
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
		{"diaisthitikos", "di-aisthi-ti-kos"},
		{"diais8htikos", "di-ais8hti-kos"},
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
		{"Αλκμήνη", "Αλκ-μή-νη"},
		{"aλkμiοnη", "aλk-μi-ο-nη"},
		{"φαστφουντάδικο", "φαστ-φου-ντά-δι-κο"},
		{"grafete Ellhnika", "gra-fe-te Ellhni-ka"},
		{`δυο
      γραμμές.`, `δυ-ο
      γραμ-μές.`},
		{"Μια και δυο. Μία και δύο", "Μι-α και δυ-ο. Μί-α και δύ-ο"},
		{"χελιδόνια", "χε-λι-δό-νι-α"},
		{"αηδόνια", "α-η-δό-νι-α"},
		{"English words are hyphenated using Greek/Greeklish grammar. Their detection and exclusion has to happen outside of this package.", "En-glish w-ords a-re hyphe-na-ted u-sing Gre-ek/Gre-e-klish gram-mar. Their de-tec-ti-on and exc-lu-si-on has to hap-pen ou-tsi-de of this pac-ka-ge."},
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

// Use rules without caching.
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
		{"άνθια", "άν-θια"},
		{"anthia", "an-thia"},
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
