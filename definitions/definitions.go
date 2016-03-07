package grhyph

import (
	"fmt"
	"regexp"
	"strings"
)

// Match groups: ["" "vowels" "consonants" "punctuation" "other"]
var SpeechSoundRe = "(?i)(?P<vowels>[ϊϋΐΰ]|[αa][ύυuy]|[εe][ύυuy]|[ηi][ύυuy]|[αa][ίιi]|[εe][ίιi]|[οo][ύυuy]|[οo][ίιi]|[άαa]|[έεe]|[ήηhi]|[ίιi]|[όοo]|[ύυyu]|[ώωwo])|(?P<consonants>(?:[μm][πp]|b)|(?:[γg][κk]|[γg])|[νn][τtj]|[νn]|(?:[τtj]h|[θ8])|(?:[δd])|(?:[τtj][ζz]|j)|[ζz]|[τtj][σs]|[σs][τtj]|[βv]|[λl]|[μm]|(?:ks|κs|kσ|[ξx3])|[ρr]|[τt]|[φf]|[χxh]|ch|(?:[pπ][σs]|[ψ4])|[πp]|[σsc]|[κk])|(?P<punctuation>[\\s\\.,\\-–—―\\/'’\":!?;&@«»])|(?P<other>.?)"

var NuTauRe = "(?i)^([νn])([τt])$"

// Valid Greek word starting consontants.
// Important: Verify the getWSCRe()'s conditions when altering.
var WordStartConsonantsRe = "(?i)^([βvb](?:[τt]h|[δdγgλlρr])|[γg](?:[τt]h|[δdκkλlνnρr])|(?:[τt]h|[δd])[νn]|(?:[τt]h|[δd])[ρr]|(?:[τt]h|[θ8])[λlνnρr]|[κk][βvb]|[κk][λlνnρrτtj]|[μm][νnπp]|[νn][τtj][^h]|[πp][λlνnρrτtj]|[πp][φf]|[σs](?:[τt][^h]|[θ8βvbγgκkλlμmνnπpφfχxh])|[τt][ζzμmρrσs]|[φf](?:[τt]h?|[θ8λlρrχxh]|ch)|[φf][κk]|(?:[χxh]|ch)(?:[θ8λlνnρr]|[τt]h?))"

var WSCReMap = map[string]*regexp.Regexp{}

// Get a WordStartConsonantsRe based on the combination options.
func GetWSCRe(combDn, combKv, combPf, combFk bool) *regexp.Regexp {
	// todo: Use an interface{}, instead of the following custom string, as a key for the WSCReMap.
	customKey := fmt.Sprintf("%s%s%s%s",
		fmt.Sprint(combDn)[0:1],
		fmt.Sprint(combKv)[0:1],
		fmt.Sprint(combPf)[0:1],
		fmt.Sprint(combFk)[0:1])

	if _, ok := WSCReMap[customKey]; ok {
		return WSCReMap[customKey]
	} else {
		wSCRe := WordStartConsonantsRe

		if !combDn {
			wSCRe = strings.Replace(wSCRe, "|(?:[τt]h|[δd])[νn]", "", 1)
		}
		if !combKv {
			wSCRe = strings.Replace(wSCRe, "|[κk][βvb]", "", 1)
		}
		if !combPf {
			wSCRe = strings.Replace(wSCRe, "|[πp][φf]", "", 1)
		}
		if !combFk {
			wSCRe = strings.Replace(wSCRe, "|[φf][κk]", "", 1)
		}

		WSCReMap[customKey] = regexp.MustCompile(wSCRe)
		return WSCReMap[customKey]
	}
}

// Vowel combinations prone to synizesis.
var SynizesisVowelsRe = "(?i)^([αάa][ηhιϊi]|[εe][ϊ]|[εe][ιi](?:[αάaοόoωώw]|[οo][υύuy])|[ιi](?:[αάaεέeοόoωώw]|[αa][ιίi]|[οo][ιίiυύuy])|[οόo](?:[ιiϊ]|[εe][ιi])|[οo][ιi](?:[αάaεέeοόoωώw]|[οo][ιίiυύuy])|[υuy][αάaιiοόoωώw])$"

var customRegexpsMap = map[string]string{ // todo: Test map records.
	"(.*)":                          "(.*)",
	"(.+)":                          "(.+)",
	"^":                             "^",
	"$":                             "$",
	"(α)":                           "([αa])",
	"(ά)":                           "(ά)",
	"(α|ο)":                         "([αaοo])",
	"(α|ε|ο|ω)":                     "([αaεeοoωw])",
	"(α|ά|ω|ώ)":                     "([αάaοόoωώw])",
	"(α|ά|ε|έ|ω|ώ)":                 "([αάaεέeωώwοόo])",
	"(α|ε|ω|ώ)":                     "([αaεeωώwοόo])",
	"(α|ε|ω)":                       "([αaεeωwοo])",
	"(ά|έ|ώ)":                       "([άέώό])",
	"(α|ε)":                         "([αaεe])",
	"(α|ά)":                         "([αάa])",
	"(α|ά|ας|άς)":                   "([αάa][σs]?)",
	"(α|ά|ας|άς|ες|ές|'ων'|'ών')":   "([αάa][σs]?|[εέe][σs]|[ωώw][νn])",
	"(α|ας|ες|ων)":                  "([αa][σs]?|[εe][σs]|[ωwοo][νn])",
	"(α|ας|ες|ων 2)":                "([αάa][σs]?|[εέe][σs]|[ωώwοόo][νn])",
	"(α|ες|οι|ους)":                 "([αa]|[εe][σs]|[οo](?:[ιi]|[υuy][σs]))",
	"('ά'|'άς'|'ές'|'ών')":          "([ά][σs]?|[έ][σs]|[ώό][νn])",
	"(α|κ)":                         "([αaκk])",
	"(α|ας|ω|ώ|ως|ώς)":              "([αa][σs]?|[ωώwοόo][σs])",
	"(α|ας)":                        "([αa][σs]?)",
	"(α|ας|ο|ου|'ως')":              "([αa][σs]?|[οo][υuy]?|[ωw][σs])",
	"(α|ο|ου|ων)":                   "([αa]|[οo][υuy]?|[ωwοo])",
	"(α|ω)":                         "([αaοoωw])",
	"(ά|ώ)":                         "([άόώ])",
	"(α|ων)":                        "([αa]|[ωwοo][νn])",
	"(α|ων|ών)":                     "([αa]|[ωώwοόo][νn])",
	"(αι|αί|α|ά|ε|έ)":               "([αa][ιίi]|[αάa]|[εέe])",
	"(Α|Ά)":                         "((?-i:[ΑΆA]))",
	"(β)":                           "([βvb])",
	"(ββ|β)":                        "([βbv][βbv]?)",
	"(Β)":                           "((?-i:[ΒVB]))",
	"(β|μπ)":                        "([βvb]|[μm][πp])",
	"(γ)":                           "([γg])",
	"(!γ)":                          "([^γg])",
	"(γ|στ|χ)":                      "([γg]|[σs][τt]|[χxh]|ch)",
	"(γ|σ)":                         "([γgσs])",
	"(γγ)":                          "([γg][γg])",
	"(γκ)":                          "([γg][κk])",
	"(γγ|γκ)":                       "([γg][γgκk])",
	"(δ)":                           "([τt]h|[δd])",
	"(Δ)":                           "((?-i:[ΤT][Hh]|[ΔD]))",
	"(δ|ντ)":                        "([τt]h|[δd]|[νn][τt])",
	"(ε)":                           "([εe])",
	"(έ)":                           "(έ)",
	"(αι|ε)":                        "([αa][ιi]|[εe])",
	"(αι|αί|ε|έ)":                   "([αa][ιίi]|[εέe])",
	"(αί|έ)":                        "([αa][ί]|[έ])",
	"(ε|έ)":                         "([εέe])",
	"(ε|έ|ο|ό)":                     "([εέeοόo])",
	"(ε|ω)":                         "([εeωwοo])",
	"(ε|έ|η|ή|ι|ί)":                 "([εέeηήhιίi])",
	"(ε|η|ι)":                       "([εeηhιi])",
	"(α|αν|ε|ες)":                   "([αa][νn]?|[εe][σs]?)",
	"(ει)":                          "([εe][ιi])",
	"(ει|ι)":                        "([εe][ιi]|[ιi])",
	"(ε|έ|ι|ί)":                     "([εέeιίi])",
	"(ει|εί|ι|ί)":                   "([εe][ιίi]|[ιίi])",
	"(ευ|β)":                        "([εe][υuy]|[βvb])",
	"(ζ)":                           "([ζz])",
	"('Ζ')":                         "((?-i:[ΖZ]))",
	"(ζ|σμ|στ)":                     "([ζz]|[σs][μm]|[σs][τt])",
	"(ζ|σ|σμ|στ)":                   "([ζzσs]|[σs][μm]|[σs][τt])",
	"(ζ|σ)":                         "([ζzσs])",
	"(τζ|τσ)":                       "([τt][ζzσs]|j)",
	"(ζ|κ|σμ|στ)":                   "([ζz]|[κk]|[σs][μmτt])",
	"(ζ|ρ|σμ|στ)":                   "([ζz]|[ρr]|[σs][μmτt])",
	"('η')":                         "([ηh])",
	"(η)":                           "([ηhιi])",
	"(η|ή)":                         "([ηήhιίi])",
	"(ης...)":                       "([εe]?[ιi][σs]|[εe][ωwοo][νnσs]|[ηhιi][σs]?)",
	"(ης|η)":                        "([ηhιi][σs]?)",
	"(ήπι-α)":                       "([αa](?:[μmτt][εe]|[νn][εe]?)?|[εe][σs]?)",
	"(ισ?)":                         "([ιi][σs]?)",
	"(ίσ?)":                         "([ιίi][σs]?)",
	"(θ)":                           "([τt]h|[θ8])",
	"(Ι)":                           "((?-i:[ΙI]))",
	"(Ι|Ί)":                         "((?-i:[ΙΊI]))",
	"(ι)":                           "([ιi])",
	"(ι|ί)":                         "([ιίi])",
	"(ί)":                           "(ί)",
	"(ι|ί|ϊ|ΐ)":                     "([ιίϊΐi])",
	"(ι|ϊ)":                         "([ιϊi])",
	"(ϊ)":                           "(ϊ)",
	"(ίς|ί)":                        "([ιίi][σs]?)",
	"(ια|ιού|ιών)":                  "([αa]|[οo][υύuy]|[ωώwοόo][νn])",
	"(ιά|ιού|ιών)":                  "([αάa]|[οo][υύuy]|[ωώwοόo][νn])",
	"(ος|ού|ό|έ...)":                "([οo](?:[υύuy][σs]?|[ιίi])|[οόo][σs]?|[εέe]|[ωώwοόo][νn])",
	"(ολ-ος)":                       "([εe]|[οo](?:[ιiσs]|[υuy][σs]?)|[ωwοo][νn])",
	"(στ|τ)":                        "([σs][τt]|[τt])",
	"(κ)":                           "([κk])",
	"(κκ|κ)":                        "([κk][κk]?)",
	"(κ|σμ|στ)":                     "([κk]|[σs][μmτt])",
	"(κος)":                         "([κk](?:[εe]|[οo](?:[ιi]|[σs]|[υuy][σs]?|[νn])?|[ωw][νn]))",
	"(λ)":                           "([λl])",
	"(μ)":                           "([μm])",
	"(μμ|μ)":                        "([μm][μm]?)",
	"(μ|ν|σ|τ)":                     "([μmνnσsτt])",
	"(λλ|λ)":                        "([λl][λl]?)",
	"(μπ)":                          "([μm][πp]|b)",
	"(ν)":                           "([νn])",
	"(λ|ρ)":                         "([λlρr])",
	"(λ|ν)":                         "([λlνn])",
	"(νν|ν)":                        "([νn][νn]?)",
	"(νν)":                          "([νn][νn])",
	"(ντ)":                          "([νn][τt]|[νn]?d)",
	"(ντ|σμ)":                       "([νn][τt]|d|[σs][μm])",
	"(ξ)":                           "(k[σs]|κs|[ξx3])", // todo: Explain why no 'κσ' ("εξ-τριμ cases").
	"(ο)":                           "([οo])",
	"(ό)":                           "(ό)",
	"(ο|ό)":                         "([οόo])",
	"(ό|ή)":                         "([όή])",
	"(ο|η)":                         "([οoηhιi])",
	"(ο|ό|ε|έ|ω|ώ)":                 "([οόoεέeωώw])",
	"(οι|ι)":                        "([οo][ιi]|[ιi])",
	"(ο|ου)":                        "([oο][υuy]?)",
	"(οι|οί|ι|ί)":                   "([οo][ιίi]|[ιίi])",
	"(ου)":                          "([οo][υuy])",
	"(ου|ού)":                       "([οo][υύuy])",
	"(ου|υ|ι)":                      "([οo]?[υuy]|[ιi])",
	"(π)":                           "([πp])",
	"(ππ|π)":                        "([πp][πp]?)",
	"(π|φ)":                         "([πpφf])",
	"(ρρ|ρ)":                        "([ρr][ρr]?)",
	"(ρ)":                           "([ρr])",
	"(σ)":                           "([σs])",
	"(σ|τ)":                         "([σsτt])",
	"(σσ|σ)":                        "([σs][σs]?)",
	"(σμ)":                          "([σs][μm])",
	"(ς?)":                          "([σs]?)",
	"(στ)":                          "([σs][τt])",
	"(σμ|μ|ζ|τ|σ|ρ)":                "([τt]|[ρr]|[μm]|[ζz]|[σs](?:[μm]|[τt])?)",
	"(ειάζω)":                       "([ζz]|[νn]|[σs](?:[τtμm]|[αaεeοoωw]))",
	"(τ)":                           "([τt])",
	"(τ|τσ)":                        "([τt][σs]?)",
	"(σ|τσ)":                        "([τt]?[σs])",
	"(ν|τ)":                         "([νnτt])",
	"(τσ)":                          "([τt][σs])",
	"(τζ)":                          "([τt][ζz]|j)",
	"(τ|σ)":                         "([τt]|[σs])",
	"(υ)":                           "([υuy])",
	"(υ|ύ)":                         "([υύuy])",
	"(υ|β)":                         "([υuyβvb])",
	"(υ|ύ|β)":                       "([υύuyβvb])",
	"(υ|φ)":                         "([υuyφf])",
	"(υ|ύ|φ)":                       "([υύuyφf])",
	"(υ|ι)":                         "([υuyιi])",
	"(υ|ύ|ι|ί)":                     "([υύuyιίi])",
	"(φ)":                           "([φf])",     // "ph" would conflict with "πη".
	"(χ)":                           "([χxh]|ch)", // "ch" doesn't confict with "ση".
	"(ψ)":                           "([pπ][σs]|[ψ4])",
	"(ω|ώ|ο|ό)":                     "([ώωwοόo])",
	"(ω)":                           "([ωwοo])",
	"('ω')":                         "([ωw])",
	"('ω|ώ')":                       "([ωώw])",
	"(ώ)":                           "(ώ)",
	"(ω|ώ)":                         "([ωώwοόo])",
	"(ος|α|ο)":                      "([αa][σs]?|[εe][σs]?|[οo](?:[ιiσs]|[υuy][σs]?)?|[ωwοo][νn])",
	"(ος|α|ο 2)":                    "([αάa][σs]?|[εέe][σs]?|[οόo][σs]?|[οo](?:[ιίi]|[υύuy][σs]?)|[ωώwοόo][νn])",
	"(ος|η|ο)":                      "([αa]|[εe][σs]?|[ηhιi][σs]?|[οo](?:[σs]|[ιi]|[υuy][σs]?|[νn])?|[ωw][νn])",
	"(ολ-ή)":                        "([εέe][σs]|[ηήh][σs]?|[ώό][νn])",
	"(ολ-έας)":                      "([εέe](?:[αa][σs]?|[ωwοo][νn])|(?:[εe][ιίi]|[ιίi])[σs])",
	"(ος...)":                       "([εe]|[οo](?:[ιi]|[σs]|[υuy][σs]?|[νn])?|[ωw][νn])",
	"(ος|ός|ους|ούς)":               "([οόo][υύuy]?[σs])",
	"(ι-ός)":                        "([εέe]|[οo](?:[ιίi]|[υύuy][σs]?)|[οόo][σs]?|[ωώwοόo][νn])",
	"(ι-άς...)":                     "((?:[αάa][σs]?)|(?:[εέe][σs])|(?:[ωώwοόo][νn]))",
	"(α|ά|ο|ό)":                     "([αάa]|[οόo])",
	"(α|ά|ο|ό|ω|ώ)":                 "([αάa]|[οόo]|[ωώw])",
	"(α|ά|ας|άς|ές|ών)":             "([αάa][σs]?|[εέe][σs]|[ωώwοόo][νn])",
	"(α|ά|ου|ού|ων|ών)":             "([αάa]|[οo][υύuy]|[ωώwοόo][νn])",
	"(α|ά|ου|ού)":                   "([αάa]|[οo][υύuy])",
	"(α|ου)":                        "([αa]|[οo][υuy])",
	"(α|ου|ού|ων|ών)":               "([αa]|[οo][υύuy]|[ωώwοόo][νn])",
	"(ε|έ|η|ή|ο|ό|ω|ώ)":             "([εέe]|[ηήhi]|[οόo]|[ωώw])",
	"(ι|ια|ου|ού|ων|ών)":            "([ιi][αa]?|[οo][υύuy]|[ωώw][νn])",
	"(αλ-ώ...)":                     "([εe]?[ιίi](?:[τt][εe]|[σs]?)|[οo][υύuy](?:[μm][εe]|[νn][εe]?|[σs][αa](?:[μmτt][εe]|[νn][εe]?)|[σs][εe][σs]?)|[ωώwοόo](?:[νn][τt][αa][σs])?)", // διαλώ
	"(ας|ες|ων|ών|α)":               "((?:[αa][σs])|(?:[εe][σs])|(?:[ωώw][νn])|[αa])",
	"(ση|σης|σεων|σεως)":            "([σs](?:[ηhιi][σs]?|[εe][ωwοo][νnσs]))",
	"(άτ-ης...)":                    "([εe][σs]|[ηhιi][σs]?|[ιi][σs]{1,2}(?:[αa][σs]?|[εe][σs]?|[ωώwοόo][νn])|[ωώwοόo][νn])",
	"(τ-ός, ή, ό)":                  "([εέe][σs]?|[οo](?:[ιίi]|[υύuy][σs]?)|[οόo][σs]?|[ωώwοόo][νn]|[αάa]|[ηήhιίi][σs]|)",
	"(ους|ας|ες|οι|ος|ου|ων|α|ε|ο)": "((?:[οo][υuy]?[σs]?)|(?:[αa][σs]?)|(?:[εe][σs]?)|(?:[οo](?:[ιi]|[σs]|[υuy])?)|(?:[ωw][νn]))",
	"(αίνω...)":                     "([ηήhιίi][κk][αa](?:[μm][εe]|[νn][εe]?|[τt][εe])|[ηhιi][κk](?:[αa][νn]?|[εe][σs]?)|(?:[αa][ιίi]|[εέe])[νn](?:[αa](?:[μm][εe]|[νn][εe]?|[τt][εe])|(?:[εe][ιi]|[ιi])[σs]?|[εe][τt][εe]|[οo](?:[μm][εe]|[νn][τt][αa][σs]|[υuy](?:[μm][εe]|[νn][εe]?))?|[ωw])|(?:[αa][ιi]|[εe])[νn](?:[αa][νn]?|[εe][σs]?)|(?:[εe][ιίi]|[ιίi])(?:[σs]|[τt][εe])?|[οo][υύuy](?:[μm][εe]|[νn][εe]?)|[ωώwοόo])",
	"(έμαι)":                        "([εέe](?:[μm](?:[αa][ιi]|[εe])|[σs](?:[αa][ιi]|[εe]|[τt][εe])|[τt](?:[αa][ιi]|[εe]))|[οo][υύuy][νn][τt](?:[εe]|[αa][ιiνn])|[οόo](?:[μm](?:[αa][σs][τt](?:[αa][νn]|[εe])|[οo][υuy][νn][αa]?)|[νn][τt](?:[αa](?:[ιi]|[νn][εe]?)|[οo][υuy][σs][αa][νn])|[σs](?:[αa][σs][τt](?:[αa][νn]|[εe])|[οo][υuy][νn][αa]?)|[τt][αa][νn][εe]?))",
	"(-έγω)":                        "([εέe](?:[γg](?:[αa](?:[μmτt][εe]|[νn][εe]?)|(?:[εe][ιi]|[ιi])[σs]?|[εe](?:[σsτt](?:[αa][ιi]|[εe])|[σs][τt][εe])|[οo](?:(?:[μm]|[νn][τt])(?:[αa][ιi]|[εe])|[νn][τt][αa][νnσs]|[υuy](?:[μm][εe]|[νn][εe]?))?|[ωw])|(?:k(?:s|σ)|κs|[ξx3])(?:[αa](?:[μmτt][εe]|[νn][εe]?)|(?:[εe][ιi]|[ιi])[σs]?|[εe][τt][εe]|[οo](?:[μm][εe]|[υuy](?:[μm][εe]|[νn][εe]?)?)?|[τt][εe]|[ωw])|(?:[χxh]|ch)[τt][ηhιi][κk](?:[αa][νn]?|[εe][σs]?))|[εe](?:[γg](?:[αa][νn]?|[εe][σs]?|[μm][εέe][νn].*|[οόo](?:[μm](?:[αa][σs][τt](?:[αa][νn]|[εe])|[οo][υuy][νn][αa]?)|[νn][τt](?:[αa][νn][εe]|[οo][υuy][σs][αa][νn])|[σs](?:[αa][σs][τt](?:[αa][νn]|[εe])|[οo][υuy][νn][αa]?)|[τt][αa][νn][εe]?))|(?:k(?:s|σ)|κs|[ξx3])(?:[αa][νn]|[εe][σs]?)|(?:[χxh]|ch)[τt](?:[ηήhιίi][κk](?:[αa](?:[μmτt][εe]|[νn][εe]?))|(?:[εe][ιίi]|[ιίi])(?:[σs]|[τt][εe])?|[οo][υύuy](?:[μm][εe]|[νn][εe]?)|[ωώwοόo])))",
	// "(ά-ζω)":         "([ζzσs](?:[αa](?:[μm][εe]|[νn][εe]?|[τt][εe])?|[εe](?:[ιi][σs]?|[σs]|[τt][εe]|)?|[οo](?:[μm][εe]|[υuy](?:[μm][εe]|[νn][εe]?))?|[ωw])|[σs][τt][εe]|[ζz][οo][νn][τt].*|[σs][μm][εέe][νn].*)",
	"(ά-ζω)":         "([ζzσs](?:[αa](?:[μm][εe]|[νn][εe]?|[τt][εe])?|[εέe](?:[ιi][σs]?|[σs]|[τt][εe])?|[οo](?:[μm][εe]|[υuy](?:[μm][εe]|[νn][εe]?))?|[ωw])|[ζz](?:[εe](?:[σs](?:[εe]|[αa][ιi]|[τt][εe])|[τt](?:[εe]|[αa][ιi]))|[οo](?:[μm][αa][ιi]|[νn][τt](?:[εe]|[αa][ιiνnσs]))|[οόo](?:[μmσs](?:[αa](?:[σs][τt](?:[αa][νn]|[εe]))|[οo][υuy][νn][αa]?)|[νn][τt](?:[αa][νn][εe]|[οo][υuy][σs][αa][νn])|[τt][αa][νn][εe]?))|[σs](?:[μm][εέe][νn].*|[οo][υuy]|[τt](?:[ηhήιίi][κk][αa](?:[μm][εe]|[νn][εe]?|[τt][εe])|[ηhιi][κk](?:[αa]|[εe][σs]?)|[εe](?:[ιίi](?:[σs]|[τt][εe])?)?|[οo][υύuy](?:[μm][εe]|[νn][εe]?)|[ωώwοόo])))",
	"(ώ-νω)":         "([νnσs](?:[αa](?:[μm][εe]|[νn][εe]?|[τt][εe])?|[εέe](?:[ιi][σs]?|[σs]|[τt][εe])?|[οo](?:[μm][εe]|[υuy](?:[μm][εe]|[νn][εe]?))?|[ωw])|[νn](?:[εe](?:[σs](?:[εe]|[αa][ιi]|[τt][εe])|[τt](?:[εe]|[αa][ιi]))|[οo](?:[μm][αa][ιi]|[νn][τt](?:[εe]|[αa][ιiνnσs]))|[οόo](?:[μmσs](?:[αa](?:[σs][τt](?:[αa][νn]|[εe]))|[οo][υuy][νn][αa]?)|[νn][τt](?:[αa][νn][εe]|[οo][υuy][σs][αa][νn])|[τt][αa][νn][εe]?))|[σs](?:[μm][εέe][νn].*|[οo][υuy]|[τt](?:[ηhήιίi][κk][αa](?:[μm][εe]|[νn][εe]?|[τt][εe])|[ηhιi][κk](?:[αa]|[εe][σs]?)|[εe](?:[ιίi](?:[σs]|[τt][εe])?)?|[οo][υύuy](?:[μm][εe]|[νn][εe]?)|[ωώwοόo])))",
	"(α...)":         "([αa](?:[μm][εe]|[νn][εe]?|[τt][εe])?|[εe](?:[ιi][σs]?|[σs]|[τt][εe])?|[οo](?:[μm][εe]|[υuy](?:[μm][εe]|[νn][εe]?))?|[τt][εe]|[ωw])",
	"(βιά-ζω)":       "([αάa](?:[ζz](?:[αa](?:[μmτt][εe]|[νn][εe]?)|[εe](?:[ιi][σs]?|[σs](?:[αa][ιi]|[τt][εe])|<τt]($6?:[αa][ιi]|[εe]))|[οόo](?:[μm](?:[αa](?:[ιi]|[σs][τt](?:[αa][νn]|[εe]))|[εe]|[οo][υuy][νn][αa]?)|[νn][τt](?:[αa](?:[ιiσs]|[νn][εe]?)|[οo][υuy][σs][αa][νn])|[σs](?:[αa][σs][τt](?:[αa][νn]|[εe])|[οo][υuy][νn][αa]?)|[τt][αa][νn][εe]?[υuy](?:[μm][εe]|[νn][εe]?))|[ωwοo])|[σs](?:[αa](?:[μmτt][εe]|[νn][εe]?)|[εe](?:[ιi][σs]?|[τt][εe])|[οo](?:[μm][εe]|[υuy](?:[μm][εe]|[νn][εe]?)?)|[τt](?:(?:[εe][ιίi]|[ιίi])(?:[σs]|[τt][εe])?|[εe]|[ηήhιίi](?:[κk](?:[αa](?:[μmτt][εe]|[νn][εe]?)?|[εe][σs]?))|[οo][υύuy](?:[μm][εe]|[νn][εe]?)|[ωώwοόo])|[ωwοo])))",
	"(ιώ-ξω)":        "([αa](?:[μmτt][eε]|[νn][εe]?)?|[εe](?:[ιi][σs]?|[σs]|[τt][εe])?|[οo](?:[μm][εe]|[υuy](?:[μm][εe]|[νn][εe]?)?)?|[τt][εe]|[ωw])",
	"(στος|στη|στο)": "([σs][τt](?:[αa]|[εe](?:[σs])?|[ηhιi](?:[σs])?|[οo](?:[ιi]|[σs]|[υuy](?:[σs])?|[νn])?|[ωw][νn]))",
	"(-ιάση)":        "([αάa][σs](?:[εe]?[ιi][σs]|[εe][ωwοo][νnσs]|[ηhιi][σs]?))",
	"(ους|ούς|ες|ές|ων|ών|ου|ού|α|ά|ε|έ|ο|ό)": "((?:[οo][ύυuy][σs])|(?:[εέe][σs])|(?:[ωώw][νn])|(?:[οo][ύυuy])|[αάa]|[εέe]|[οόo])",
}

type GrhyphRule struct {
	CompiledCustomRe *regexp.Regexp
	Replacement      string
}

// Transforms the custom regexps into the canonical format, and compiles the result.
func customRegexpCompile(grhyphRegexps []string) *regexp.Regexp {
	var canonRegexp []byte
	canonRegexp = append(canonRegexp, "(?i)"...)

	for _, grhyphRegexp := range grhyphRegexps {
		if mappedRe, ok := customRegexpsMap[grhyphRegexp]; ok {
			canonRegexp = append(canonRegexp, mappedRe...)
		} else {
			canonRegexp = append(canonRegexp, grhyphRegexp...)
			// debug
			fmt.Println("warning: custom regexp not included in map", grhyphRegexp) // todo: Error.
		}
	}

	return regexp.MustCompile(string(canonRegexp[:])) // todo: Error instead of panic.
}

// This list catalogs whole words, or parts of words, for which synizesis is most likely to occur.
// Rules to seperate vowels may also exist, useful for words that miss accents/diacritics or are
// written in Greeklish.
var GrhyphRules = []GrhyphRule{
	// GrhyphRule{customRegexpCompile([]string{}), ""},

	// αβάισσα, αβάισα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ββ|β)", "(α)", "(ι)", "(σσ|σ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαισ*&dq=

	// Αγλαΐα, αγλαΐζω, αγλάισμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γ)", "(λ)", "(α)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγλαί*&dq=

	// Αδελαΐδα (adela-ida)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε)", "(λ)", "(α)", "(ι|ί)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δελαιδ*&dq=

	// αδενοϋπόφυση
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(υ)", "(π)", "(ο)", "(φ)", "(.*)"}), "$1$2>-<$3$4$5$6$7"},
	// http://multilingual.sensegates.com/ΟΥΠΟΦ/string.html

	// αδιανόητα (adiano-ita)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ο)", "(ι)", "(τ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανοιτ*&dq=

	// ναυσιπλοίαρχος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(λ)", "(ο)", "(ι|ί)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4>$5-<$6$7$8"},

	// ποντοπλοΐα (ποντοπλο-ια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ντ)", "(ο)", "(π)", "(λ)", "(ο)", "(ι|ί)", "(.*)"}),
		"$1$2$3$4$5$6$7>-<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αυσιπλο*&dq=

	// ακτοπλοΐα (ακτοπλο-ια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ο)", "(π)", "(λ)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τοπλοι*&dq=

	// αξιοπλοΐα (αξιοπλο-ια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ξ)", "(ι)", "(ο)", "(π)", "(λ)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6$7>-<$8$9"},
	// *ξιοπλοι*

	// ιστιοπλοΐα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ι)", "(ο)", "(π)", "(λ)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6$7>-<$8$9"},
	// todo: όπλοια rule?

	// ναυσιπλοΐα (ναυσιπλο-ια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(υ|φ)", "(σ)", "(ι)", "(π)", "(λ)", "(ο)", "(ι|ί)", "(.*)"}),
		"$1$2$3$4$5$6$7$8>-<$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αυσιπλο*&dq=

	// αεροπλοΐα (αεροπλο-ια, όμως αεροπλοι-αρχος και διαστημόπλοι-α)
	// GrhyphRule{customRegexpCompile(
	// []string{"(.+)", "(π)", "(λ)", "(ο)", "(ί)", "(α)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πλοια*&dq=

	// αεροπλοϊκός (αεροπλο-ικος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(λ)", "(ο)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ΠΛΟΙΚ*&dq=

	//Αζερμπαϊτζάν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(α)", "(ι)", "(τ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπαϊτ*&dq=

	// αζωικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ο)", "(ι)", "(κ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζωικ*&dq=

	// άηχος (a-ixos)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ι)", "(χ)", "(ο)", "(.*)"}), "$1$2>-$3-$4<$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αιχο*&dq=

	// αθεΐα (αθε-ια)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(θ)", "(ε)", "(ι)", "(α|ά)", "(φ)", "(.*)"}), "$1-$2$3$4$5-<$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(θ)", "(ε)", "(ι|ί)", "(α)", "(.*)"}), "$1$2$3>-<$4$5$6"}, //"$1-$2$3-<$4$5$6"
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αθεια*&dq=

	// αθεΐζω, αθεϊσμός, αθεϊστής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε)", "(ι|ί)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θειζ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θεισμ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θειστ*&dq=

	// αθηναϊκός (αθηνα-ικος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(ν)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηναικ*&dq=

	// αϊκιού (α-ϊ-κιού)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ι|ί|ϊ|ΐ)", "(κ)", "(ι)", "(ο)", "(.*)"}), "$1$2>-$3-$4$5<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αικιου*&dq=

	// αϊλάινερ (α-ι-λ-α-ι-νερ)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ι|ί|ϊ|ΐ)", "(λ)", "(α|ά)", "(ι|ί|ϊ|ΐ)", "(.*)"}), "$1$2>-$3-$4$5-<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αιλάι*&dq=

	// Αϊνστάιν, αϊνσταΐνιο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ι|ϊ)", "(ν)", "(στ)", "(α)", "(ι|ί|ϊ|ΐ)", "(.*)"}), "$1$2>-$3$4-$5$6-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αϊνστα*&dq=

	// αϊράνι (α-ιρανι)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ι|ϊ)", "(ρ)", "(α)", "(ν)", "(ι)", "(α|ά|ο|ό|ω|ώ)", "(.*)"}), "$1$2>$3-$4$5-$6$7<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αιράν*&dq=

	// Αϊτή
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ι|ϊ)", "(τ)", "(η|ή)", "(ς?)", "$"}), "$1-$2-$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Αϊτήσ*&dq=

	// Αϊτινός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ι|ϊ)", "(τ)", "(ι)", "(ν)", "(ε|έ|η|ή|ο|ό|ω|ώ)", "(.*)"}), "$1$2>-<$3$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Αϊτιν*&dq=

	// Αϊ (Αϊ-βαλί), άι (άι-Γιάννης)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ά)", "(ι)", "(.*)"}), "$1<$2$3"}, //^άι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ϊ)", "(.*)"}), "$1<$2$3"}, //^αϊ
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αι*&dq=

	// αϊβαλιώτης, αϊβαλιώτικος, αϊβαλιώτισσα, ...
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι|ί|ϊ|ΐ)", "(β)", "(α)", "(λ)", "(ι)", "(α|ά|ο|ό|ω|ώ)", "(.*)"}), "$1$2-$3$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αιβα*&dq=

	// αϊδίνι (α-ιδινι)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ι)", "(δ)", "(ι|ί)", "(ν)", "(.*)"}), "$1-<$2$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αϊδίν*&dq=

	// αΐδιος (α-ιδιος)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ι|ί)", "(δ)", "(ι)", "(.*)"}), "$1-$2-<$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αιδι*&dq=

	// αιμορροϊδες (αιμορρο-ιδες)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(ο)", "(ρρ|ρ)", "(ο)", "(ι|ί)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μοροϊδ*&dq=

	// αϊνάς (α-ινας)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ι)", "(ν)", "(α|ά)", "(.*)"}), "$1-<$2$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αϊνά*&dq=

	// ακαδημαϊκός (ακαδημα-ικος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(δ)", "(η)", "(μ)", "(α)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αδημαϊ*&dq=

	// ακαΐα (ακα-ια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(α)", "(ι)", "(α)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ακαια*&dq=

	// ακατανόητα (akatano-ita)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ο)", "(ι)", "(τ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*ανόητ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανόιτ*&dq=

	// ακορόιδευτος (ακο-ρόι-δευτος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ρ)", "(ο|ό)", "(ι|ϊ)", "(δ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ορόιδ*&dq=

	// ακρολεΐνη (ακρολε-ινη)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ο)", "(λ)", "(ε)", "(ι|ί)", "(ν)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρολεΐν*&dq=

	// αλαφροΐσκιωτος (αλαφρο-ισκιωτος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ι|ί|ϊ|ΐ)", "(σ)", "(κ)", "(ι)", "(ω|ώ)", "(.*)"}), "$1$2>-$3-$4$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οΐσκιω*&dq=

	// ακτινοϋποδοχέας (ακτινο-υποδοχέας)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ο)", "(υ|ι)", "(π)", "(ο)", "(δ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νοϋποδ*&dq=

	// αλδεΰδη (αλδε-υδη)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(δ)", "(ε)", "(υ|ύ)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λδευ*&dq=

	// αλκαϊκός (αλκα-ικος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(α)", "(ι|ί|ϊ|ΐ)", "(κ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*καϊκ*&dq=

	// αλληλοϋποστήριξη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ο)", "(υ|ι)", "(π)", "(ο)", "(στ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λοϋποστ*&dq=

	// αλταϊκός (αλτα-ικος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(τ)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λταικ*&dq=

	// Αλτσχάιμερ (Αλτσ-χάι-μερ)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(τσ)", "(χ)", "(α|ά)", "(ι)", "(μ)", "(.*)"}), "$1$2$3>-$4$5$6-<$7$8"},
	// http://lexilogia.gr/forum/showthread.php?13995-Ο-συλλαβισμός-των-λέξεων&p=249383&viewfull=1#post249383

	// αμινοξεϊκός (aminokse-ikos)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ξ)", "(ε)", "(ι)", "(κ)", "(.*)"}), "$1$2>-$3$4-$5<$6$7$8"},
	// *οξεϊκ*

	// αμφισσαϊκός (αμφισσα-ικός)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(σσ|σ)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισαϊκό*&dq=

	// αναπνοή (anapno-i)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(α)", "(π)", "(ν)", "(ο)", "(ι)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ναπνοι&dq=

	// ανατολικοευρωπαϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(π)", "(α)", "(ι)", "(κ|σμ|στ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρωπαϊκ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ωπαισ*&dq=

	// ανευθυνοϋπεύθυνος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(υ|ι)", "(π)", "(ε)", "(υ|ύ|φ)", "(θ)", "(.*)"}), "$1$2>-<$3$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουπέυ*&dq=

	// ανισοϋψής (ανισο-υψης)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(σ)", "(ο)", "(υ|ι)", "(ψ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισοϋψ*&dq=

	// ανοησία (ano-isia)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ο)", "(ι)", "(σ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανοιτ*&dq=

	// ανοσοϊστοχημεία (ανασο-ιστοχημεια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ο)", "(σ)", "(ο)", "(ι)", "(στ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νοσοϊστ*&dq=

	// αντιλαϊκά (αντιλα-ικά)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*λαϊκ*&dq=

	// αντινατοϊκός (αντινατο-ικός)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(α)", "(τ)", "(ο)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νατοικ*&dq=

	// αντιηρωικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(ρ)", "(ο)", "(ι)", "(κ|σμ|στ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηρωικ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιροικ*&dq=

	// αουτσάιντερ (αουτσα-ιντερ)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(α)", "(ι)", "(ντ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σαιντερ*&dq=

	// απαλλιώς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(α)", "(λλ|λ)", "(οι|ι)", "(ω|ώ)", "(σ)", "$"}), "$1$2$3$4$5$6><$7$8"},

	// απαρτχάιντ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(χ)", "(α)", "(ι)", "(ντ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τχάιντ*&dq=

	// αποηχηροποίηση (apo-ixiropoiisi)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ο)", "(ι)", "(χ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ποηχ*&dq=

	// αποϊδρυματισμός (απο-ιδρυματισμός)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ο)", "(ι)", "(δ)", "(ρ)", "(.*)"}), "$1$2$3>-<$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ποϊδρ*&dq=

	// αποκαΐδια (αποκα-ι-δια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(κ)", "(α)", "(ι|ί|ϊ|ΐ)", "(δ)", "(ι)", "(α)", "(.*)"}), "$1$2$3$4>-$5-$6$7<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οκαίδια*&dq=

	// απονενοημένος (απενο-ιμενος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ο)", "(ι)", "(μ)", "(ε)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ενοιμέ*&dq=

	// αποϋλοποίηση (απο-υλοποιηση)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο)", "(υ|ι)", "(λ)", "(ο)", "(π)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απουλοπ*&dq=

	// απράυντος (απρα-υντος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ρ)", "(α)", "(υ|ύ|ι|ί)", "(ντ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πράυντ*&dq=

	// απρονόητος (aprono-itos), απρονοησία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ν)", "(ο)", "(ι)", "(τ|σ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ονόιτ*&dq=

	// απτόητα (apto-ita)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(τ)", "(ο)", "(ι|ί)", "(τ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απτοιτ*&dq=

	// αραμαϊκός (αραμα-ικός)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(α)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμαϊκ*&dq=

	// αραχνοΰφαντος (αραχνο-ύφαντος), αραχνοΰφασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ν)", "(ο)", "(υ|ύ|ι|ί)", "(φ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χνουφ*&dq=

	// αργατολόι (αργατολο-ι)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(γ)", "(α)", "(τ)", "(ο)", "(λ)", "(ο)", "(ι)", "(.*)"}), "$1$2$3$4$5$6$7$8>-<$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ργατολόι*&dq=

	// αργυροχοΐα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(χ)", "(ο)", "(ι|ί)", "(α)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οχοΐα*&dq=

	// αρνησιθεΐα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(σ)", "(ι)", "(θ)", "(ε)", "(ι|ί)", "(α)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ησιθεΐα*&dq=

	// Άρπυια (Άρπυι-α)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(π)", "(υ)", "(ι)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρπυι*&dq=

	// Αρσινόη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(σ)", "(ι)", "(ν)", "(ο)", "(ι)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρσινόη*&dq=

	// αρχαΐζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(χ)", "(α)", "(ι|ί)", "(ζ|κ|σμ|στ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρχαιζ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρχαικ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρχαϊσ*&dq=

	// αρχαιοϊστορικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(χ)", "(αι|ε)", "(ο)", "(ι)", "(στ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// ρχαιοϊστ

	// αρχοντολόι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ο)", "(λ)", "(ο)", "(ι)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντολόι*&dq=

	// ασυνεννόητα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ε)", "(νν|ν)", "(ο)", "(ι|ί)", "(σ|τ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νεννόητ*&dq=

	// αστεΐζομαι, αστεϊσμός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ε)", "(ι|ί)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*στεΐσ*&dq=

	// ατάιστος (ατα-ιστος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(α)", "(ι)", "(στ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αταιστ*&dq=

	// ατμοηλεκτρικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(μ)", "(ο)", "(ι)", "(λ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τμοηλ*&dq=

	// ατμόιππος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(μ)", "(ο)", "(ι)", "(π)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τμοιπ*&dq=

	// αϋπνία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(υ|ι)", "(π)", "(ν)", "(.*)"}), "$1$2>-<$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αυπν*&dq=

	// αυτηκοΐα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(κ)", "(ο)", "(ι|ί)", "(α)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηκοΐα*&dq=

	// αυτοϊκανοποίηση
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ο)", "(ι)", "(κ)", "(α)", "(ν)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τοϊκαν*&dq=

	// Αχαΐα, αχαϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(α)", "(ι|ί)", "(α|κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχαια*&dq=

	// αχάιδευτος (α-χάι-δευτος, α-χαϊ-δευτος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(α|ά)", "(ι|ϊ)", "(δ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχάιδε*&dq=

	// αχαΐρευτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(α)", "(ι|ί)", "(ρ)", "(ε)", "(υ|φ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χαΐρευ*&dq=

	// αχολόημα (axolo-ima)
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ο)", "(λ)", "(ο)", "(ι|ί)", "(μ)", "(α)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ολοιμα*&dq=

	// αχολόι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(ο)", "(λ)", "(ο)", "(ι)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχολόι*&dq=

	// αβανιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(β)", "(α)", "(ν)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αβανια*&dq=

	// αβάν πρεμιέρ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ρ)", "(ε|έ)", "(μ)", "(ι)", "(ε|έ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πρεμιέ*&dq=

	// αβαριάτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α|ά)", "(ρ)", "(ι)", "(α|ά)", "(τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαριάτ*&dq=

	// αβροχιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(ο|ό)", "(χ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βροχιά*&dq=

	// αγαθιάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(θ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αθιάρ*&dq=

	// αγάλια, αγάλλιασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γ)", "(α|ά)", "(λ)", "(ι)", "(α)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γ)", "(ά)", "(λλ|λ)", "(ι)", "(α)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγάλια*&dq=

	// αγαπησιάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(η)", "(σ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πησιάρ*&dq=

	// αγγελιάζομαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε|έ)", "(λ)", "(ι)", "(α|ά)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γελιάζ*&dq=

	// αγγουριά, αγγούρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ου|ού)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γουρια*&dq=
	// todo: ολιγουρία ?

	// αγδίκιωτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(δ)", "(ι|ί)", "(κ)", "(ι)", "(ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γδίκιω*&dq=

	// αγελαδήσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(δ)", "(η|ή)", "(σ)", "(ι)", "(α|ε|ο|ω)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αδίσιο*&dq=

	// ους|ας|ες|οι|ος|ου|ων|α|ε|ο
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α|ά)", "(γ)", "(ι)", "(ους|ας|ες|οι|ος|ου|ων|α|ε|ο)", "$"}), "$1-$2$3$4"},

	// αγιάζι
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(α)", "(γ)", "(ι)", "(α|ά)", "(ζ)", "(ι)", "(α)", "(.*)"}), "$1$2>-$3$4$5-$6$7<$8$9"},
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(α)", "(γ)", "(ι)", "(α|ά)", "(ζ)", "(ι)", "(.*)"}), "$1$2$3$4><$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγιάζι*&dq=

	// αγιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α|ά)", "(γ)", "(ι)", "(α|ά)", "(ά-ζω)", "$"}), "$1$2$3><$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγιάζω*&dq=

	// άγιασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ά)", "(γ)", "(ι)", "(α)", "(σ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// todo: fix: άγιασμα, αγιάσματος

	// αγιαστήρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α)", "(στ)", "(η|ή)", "(ρ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαστήρ*&dq=

	// αγιαστούρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α)", "(στ)", "(ου|ού)", "(ρ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαστούρ*&dq=

	// αγιατολάχ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γ)", "(ι)", "(α)", "(τ)", "(ο)", "(λλ|λ)", "(.*)"}), "$1$2$3$4><$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγιατολ*&dq=

	// αγιάτρευτα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(τ)", "(ρ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*γιάτρ*&dq=

	// αγιοβασιλιάτικα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(β)", "(α)", "(σ)", "(ι)", "(λ)", "(ι)", "(α|ά)", "(.*)"}),
		"$1>-$2$3$4-$5$6-$7$8-$9$10<$11$12"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιοβασιλιάτ*&dq=

	// Αγιοβασίλης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(β)", "(α)", "(σ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιοβασ*&dq=

	// αγιοβότανο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(β)", "(ο|ό)", "(τ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιοβότ*&dq=

	// αγιογδύτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(γ)", "(δ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιογδ*&dq=

	// αγιοδημητριακός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(δ)", "(η)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγιοδημ*&dq=

	// αγιοκέρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(κ)", "(ε|έ)", "(ρ)", "(ι)", "(α)", "(.*)"}), "$1$2$3>$4-$5$6-$7$8<$9$10"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(κ)", "(ε|έ)", "(ρ)", "(.*)"}), "$1$2$3>$4-<$5$6$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(κ)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά)", "$"}), "$1$2$3$4$5$6>$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*κέρι&dq=
	// todo kalokairiwn, kakokairi-wn "(αι|αί|ε|έ)"

	// αγιόκλημα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο|ό)", "(κ)", "(λ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιόκλ*&dq=

	// Αγιονορείτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(ν)", "(ο)", "(ρ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιονορ*&dq=

	// Αγιορειτικός (Αγι-ορειτικός)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(ρ)", "(ει|ι)", "(τ)", "(ι)", "(κ)", "(ό|ή)", "(.*)"}),
		"$1$2$3>-<$4$5$6$7$8$9$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιοριτικο*&dq=

	// Αγιορείτικος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(ρ)", "(ει|εί|ι|ί)", "(τ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιορειτ*&dq=

	// αγιοστράτηγος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(στ)", "(ρ)", "(α|ά)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// *γιοστρά*

	// αγιούπας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(π)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιούπ*&dq=

	// αγκαθένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(α)", "(θ)", "(ε|έ)", "(ν)", "(ι)", "(α|ε|ο|ω)", "(.*)"}),
		"$1$2$3$4$5$6$7><$8$9"},
	// http://www.neurolingo.gr/en/online_tools/lexiscope.htm?term=αγκαθένιος

	// αγκαθιά, αγκάθια, αγκαθιώνας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(α|ά)", "(θ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γκαθι*&dq=

	// αγκαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αγκαλιά*&dq=

	// αγκιναριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(ι)", "(ν)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκιναριά*&dq=

	// αγκιναροκούκια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ου|ού)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κούκι*&dq=

	// αγκιστριά, αγκίστρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(ι|ί)", "(στ)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκίστρια*&dq=

	// αγκωνιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γκ)", "(ω)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γκωνι*&dq=

	// αγνάντια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ν)", "(α|ά)", "(ντ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γνάντι*&dq=

	// αγουστιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ου)", "(στ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γουστι*&dq=

	// Αγραφιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(α)", "(φ)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γραφιώτ*&dq=

	// αγριελιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ι)", "(ε)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γριελι*&dq=

	// Αγρινιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ι)", "(ν)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γρινιώτ*&dq=

	// αγροτιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ο)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γροτιά*&dq=

	// αγρυπνιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(υ)", "(π)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγρύπνια*&dq=

	// αγυιά, αγυιόπαιδο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γ)", "(υ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγυι*&dq=

	// αγωγιάτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ω)", "(γ)", "(ι)", "(α|ά)", "(στ|τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γωγιάτ*&dq=

	// άδεια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(δ)", "(ει)", "(ος|α|ο)", "$"}), "$1$2$3$4><$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δειά*&dq=

	// αδειάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(δ)", "(ει)", "(α|ά)", "(ειάζω)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	//http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δειάζ*&dq=

	// αδελφοξαδέλφια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ξ)", "(α)", "(δ)", "(ε|έ)", "(λ|ρ)", "(φ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ξαδέλφι*&dq=

	// αδιάβαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(δ)", "(ι)", "(α|ά)", "(β)", "(α)", "(στ)", "(.*)"}), "$1$2$3$4><$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διάβαστ*&dq=

	// αδιαγούμιστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(γ)", "(ου|ού)", "(μ)", "(.*)"}), "$1$2><$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διαγούμ*&dq=

	// αδιαντροπιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(ντ)", "(ρ)", "(ο)", "(π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6$7-$8$9<$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντροπιά*&dq=

	// αδιάντροπα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α|ά)", "(ντ)", "(ρ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*διάντρ*&dq=

	// αδραξιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ρ)", "(α)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δραξι*&dq=

	// αετήσιος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ε)", "(τ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αετήσι*&dq=
	// εικοσαετήσι-ου!

	// δονάκιο
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ο)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ονάκι*&dq=

	// αηδόνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(η)", "(δ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2>$3-$4$5-$6$7<$8$9"},

	// αηδόνι (αη-δόνι)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(η)", "(δ)", "(ο|ό)", "(ν)", "(.*)"}), "$1$2><$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αηδον*&dq=

	// αηδονάκια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// αηδονολαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ν)", "(ο)", "(λ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ονολαλιά*&dq=

	// αηδονοφωλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ν)", "(ο)", "(φ)", "(ω)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*φωλι*&dq=
	// todo: φωλιά
	// todo: αηδονήσιος

	// αητός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "('η')", "(τ)", "(ο|ό)", "(.*)"}), "$1$2><$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αητο*&dq=

	// αθεμέλιωτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε)", "(μ)", "(έ)", "(λ)", "(ι)", "('ω')", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θεμέλιω*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θεμέλιο*&dq=
	// θεμελιωτής

	// Αθηνιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(η)", "(ν)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θηνιώτ*&dq=

	// Αιγαλιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(λ)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γαλιοτ*&dq=

	// αιμασιά, αιματιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(αι|ε)", "(μ)", "(α|ά)", "(τ|σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αιματια*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*εματια*&dq=

	// ακαμάκιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμάκια*&dq=

	// ακαταδεξιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(α)", "(δ)", "(ε)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ταδεξιά*&dq=

	// ακέριος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ακέρι*&dq=

	// ακεφιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ε)", "(φ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ακεφι*&dq=

	// ακοινωνισιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(οι|ι)", "(ν)", "(ω)", "(ν)", "(ι)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κινωνισι*&dq=

	// ακομπανιαμέντο, ακομπανιάρισμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(α|ά)", "(ν)", "(ι)", "(α|ά)", "(σμ|μ|ζ|τ|σ|ρ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	//http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπανιαμ*&dq=
	//http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπανιαρ*&dq=

	// ακόμπιαστα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ο|ό)", "(μπ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κόμπια*&dq=

	// ακόπιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ο|ό)", "(π)", "(ι)", "(α|ά)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κόπιαζ*&dq=

	// ακρίβια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ρ)", "(ι|ί)", "(β)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κρίβια*&dq=

	// ακριβογιός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(β)", "(ο)", "(γ)", "(ι)", "(ο|ό|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιβογι*&dq=

	// ακριτομύθεια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(τ)", "(ο)", "(μ)", "(υ|ύ|ι|ί)", "(θ)", "(ει)", "(α|ά|ο|ό|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// *ιτομύθεια*

	// ακρογιαλιά, ακρογιάλια
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ο)", "(γ)", "(ι)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4>$5-$6$7<$8$9"},

	// ακρογιάλι
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ο)", "(γ)", "(ι)", "(α|ά)", "(λ)", "(ι)", "(.*)"}), "$1$2$3$4><$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ογιαλι*&dq=

	// ακροδεσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ο)", "(δ)", "(ε)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ροδεσι*&dq=

	// ακροθαλασσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(α)", "(λ)", "(α)", "(σσ|σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θαλασσια*&dq=

	// ακροποταμιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(π)", "(ο)", "(τ)", "(α)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οποτάμι*&dq=
	// todo: ποτάμια

	// ακρορεματιά
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ε)", "(μ)", "(α)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*εματια*&dq=

	// αλαλιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λαλιάσ*&dq=

	// αλαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(λ)", "(ι)", "(ά|έ|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλαλιά*&dq=

	// αλανιάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*αλανι*&dq=

	// αλατενιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α)", "(τ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λατένι*&dq=

	// αλατζαδένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τζ)", "(α)", "(δ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τζαδένι*&dq=

	// αλατζένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(τζ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λατζένιος*&dq=

	// αλατιέρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α)", "(τ)", "(ι)", "(ε|έ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλατιε*&dq=

	// αλαφιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α|ά)", "(φ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλαφιά*&dq=

	// αλαφριός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α|ά)", "(φ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλαφρι*&dq=

	// αμυαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(υ|ι)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*μυαλι*&dq=

	// αλαφρομυαλιά
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ο)", "(μ)", "(υ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4>$5-$6$7<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ομυαλι*&dq=

	// άμυαλα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(υ|ι)", "(α)", "(λ)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μυαλα*&dq=

	// αλαφρόμυαλος
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(α)", "(φ)", "(ρ)", "(ο|ό)", "(μ)", "(υ|ι)", "(α)", "(λ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αφρόμυαλ*&dq=

	// αλειματουργησιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ου)", "(ρ)", "(γ)", "(η)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// *τουργισι* α|ά|ε|έ|ω|ώ

	// αλεπουδίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(π)", "(ου)", "(δ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*επουδίσι*&dq=

	// αποφώλιος (αποφώλι-ος)
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(α)", "(π)", "(ο)", "(φ)", "(ω|ώ)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
	// "$1$2$3$4$5$6$7$8>-<$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αποφωλι*&dq=

	// αλεποφωλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(π)", "(ο)", "(φ)", "(ω)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*εποφωλι*&dq=

	// αλεσιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(λ)", "(ε)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλεσιά*&dq=

	// αλετριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ε|έ)", "(τ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλετρι*&dq=

	// αλευρένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ε)", "(υ|β)", "(ρ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ευρένιο*&dq=

	// αλευριά, αλεύρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ε)", "(υ|ύ|β)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλευριά*&dq=

	// αλήθεια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(η|ή)", "(θ)", "(ει|ι)", "(α|ε|ω)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αληθει*&dq=

	// αλησμονησιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(σμ)", "(ο)", "(ν)", "(η)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ησμονησι*&dq=

	// αλιάδα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(λ)", "(ι)", "(α|ά)", "(δ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αλιάδ*&dq=

	// αλισφακιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(φ)", "(α)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σφακι*&dq=

	// αλιφασκιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(φ)", "(α)", "(σ)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λιφασκι*&dq=

	// αλλαξιά (αλλα-ξια)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(λλ|λ)", "(α)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αλλαξι*&dq=

	// αλλαξοκαιριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ξ)", "(ο)", "(κ)", "(αι|ε)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ξοκαιρι*&dq=

	// αλλαξοκωλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ξ)", "(ο)", "(κ)", "(ω)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ξοκωλ*&dq=

	// αλλιώς
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(λλ|λ)", "(οι|ι)", "(ω|ώ)", "(σ)", "$"}), "$1$2$3><$4$5"},
	// ^alios$

	// αλλιώτικος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(λλ|λ)", "(οι|ι)", "(ω|ώ)", "(τ)", "(ι)", "(κ)", "(ος|η|ο)", "$"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλλιώ*&dq=
	// todo: i-kia

	// αλλοκοτιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λλ|λ)", "(ο)", "(κ)", "(ο)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλλοκοτιά*&dq=

	// αλλοχωριανός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ω)", "(ρ)", "(ι)", "(α)", "(ν|τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χωριαν*&dq=
	// todo: αρχοντοχωριατιά

	// αλόγιαστα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ο|ό)", "(γ)", "(ι)", "(α)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λόγιαστ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*λόγιαζ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λόγιασμ*&dq=

	// αλογίσια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ο)", "(γ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λογίσι*&dq=

	// Αλοννησιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(!γ)", "(ν)", "(η)", "(σ)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νησιότ*&dq=

	// νησιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ν)", "(η)", "(σ)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νησιότ*&dq=

	// αλουμινένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(ε|έ)", "(ν)", "(ι)", "(α|ε|ω)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινένι*&dq=

	// αλουσιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(λ)", "(ου)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλουσιά*&dq=

	// αλυπησιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(υ|ι)", "(π)", "(η)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λυπησιά*&dq=

	// αλφαδιά, αλφάδια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(φ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*λφάδι*&dq=

	// αμαξιάτικος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(α|ά)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμαξι*&dq=

	// αμερικανιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(ι)", "(κ)", "(α)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*ερικανι*&dq=

	// αμμουδιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μμ|μ)", "(ου|ού)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμμουδι*&dq=

	// Αμμοχωστιανός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(στ)", "(ι)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ωστιαν*&dq=

	// άμοιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(μ)", "(οι|ι)", "(α)", "(στ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*άμοιαστ*&dq=

	// αμοργιανός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(ο|ό)", "(ρ)", "(γ)", "(ι)", "(α)", "(ν)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μοργιαν*&dq=

	// αμορφωσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(ο)", "(ρ)", "(φ)", "(ω)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μορφωσι*&dq=

	// αμπελίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(η|ή)", "(σσ|σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελίσι*&dq=
	// todo: αμπελοειδή

	// αμπιγιέ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(ι)", "(γ)", "(ι)", "(ε|έ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπιγι*&dq=

	// αμπόλιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(ο|ό)", "(λ)", "(ι)", "(α|ά)", "(ζ|ρ|σμ|στ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπόλιασ*&dq=

	// αμπραγιάζ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(ρ)", "(α)", "(γ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπραγιά*&dq=

	// αμυγδαλένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(δ)", "(α)", "(λ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γδαλένι*&dq=

	// αμυγδαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(δ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γδαλι*&dq=
	// todo: agriamugdalia

	// αναγαλλιάζω (αναγαλ-λιά-ζω)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(α)", "(γ)", "(α)", "(λλ|λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ναγαλλι*&dq=

	// Μαγούλιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(γ)", "(ου|ού)", "(λλ|λ)", "(ι|ί)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7>-<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Μαγούλι*&dq=

	// αναγουλιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ου|ού)", "(λλ|λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/gr eekLang/modern_greek/tools/lexica/search.html?start=20&lq=*γουλι*&dq=

	// αναδεξιμιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε)", "(ξ)", "(ι)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δεξιμι*&dq=

	// ιεροδουλία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(ο)", "(δ)", "(ου)", "(λ)", "(ι)", "(α)", "(.*)"}),
		"$1$2$3$4$5$6$7$8>-<$9$10"},
	// ξενοδουλία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ο)", "(δ)", "(ου)", "(λ)", "(ι)", "(α|ε|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8>-<$9$10"},
	// εθελοδουλία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε)", "(λ)", "(ο)", "(δ)", "(ου)", "(λ)", "(ι)", "(α|ε|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9>-<$10$11"},
	// αναδουλειά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ου)", "(λ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*δουλι*&dq=

	// αναμαλλιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α|ά)", "(λλ|λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(ζ|ρ|σμ|στ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μαλλιαζ*&dq=

	// αναπαραδιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(α)", "(ρ)", "(α)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απαραδι*&dq=

	// Αναπλιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(α|ά)", "(π)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ναπλι*&dq=

	// αναποδιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(α)", "(π)", "(ο)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ναποδι*&dq=

	// ανάριος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(α|ά)", "(ρ)", "(ι)", "(ος|α|ο)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=0&lq=*ανάρι*&dq=

	// αναρωτιέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ω)", "(τ)", "(ι)", "(ε|έ|ο|ό)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρωτιέ*&dq=

	// ανατριχιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ρ)", "(ι)", "(χ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*τριχι*&dq=

	// Αναφιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(φ)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αφιώτ*&dq=

	// ανεγκεφαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(ε)", "(φ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκεφαλι*&dq=

	// ανεμελιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ε)", "(μ)", "(ε)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νεμελι*&dq=

	// ανεμοβλογιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ευ|β)", "(λ)", "(ο)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βλογι*&dq=

	// ανεψιά, ανηψιά, ανιψιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ε|έ|η|ή|ι|ί)", "(ψ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανιψι*&dq=

	// ανήλιαγος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(η|ή)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ανήλι*&dq=

	// ανημποριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(μπ)", "(ο)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ημπορι*&dq=

	// λαμπαδηφορία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(η)", "(φ)", "(ο)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7>-<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*ηφορι*&dq=

	// ανηφοριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(φ)", "(ο)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*ηφορι*&dq=

	// ανθογυάλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ)", "(α|ά)", "(λ)", "(ι)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γυάλι*&dq=
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ|ι)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ|ι)", "(α|ά)", "(λ)", "(ι)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(γ)", "(υ)", "(α)", "(λ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*γιαλι*&dq=

	// ανθρακιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*θρακι*&dq=

	// ανθρωπιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(θ)", "(ρ)", "(ω)", "(π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νθρωπι*&dq=

	// Άννια, Αννιώ
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(Α|Ά)", "(νν|ν)", "(ι)", "(α|ας|ω|ώ|ως|ώς)", "$"}), "$1$2$3>$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Άννιο*&dq=

	// ανοιξιάτικα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(οι|ι)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανοιξι*&dq=

	// ανοιχτωσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(οι|ι)", "(χ)", "(τ)", "(ω)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νιχτοσι*&dq=

	// ανοργανωσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(ν)", "(ω)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γανωσι*&dq=

	// ανορεξιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ρ)", "(ε)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ορεξι*&dq=

	// ανοστιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ο)", "(στ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανοστι*&dq=

	// ανταριάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ντ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά)", "(ζ|σ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανταριάζ*&dq=

	// αντζούγια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(τζ|τσ)", "(ου|ού)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντζούγια*&dq=

	// αντηλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ντ)", "(η)", "(λ)", "(ι)", "(α|ά|ας|άς)", "$"}), "$1$2$3$4$5$6><$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντηλι*&dq=

	// αντιζυγιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(υ|ύ)", "(γ)", "(ι)", "(α|ά)", "(ά-ζω)", "$"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζυγιά*&dq=
	// todo: remove?

	// αντιζύγιασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(υ|ύ)", "(γ)", "(ι)", "(α|ά)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζυγιάσ*&dq=

	// αντιλαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(λ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιλαλι*&dq=

	// αντιμετριέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(ε)", "(τ)", "(ρ)", "(ι)", "(ε|έ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μετριέ*&dq=

	// αντιμιλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ι)", "(μ)", "(ι)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντιμιλι*&dq=

	// αντιφεγγιά, αντιφέγγιασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ι)", "(φ)", "(ε|έ)", "(γγ|γκ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*φέγγι*&dq=

	// αντρειεύομαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ρ)", "(ει|ι)", "(ε)", "(υ|ύ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντρειεύ*&dq=

	// αντρίκεια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ρ)", "(ι|ί)", "(κ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντρίκει*&dq=

	// ανωμαλιάρης
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ω)", "(μ)", "(α)", "(λ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ωμαλιάρ*&dq=

	// αξιώτικος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ξ)", "(ι)", "(ώ)", "(τ)", "(ι)", "(κ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αξιώτικ*&dq=

	// αξόμπλιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(λ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*μπλια*&dq=

	// αξουρισιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ)", "(ρ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υρισι*&dq=

	// απαγκιάζω, απάγκιο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(α|ά)", "(γγ|γκ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απαγκι*&dq=

	// απανεμιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(α)", "(ν)", "(ε)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανεμιά*&dq=

	// απανωσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(α)", "(ν)", "(ω)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πανωσι*&dq=

	// απαρνησιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(ν)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*αρνησι*&dq=

	// απατεωνιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(ε)", "(ω|ώ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ατεωνι*&dq=

	// απελπισιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ε)", "(λ)", "(π)", "(ι)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελπισι*&dq=

	// απηλιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(η)", "(λ)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πηλιώτ*&dq=

	// απηλογιέμαι, απολογιέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο|η)", "(λ)", "(ο)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απολογι*&dq=

	// άπιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(π)", "(ι)", "(α|ά)", "(στ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απιαστ*&dq=

	// απιδιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*απιδι*&dq=

	// άπιοτος -η -ο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(π)", "(ι)", "(ω)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πιοτό*&dq=
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πιοτ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απιοτ*&dq=

	// απλανιάριστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(λ)", "(α)", "(ν)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πλανιάρ*&dq=

	// ευχέρεια (ευχέρι-α)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(υ|φ)", "(χ)", "(ε|έ)", "(ρ)", "(ι)", "(α)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8$9"},
	// δυσχέρεια (δυσχέρι-α)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(υ)", "(σ)", "(χ)", "(ε|έ)", "(ρ)", "(ι)", "(α)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9$10"},
	// απλοχεριά (απλοχε-ριά)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χερια*&dq=

	// απλοχωριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(χ)", "(ω|ώ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*οχωρι*&dq=

	// απλυσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(λ)", "(υ)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πλυσι*&dq=

	// απλωσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(λ)", "(ω)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πλωσι*&dq=

	// απόειδα
	GrhyphRule{customRegexpCompile(
		[]string{"(α)", "(π)", "(ο|ό)", "(ει|ι)", "(δ)", "(α|ε)", "(.*)"}), "$1$2$3><$4$5$6$7"},

	// απόγιομα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ο|ό)", "(γ)", "(ι)", "(ο)", "(μ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*όγιομ*&dq=

	// αποδιωγμός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(ω)", "(γ)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διωγμ*&dq=

	// αποδιώχνω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(ω|ώ)", "(χ)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διωχν*&dq=

	// απέδιωξα, αποδιώξουμε
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(ω|ώ)", "(ξ)", "$"}), "$1$2$3$4"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(ω|ώ)", "(ξ)", "(ιώ-ξω)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διώξ*&dq=
	// (α)(π)(ο|ό)(δ)(ι)(ω|ώ)(ξ)(...)

	// αποθυμιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο)", "(θ)", "(υ|ι)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αποθυμι*&dq=

	// αποκοιμιέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο)", "(κ)", "(οι|ι)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αποκοιμι*&dq=

	// αποκοτιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(κ)", "(ο)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οκοτι*&dq=

	// αποκουτιαίνω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(κ)", "(ου)", "(τ)", "(ι)", "(αι|αί|ε|έ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οκουτιαί*&dq=

	// απολησμονιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(σμ)", "(ο)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ησμονι*&dq=

	// απονέρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ν)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ονέρι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νέρι*&dq=
	// todo: αγιονέρια

	// απονιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*απονι*&dq=

	// αποξεχνιέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ξ)", "(ε)", "(χ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ξεχνι*&dq=

	// πιόμα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(π)", "(ι)", "(ω|ώ)", "(μ)", "(α)", "(.*)"}), "$1$2><$3$4$5$6"},
	// απόπιομα, απόπιωμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ο|ό)", "(π)", "(ι)", "(ω|ώ)", "(μ)", "(α)", "(.*)"}), "$1$2$3$4$5><$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πιομα*&dq=

	// αποσκίασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ο)", "(σ)", "(κ)", "(ι|ί)", "(α)", "(σμ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// απόσκιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ο|ό)", "(σ)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απόσκι*&dq=

	// αποτελείωση (αποτελει-ω-ση)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(π)", "(ο)", "(τ)", "(ε)", "(λ)", "(ει|ι)", "(ω|ώ)", "(ση|σης|σεων|σεως)", "$"}),
		"$1$2$3$4$5$6$7>-<$8$9"},
	// αποτέλειωμα (αποτέλειω-μα)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο)", "(τ)", "(ε|έ)", "(λ)", "(ει|ι)", "(ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αποτέλειω*&dq=

	// αποτέτοιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ε|έ)", "(τ)", "(οι|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τέτοι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τέτι*&dq=

	// αποτρύγια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ρ)", "(υ|ύ|ι|ί)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τρύγι*&dq=

	// αποτρυγίδια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γίδι*&dq=

	// αποφάγια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο)", "(φ)", "(α|ά)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αποφάγι*&dq=

	// απριλιάτικα, απριλιάτικος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ρ)", "(ι)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απριλιά*&dq=

	// απροθυμιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ρ)", "(ο)", "(θ)", "(υ|ι)", "(μ)", "(ι)", "(ά|έ|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απροθυμι*&dq=

	// αραδιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*αραδι*&dq=

	// αραποσιτιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ο)", "(σ)", "(ι|ί)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αποσιτι*&dq=

	// αραχνιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(α)", "(χ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αραχνιά*&dq=

	// Αρβανιτιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(ν)", "(ι|ί)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βανιτι*&dq=

	// αργαλειός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(γ)", "(α)", "(λ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αργαλει*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αργαλι*&dq=

	// αρκουδιάρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(δ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουδιάρ*&dq=

	// αρκουδίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(δ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουδίσι*&dq=

	// αρμαθιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(μ)", "(α|ά)", "(θ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρμαθι*&dq=

	// αρματωσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(ω)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ατωσι*&dq=

	// αρμεξιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(μ)", "(ε)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρμεξι*&dq=

	//αρνιέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=0&lq=*αρνι*&dq=15

	// αρπαξιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(π)", "(α)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρπαξι*&dq=

	// αρραβωνιαστικιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*στικι*&dq=

	// αρραβωνιάζω, αρρεβωνιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρρ|ρ)", "(α|ε)", "(β)", "(ω)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ραβωνι*&dq=

	// αρρώστια, αρρώστεια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρρ|ρ)", "(ω|ώ)", "(στ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ροστι*&dq=

	// αρτιφισιέλ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φισι*&dq=

	// αρχιδιάκος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ι)", "(δ)", "(ι)", "(α|ά)", "(κος)", "$"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χιδιάκ*&dq=

	// αρχιμαφιόζος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(φ)", "(ι)", "(ο|ό)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μαφιό*&dq=

	// αρχιμηνιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(η)", "(μ)", "(η)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χημηνι*&dq=

	// αρχιχρονιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(χ)", "(ρ)", "(ο)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιχρονι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηχρονι*&dq=

	// αρχοντιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(χ)", "(ο)", "(ντ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρχοντι*&dq=

	// ασβολιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(β)", "(ο)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σβολι*&dq=
	// todo: greeklish azvolia ?

	// άσκιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(σ)", "(κ)", "(ι)", "(α)", "(στ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// άσκιαχτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(σ)", "(κ)", "(ι)", "(α)", "(χ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σκιαστ*&dq=

	// ασπρουλιάρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(λ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουλιάρ*&dq=

	// αστοχασιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ο)", "(χ)", "(α)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*στοχασι*&dq=

	// αστραποφεγγιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ε|έ)", "(γγ|γκ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*φεγγι*&dq=

	// Αστροπαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ο)", "(π)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ροπαλι*&dq=

	// Αστυπαλιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(υ|ι)", "(π)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*στυπαλι*&dq=

	// ασυλλογισιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ο)", "(γ)", "(ι)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λογισι*&dq=
	// todo: (σ)->(σσ|σ) ?

	// ασχήμια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(σ)", "(χ)", "(η|ή)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ασχήμι*&dq=

	// ματεριαλισμός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(ε|έ)", "(ρ)", "(ι)", "(α)", "(λ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// αταίριαστα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(αι|αί|ε|έ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ατερι*&dq=
	//todo: αϊταίρι ?

	// ατέλειωτα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(ε|έ)", "(λ)", "(ει|ι)", "('ω|ώ')", "(μ|ν|σ|τ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τέλειωτ*&dq=
	// todo: τελειώνω

	// ατελιέ
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(τ)", "(ε)", "(λ)", "(ι)", "(ε|έ)", "$"}), "$1$2$3$4$5><$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ατελιε*&dq=

	// ατημελησιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(ε)", "(λ)", "(η)", "(σσ|σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μελισσι*&dq=

	// ατλαζένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ζ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αζένι*&dq=

	// ατόφιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(ο|ό)", "(φ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ατόφι*&dq=

	// ατραξιόν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ρ)", "(α|ά)", "(ξ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τραξι*&dq=

	// ατσαλένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλένι*&dq=
	// todo: γυαλένιος

	// ατσαλιά, ατσάλια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τσ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*τσαλι*&dq=

	// αυλακιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(υ|β)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αυλακι*&dq=

	// φιλαυτία, ναυτία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ|ν)", "(α)", "(υ|φ)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// αυτιά, αφτιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(υ|φ)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=50&lq=*αυτι*&dq=

	// αφαγιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(φ)", "(α)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αφαγιά*&dq=

	// αφεντιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ε)", "(ντ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φεντι*&dq=

	// αφιόνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ι)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φιόνι*&dq=

	// αφιονίζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ι)", "(ο|ό)", "(ν)", "(ι|ί)", "(.*)"}), "$1$2$3><$4$5$6$7"},

	// άφκιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(φ)", "(κ)", "(ι)", "(α|ά)", "(γ|στ|χ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*άφκια*&dq=
	// todo: φκιάνω

	// αφοβιά (αφο-βιά, αφοβι-α)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(φ)", "(ο)", "(β)", "(ι)", "(ά|έ|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αφοβι*&dq=

	// αφραγκιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ρ)", "(α)", "(γγ|γκ)", "(ι)", "(ά|έ|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φραγκι*&dq=

	// αφροντησιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ρ)", "(ο)", "(ντ)", "(η)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φροντησι*&dq=

	//υδροξύλιο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ρ)", "(ο)", "(ξ)", "(υ|ύ|ι|ί)", "(λ)", "(ι)", "(α|ο)", "(.*)"}),
		"$1$2$3$4$5$6$7$8>-<$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δροξυλι*&dq=
	// αφροξυλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ξ)", "(υ|ύ|ι|ί)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οξυλι*&dq=

	// αχιόνιστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(ι)", "(ο|ό)", "(ν)", "(ι)", "(στος|στη|στο)", "(.*)"}), "$1$2$3$4><$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχιόνιστ*&dq=

	//αχλαδιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(λ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χλαδι*&dq=

	// αχορταγιά, αχορτασιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(ο)", "(ρ)", "(τ)", "(α)", "(γ|σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχορταγι*&dq=

	// αχρόνιαστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ρ)", "(ο|ό)", "(ν)", "(ι)", "(α)", "(στ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χρόνιαστ*&dq=
	// todo χρωνιάζω...

	// αχυρένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(υ|ι)", "(ρ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χυρένι*&dq=

	// αψυλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ψ)", "(υ|ι)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αψυλι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αψιλι*&dq=

	// αβάντι (αβά-ντια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α|ά)", "(ντ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*βάντι*&dq=

	// αβγουλάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ουλακι*&dq=

	// αγαλματάκι (αγαλματά-κια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ματάκι*&dq=

	// αγγελάκι
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(γ)", "(ε)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γελάκι*&dq=

	// αγγελουδάκια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(δ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουδάκι*&dq=

	// αγγελούδι
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(λ)", "(ου|ού)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*λουδι*&dq=

	// αγγόνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γγ|γκ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγγόνι*&dq=

	// αγέρι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(γ)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγέρι*&dq=

	// αγκαθάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(θ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αθάκι*&dq=

	// αγκιδάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(δ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ιδάκι*&dq=

	// αγκίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκίδι*&dq=

	// αγκιστράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*στράκι*&dq=

	// αγκύλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γγ|γκ)", "(υ|ύ|ι|ί)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγκυλι*&dq=

	// *ωναριάτ*
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(ν)", "(α)", "(ρ)", "(ι)", "(α|ά)", "(τ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ωναριάτ*&dq=
	// αγκωνάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(ω)", "(ν)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκωνάρι*&dq=

	// αγοράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*οράκι*&dq=

	// αγόρια, αγώρια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(γ)", "(ω|ώ)", "(ρ)", "(ι)", "(ια|ιού|ιών)", "$"}), "$1$2$3$4$5>$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*οράκι*&dq=

	// αγριμάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ι)", "(μ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ριμάκι*&dq=

	// αγρίμι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ι|ί)", "(μ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γρίμι*&dq=

	// αγριοβόρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ι)", "(ο)", "(β)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ριοβόρι*&dq=

	// αγώγι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(γ)", "(ω|ώ)", "(γ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=αγώγι*&dq=

	// αδελφάκι, αδερφάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ|ρ)", "(φ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρφάκι*&dq=

	// μισαδελφία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(α)", "(δ)", "(ε)", "(λ)", "(φ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8>-<$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σαδέλφι*&dq=
	// φιλαδέλφια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ι)", "(λ)", "(α)", "(δ)", "(ε|έ)", "(λ)", "(φ)", "(ι)", "(α|ά)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9$10>-<$11$12"},
	// αδέλφι, αδέρφι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε|έ)", "(λ|ρ)", "(φ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*δέλφι*&dq=

	// αδερφομεράδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(μ)", "(ε)", "(ρ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ομεράδ*&dq=

	// αδερφομοιράδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(οι|ι)", "(ρ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μοιράδι*&dq=
	// todo: *μοιράσι*

	// ισομοιρία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(ο)", "(μ)", "(οι|ι)", "(ρ)", "(ι)", "(α|ε|ω)", "(.*)"}), "$1$2$3$4$5$6$7>-<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σομοίρι*&dq=

	// αδερφομοίρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(μ)", "(οι|οί|ι|ί)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ομοίρι*&dq=

	// αδραχτάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χτάκι*&dq=

	// αδράχτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ρ)", "(α|ά)", "(χ)", "(τ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δράχτι*&dq=

	// αεράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*εράκι*&dq=

	// αεροπλανάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*ανακι*&dq=

	// διακονιάρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(κ)", "(ο)", "(ν)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8<$9$10$11"},
	// διακονιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α|ά)", "(κ)", "(ο)", "(ν)", "(ι)", "(ά|έ|ώ)", "(.*)"}), "$1$2$3-$4$5-$6$7<$8$9"},
	// todo: αρχιδιακονία
	// ακόνι (ακό-νια, ακο-νιά), εξαιτίας αρχιδιακονίας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ό)", "(ν)", "(ι)", "(α|ω)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ο)", "(ν)", "(ι)", "(ά|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=30&lq=*ακόνι*&dq=

	// ακουμπιστήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(ι)", "(στ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπιστήρι*&dq=

	// ακρομόλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ρ)", "(ο)", "(μ)", "(ο|ό)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// κρομόλια

	// ακρωτήρι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(κ)", "(ρ)", "(ω)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(ια|ιού|ιών)", "$"}),
		"$1$2$3$4$5$6$7$8>$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ακρωτήρι*&dq=

	// αλατάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λατάκι*&dq=

	// αλάτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*αλάτι*&dq=

	// αλειφτήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(φ)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιφτήρι*&dq=

	// άλειωτος, άλιωτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(λ)", "(ει|ι)", "(ω)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// ! http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*άλιοτ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*άλιωτ*&dq=

	// μωρουδιακός
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ου)", "(δ)", "(ι)", "(α)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουδιακ*&dq=
	// αλεπούδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=160&lq=*ούδι*&dq=

	// αλεστήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ε)", "(στ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λεστήρι*&dq=

	// αλετροπόδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ο)", "(π)", "(ο|ό)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ροπόδι*&dq=

	// αλευράκι (υράκι)
	// todo(end)
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*υράκι*&dq=
	// αλευράκι (ευράκι)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(υ|β)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ευράκι*&dq=

	// αλητάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ητάκι*&dq=

	// ολισβερίσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ε)", "(ρ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ερίσι*&dq=

	// αλκολίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(λ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ολίκι*&dq=

	// άλλοθι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ο)", "(θ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λοθι*&dq=

	// αλογάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(γ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ογάκι*&dq=

	// αλσάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*σάκι*&dq=

	// αλώνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ω|ώ)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αλωνι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αλονι*&dq=

	// αμανάτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(ν)", "(α|ά)", "(τ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μανάτι*&dq=

	// αμαξάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ξ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αξακι*&dq=

	// αμέτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(ε|έ)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μέτι*&dq=

	// αμόνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμονι*&dq=

	// σνομπαρία (σνομπαρι-α)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(ν)", "(ο)", "(μπ)", "(α)", "(ρ)", "(ι)", "(α|ε|ω)", "(.*)"}),
		"$1$2$3$4$5$6$7$8>-<$9$10"},
	// αμπάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*μπαρι*&dq=

	// αμπελάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελάκι*&dq=

	// αμπέλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(ε|έ)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=30&lq=*μπέλι*&dq=

	// κάμπριο (καμπρι-ο)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(α|ά)", "(μπ)", "(ρ)", "(ι)", "(ω|ώ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// αμπρί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μπ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμπρί*&dq=

	// αμυγδαλάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαλακι*&dq=

	// ανεκδοτάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οτάκι*&dq=

	// ανεμίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ε)", "(μ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εμίδι*&dq=

	// ανεμιστηράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(η)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*στηράκι*&dq=

	// ανεμοβόρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ε)", "(μ)", "(ο)", "(β)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μοβόρι*&dq=

	// ανεμογκάστρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(α|ά)", "(στ)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκάστρι*&dq=

	// ανθί, άνθι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α|ά)", "(ν)", "(θ)", "(ι)", "(ια|ιού|ιών)", "$"}), "$1$2$3$4>$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ανθί*&dq=

	// ανθοτύρι, βουλλωτύρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(τ)", "(υ|ύ|ι|ί)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οτύρι*&dq=
	// todo: fix εκδοτήρια (ekdotiria)

	// ανθρωπάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(π)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ωπάκι*&dq=

	// ανιψάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ψ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιψάκι*&dq=

	// ανοιχτήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(χ)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιχτήρι*&dq=

	// αντερί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ντ)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντερί*&dq=

	// αντιζύγι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ζ)", "(υ|ύ|ι|ί)", "(γ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιζύγι*&dq=

	// αντιστύλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ι)", "(στ)", "(υ|ύ|ι|ί)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντιστύλι*&dq=

	// αντράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντρακι*&dq=

	// αντριλίκι, υπαλληλίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(λ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*ιλίκι*&dq=
	// todo: book κοροϊδιλίκια

	// ανώγι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ω|ώ)", "(γ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νώγι*&dq=

	// ανώφλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω|ώ)", "(φ)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ώφλι*&dq=

	// αξόνι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ξ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αξόνι*&dq=

	// απακούμπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ου|ού)", "(μπ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=30&lq=*κούμπι*&dq=

	// απανωπροίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ρ)", "(οι|οί|ι|ί)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*προίκι*&dq=

	// απλάδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(λ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πλάδι*&dq=
	// todo: *λάδι* ?

	// αποβόρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ο)", "(β)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ποβόρι*&dq=
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οβόρι*&dq=

	// απολειφάδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(φ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιφάδι*&dq=

	// απομεινάδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινάδι*&dq=

	// απομεινάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(μ)", "(ει|ι)", "(ν)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μεινάρι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μινάρι*&dq=

	// αποπαίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο)", "(π)", "(αι|αί|ε|έ)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αποπαίδι*&dq=

	// αποπότι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(π)", "(ο|ό)", "(τ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οπότι*&dq=

	// αποσαρίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(α)", "(ρ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σαρίδι*&dq=

	// αποσπόρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(π)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*σπόρι*&dq=
	// with and without = with synizesis

	// αποτηγανίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(ν)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γανίδι*&dq=

	// αποφόρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ο)", "(φ)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αποφόρι*&dq=

	// αποχτενίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ενίδι*&dq=

	// αραλίκι
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ρ)", "(α)", "(λ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ραλίκι*&dq=

	// αραξοβόλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(α)", "(ξ)", "(ο)", "(β)", "(ο|ό)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9$10><$11$12"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αξοβόλι*&dq=

	// Αραπάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απάκι*&dq=

	// Αραπιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(α)", "(π)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Αραπιά*&dq=

	//αραχνάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χνάκι*&dq=

	// αρίδι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(ρ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αρίδι*&dq=

	// αρμάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(μ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρμάρι*&dq=

	// αρμίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(μ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρμίδι*&dq=

	// αρμίθι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(μ)", "(ι|ί)", "(θ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρμίθι*&dq=

	// αρμυρίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(μ)", "(υ|ι)", "(ρ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μυρίκι*&dq=

	// αρνάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρνάκι*&dq=

	// αρνί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*αρνί*&dq=

	// αρχίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(χ)", "(ι|ί)", "(δ)", "(ι)", "(ιά|ιού|ιών)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρχίδι*&dq=

	// αρχονταρίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρίκι*&dq=

	// ασημένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(η)", "(μ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σιμενι*&dq=

	// ασήμι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(σ)", "(η|ή)", "(μ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ασήμι*&dq=

	// ασκαύλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(α)", "(υ|ύ|β)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*καύλι*&dq=

	// ασκέρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(κ)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σκέρι*&dq=

	// ασκί
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α)", "(σ)", "(κ)", "(ι)", "(ιά|ιού|ιών)"}), "$1$2$3$4>$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ασκί*&dq=

	// ασπραδάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(δ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αδάκι*&dq=

	// ασπράδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ρ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πράδι*&dq=

	// αστάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(στ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αστάρι*&dq=

	// αστειάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ει|ι)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*στειάκι*&dq=

	// αστέρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ε)", "(ρ)", "(ι)", "(ια|ιού|ιών)", "$"}), "$1$2$3$4$5>$6"},
	// todo: ξαστεριά
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=αστερι*&dq=

	// αστρί
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(α|ά)", "(στ)", "(ρ)", "(ι)", "(ιά|ιού|ιών)", "$"}), "$1$2$3$4>$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=αστρί*&dq=

	// αστροπελέκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(ε|έ)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ελέκι*&dq=

	// ασφοδέλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(δ)", "(ε|έ)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οδέλι*&dq=

	// ασφοδίλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(δ)", "(ι|ί)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οδίλι*&dq=

	// αφάλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(φ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αφάλι*&dq=

	// αχείλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(ει|εί|ι|ί)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχείλι*&dq=

	// αχνάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ν)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χνάρι*&dq=

	// αχούρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(ου|ού)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχούρι*&dq=

	// αψέντι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ψ)", "(ε|έ)", "(ντ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ψέντι*&dq=

	// βαβούλι (βαβούλια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(β)", "(ου|ού)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαβούλι*&dq=
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αβούλι*&dq=
	// todo: αβούλιαχτος, στραβούλιακας

	// Βαγγέλιο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(γγ|γκ)", "(ε|έ)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαγγέλιο*&dq=

	// βαγένι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(γ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Βαγένι*&dq=

	// Βάγια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α|ά)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Βάγια*&dq=

	// βαγόνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(γ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαγόνι*&dq=

	// βαζάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(ζ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αζάκι*&dq=

	// βαθιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(θ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαθι*&dq=

	// βαθύσκιωτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(υ|ύ|ι|ί)", "(σ)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θυσκιω*&dq=
	// todo: ανήσκιωτος
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισκιω*&dq=
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισκιο*&dq=

	// βάι (βαι)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ά)", "(ι)", "$"}), "$1$2$3"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βάι&dq=

	// Βαϊοφόρος (βα-ι-οφόρος)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(ι)", "(ο)", "(φ)", "(.*)"}), "$1$2$3>-<$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαϊοφ*&dq=

	// Βάιος (Βα-ι-ος)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(Β)", "(α|ά)", "(ι|ί)", "(ος...)", "$"}), "$1$2-$3-$4"},
	// todo: Βαΐων (βα-ί-ων)

	// Βάια, βάιο (βάι-ο)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α|ά)", "(ι|ϊ)", "(α|ας|ο|ου|'ως')", "$"}), "$1$2-$3$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βαΐα&dq=

	// Βάγια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α|ά)", "(γ)", "(ι)", "(α|ας)", "$"}), "$1$2$3$4>$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βαΐα&dq=

	// Βαΐα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α)", "(ί)", "(α|ας)", "$"}), "$1$2-$3-$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βαΐα&dq=

	// Βαϊκάλη (Βα-ικάλη)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαικ*&dq=
	// todo: .αικ

	// Βαϊμάρη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(ι)", "(μ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Βαϊμ*&dq=

	// βακούφι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ου|ού)", "(φ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ακούφι*&dq=

	// βαλανίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(ν)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λανίδι*&dq=

	// Βάλια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=Βάλι*&dq=

	// Βαλτέτσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ε|έ)", "(τσ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τέτσι*&dq=

	// βαλτοτόπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(τ)", "(ο)", "(τ)", "(ο|ό)", "(π)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λτοτόπι*&dq=

	// ανθρακένιο (ανθρακένι-ο)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(θ)", "(ρ)", "(α)", "(κ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9$10>-<$11$12"},

	// βαμβακένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ακένι*&dq=

	// βαμβάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(β)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μβάκι*&dq=

	// βαμπάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*μπάκι*&dq=

	// βαπόρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(π)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαπόρι*&dq=

	// βαπορίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ρ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ορίσι*&dq=

	// βαράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=30&lq=*αράκι*&dq=

	// βαρβαριστί
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αριστί*&dq=

	// βαρβατιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(β)", "(α|ά)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρβατι*&dq=

	// βαρδιάνος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(δ)", "(ι)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρδιάν*&dq=

	// Βαρδούσια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(δ)", "(ου|ού)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρδούσι*&dq=

	// βαρέλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(ρ)", "(ε|έ)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βαρέλι*&dq=

	// βαρελίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ε)", "(λ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρελίσι*&dq=

	// βάριο (βάρι-ο)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α)", "(ρ)", "(ι)", "(ο)", "$"}), "$1$2-$3$4-$5"},
	// βαριούχος (βαρι-ού-χος)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α)", "(ρ)", "(ι)", "(ου|ού)", "(χ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// βαριά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=βαρι*&dq=

	// βαρκάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(κ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ρκάρι*&dq=

	// βαρυγγώμια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(γγ|γκ)", "(ω|ώ)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υγγώμι*&dq=

	// βαρυστομαχιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ο)", "(μ)", "(α|ά)", "(χ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*στομαχι*&dq=
	// todo: στομαχάκι

	// βαρυχειμωνιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ει|ι)", "(μ)", "(ω|ώ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χειμωνι*&dq=

	// Βάσια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(α|ά)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βασι*&dq=

	// βασιλιάς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(α)", "(σ)", "(ι)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=30&lq=*βασιλι*&dq=

	// βατομουριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(μ)", "(ου)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ομουρι*&dq=

	// βατραχάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχάκι*&dq=

	// βατράχι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ρ)", "(α)", "(χ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τράχι*&dq=

	// βατραχοειδής (βατραχο-ι-δης)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ρ)", "(α)", "(χ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τραχοιδ*&dq=

	// βατσινιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τσ)", "(ι|ί)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τσινι*&dq=

	// βαφτιστήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π|φ)", "(τ)", "(ι)", "(στ)", "(η|ή)", "(ρ)", "(ι)", "(ι|ια|ου|ού|ων|ών)", "$"}),
		"$1$2$3$4$5$6$7$8>$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φτιστήρι*&dq=

	// βδομαδιάτικος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(δ)", "(ι)", "(ά)", "(τ)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4><$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αδιατικ*&dq=

	// βεγόνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(γ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εγόνι*&dq=

	// βελόνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ελόνι*&dq=

	// βελουδένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(δ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουδένι*&dq=

	// βελούχι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(ου|ού)", "(χ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελούχι*&dq=
	// todo: βελούκια ??

	// βενετιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ε)", "(τ|τσ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ενετι*&dq=

	// βεντάγια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ντ)", "(α|ά)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εντάγι*&dq=

	// βεντάλια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ντ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εντάγι*&dq=

	// βεργί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ε)", "(ρ)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βεργί*&dq=

	// βερεσέδια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(σ)", "(ε|έ)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εσέδι*&dq=

	// βερικοκιά, βερυκοκιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(κ)", "(ο|ό)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ικοκι*&dq=

	// βερμιγιόν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(υ|ύ|ι|ί)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μιγι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μυγι*&dq=

	// βερνίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ν)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρνίκι*&dq=

	// βερσιόν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ε)", "(ρ)", "(σ)", "(ι)", "(ο|ό)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βερσι*&dq=

	// βετούλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ου|ού)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*τούλι*&dq=

	// βια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ι)", "(α|ά|ας|άς|ες|ές|'ων'|'ών')", "$"}), "$1$2$3"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βια&dq=

	// έβιαζα (έ-βια-ζα), έβιασα (έ-βια-σα)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε|έ)", "(β)", "(ι)", "(α)", "(ζ|σ)", "(α|αν|ε|ες)", "$"}), "$1$2$3$4><$5$6$7"},

	// εβιασμένος (ε-βια-σμένος)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(β)", "(ι)", "(α)", "(σμ)", "(ε|έ)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// todo: αλισίβιασμα http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιασμ*&dq=

	// βιάζε, βιάσε
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ι)", "(α|ά)", "(ζ|σ)", "(ε)", "$"}), "$1$2><$3$4$5"},

	// βιάζω, βιάσω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ι)", "(βιά-ζω)", "$"}), "$1$<><$3$6"},

	// βιάση
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ι)", "(-ιάση)", "$"}), "$1$2><$3$4"},

	// εκβιαστικά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(κ)", "(β)", "(ι)", "(α)", "(στ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// βιαστηκά, βιαστικά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(α)", "(στ)", "(ι)", "(κ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιαστικ*&dq=

	// βιασύνη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(α)", "(σ)", "(υ|ύ|ι|ί)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιασυν*&dq=

	// βιβάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(β)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιβάρι*&dq=

	// βιβλιαράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(λ)", "(ι)", "(α)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7$8$9><$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βλιαράκι*&dq=

	// βιγόνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(γ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιγόνι*&dq=

	// βιδωτήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(δ)", "(ω)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(κ)", "(λ)", "(ει|ι)", "(δ)", "(ω)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
	// "$1$2$3$4$5$6$7$8$9$10><$11$12"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κλειδωτήρι*&dq=
	// todo: 'ω'τήρια*
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δοτήρι*&dq=

	// βιλαέτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(ε|έ)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λαέτι*&dq=

	// βιλάι (βαλα-ι)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(λ)", "(α)", "(ι)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιλάι*&dq=

	// Βίλλια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι|ί)", "(λλ|λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βίλι*&dq=

	// βινιέτα, ινιέστα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(ι)", "(ε|έ)", "(στ|τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινιέτ*&dq=

	// βιντεογκέιμ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(γγ|γκ)", "(ε)", "(ι)", "(μ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// todo: γκέιμ, γκέιμερ
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ογκέιμ*&dq=

	// βιντεοπαιχνίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ν)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χνίδι*&dq=

	// βίντσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι|ί)", "(ν)", "(σ|τσ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ίντσι*&dq=

	// βιοηθική
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο)", "(ι)", "(θ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιοηθ*&dq=

	// βιοηλεκτρισμός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο)", "(ι)", "(λ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιοιλ*&dq=

	// βιοϊατρικά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο)", "(ι)", "(α)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιοϊα*&dq=

	// βιολογία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο)", "(λ)", "(ο|ό)", "(γ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// βιόλυση (βι-όλυση)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο|ό)", "(λ)", "(υ|ι)", "(σ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*βιολογ*&dq=
	// βιολιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο|ό)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// βιολάκια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο|ό)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8<$9$10"},
	// βιόλα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο|ό)", "(λ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=30&lq=*βιολ*&dq=

	// βιος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ι)", "(ος|ός|ους|ούς)", "$"}), "$1$2$3"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιος&dq=

	// βιοϋλικό
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι)", "(ο)", "(υ|ι)", "(λ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βιοϋλ*&dq=

	// Βισκαϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(κ)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισκαϊκ*&dq=
	// todo: ska-i

	// Βίτσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ι|ί)", "(τσ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*Βίτσι*&dq=

	// βλαστήμια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(στ)", "(η|ή)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αστήμι*&dq=

	// Βλαχία (Βλαχι-α)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(Β)", "(λ)", "(α)", "(χ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// βλαχιά (βλα-χιά)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(λ)", "(α)", "(χ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βλαχι*&dq=

	// βλεννορροϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ρρ|ρ)", "(ο)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ορροϊκ*&dq=

	// βογιάρικος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(γ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ογιάρ*&dq=

	// βόδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο|ό)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*βόδι*&dq=

	// βόειος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο)", "(ι)", "(ος|α|ο 2)", "$"}), "$1$2$3>-$4-$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βοιο*&dq=

	// βοήθεια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ο)", "(η|ή)", "(θ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2-$3-$4$5<$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο)", "(ι|ί)", "(θ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βοιθ*&dq=

	// Βοημή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο)", "(ι)", "(μ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βοιμ*&dq=

	// Βοϊβοβίνα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο)", "(ι)", "(β)", "(ο)", "(.*)"}), "$1$2$3>-<$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Βοϊβο*

	// βοϊδάμαξα, βοϊδήσιος, βόιδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο|ό)", "(ι|ϊ)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο)", "(ι|ϊ)", "(δ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8<$9$10"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο|ό)", "(ι|ϊ)", "(δ)", "(.*)"}), "$1$2$3>$4-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βοϊδ*&dq=

	// βολάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ολάκι*&dq=

	// βόλεϊ (βόλε-ι)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ό)", "(λλ|λ)", "(ε)", "(υ|ι)", "$"}), "$1$2$3$4>-$5"},
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ο|ό)", "(λ)", "(ε)", "(υ|ι)", "$"}), "$1$2$3$4>-$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βόλω&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βόλεϊ&dq=

	// βολεϊμπολίστας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ο)", "(λ)", "(ε)", "(υ|ι)", "(μπ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βολεϊμπ*&dq=

	// βολεψάκιας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε|η|ι)", "(ψ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εψάκι*&dq=

	// βόλι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ο|ό)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βόλι*&dq=

	// Βομβάη
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ο)", "(μ)", "(β)", "(α|ά)", "(ι)", "$"}), "$1$2$3$4$5>-<$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Βομβάη*&dq=

	// Βονιτσιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(τσ)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιτσιώτ*&dq=

	// βοριαδάκι, βοριάς
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ο)", "(ρ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βορια*&dq=

	// βοτάνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(τ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*οτάνι*&dq=

	// βοτρυοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(υ|ι)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ριοειδ*&dq=

	// δαμαλάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(α)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// *αμαλάκι*

	// δαμάλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(μ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαμάλι*&dq=

	// μαλάκιο (μαλάκι-ο)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(ο|ου)", "$"}), "$1$2$3$4$5$6$7>-$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(λ)", "(ά)", "(κ)", "(ι)", "(α|ων)", "$"}), "$1$2$3$4$5$6$7>-$8"},
	// βοτσαλάκι (βοτσαλά-κια, βοτσαλα-κιού, βοτσαλα-κιών)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αλάκι*&dq=

	// βουβάλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(β)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ουβάλι*&dq=

	// βουβαλίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(υ|ύ|ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλίσι*&dq=
	// todo: αλυσιδιάζω, αλυσιδιασμένος κλπ.

	// βουζούνι
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ζ)", "(ου|ού)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζουνι*&dq=

	// βουκαμβίλια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(β)", "(ι|ί)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμβίλι*&dq=

	// βούλιαγμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(λ)", "(ι)", "(α)", "(γ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ούλιαγ*&dq=

	// βουλιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ου|ού)", "(λ)", "(ι)", "(α|ά)", "(ζ|σ|σμ|στ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βουλιάζω&dq=
	// todo: αβούλιαχτος

	// βουλωτήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(λλ|λ)", "(ω)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουλλωτηρι*&dq=

	// βουνάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουνάκι*&dq=

	// Iούνιος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι)", "(ου|ού)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4>-<$5$6"},
	// ουνία
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ου|ού)", "(ν)", "(ι)", "(α|ας|ες|ων)", "(.*)"}), "$1$2$3>-<$4$5"},
	// βουνί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=200&lq=*ουνι*&dq=
	// todo: μαυροβούνι-ε

	// βουνίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(ν)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουνίσι*&dq=

	// βουνοπλαγιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ο)", "(π)", "(λ)", "(α)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νοπλαγι*&dq=

	// βουρκοτόπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ο)", "(τ)", "(ο|ό)", "(π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κοτοπι*&dq=

	// κονσόρτσιουμ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(ο|ό)", "(ρ)", "(τσ)", "(ι)", "(ου)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σόρτσι*&dq=
	// βουρτσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(τσ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ρτσι*&dq=

	// βουτιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ου|ού)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βουτι*&dq=

	// βουτσί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(τσ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=80&lq=*ουτσι*&dq=
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*ουτσι*&dq=
	//  αποταχιούτσικα
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιουτσι*&dq=

	// βουτυράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(υ|ι)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τυράκι*&dq=

	// βουτυρένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(υ|ι)", "(ρ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τυρενι*&dq=

	// βουτυριέρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(ρ)", "(ι)", "(ε|έ)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υριέρ*&dq=

	// βραγιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(α)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βραγι*&dq=

	// βραδιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(α|ά)", "(δ)", "(υ|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*βραδι*&dq=

	// βραδυκαής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(α)", "(δ)", "(υ|ι)", "(κ)", "(α)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6$7$8>-<$9$10"},

	// Βραζιλιάνα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ζ)", "(ι)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αζιλι*&dq=

	// Βραΐλα, Μπραΐλα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β|μπ)", "(ρ)", "(α)", "(ι)", "(λ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},

	// βρακάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(α)", "(κ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρακάκι*&dq=

	// βρακί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ου|ού|ων|ών)", "$"}), "$1$2$3$4$5$6>$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βρακί&dq=

	// βρασιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(α|ά)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βρασι*&dq=

	// βράχια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ρ)", "(α|ά)", "(χ)", "(ι)", "(α|ά)", "$"}), "$1$2$3$4$5>$6"},
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βράχι*&dq=
	//  βράχια

	// βραχιόλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(ι)", "(ο|ό)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4>$5-$6$7<$8$9"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(ι)", "(ο|ό)", "(λ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχιόλι*&dq=

	// βραχνιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(α|ά)", "(χ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*ραχνι*&dq=
	// todo: remove αραχνιάζω rule

	// βραχονήσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ο)", "(ν)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χονησι*&dq=

	// βραχοτόπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ο)", "(τ)", "(ο|ό)", "(π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χοτοπι*&dq=

	// βρετίκια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(τ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ετίκι*&dq=

	// βρεχτάδια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(τ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*χτάδι*&dq=

	// βρισίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(υ|ι)", "(σ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// βρισιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(υ|ι)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βρισι*&dq=

	// βρομιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(ω)", "(μ)", "(ι)", "(α|ά|ας|άς|ές|ών)", "$"}), "$1$2$3$4$5$6>$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*βρομι*&dq=

	// βρομιάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ρ)", "(ω)", "(μ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βρομιάρ*&dq=
	// todo: βρόμιασα, βρόμιαζε

	// βρυχιέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(β)", "(ρ)", "(υ|ι)", "(χ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=βρυχι*&dq=
	// todo: *βρυχι*

	// βυζάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(ζ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υζάκι*&dq=

	// βυζανιάρικο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ι)", "(α|ά)", "(ρ)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4><$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανιαρικ*&dq=

	// γιουροβίζιον
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου)", "(ρ)", "(ο)", "(β)", "(υ|ύ|ι|ί)", "(ζ)", "(ι)", "(ο)", "(ν)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8-$9$10-<$11$12$13"},
	// todo: γιούρο
	// βυζί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(υ|ι)", "(ζ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*βιζί*&dq=

	// βυσσινιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(σσ|σ)", "(ι)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υσσινι*&dq=

	// γαβιάλης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(β)", "(ι)", "(α|ά)", "(λ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αβιάλ*&dq=

	// γαγγραινιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γγ|γκ)", "(ρ)", "(αι|ε)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γγραινι*&dq=

	// γαζί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α|ά)", "(ζ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γαζί*&dq=
	// todo: μπαγκάζια

	// γαϊδάρα, γάιδαρος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α|ά)", "(ι|ϊ)", "(δ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=30&lq=*γαϊδ*&dq=

	// γαϊδουράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ου)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουράκι*&dq=

	// γαϊδουριά, γαϊδούρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α|ά|ου|ού)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*δουρι*&dq=

	// γαϊδουρίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α|ου)", "(ρ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δουρισι*&dq=

	// γαϊτάνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(ι|ϊ)", "(τ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8<$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γαϊτανι*&dq=

	// γαϊτανοφρυδάτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(ϊ)", "(τ)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γαϊταν*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αϊτα*&dq=

	// γαλάζιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(λ)", "(α|ά)", "(ζ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γαλάζι*&dq=

	// Γαλαξίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ξ)", "(ει|εί|ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αξιδει*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αξιδι*&dq=

	// γαλάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(λ)", "(α|ά)", "(ρ)", "(ι)", "(α|ου|ού|ων|ών)", "$"}), "$1$2$3$4$5$6$7>$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γαλάρι*&dq=
	// todo: book example

	// Γαλάτσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α|ά)", "(τσ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλάτσι*&dq=

	// γαλί
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=γαλί*&dq=

	// γαλιφιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ι)", "(φ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλιφι*&dq=

	// γαμήσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμήσι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμισι*&dq=

	// γαμιόλα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(ι)", "(ο|ό)", "(λ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμιολ*&dq=

	// γαμιάς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(ι)", "(α|ά)", "(δ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(ι)", "(ά)", "(.*)"}), "$1$2$3$4><$5$6"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(μ)", "(ι)", "(έμαι)", "$"}), "$1$2$3$4$5><$6"},
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(α)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γαμι*&dq=

	// γάμπια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α|ά)", "(μπ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γάμπι*&dq=

	// γανιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α|ά)", "(ν)", "(ι)", "(α|ά)", "(ζ|σ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// todo: check->
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γανι*&dq=

	// γαντάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(ντ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντάκι*&dq=

	// γαντζάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ντ)", "(ζ|σ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ζ|σ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(ζ|σ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανζακι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανσακι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ατζάκι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ατσάκι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντζακι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντσακι*&dq=

	// γάντι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(α|ά)", "(ντ)", "(ι)", "(α|ου|ού|ων|ών)", "$"}), "$1$2$3$4>$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γάντι*&dq=

	// γαρδέλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(δ)", "(ε|έ)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρδέλι*&dq=

	// γαρδένια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(δ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρδένι*&dq=

	// γαριάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(ρ)", "(ι)", "(α|ά)", "(ά-ζω)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αριάζω*&dq=
	// todo: ιασμένος -> http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αριάζω*&dq=

	// γάριασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α|ά)", "(ρ)", "(ι)", "(α|ά)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*άριασμ*&dq=
	// todo: -αριας

	// γαριφαλιά, γαρουφαλιά, γαρυφαλλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|υ|ι)", "(φ)", "(α)", "(λλ|λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιφαλι*&dq=

	// γατάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ατάκι*&dq=
	// todo: γιατάκι

	// γατί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(α)", "(τ)", "(ι|ί)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*γατι*&dq=

	// γατήσιος, γατίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ατισι*&dq=

	// γαυράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(υ|β)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αυράκι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αβράκι*&dq=

	// γαύριασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(υ|ύ|β)", "(ρ)", "(ι)", "(α)", "(σμ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αύριασμ*&dq=

	// γδικιέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(δ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γδικι*&dq=
	// todo: remove αγδίκιωτος rule?

	// γειτονιάρχης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ει|ι)", "(τ)", "(ο)", "(ν)", "(ι)", "(α|ά)", "(ρ)", "(χ)", "(.*)"}),
		"$1$2$3$4$5$6$7>-<$8$9$10$11"},
	// γειτονιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ει|ι)", "(τ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γειτονι*&dq=

	// γελάδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ελάδι*&dq=

	// γελεκάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ε)", "(κ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λεκάκι*&dq=

	// γέλιο
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ε|έ)", "(λ)", "(ι)", "(α|ο|ου|ων)", "$"}), "$1$2$3$4>$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γέλι*&dq=

	// γεμοφεγγαριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(γγ|γκ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εγγαρι*&dq=
	// todo: single greeklish "g"?

	// Δερβενακίων
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(Δ)", "(ε)", "(ρ)", "(β)", "(ε)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(ε|ω)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9>-<$10$11"},
	// γενάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ενακι*&dq=

	// γεναριάτικος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(ι)", "(α|ά)", "(τ)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4><$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αριάτικ*&dq=

	// γενειάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ει|ι)", "(α|ά)", "(ά-ζω)", "$"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ενειάδ*&dq=

	// γενειάδα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ε)", "(ν)", "(ει|ι)", "(α|ά)", "(δ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γενειά*&dq=

	// γένι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ε|έ)", "(ν)", "(ει|ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=γενι*&dq=

	// γενιά, γεννιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ε|έ)", "(νν|ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γενια&dq=

	// γεννησιμιό
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(σ)", "(ι)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ησιμι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισιμι*&dq=

	// γεννητούρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(τ)", "(ου|ού)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ητούρι*&dq=

	// γεννοφάσκια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(α|ά)", "(σ)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φασκι*&dq=

	// γεντιανή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ε|έ)", "(ντ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γεντι*&dq=

	// Ιθακίσιος (Ιθακίσι-ος)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(Ι)", "(θ)", "(α)", "(κ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7>-<$8$9"},
	// γερακίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κίσι*&dq=

	// γεράνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ε)", "(ρ)", "(α|ά)", "(ν)", "(ι)", "(α|ου|ού|ων|ών)", "$"}), "$1$2$3$4$5$6$7>$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εράνι*&dq=

	// γερατειά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ε|η|ι)", "(ρ)", "(α|ά)", "(τ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γερατει*&dq=

	// γερμανοτσολιάς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τσ)", "(ο|ό)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τσολι*&dq=

	// διαολιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α|ά)", "(ο)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2>$3-$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διαολι*&dq=

	// γεροδιάολος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α|ά)", "(ο)", "(λ)", "(.*)"}), "$1$2><$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιαολ*&dq=
	// todo: διαολιά, διαόλια

	// κοντάκιο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ο)", "(ντ)", "(α|ά)", "(κ)", "(ι)", "(ο)", "$"}), "$1$2$3$4$5$6$7>-<$8"},

	// κοντακιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ο)", "(ντ)", "(α)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},

	// γεροντάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ο)", "(ντ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ροντάκι*&dq=
	// todo: κοντάκι-ο, κοντα-κιά

	// γερόντι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(ο|ό)", "(ντ)", "(ι)", "(α|ων|ών)", "$"}), "$1$2$3$4$5$6>$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(ο)", "(ντ)", "(ι)", "(ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ερόντι*&dq=

	// γεροντοκοριάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(κ)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*οκορι*&dq=

	// γεροφτιαγμένος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(τ)", "(ει|ι)", "(α)", "(γ)", "(μ)", "(.*)"}), "$1$2$3$4><$5$6$7$8"},
	// todo: φτιάχνω, φτιάχνω
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φτιαγμ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φτειαγμ*&dq=

	// γέτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ε|έ)", "(τ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γέτι*&dq=

	// γεφυράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(υ|ι)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// todo: γιοφυράκι
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φυράκι*&dq=

	// Ζαφείριος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "('Ζ')", "(α)", "(φ)", "(ει|εί|ι|ί)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7>-<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αφειρι*&dq=
	// γεφύρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(υ|ύ|ι|ί)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εφύρι*&dq=

	// γεωοικονομία, γεωηλεκτρικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ε)", "(ο)", "(ι)", "(.*)"}), "$1$2$3$4>-<$5$6"},

	// γηπεδάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(δ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πεδάκι*&dq=

	// γητειά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(η|ή)", "(τ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γιτι*&dq=

	// γεια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ει|ι)", "(α|ά)", "$"}), "$1$2$3"},
	// []string{"^", "(γ)", "(ει|ι)", "(α|ά)", "(.*)"}), "$1$2<$3$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γεια*&dq=

	// γιαβάς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(β)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαβ*&dq=

	// γιαγιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαγι*&dq=

	// γιαίνω, γιάνω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ει|ι)", "(αι|αί|α|ά|ε|έ)", "(ν)", "(.*)"}), "$1$2<$3$4$5"},
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε|έ)", "(γ)", "(ει|ι)", "(αι|αί|α|ά|ε|έ)", "(ν)", "(.*)"}), "$1-$2$3<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαίν*&dq=

	// γιακάς
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ι)", "(α)", "(κ)", "(α|ά)", "(.*)"}), "$1$2<$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γιακ*&dq=

	// Γιακουμής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(κ)", "(ου|ού)", "(μ)", "(.*)"}), "$1$2><$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιακούμ*&dq=

	// γιαλός
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ι)", "(α|ά)", "(λ)", "(.*)"}), "$1$2<$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαλ*&dq=

	// Γιάλτα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(λ)", "(τ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιάλτ*&dq=

	// γιάμπολη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιάμ*&dq=

	// Γιάνκης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α|ά)", "(ν)", "(κ)", "(.*)"}), "$1$2><$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιάνκ*&dq=

	// Γιαννιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α)", "(ν)", "(ν)", "(ι)", "(ω|ώ)", "(.*)"}), "$1$2$3>$4$5-$6$7<$8$9"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α)", "(ν)", "(ι)", "(ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαννι*&dq=

	// Γιάννενα, Γιάννινα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(νν|ν)", "(ε|η|ι)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},

	// Γιαννινιώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(νν|ν)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινιώτ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινιότ*&dq=

	// Γιάννης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(νν|ν)", "(η|ή)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαννι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γιανι*&dq=

	// γιάντες
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(ντ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαντ*&dq=

	// γιαούρτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(ου|ού)", "(ρ)", "(.*)"}), "$1$2><$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ιαούρ*&dq=
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(ρ)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ούρτι*&dq=

	// γιαουρτάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(ρ)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουρτάκι*&dq=

	// γιαπί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(π)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαπι*&dq=

	// γιαπράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(π)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*απράκι*&dq=

	// γιάπης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α|ά)", "(π)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαπ*&dq=

	// Γιαπωνέζα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(π)", "(ω)", "(ν)", "(.*)"}), "$1$2><$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιαπων*&dq=

	// γιαραμπής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α)", "(ρ)", "(α)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαρα*&dq=

	// γιάρδα, γυάρδα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(α|ά)", "(ρ)", "(δ)", "(.*)"}), "$1$2><$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιαρδ*&dq=

	// γιαρμάς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α)", "(ρ)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιαρμ*&dq=

	// γιασεμάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(μ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εμάκι*&dq=

	// γιασεμί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(σ)", "(ε|έ)", "(μ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2>$3-$4$5-$6$7<$8$9"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α)", "(σ)", "(ε|έ)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιασεμ*&dq=

	// γιασεμένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(ε)", "(μ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σεμένι*&dq=

	// γιασμάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(σ)", "(μ)", "(α|ά)", "(κ)", "(.*)"}), "$1$2><$3$4$5$6$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(μ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σμάκι*&dq=

	// γιαταγάνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(α)", "(τ)", "(α|ά)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιατα*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αγάνι*&dq=

	// γιατί
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ι)", "(α)", "(τ)", "(ι|ί)", "$"}), "$1$2><$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γιατί*&dq=

	// γιατροσόφι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ο)", "(σ)", "(ο|ό)", "(φ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ροσόφι*&dq=

	// γιάφκα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α|ά)", "(φ)", "(κ)", "(.*)"}), "$1$2><$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιάφκ*&dq=

	// γιαχνί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α|ά)", "(χ)", "(ν)", "(.*)"}), "$1$2><$3$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιαχνι*&dq=

	// για*
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2><$3$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=για*&dq=

	// γιγαντένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντενι*&dq=

	// γιδίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(δ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιδίσι*&dq=

	// γιδοτόπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(τ)", "(ο|ό)", "(π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οτόπι*&dq=

	// γιεγιές
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ε)", "(γ)", "(ι)", "(ε|έ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιεγιέ*&dq=

	// γιεν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ε|έ)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6"},

	// γιες
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ι)", "(ε|έ)", "(σ)", "$"}), "$1$2$3$4"},

	// γιέσμαν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ε|έ)", "(σ)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6$7"},

	// γινάτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(α|ά)", "(τ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινάτι*&dq=

	// γιογιό
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(γ)", "(ι)", "(ο|ό)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιογιό*&dq=

	// γιόγκα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο|ό)", "(γγ|γκ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιογκ*&dq=

	// γιοκ
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ι)", "(ο|ό)", "(κ)", "$"}), "$1$2$3$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γιοκ*&dq=
	// todo: test πλαγι-οκόπηση, αγι-οκατάταξη

	// γιόκας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο|ό)", "(κ)", "(α|ας)", "$"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιοκας&dq=

	// γιοματάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(μ)", "(α)", "(τ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οματάρι*&dq=

	// γιόμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο|ό)", "(μ)", "(α|ά)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιόμα*&dq=

	// γιομίζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο|ό)", "(μ)", "(ι|ί)", "(ζ|σ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιομί*&dq=

	// γιορντάνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(ρ)", "(ντ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ρ)", "(ντ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ορντάνι*&dq=

	// γιορτάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο|ό)", "(ρ)", "(τ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*γιορτ*&dq=

	// γιορτάσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(τ)", "(α|ά)", "(σ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρτάσι*&dq=

	// γιορτιάτικος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(τ)", "(ι)", "(α|ά)", "(τ)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4><$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρτιατικ*&dq=

	// γιός
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ι)", "(ι-ός)", "$"}), "$1$2$3"},
	// todo: καλογιός, παραγιός, ψυχογιός, μοναχογιός
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*γιος&dq=

	// γιοτ, γιωτ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ω|ώ)", "(τ)", "$"}), "$1$2$3>$4$5"},

	// γιουβαρελάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου)", "(β)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιουβ*&dq=

	// γιουβαρλάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρλάκι*&dq=

	// γκιουβέτσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ι)", "(ου)", "(β)", "(.*)"}), "$1$2$3><$4$5$6"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(β)", "(ε|έ)", "(τσ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκιουβ*&dq=

	// γιούκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ου|ού)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2>$3-$4$5<$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ου|ού)", "(κ)", "(.*)"}), "$1$2><$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιουκ*&dq=

	// γιουκαλέλι, γιουκαλίλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ε|έ|ι|ί)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλίλι*&dq=

	// γιούλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},

	// Γιούλης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(λ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιουλ*&dq=

	// Γιούνης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6"},

	//γιούρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιούρι*&dq=

	// γιούργια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(ρ)", "(γ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3>$4$5-$6$7<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιούργι*&dq=
	// todo: test δημιουργια

	// γιούρο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(ρ)", "(.*)"}), "$1$2$3><$4$5$6"},

	// γιουρούσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(ρ)", "(ου|ού)", "(σ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουρουσι*&dq=

	// γιουσουρούμ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(σ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιουσ*&dq=

	// γιούσουρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(σ)", "(ου|ού)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουσούρι*&dq=
	// todo: οσούρι (ανεμοσούρι)

	// ομολογιούχα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(γ)", "(ι)", "(ου|ού)", "(χ)", "(α)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8"},
	// γιουχαΐζω (γιουχα-ίζω)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(χ)", "(α)", "(ι|ί)", "(.*)"}), "$1$2$3>$4-$5$6-<$7$8"},
	// γιούχα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ου|ού)", "(χ)", "(α|ά)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιουχα*&dq=

	// γιοφύλλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(φ)", "(υ|ύ|ι|ί)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8<$9$10"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(φ)", "(υ|ύ|ι|ί)", "(λ)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6$7-$8$9<$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιοφυλλι*&dq=

	// γιοφύρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ο)", "(φ)", "(υ|ύ|ι|ί)", "(ρ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},

	// Γιώργος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ω|ώ)", "(ρ)", "(γ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιωργ*&dq=

	// Γιώτα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "('ω|ώ')", "(τ)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// GrhyphRule{customRegexpCompile(
	// []string{"^", "(γ)", "(ι)", "(ω|ώ)", "(τ)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2><$3$4$5$6"},

	// γιωτακισμός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ω|ώ)", "(τ)", "(α|ά)", "(κ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιωτακ*&dq=

	// γιωταχί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ι)", "(ω|ώ)", "(τ)", "(α)", "(χ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γιωταχ*&dq=

	// ^γιο, ^γιου, ^γιω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ι)", "(ω|ώ)", "(.*)"}), "$1$2<$3$4"},

	// γκαβούλιακας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(λ)", "(ι)", "(α|ά)", "(κ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ούλιακ*&dq=

	// γκάζι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(α|ά)", "(ζ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκάζι*&dq=

	// δολοφονία (δολοφονια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ο)", "(φ)", "(ο|ό)", "(ν)", "(ι)", "(α|ε|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7>-<$8$9"},
	// γκαζοφονιάς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φονι*&dq=
	// todo: ιπποφονία, μιαιφονία (μόλυνση τών χεριών που οφείλεται σε μιαρό φόνο) ?

	// γκάιντα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γκ)", "(α)", "(ι)", "(ντ)", "(.*)"}), "$1$2-<$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκαιντ*&dq=

	// γκαϊντατζής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(α)", "(ι)", "(ντ)", "(α)", "(τζ|τσ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},

	// Γκαϊτατζής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(α)", "(ι)", "(τ)", "(α)", "(τζ|τσ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},

	// μαγκάνιο (λανθ. του μαγγάνιο)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(γκ)", "(α|ά)", "(ν)", "(ι)", "(ο)", "$"}), "$1$2$3$4$5$6$7>-$8"},

	// γκανιάν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκανι*&dq=

	// γκαντεμιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ντ)", "(ε|έ)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντεμι*&dq=

	// γκαράζι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(α|ά)", "(ζ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αράζι*&dq=

	// γκαρδιακός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(α)", "(ρ)", "(δ)", "(ι)", "(α|ά)", "(κ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκαρδιακ*&dq=

	// γκαρσόνι, γκαρσονιέρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(σ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρσόνι*&dq=

	// γκέι (γκε-ι)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ε)", "(ι)", "$"}), "$1$2$3>-$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκει&dq=

	// γκέισα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ε)", "(ι)", "(σ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκεισ*&dq=

	// γκέμι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ε|έ)", "(μ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκέμι*&dq=

	// γκεσέμι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(σ)", "(ε|έ)", "(μ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εσέμι*&dq=

	// γκιαούρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ι)", "(α|ά)", "(ου|ού)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκιαού*&dq=
	// todo: "(ι)", "(α)" -> "(ι)", "(α|ά)" where synizesis occurs.

	// γκίνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ι|ί)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκίνι*&dq=

	// γκινισιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// *ινισι*

	// Γκιόνα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ι)", "(ω|ώ)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*Γκιόν*&dq=

	// γκιόσα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ι)", "(ο|ό)", "(σ)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκιόσα*&dq=

	// γκιουλέκας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ι)", "(ου|ού)", "(λ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκιουλ*&dq=

	// γκλαμουριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(ου|ού)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*αμουρι*&dq=

	// γκόλφι, γκόρφι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ο|ό)", "(λ|ρ)", "(φ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκόλφι*&dq=

	// γκολφάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(λ)", "(φ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// *ολφάκια*

	// γκομενιάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ενιάρ*&dq=

	// γκοριτσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ο)", "(ρ)", "(ι)", "(τσ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκοριτσι*&dq=

	// γκρέιντερ (γρε-ιντερ)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ρ)", "(ε)", "(ι)", "(ντ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρέιντ*&dq=
	// todo: +γ, +γγ ...

	// γκρέιπφρουτ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ρ)", "(ε)", "(ι)", "(π)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκρέιπ*&dq=

	// γκρίνια, γρίνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(ρ)", "(ι|ί)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ρ)", "(ι|ί)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γκρίνι*&dq=

	// γκρουπάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(π)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουπάκι*&dq=

	// γλαδίολος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(δ)", "(ι)", "(ο|ό)", "(λ)", "(ου)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8$9"},
	// γλαδιόλα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(δ)", "(ι)", "(ο|ό)", "(λ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λαδιολ*&dq=
	// todo: improve γλαδίολος

	// γλαρόνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρόνι*&dq=

	// γλαροπούλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(π)", "(ου|ού)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*οπούλι*&dq=

	// Μαντζουρία, Ματζουρία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(ν)", "(τζ|τσ)", "(ου|ού)", "(ρ)", "(ι)", "(α)", "(.*)"}),
		"$1$2$3$4$5$6$7$8>-<$9$10"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(τζ|τσ)", "(ου|ού)", "(ρ)", "(ι)", "(α)", "(.*)"}),
		"$1$2$3$4$5$6$7>-<$8$9"},
	// γλειφιτζούρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τζ|τσ)", "(ου|ού)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τζούρι*&dq=

	// γλειφτρόνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(τ)", "(ρ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φτρόνι*&dq=

	// γλεντάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ντ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εντάκι*&dq=

	// γλέντι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(λ)", "(ε|έ)", "(ντ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γλέντι*&dq=

	// γλεντοκόπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ο)", "(κ)", "(ο|ό)", "(π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντοκόπι*&dq=
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οκόπι*&dq=

	// γλινί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(λ)", "(ι|ί)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γλινι*&dq=

	// γλιτσιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι|ί)", "(τσ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// *ιτσι*

	// γλυκάδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(κ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υκάδι*&dq=

	// γλυκάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(κ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λυκάκι*&dq=

	// γλυκόπιοτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο|ό)", "(π)", "(ι)", "(ο)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*όπιοτ*&dq=

	// γλυκούλι :'^) todo: book example
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ου|ού)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*κούλι*&dq=

	// γλωσσάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(λ)", "(ω)", "(σσ|σ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ου|ού|ων|ών)", "$"}),
		"$1$2$3$4$5$6$7$8>$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γλωσσάρι*&dq=

	// γλωσσίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ω)", "(σσ|σ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λωσσίδι*&dq=

	// γλωσσοϋφολογία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ω)", "(σσ|σ)", "(ο)", "(υ|ι)", "(φ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λοσουφ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λωσσοϋφ*&dq=

	// γνεψιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ν)", "(ε)", "(ψ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γνεψι*&dq=

	// γνωριμιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(ρ)", "(ι)", "(μ)", "(ι)", "(ά|έ|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ωριμι*&dq=
	// todo: book, "γνωριμί-α" is more used, but -ών still gets synizesis.

	// Σλοβακία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ο)", "(β)", "(α|ά)", "(κ)", "(ι)", "(α|ε|ω)", "(.*)"}), "$1$2$3$4$5$6$7>-<$8$9"},
	// γοβάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(β)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οβάκι*&dq=

	// γογγύλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(γγ)", "(υ|ύ|ι|ί)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ογγύλι*&dq=

	// γόης
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ο)", "(ι)", "(.*)"}), "$1$2-<$3$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γοι*&dq=

	// γοητευτικά (γο-ιτευτικά)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ο)", "(ι)", "(τ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ογοιτ*&dq=

	// γομάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(μ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ομάρι*&dq=

	// γονατιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ο)", "(ν)", "(α)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=γονάτι*&dq=

	// γονδολιέρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(λ)", "(ι)", "(ε|έ)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ολιέρ*&dq=
	// todo: οδιερ, ζ, ν, τ,

	// γονεϊκός (γονε-ικός)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ν)", "(ε)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ονεικ*&dq=

	// γονιός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ο)", "(ν)", "(ι)", "(ος|ού|ό|έ...)", "$"}), "$1$2$3$4$5>$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=γονι*&dq=

	// γοριλάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι|ϊ)", "(λλ|λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιλάκι*&dq=
	// todo: "(.*)", "(ι)" -> "(.*)", "(ι|ϊ)" ?

	// γουβιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ου|ού)", "(β)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γουβια*&dq=
	// todo: κλούβιασμα etc.

	// γουδί
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ου|ού)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γουδι*&dq=

	// Γουοτεργκέιτ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(γκ)", "(ε)", "(ι)", "(τ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ργκειτ*&dq=

	// γουρσουζιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(ου|ού)", "(ζ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σούζι*&dq=

	// γοφάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ο)", "(φ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γοφάρι*&dq=

	// Γραβιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(α|ά)", "(β)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γραβι*&dq=

	// γραΐδιο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(α)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γραΐδι*&dq=
	// todo: test with "ίδιος", "ολόιδιος" etc.

	// γρανάζι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(α|ά)", "(ζ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αναζι*&dq=

	// γρανιτένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(τ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιτενι*&dq=

	// γρασίδι, γρασσίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(σσ|σ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// todo: investigate *ασίδι (μοιασίδια, χρειασίδια, etc)

	// γραφιάς
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ρ)", "(α)", "(φ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γραφια*&dq=

	// γρέζι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ε|έ)", "(ζ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γρεζι*&dq=

	// γρεναδιέρος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(α)", "(δ)", "(ι)", "(ε|έ)", "(ρ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ναδιέρ*&dq=

	// γριβάδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ι)", "(β)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ριβάδι*&dq=

	// γρίλια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ι|ί)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γρίλι*&dq=

	// γριπιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ι)", "(ππ|π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γριπι*&dq=

	// γροθιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ο)", "(θ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γροθι*&dq=

	// γροθοκοπανιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ο)", "(π)", "(α)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κοπανι*&dq=

	// γρόσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ρ)", "(ο|ό)", "(σ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γρόσι*&dq=

	// γρουμπούλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(μπ)", "(ου|ού)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουμπούλι*&dq=

	// γυαλά, γιάλα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ|ι)", "(α|ά)", "(λ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουμπούλι*&dq=

	// Γυάρος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ|ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=*Γιάρ*&dq=

	// γυιόκας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ)", "(ι)", "(ο|ό)", "(.*)"}), "$1$2$3>$4<$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γυιοκ*&dq=

	// γυλιός
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(υ|ι)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},

	// γύμνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ|ύ|ι|ί)", "(μ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γύμνι*&dq=

	// γυμνοσάλιαγκας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(α|ά)", "(γκ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λιαγκ*&dq=

	// γυναίκειος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(ν)", "(αί|έ)", "(κ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υναίκει*&dq=

	// γυναικολόγια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(αι|ε)", "(κ)", "(ο)", "(λ)", "(ο|ό)", "(γ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αικολόγι*&dq=

	// γυναικολόι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(αι|ε)", "(κ)", "(ο)", "(λ)", "(ο)", "(ι)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},

	// γυναικομάνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(αι|ε)", "(κ)", "(ο)", "(μ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// ẗodo: ομάνι

	// γυροβολιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ|ι)", "(ρ)", "(ο)", "(β)", "(ο|ό)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γυροβολι*&dq=

	// Γυφτάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φτάκι*&dq=

	// γυφταριό
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(τ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φτάρι*&dq=

	// γυφτιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(υ|ι)", "(φ)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γυφτι*&dq=

	// γωβιός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ω)", "(β)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γοβι*&dq=

	// γωνιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ω)", "(ν)", "(ι)", "(α|ας|ες|ων 2)", "$"}), "$1$2-$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=γωνιά*&dq=

	// γωνιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(γ)", "(ω|ώ)", "(ν)", "(ι)", "(α|ά)", "(ά-ζω)", "$"}), "$1$2$3$4><$5$6"},

	// γώνιασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ω|ώ)", "(ν)", "(ι)", "(α|ά)", "(σ)", "(μ)", "(.*)"}), "$1$2$3$4$5><$6$7$8$9"},

	// γωνιαστός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ω|ώ)", "(ν)", "(ι)", "(α|ά)", "(σ)", "(τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8$9"},

	// γραμμοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μμ|μ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμμοειδ*&dq=

	// γριφοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(φ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιφοειδ*&dq=

	// γυψοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(ψ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},

	// γωνιοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ι)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},

	// δαγκαματιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γκ)", "(α|ω)", "(μ)", "(α)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*γκαματι*&dq=

	// δαδί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=40&lq=*δαδι*&dq=

	// δαημοσύνη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(ι)", "(μ)", "(ο)", "(σ)", "(υ|ύ|ι|ί)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαιμοσι*&dq=

	// δαιδαλοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(λ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαλοειδ*&dq=

	// δακτυλάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(υ|ι)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},

	// δακτυλίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(λ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιλιδι*&dq=

	// δακτυλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(κ)", "(τ)", "(υ|ι)", "(λ)", "(ι)", "(ι-άς...)", "$"}), "$1$2$3$4$5$6$7$8>$9"},
	// δαχτυλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(χ)", "(τ)", "(υ|ι)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},

	// δαχτυλιδένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(δ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λιδενι*&dq=

	// δακτυλιοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	//

	//Δαλάι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(λ)", "(α)", "(ι)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	//http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαλάι*&dq=

	// δαμασκηνιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(σ)", "(κ)", "(η)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μασκηνι*&dq=

	// Δαμιανός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(μ)", "(ι)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαμιαν*&dq=

	// Δανάη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(ν)", "(α)", "(ι)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαναι*&dq=

	// δαντελένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ε)", "(λ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντελένι*&dq=

	// Δαρδανέλια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(ν)", "(ε|έ)", "(λ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δανελι*&dq=

	// δασκαλίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(α)", "(λ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*καλίκι*&dq=

	// δασκαλοπαίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ο)", "(π)", "(α)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλοπαίδι*&dq=

	// δασοτόπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(ο)", "(τ)", "(ο|ό)", "(π)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σοτόπι*&dq=

	// δαυλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(υ|β)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαυλί*&dq=

	// Δαφνιάς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(α)", "(φ)", "(ν)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δαφνί*&dq=

	// δαφνοκερασιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(α|ά)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ερασί*&dq=

	// δέηση (de-isi)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ε)", "(ι)", "(σ)", "(ης...)", "$"}), "$1$2-<$3$4$5"},
	// *δεισι*

	// δέησα (δε-ισα)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ε)", "(ι)", "(σσ|σ)", "(α...)", "$"}), "$1$2-<$3$4$5"},
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(δ)", "(ε)", "(ι)", "(σσ|σ)", "(α...)", "$"}), "$1$2$3>-<$4$5$6"},

	// δεητικά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε)", "(ι)", "(τ)", "(ι)", "(κ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δειτικ*&dq=

	// δείλι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ει|εί|ι|ί)", "(λ)", "(ι)", "(α|ου|ού|ων|ών)", "$"}), "$1$2$3$4>$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=δειλι&dq=

	// δειλιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε)", "(ι)", "(λ)", "(ι)", "(α|ά)", "(ά-ζω)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δειλι*&dq=

	// δείλιασμα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ει|εί|ι|ί)", "(λ)", "(ι)", "(α|ά)", "(σ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δειλιασ*&dq=

	// δεκαεννιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ν)", "(ι)", "(α|ά)", "$"}), "$1$2$3$4$5>$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αεννια*&dq=

	// δεκαήμερος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ι|ί)", "(μ)", "(ε)", "(ρ)", "(.*)"}), "$1$2>-<$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αιμερ*&dq=

	// δεκανίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(α)", "(ν)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κανικι*&dq=

	// δεκαπενταριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(α)", "(π)", "(ε)", "(ντ)", "(α)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// todo: ...εξαριά, εφταριά...

	// δεκάρι, δεκαριά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ε)", "(κ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// δεκατιανό
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(κ)", "(α)", "(τ)", "(ι)", "(α)", "(ν)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εκατιαν*&dq=

	// δεκατριάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ρ)", "(ι)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4-$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τριάρι*&dq=
	// todo: book example (quick-synizesis + rules)

	// δελτοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τοιδ*&dq=

	// δελφινάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ι)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φινακι*&dq=

	// δελφινοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ι)", "(ν)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φινοιδ*&dq=

	// δενδροειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ|ντ)", "(ρ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δροιδ*&dq=

	// δενδράκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ρ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δράκι*&dq=

	// δεντρί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε)", "(ντ)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δεντρί*&dq=

	// δεντρογαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ο)", "(γ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρογαλι*&dq=

	// δεξίμι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε)", "(ξ)", "(ι|ί)", "(μ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δεξίμι*&dq=

	// δερβένι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(β)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ερβένι*&dq=

	// δερματένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(τ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ματένι*&dq=

	// δερνοκοπιέμαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ο)", "(π)", "(ι)", "(έμαι)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κοπιέμαι*&dq=

	// δεσιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ε)", "(σ)", "(ι)", "(ι-άς...)", "$"}), "$1$2$3$4>$5"},

	// δεσίδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε)", "(σ)", "(ι|ί)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δεσίδι*&dq=

	// δευτεριάτικα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ε)", "(ρ)", "(ι)", "(α|ά)", "(τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τεριατ*&dq=

	// δεφτέρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(φ)", "(τ)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εφτέρι*&dq=

	// δημοσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(η)", "(μ)", "(ο)", "(σ)", "(ι)", "(α|ά|ας|άς|ες|ές|'ων'|'ών')", "$"}),
		"$1$2$3$4$5$6$7>$8"},

	// δημοσιοϋπαλληλικά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ο)", "(υ|ι)", "(π)", "(α)", "(λ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιουπαλ*&dq=

	// διάβα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α|ά)", "(β)", "(α)", "$"}), "$1$2$3-$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=διαβα&dq=

	// διαβάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α|ά)", "(β)", "(α|ά)", "(ά-ζω)", "$"}), "$1$2$3><$4$5$6$7"},

	// διαβαίνω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α|ά)", "(β)", "(αίνω...)", "$"}), "$1$2$3><$4$5$6"},

	// διάβασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α|ά)", "(β)", "(α)", "(σμ)", "(.*)"}), "$1$2><$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιάβασμ*&dq=

	// διαβαστερός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α|ά)", "(β)", "(α|ά)", "(στ)", "(.*)"}), "$1$2><$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιάβαστ*&dq=

	// διαβατάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(β)", "(α)", "(τ)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2><$3$4$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιαβατάρ*&dq=

	// διαβάτης, διαβάτισσα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(β)", "(α|ά)", "(τ)", "(άτ-ης...)", "$"}), "$1$2><$3$4$5$6$7"},
	// "(άτ-ης...)"

	// διαβατικός
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α)", "(β)", "(α|ά)", "(τ)", "(ι)", "(κ)", "(.*)"}), "$1$2><$3$4$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=διαβατικ*&dq=

	// διαβολή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(β)", "(ο)", "(λ)", "(ολ-ή)", "$"}), "$1$2$3>-<$4$5$6$7$8"},
	// διαβολέας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(β)", "(ο)", "(λ)", "(ολ-έας)", "$"}), "$1$2$3>-<$4$5$6$7$8"},
	// διαβολιά, διαβόλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(β)", "(ο|ό)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6><$7$8"},
	// διαβολάκος, διαβόλισσα, διάβολος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α|ά)", "(β)", "(ο|ό)", "(λ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διαβολ*&dq=

	// διαϊδρυματικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(ι)", "(δ)", "(ρ)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διαιδρ*&dq=

	// διακαής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(κ)", "(α)", "(ισ?)", "$"}), "$1$2$3$4$5$6>-$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διακαι*&dq=

	// διακόνεμα, διακονεύω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(κ)", "(ο|ό)", "(ν)", "(ε|έ)", "(.*)"}), "$1$2$3><$4$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διακόνε*&dq=

	// διάκος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α|ά)", "(κ)", "(ος...)", "$"}), "$1$2$3><$4$5$6"},

	// διακόσια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ο|ό)", "(σ)", "(ι)", "(α|ες|οι|ους)", "$"}),
		"$1$2$3$4$5$6>$7"},
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(δ)", "(ι)", "(α)", "(κ)", "(ο|ό)", "(σ)", "(ι)", "(α|ες|οι|ους)", "$"}),
	// "$1$2$3>$4-$5$6-$7$8$9"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(κ)", "(ο|ό)", "(σ)", "(ι)", "(.*)"}), "$1$2$3>$4-<$5$6$7$8$9"},
	// διακοσάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(κ)", "(ο)", "(σ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8-$9$10<$11$12"},
	// διακοσάρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(κ)", "(ο)", "(σ)", "(α|ά)", "(ρ)", "(.*)"}),
		"$1$2$3>$4-$5$6-<$7$8$9$10"},

	// διαλαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(λ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8<$9$10"},

	// διαλάλημα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(λ)", "(α|ά)", "(λ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},

	// διάλεγμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α|ά)", "(λ)", "(ε|έ)", "(γ)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6$7$8$9"},

	// συνδιαλέγομαι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(δ)", "(ι)", "(α|ά)", "(λ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// νδιαλ

	// διαλέγω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α|ά)", "(λ)", "(-έγω)", "$"}), "$1$2$3><$4$5$6"},
	// todo: book->διαλέξεις

	// διαλεκτός
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α)", "(λ)", "(ε)", "(κ)", "(τ)", "(τ-ός, ή, ό)", "$"}), "$1$2><$3$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διάλεκτ*&dq=

	// διαλεχτός
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α)", "(λ)", "(ε)", "(χ)", "(τ)", "(τ-ός, ή, ό)", "$"}), "$1$2><$3$4$5$6$7$8"},
	// todo: αδιάλεχτος...

	// διάλος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α|ά)", "(λ)", "(ολ-ος)", "$"}), "$1$2$3-$4$5"},
	// "(ολ-ος)"

	// διαλώ
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α)", "(λ)", "(αλ-ώ...)", "$"}), "$1$2$3-<$4$5"},
	// todo: διαλάω, διαλώ http://www.slang.gr/lemma/11630-a

	// διαμαντάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ντ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αντάκι*&dq=

	// διαμάντι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(μ)", "(α|ά)", "(ντ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2>$3-$4$5-$6$7<$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιαμάντι*&dq=

	// διαμαντέ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(μ)", "(α|ά)", "(ντ)", "(.*)"}), "$1$2><$3$4$5$6$7"},

	// διαμερισματάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σμ)", "(α)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σματάκι*&dq=

	// διαμιάς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α)", "(μ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διαμιά*&dq=

	// διάνα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α|ά)", "(ν)", "(α)", "$"}), "$1$2$3-$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διάνα*&dq=

	// διανόημα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(ν)", "(ο)", "(ι)", "(μ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιανοιμ*&dq=

	// διανοήτρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ο)", "(ι|ί)", "(τ)", "(ρ)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανοιτρ*&dq=

	// διάνος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α|ά)", "(ν)", "(ος...)", "$"}), "$1$2$3-$4$5"},

	// διαπνοή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(π)", "(ν)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιαπνοι*&dq=

	// διαρροή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(ρ)", "(ρ)", "(ο)", "(ισ?)", "$"}), "$1$2$3$4$5$6>-$7"},

	// διαρροϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(α)", "(ρρ|ρ)", "(ο)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιαρροϊκ*&dq=

	// διασκελιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(σ)", "(κ)", "(ε)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ασκελι*&dq=

	// διάτα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(α|ά)", "(τ)", "(α|ας|ες|ων 2)", "$"}), "$1$2$3><$4$5$6"},

	// διαφέντεμα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α)", "(φ)", "(ε|έ)", "(ντ)", "(.*)"}), "$1$2<$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διαφεντ*&dq=

	// διάφορο
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(α|ά)", "(φ)", "(ο|ό)", "(ρ)", "(α|ο|ου|ων)", "$"}), "$1$2$3-$4$5-$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*διαφορ*&dq=

	// διβάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(β)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},

	// δικαιικά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(κ)", "(ε)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ικεικ*&dq=

	// δικαϊκά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(κ)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δικαικ*&dq=

	// δικέλλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(κ)", "(ε|έ)", "(λλ|λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δικελλι*&dq=

	// δικέρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(κ)", "(ε|έ)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δικερι*&dq=

	// δίκιο
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι|ί)", "(κ)", "(ι)", "(ος|α|ο)", "$"}), "$1$2-$3$4$5"},

	// δικράνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(κ)", "(ρ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},

	// διμήνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ι)", "(μ)", "(η|ή)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},

	// διόλου
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(ο|ό)", "(λ)", "(ου)", "$"}), "$1$2$3-$4$5"},

	// Διονύσης
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(ο)", "(ν)", "(υ|ύ|ι|ί)", "(σ)", "(ης|η)", "$"}), "$1$2<$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=Διονύσ*&dq=

	// Διονυσάκης
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ι)", "(ο)", "(ν)", "(υ|ι)", "(σ)", "(α|ά)", "(κ)", "(.*)"}), "$1$2<$3$4$5$6$7$8$9"},
	// todo: create rules for kvanto anyway:
	// http://lexilogia.gr/forum/showthread.php?13995-Ο-συλλαβισμός-των-λέξεων&p=249418&viewfull=1#post249418

	// διπλοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(π)", "(λ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιπλοιδ*&dq=

	// διπλοπενιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ο)", "(π)", "(ε)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λοπενι*&dq=

	// διπλοσκοπιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(λ)", "(ο)", "(σ)", "(κ)", "(ο)", "(π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πλοσκοπι*&dq=

	// δισκάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(σ)", "(κ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// δισκάκι

	// δισκοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(σ)", "(κ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισκοιδ*&dq=

	// διφραγκάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γγ|γκ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αγκάκι*&dq=

	// διχόνοια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(χ)", "(ο|ό)", "(ν)", "(οι|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιχόνι*&dq=

	// διχτυωτός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(υ|ύ|ι|ί)", "(χ)", "(τ)", "(υ|ι)", "(ω)", "(τ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// δίχτυ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(υ|ύ|ι|ί)", "(χ)", "(τ)", "(υ|ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δίχτυ*&dq=

	// δοβλέτι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(β)", "(λ)", "(ε|έ)", "(τ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οβλέτι*&dq=

	// Δοϊράνη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ο)", "(ι)", "(ρ)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δοιραν*&dq=

	// δοκάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ο)", "(κ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δοκαρι*&dq=

	// Δολιανά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(λ)", "(ι)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ολιαν*&dq=

	// δόλιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ο|ό)", "(λ)", "(ι)", "(ος|α|ο)", "$"}), "$1$2$3>-$4$5$6"},

	// δοντάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ντ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// δοντιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ο|ό)", "(ντ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2-$3$4<$5$6"},
	// todo: *δοντια*

	// δοξάρι, δοξαριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ο)", "(ξ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},

	// Δραγατσάνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(τσ)", "(α)", "(ν)", "(ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ατσάνι*&dq=

	// δρακοντιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(ρ)", "(α)", "(κ)", "(ο)", "(ντ)", "(ι)", "('ά'|'άς'|'ές'|'ών')", "$"}),
		"$1$2$3$4$5$6$7>$8"},

	// δρακουλίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(λ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουλισι*&dq=

	// δρακουλίνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(λ)", "(ι|ί)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// εμπαλία (εμπαλι-α)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(μπ)", "(α)", "(λ)", "(ι)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εμπαλια*&dq=

	// δράμι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ρ)", "(α|ά)", "(μ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// δρεπάνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(π)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*επανι*&dq=

	// δρεπανοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ν)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ανοιδ*&dq=

	// δρολάπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(λ)", "(α|ά)", "(π)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ολαπι*&dq=

	// δρομάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(μ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ομακι*&dq=

	// δροσιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ρ)", "(ο)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δροσι*&dq=

	// δροσοσταλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*σταλι*&dq=

	// δυάρι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(υ|ι)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(υ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// todo: τεσσάρια, πεντάρια ...

	// δυναμάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(α)", "(μ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ναμάρι*&dq=

	// δυναμοηλεκτρικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ι)", "(λ)", "(ε)", "(κ)", "(.*)"}), "$1$2>-<$3$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οιλεκ*&dq=

	// δυο
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(δ)", "(υ|ι)", "(ο|ό)", "$"}), "$1$2$3"},

	// δυόμισι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(υ|ι)", "(ο|ό)", "(μ)", "(ι)", "(σ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},

	// δυοσμαρίνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(υ|ι)", "(ο|ό)", "(σ)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(ρ)", "(ι|ί)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},

	// δυονών
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(υ|ι)", "(ο)", "(ν)", "(ω|ώ)", "(ν)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// ^δυο.*
	// GrhyphRule{customRegexpCompile(
	// []string{"^", "(δ)", "(υ)", "(ο|ό)", "(.*)"}), "$1$2<$3$4"},

	// δυσνόητα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(υ|ι)", "(σ)", "(ν)", "(ο)", "(ι)", "(τ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υσνοητ*&dq=

	// δυσπνοϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(ν)", "(ο)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πνοικ*&dq=

	// δωματιάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(τ)", "(ι)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ματιακι*&dq=

	// εβδομηνταριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(ντ)", "(α)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ηνταρι*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινταρι*&dq=

	// Εβραίικα (εβρε-ικα)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(β)", "(ρ)", "(ε)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εβρεικ*&dq=

	// Εβραϊκά (εβρα-ικα), εβραϊστής, εβραϊσμός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(β)", "(ρ)", "(α)", "(ι|ί)", "(κ|σμ|στ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},

	// έγια μόλα/λέσα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε|έ)", "(γ)", "(ι)", "(α|ά)", "$"}), "$1-$2$3$4"},

	// έγνοια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε|έ)", "(γ)", "(ν)", "(οι|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// todo: reduce?

	// εγωίστρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(γ)", "(ο)", "(ι|ί)", "(κ|σμ|στ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γοισμ*&dq=
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γοιστ*&dq=
	// todo: remove κ

	// εδεσσαϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(δ)", "(ε)", "(σσ|σ)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*δεσαικ*&dq=

	// ειδαλλιώς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ει|ι)", "(δ)", "(α)", "(λλ|λ)", "(οι|ι)", "(ω|ώ)", "(σ)", "$"}), "$1$2$3$4$5$6><$7$8"},

	// έι
	// GrhyphRule{customRegexpCompile(
	// []string{"^", "(έ)", "(ι)", "$"}), "$1$2"},

	// εικοσαριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(σ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// Ειλείθυια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ει|ι)", "(λ)", "(ει|εί|ι|ί)", "(θ)", "(υ)", "(ι)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// Έιρε
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(ι)", "(ρ)", "(ε)", "$"}), "$1-$2-$3$4"},

	// εισροή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ει|ι)", "(σ)", "(ρρ|ρ)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισροι*&dq=

	// έιτζ
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(ι)", "(τζ)", "$"}), "$1-$2$3"},

	// οστάριο
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ο)", "(στ)", "(α|ά)", "(ρ)", "(ι)", "(α|ε|ω)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	// εκατοσταριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(στ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// εκιού
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(κ)", "(ι)", "(ου|ού)", "$"}), "$1-$2$3$4"},

	// εκκλησιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κκ|κ)", "(λ)", "(η)", "(σ)", "(ι)", "(ι-άς...)", "$"}), "$1$2$3$4$5$6>$7"},

	// εκουαλάιζερ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(α)", "(λ)", "(α)", "(ι)", "(ζ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουαλάιζ*&dq=

	// εκπνοή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(π)", "(ν)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κπνοι*&dq=

	// εκροή
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(κ)", "(ρ)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4>-<$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εκροι*&dq=

	// εκχριστιανίζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(χ)", "(ρ)", "(ι)", "(στ)", "(ι)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3$4$5$6><$7$8$9"},

	// εκχυδαΐζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(χ)", "(υ|ι)", "(δ)", "(α)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κχυδαι*&dq=

	// ελαϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(α)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελαϊκ*&dq=
	// todo: book example πετρελαικοσ, πετρέλαιο

	// ακρελαΐνη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ρ)", "(ε)", "(λ)", "(α)", "(ι|ί)", "(ν)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// ελαϊνή
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(λ)", "(α)", "(ι|ί)", "(ν)", "(.*)"}), "$1$2$3>-<$4$5$6"},

	// ελαιοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(αι|ε)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελαιοιδ*&dq=

	// ελατηριάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(υ|ι)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ριάκι*&dq=

	// ελαφάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(φ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αφάκι*&dq=

	// ελάφι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(α|ά)", "(φ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελαφι*&dq=

	// ελαφοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(α)", "(φ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λαφοειδ*&dq=

	// ελαφροκεφαλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ρ)", "(ο)", "(κ)", "(ε)", "(φ)", "(α)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9$10><$11$12"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φροκεφ*&dq=

	// ελεεινά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(ε)", "(ι)", "(ν)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελειν*&dq=
	// todo: ανελέητος etc.

	// ελεημοσύνη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(ε)", "(ι)", "(μ)", "(ο)", "(σ)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελειμοσ*&dq=

	// ελεήμων
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(ε)", "(ι|ί)", "(μ)", "(ο)", "(ν)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8$9"},

	// ελέησον
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(ε)", "(ι)", "(σ)", "(ο)", "(ν)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8$9"},

	// ελέηση
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(λ)", "(ε)", "(ι|ί)", "(σ)", "(.*)"}), "$1$2$3>-<$4$5$6"},

	// ελεητικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(ε)", "(ι|ί)", "(τ)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελειτικ*&dq=

	// ελευθεροπλοΐα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε)", "(ρ)", "(ο)", "(π)", "(λ)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6$7$8>-<$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θεροπλοΐ*&dq=

	// ελευτεριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ε)", "(υ|φ)", "(τ)", "(ε)", "(ρ)", "(ι)", "(α|ά|ας|άς|ες|ές|'ων'|'ών')", "$"}),
		"$1$2$3$4$5$6$7$8>$9"},

	// ελιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(λ)", "(η)", "('ά'|'άς'|'ές'|'ών')", "$"}), "$1-$2$3$4"},
	// todo: αγρελιάς etc.

	// ελικοειδώς
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(κ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λικοιδ*&dq=

	// ελιόψωμο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(ο|ό)", "(ψ)", "(ω)", "(μ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// todo: ελιό*

	// ελισαβετιανός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ε|έ)", "(τ)", "(ι)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*βετιαν*&dq=

	// ελκοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(κ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελκοιδ*&dq=

	// ελληνορωμαϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ω)", "(μ)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρωμαικ*&dq=

	// ελμινθοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(θ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινθοιδ*&dq=

	// ελυτροειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(τ)", "(ρ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υτροιδ*&dq=

	// εμβοή
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(μ)", "(β)", "(ο)", "(ίς|ί)", "$"}), "$1$2$3$4>-$5"},
	// todo: expand?

	// εμορφιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(μ)", "(ο)", "(ρ)", "(φ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εμορφι*&dq=

	// εμουλσιόν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(λ)", "(σ)", "(ι)", "(ο|ό)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουλσιο*&dq=

	// εμπασιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(μπ)", "(α)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εμπασι*&dq=

	// εμποροϋπαλληλικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(ρ)", "(ο)", "(υ|ι)", "(π)", "(α)", "(λ)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ορουπαλ*&dq=

	// ενδεής
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(ν)", "(δ)", "(ε)", "(ίς|ί)", "$"}), "$1$2$3$4>-$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=ενδει*&dq=

	// εννιαθέσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α)", "(θ)", "(.*)"}), "$1$2$3$4><$5$6$7"},

	// εννιάκιλος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(κ)", "(.*)"}), "$1$2$3$4><$5$6$7"},

	// Βενιαμίν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(β)", "(ε)", "(ν)", "(ι)", "(α)", "(μ)", "(η|ή)", "(ν)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8$9$10"},

	// εννιαμελής, εννιάμηνος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(μ)", "(ε|η|ι)", "(.*)"}), "$1$2$3$4><$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εννιαμη*&dq=

	// εννιαπλάσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(π)", "(λ)", "(.*)"}), "$1$2$3$4><$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εννιαπλ*&dq=

	// εννιάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4>$5-$6$7<$8$9"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ενιάρι*&dq=

	// εννιατάξιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},

	// ενιαύσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ι)", "(α|ά)", "(υ|ύ|φ)", "(σ)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8"},
	// εννιάφυλλος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(φ)", "(.*)"}), "$1$2$3$4><$5$6$7"},

	// ενιαχού
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ι)", "(α)", "(χ)", "(ου)", "(.*)"}), "$1$2$3$4>-<$5$6$7$8"},
	// εννιάχρονα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(χ)", "(.*)"}), "$1$2$3$4><$5$6$7"},

	// εννιαψήφιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(ψ)", "(.*)"}), "$1$2$3$4><$5$6$7"},

	// εννιάωρος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(νν|ν)", "(ι)", "(α|ά)", "(ω)", "(.*)"}), "$1$2$3$4><$5$6$7"},

	// ενιαίως
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ν)", "(ι)", "(αι|αί|ε|έ)", "(.*)"}), "$1$2$3$4>-<$5$6"},

	// ^εννι[αά]
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(νν)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3><$4$5"},
	// todo: νν|ν?

	// έννοια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(νν|ν)", "(οι|ι)", "(ω|ώ)", "(ν)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε|έ)", "(νν|ν)", "(οι|ι)", "(α|ας|ες|ων)", "$"}), "$1$2$3>$4"},

	// εντεκάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ντ)", "(ε)", "(κ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ντεκάρι*&dq=

	// εντελβάις
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(β)", "(α)", "(ι)", "(σ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ελβαισ*&dq=

	// εξάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ξ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εξάρι*&dq=

	// εξαϋλώνω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ξ)", "(α)", "(υ|ύ|ι|ί)", "(λ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εξαυλ*&dq=

	// εξυπνάκιας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(π)", "(ν)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υπνακι*&dq=

	// εξωκλήσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(κ)", "(λ)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ωκλήσι*&dq=
	// todo: Εξωκλησί-ου, Μαυροκλησί-ου, Πρωτοκλησί-ου, Παλαιοκλησί-ου

	// εξωραΐζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ξ)", "(ω)", "(ρ)", "(α)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ξωραι*&dq=

	// επαΐουσα, επαΐων
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(π)", "(α)", "(ι|ί)", "(ω|ώ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*επαιο*&dq=

	// επανωφόρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ν)", "(ω)", "(φ)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*νωφόρι*&dq=
	// todo: ανηφόρι etc.

	// επινόημα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(ο)", "(ι|ί)", "(μ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινοιμ*&dq=
	// todo: ^νόημα, ^νόηση etc.

	// επινόηση
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(ο)", "(ι|ί)", "(μ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινοισ*&dq=

	// επινοητής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(ο)", "(ι|ί)", "(τ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινοιτ*&dq=

	// επιπλάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(λ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*πλακι*&dq=

	// επιρροή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(π)", "(ι)", "(ρρ|ρ)", "(ο)", "(ίς|ί)", "$"}), "$1$2$3$4$5$6>-$7"},

	// εργάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(γ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εργάκι*&dq=

	// εργατιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(γ)", "(α)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ργατι*&dq=

	// εργατοϋπαλληλικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ο)", "(υ|ι)", "(π)", "(α|ά)", "(λ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τουπαλ*&dq=

	// ερημιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(ρ)", "(η)", "(μ)", "(ι)", "(ι-άς...)", "$"}), "$1$2$3$4$5>$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ερημι*&dq=

	// ερημονήσι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(η)", "(μ)", "(ο)", "(ν)", "(η|ή)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9$10><$11$12"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ερημονήσι*&dq=

	// ερμαϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(μ)", "(α)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρμαικ*&dq=

	// ερμάρι
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(α|ε)", "(ρ)", "(μ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρμαρι*&dq=

	// ερυθραϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(θ)", "(ρ)", "(α)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υθραικ*&dq=

	// ερωτιάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ω)", "(τ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρωτιαρ*&dq=

	// εσαεί
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(σ)", "(α)", "(ι|ί)", "$"}), "$1-$2$3-$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εσαεί*&dq=

	// εσκιμωικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ι)", "(μ)", "(ο)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κιμοικ*&dq=

	// εσπεριδοειδές
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(δ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιδοιδ*&dq=

	// εσχατιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(χ)", "(α)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σχατι*&dq=

	// ετεροειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ε)", "(ρ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*τεροιδ*&dq=

	// ευβοϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(υ|β)", "(β)", "(ο)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},

	// ευζωία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(υ|β)", "(ζ)", "(ο)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5>-<$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ευζωι*&dq=
	// todo: book example evzoias

	// ευκλεής
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ε)", "(υ|φ)", "(κ)", "(λ)", "(ε)", "(ίς|ί)", "$"}), "$1$2$3$4$5>-$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=ευκλε*&dq=

	// ευνόητος, ευνοϊκά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(υ|β)", "(ν)", "(ο)", "(ι|ί)", "(τ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(υ|β)", "(ν)", "(ο)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ευνοι*&dq=

	// ευρωπαΐζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(υ|β)", "(ρ)", "(ω)", "(π)", "(α)", "(ι|ί)", "(.*)"}), "$1$2$3$4$5$6$7>-<$8$9"},

	// ευωδιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(υ)", "(ω)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ευωδι*&dq=

	// ζαγάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(γ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=30&lq=*αγαρι*&dq=

	// Ζαγοριανός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(α)", "(γ)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζαγορι*&dq=

	// ζαϊνισμός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(α)", "(ι)", "(ν)", "(ι)", "(σ)", "(μ)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3>-<$4$5$6$7$8$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζαινισμ*&dq=

	// Ζαΐρ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(α)", "(ι|ί)", "(ρ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζαιρ*&dq=

	// ζακετάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ε)", "(τ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ετακι*&dq=

	// ζακόνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(α)", "(κ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζακόνι*&dq=

	// Ζακυθινιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(θ)", "(ι)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιθινι*&dq=

	// Ζακυνθινιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(ν)", "(θ)", "(ι)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},

	// ζαλίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αλίκι*&dq=

	// ζαμάνι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(μ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αμανι*&dq=

	// ζάπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(α|ά)", "(π)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζάπι*&dq=

	// ζάρι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ζ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2-$3$4<$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζαρι*&dq=

	// ζαρκάδι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(α)", "(ρ)", "(κ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ζαρκαδι*&dq=

	// ζαρντινιέρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ν)", "(ι)", "(ε|έ)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ινιερ*&dq=

	// ζαρτιέρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(τ)", "(ι)", "(ε|έ)", "(ρ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρτιέρ*&dq=

	// ζαρωματιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(ω)", "(μ)", "(α)", "(τ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αρωματι*&dq=

	// ζαφειρένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(φ)", "(ει|ι)", "(ρ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*φειρένι*&dq=

	// ζάφτι
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(ζ)", "(α|ά)", "(φ)", "(τ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζάφτι*&dq=

	// ζαχαρένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αρενι*&dq=

	// ζαχαριέρα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(ρ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αρενι*&dq=

	// ζεβζεκιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(β)", "(ζ)", "(ε|έ)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εβζεκι*&dq=

	// ζεϊμπεκιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(ε|έ)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ε)", "(ι)", "(μπ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπεκι*&dq=

	// ζεμπίλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(ι|ί)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},

	// ζέρσεϊ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ε|έ)", "(ρ)", "(σ)", "(ε)", "(ι)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζερσει*&dq=

	// ζεστασιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(στ)", "(α)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εστασι*&dq=

	// ζευγάρι, ζυγαριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|β)", "(γ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*υγάρι*&dq=

	// ζευγολατειό
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(γ)", "(ο)", "(λ)", "(α)", "(τ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*γολατειό*&dq=

	// ζήλια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ζ)", "(η|ή)", "(λ)", "(ει|ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=ζήλια*&dq=

	// ζηλιάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(η|ή)", "(λ)", "(ει|ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζηλιάρ*&dq=

	// ζημιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ζ)", "(η)", "(μ)", "(ι)", "(ι-άς...)", "$"}), "$1$2-$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=ζημιά&dq=

	// ζημιάρης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(η)", "(μ)", "(ι)", "(α|ά)", "(ρ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζημιάρ*&dq=

	// Ζήρεια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(η|ή)", "(ρ)", "(ει|ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},

	// ζητιανιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(τ)", "(ι)", "(α)", "(ν)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*τιανι*&dq=
	// todo: *ιανι*

	// ζητιάνα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(τ)", "(ι)", "(α|ά)", "(ν)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ητιαν*&dq=
	// todo: *τιαν*

	// ζιμπούλι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(μπ)", "(ου|ού)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιμπουλι*&dq=

	// ζίνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ι|ί)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζινι*&dq=

	// ζιπάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(π)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιπακι*&dq=

	// ζόμπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ο|ό)", "(μπ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},

	// ζόρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ο|ό)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζορι*&dq=

	// ζορμπαλιδίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(δ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λιδικι*&dq=

	// ζουλάπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(λ)", "(α|ά)", "(π)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουλαπι*&dq=

	// ζουμάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου)", "(μ)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ουμακι*&dq=

	// ζουμί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ου|ού)", "(μ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ζουμι*&dq=

	// ζούρλια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ου|ού)", "(ρ)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*ουρλι*&dq=

	// ζούφιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ου|ού)", "(φ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζουφι*&dq=

	// ζοχαδιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(χ)", "(α|ά)", "(δ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αχαδι*&dq=

	// ζύγι, ζυγιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ζ)", "(υ|ύ|ι|ί)", "(γ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=ζύγι*&dq=

	// ζυγολούρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(λ)", "(ου|ού)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ολούρι*&dq=
	// todo: Μεσολουριου, Μοσχολουριου

	// ζυμάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(υ|ι)", "(μ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζυμαρι*&dq=
	// todo: Ζυμαρι-ου

	// ζυμοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(υ|ι)", "(μ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζυμοι*&dq=

	// ζωάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ω)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},

	// ζωγραφιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ω)", "(γ)", "(ρ)", "(α)", "(φ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// todo: Ζωγραφία

	// ζωή
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ζ)", "(ω)", "(ίσ?)", "$"}), "$1$2-$3"},

	// ζωηρά, ζωηρεύω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ω)", "(η|ή)", "(ρ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// todo: Λατζόιο

	// ζωηφόρος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ω)", "(η)", "(φ)", "(ο|ό)", "(ρ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// ζοιφορ

	// -ζωία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ω)", "(ι|ί)", "(α)", "(.*)"}), "$1$2$3>-<$4$5$6"},

	// ζωίτσα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(ι|ί)", "(τσ)", "(.*)"}), "$1$2>-<$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οιτσ*&dq=

	// ζωνάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ω)", "(ν)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ζωναρι*&dq=

	// ζωντάνια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(ντ)", "(α|ά)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*οντανι*&dq=

	// ζωοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ω)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3>-<$4$5$6"},

	// ζωύφιο
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ζ)", "(ω)", "(υ|ύ)", "(φ)", "(ι|ί)", "(.*)"}), "$1$2>-<$3$4$5$6"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ζ)", "(ω)", "(ι|ί)", "(φ)", "(ι|ί)", "(.*)"}), "$1$2$3>-<$4$5$6$7"},

	// ηθμοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(μ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},

	// ηλιοβασίλεμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(ο)", "(β)", "(α)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηλιοβα*&dq=

	// ηλιοθεραπεία
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(λ)", "(ι)", "(ο)", "(θ)", "(ε)", "(.*)"}), "$1$2$3$4><$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηλιοθε*&dq=

	// ηλιοκαμένος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(ο)", "(κ)", "(α|ά)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λιοκα*&dq=

	// ηλιόλουστος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(ο|ό)", "(λ)", "(ου|ού)", "(.*)"}), "$1$2$3><$4$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λιολου*&dq=

	// Ίλιον
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(Ι|Ί)", "(λ)", "(ι)", "(ο)", "(ν)", "$"}), "$1-$2$3-$4$5"},
	// ήλιος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(η|ή)", "(λ)", "(ι)", "(ος...)", "$"}), "$1-$2$3$4"},

	// ηλιόσπορος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(ο|ό)", "(σ)", "(π)", "(ο|ό)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λιοσπο*&dq=

	// ηλιοφώτιστος, ηλιόφωτος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(λ)", "(ι)", "(ο|ό)", "(φ)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3><$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*λιοφω*&dq=

	// ηλιόχαρος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(λ)", "(ι)", "(ο|ό)", "(χ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηλιόχ*&dq=

	// ηλιοψημένος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(λ)", "(ι)", "(ο)", "(ψ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηλιοψ*&dq=

	// ημερολογιάκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(γ)", "(ι)", "(α|ά)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ογιακι*&dq=

	// ημηισεληνοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(λ)", "(η)", "(ν)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*εληνοι*&dq=

	// ήπια
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(η|ή)", "(π)", "(ι)", "(ήπι-α)", "$"}), "$1$2$3><$4"},
	// todo: πιες

	// Ηρακλειώτης
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(λ)", "(ει|ι)", "(ω|ώ)", "(τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ακλειώτ*&dq=

	// ηρωίδα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(ρ)", "(ο)", "(ι|ί)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},

	// ηρωίνη
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ο)", "(ι|ί)", "(ν)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ροιν*&dq=

	// Ησαΐας
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(σ)", "(α)", "(ι|ί)", "(α)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ησαΐα*&dq=

	// θαλάμι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(α)", "(λ)", "(α|ά)", "(μ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θαλάμι*&dq=

	// θαλασσαϊτός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(σσ|σ)", "(α)", "(ϊ)", "(τ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// todo: σταυραϊτός χαρταϊτός, χρυσαϊτός, θαλασσαητός

	// θαλασσοπλοΐα (θαλασσοπλο-ία)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(α)", "(σσ|σ)", "(ο)", "(π)", "(λ)", "(ο)", "(ί)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9>-<$10$11"},
	// ẗodo: ποταμόπλοια vs ποταμοπλοία

	// θαλασσοσπηλιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(π)", "(η)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σπηλιά*&dq=

	// Θάλεια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(α|ά)", "(λ)", "(ει|ι)", "(α|ας|ες|ων)", "$"}), "$1$2$3$4$5>$6"},

	// θαμνοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(ν)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μνοιδ*&dq=

	// θειός
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(ει|ι)", "(ος|α|ο 2)", "$"}), "$1$2$3"},

	// θειαφένιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ει|ι)", "(α)", "(φ)", "(ε|έ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3>$4-$5$6-$7$8<$9$10"},

	// θειάφι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ει|ι)", "(α|ά)", "(φ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3>$4-$5$6<$7$8"},

	// θειαφιστήρια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(στ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αφιστήρι*&dq=

	// θειαφίζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ει|ι)", "(α|ά)", "(φ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θειαφ*&dq=

	// θεϊκά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε)", "(ι|ί)", "(κ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θεϊκ*&dq=

	// θελιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θελι*&dq=

	// θεοείδεια (θεο-ίδεια)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε)", "(ο)", "(ι|ί)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θεοίδ*&dq=

	// θεριακλίκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(λ)", "(ι|ί)", "(κ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κλικι*&dq=

	// θεριακή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε)", "(ρ)", "(ι)", "(α)", "(κ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θεριακ*&dq=

	// θέριεμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ε|έ)", "(ρ)", "(ι)", "(ε|έ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θέριε*&dq=

	// θεριό
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(ε)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=20&lq=θερι*&dq=

	// θερμοηχομονωτικός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(μ)", "(ο)", "(ι)", "(χ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρμοιχ*&dq=

	// θερμοπαρακάλια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(α)", "(κ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρακαλι*&dq=

	// Θεσσαλονικιά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(λ)", "(ο)", "(ν)", "(ι)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7$8><$9$10"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*αλονικι*&dq=
	// todo: ικι'ά'...

	// θηβαίικος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(β)", "(ε)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηβεικ*&dq=

	// θηκάρι
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(η)", "(κ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θηκάρι*&dq=

	// θηλιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(η)", "(λ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// todo: συρτοθηλιά

	// θηλοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(λ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηλοιδ*&dq=

	// θηλύκι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(η)", "(λ)", "(υ|ύ|ι|ί)", "(κ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θηλύκι*&dq=

	// θηλυκωτάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ω)", "(τ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κωτάρι*&dq=

	// θηλυκωτήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(κ)", "(ω)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*κωτήρι*&dq=

	// θημωνιά, θημώνιασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(η)", "(μ)", "(ω|ώ)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θημωνι*&dq=

	// θηραϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(η)", "(ρ)", "(α)", "(ι)", "(κ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ηραικ*&dq=

	// Θιακή
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(ι)", "(α|ά)", "(κ)", "(.*)"}), "$1$2><$3$4$5"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=θιακ*&dq=

	// θολοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(λ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ολοιδ*&dq=

	// θροΐζω, θρόισμα
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(ρ)", "(ο)", "(ι|ί)", "(ζ|σ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=θροι*&dq=

	// θρονί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ρ)", "(ο|ό)", "(ν)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*θρονί*&dq=

	// θρούμπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(ρ)", "(ου|ού)", "(μπ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θρούμπι*&dq=

	//	θρυψαλιάζω, θρυψάλιασμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(υ|ι)", "(ψ)", "(α|ά)", "(λ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρυψαλι*&dq=

	// θυλακοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(α)", "(κ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ακοιδ*&dq=

	// θυμαρίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(α)", "(ρ)", "(ι|ί)", "(σ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μαρισι*&dq=
	// todo: http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*αρισι*&dq=

	// θυμάρι, θυμαριά
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(υ|ι)", "(μ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θυμαρι*&dq=

	// θυμητάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(η)", "(τ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μητάρι*&dq=

	// θυμιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(υ|ι)", "(μ)", "(ι)", "(α|ά)", "(ά-ζω)", "$"}), "$1$2$3$4><$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θυμιάζω*&dq=

	// θυμιατήριο
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(υ|ι)", "(μ)", "(ι)", "(α)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(ο)", "$"}),
		"$1$2$3$4$5><$6$7$8$9$10$11"},
	// θυμιατήρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(υ|ι)", "(μ)", "(ι)", "(α)", "(τ)", "(η|ή)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}),
		"$1$2$3$4$5>$6-$7$8-$9$10<$11$12"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θυμιατήρι*&dq=

	// θυμιατίζω
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(υ|ι)", "(μ)", "(ι)", "(α|ά)", "(τ)", "(.*)"}), "$1$2$3$4$5><$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?start=10&lq=*θυμιατ*&dq=

	// θυμοειδής
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(υ|ι)", "(μ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},

	// θυρεοειδεκτομή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(υ|ι)", "(ρ)", "(ε)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8$9"},
	// "θυροιδής"
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(θ)", "(υ|ι)", "(ρ)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*θυρεοειδ*&dq=

	// Θωμαή
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(ω)", "(μ)", "(α)", "(ι|ί)", "(.*)"}), "$1$2$3$4>-<$5$6"},
	// Θωμαη*

	// θώρι, θωριά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(θ)", "(ω|ώ)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},
	// todo: ξεθωριάζω

	// ίδιος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι|ί)", "(δ)", "(ι)", "(ος|α|ο)", "$"}), "$1-$2$3$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=ιδιος*&dq=

	// ιδροκόπι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(δ)", "(ρ)", "(ω)", "(κ)", "(ο|ό)", "(π)", "(ι)", "(α|ά)", "(.*)"}),
		"$1$2$3$4$5$6$7$8$9><$10$11"},

	// ιδρωτάρι
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ρ)", "(ω)", "(τ)", "(α|ά)", "(ρ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4$5$6$7><$8$9"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ρωταρι*&dq=

	// Ικάριος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι)", "(κ)", "(α|ά)", "(ρ)", "(ι)", "(ος|α|ο)", "$"}), "$1$2$3$4$5>-<$6"},
	// Ικαριά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι)", "(κ)", "(α)", "(ρ)", "(ι)", "(α|ά|ε|έ|ω|ώ)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=ικαρι*&dq=

	// ιλουστρασιόν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(στ)", "(ρ)", "(α)", "(σ)", "(ι)", "(ο|ό)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*στρασι*&dq=

	// Ιμαλάια
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(μ)", "(α)", "(λ)", "(α)", "(ι)", "(.*)"}), "$1$2$3$4$5$6>-<$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιμαλαι*&dq=

	// μπαϊλντί
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μπ)", "(α)", "(ι)", "(λ)", "(.*)"}), "$1$2$3>-<$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπαιλ*&dq=

	// μέικερ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(μ)", "(ε)", "(ι)", "(κ)", "(ε)", "(ρ)", "(.*)"}), "$1$2$3>-<$4$5$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μεικερ*&dq=

	// ιμιτασιόν
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(τ)", "(α)", "(σ)", "(ι)", "(ο|ό)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιτασιο*&dq=

	// ιντερβιού
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ε)", "(ρ)", "(β)", "(ι)", "(ου|ού)", "(.*)"}), "$1$2$3$4$5><$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ερβιού*&dq=

	// Νησιά
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ν)", "(η)", "(σ)", "(ι)", "(α|ά|ω|ώ)", "(.*)"}), "$1$2$3$4><$5$6"},

	// ιουδαϊκός
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ου)", "(δ)", "(α)", "(ι)", "(κ|σμ|στ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιουδαι*&dq=

	// ιπποειδή
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ππ|π)", "(ο)", "(ι)", "(δ)", "(.*)"}), "$1$2$3$4>-<$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιπποιδ*&dq=

	// ίσιος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι|ί)", "(σ)", "(ι)", "(ος|α|ο)", "$"}), "$1-$2$3$4"},
	// todo: ολόισια

	// ισιάδα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(σ)", "(ι)", "(α|ά)", "(δ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισιαδ*&dq=

	// ισιάζω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι|ί)", "(σ)", "(ι)", "(α|ά)", "(ζ|σμ|στ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισιαζ*&dq=
	// todo: ίσιασμα

	// ισιώνω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι|ί)", "(σ)", "(ι)", "(ω|ώ)", "(μ|ν|σ|τ)", "(.*)"}), "$1$2$3><$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισιαζ*&dq=

	// ίσκιος
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι|ί)", "(σ)", "(κ)", "(ι)", "(ος...)", "$"}), "$1$2$3$4>$5"},
	// todo: ανίσκιος, περίσκιος, ήσκιος

	// ίσκιωμα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(σ)", "(κ)", "(ι)", "(ω|ώ)", "(μ)", "(.*)"}), "$1$2$3$4><$5$6$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*σκιωμ*&dq=

	// ισκιώνω
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(ι|ί)", "(σ)", "(κ)", "(ι)", "(ω|ώ)", "(ώ-νω)", "$"}), "$1$2$3$4><$5$6"},

	// Ισμαήλ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ζ|σ)", "(μ)", "(α)", "(ι|ί)", "(λ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},

	// Ιζραήλ
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(ζ|σ)", "(ρ)", "(α)", "(ι|ί)", "(λ)", "(.*)"}), "$1$2$3$4$5>-<$6$7$8"},

	// Ιταλιάνα
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ι)", "(τ)", "(α)", "(λ)", "(ι)", "(α|ά)", "(.*)"}), "$1$2$3$4$5$6><$7$8"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ιταλιαν*&dq=

	// ...

	// ισοπαλία (ισοπαλι-α)
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(ο)", "(π)", "(α)", "(λ)", "(ι)", "(ας|ες|ων|ών|α)", "$"}), "$1$2>-$3$4-$5$6-$7"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*ισοπαλ*&dq=

	// μπαλιά
	// GrhyphRule{customRegexpCompile(
	// []string{"(.*)", "(μπ)", "(α)", "(λ)", "(ι)", "(ι-άς...)", "$"}), "$1$2$3>-$4$5$6"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=*μπαλια*&dq=

	// παλιά (πα-λια)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(π)", "(α)", "(λ)", "(ι)", "(ους|ούς|ες|ές|ων|ών|ου|ού|α|ά|ε|έ|ο|ό)", "$"}), "$1$2-$3$4$5"},
	GrhyphRule{customRegexpCompile(
		[]string{"(.*)", "(π)", "(α)", "(λ)", "(ι)", "(α|ά|ο|ό)", "(.*)"}), "$1$2$3>-$4$5<$6$7"}, //[^μm](π)
	// todo: οπάλιο (οπαλι-ο), Μπαλίου (μπαλι-ου)

	GrhyphRule{customRegexpCompile(
		[]string{"^", "(π)", "(ι)", "(ο)", "$"}), "$1$2$3"},

	// ...

	// πιάνω (πιάσ' το)
	GrhyphRule{customRegexpCompile(
		[]string{"^", "(π)", "(ι)", "(α|ά)", "(σ)", "$"}), "$1$2$3$4"},
	// http://www.greek-language.gr/greekLang/modern_greek/tools/lexica/search.html?lq=πιασ&dq=
	// todo: πιάνω

	//todo: test arguments count with queries
	//todo: test replacements number sequence (the rules are dynamic either anyway)
	//todo: test >< placements, multiple occurences etc.

	// ...
}
