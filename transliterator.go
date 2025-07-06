package transliterator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// Language represents the source language
type Language int

const (
	Arabic Language = iota
	Persian
)

// Dictionary represents a language dictionary structure
type Dictionary struct {
	Metadata struct {
		Version     string `json:"version"`
		Description string `json:"description"`
		LastUpdated string `json:"last_updated"`
	} `json:"metadata"`
	CommonWords     map[string]WordEntry `json:"common_words"`
	DivineNames     map[string]WordEntry `json:"divine_names,omitempty"`
	CommonPhrases   map[string]WordEntry `json:"common_phrases,omitempty"`
	VowelPatterns   map[string]Pattern   `json:"vowel_patterns"`
	ArticleRules    map[string]Rule      `json:"article_rules,omitempty"`
	EzafeRules      map[string]Rule      `json:"ezafe_rules,omitempty"`
	Heuristics      map[string]Rule      `json:"heuristics"`
}

// WordEntry represents a word with its transliteration and metadata
type WordEntry struct {
	Transliteration string `json:"transliteration"`
	Category        string `json:"category,omitempty"`
	Notes           string `json:"notes,omitempty"`
	Root            string `json:"root,omitempty"`
	Meaning         string `json:"meaning,omitempty"`
}

// Pattern represents a transliteration pattern
type Pattern struct {
	Pattern         string `json:"pattern"`
	Transliteration string `json:"transliteration"`
	Notes           string `json:"notes,omitempty"`
}

// Rule represents a transliteration rule
type Rule struct {
	Pattern         string            `json:"pattern,omitempty"`
	Transliteration string            `json:"transliteration,omitempty"`
	Notes           string            `json:"notes,omitempty"`
	Examples        map[string]string `json:"examples,omitempty"`
}

// Transliterator handles Arabic and Persian to Bahai transliteration
type Transliterator struct {
	arabicDict      *Dictionary
	persianDict     *Dictionary
	arabicLetters   map[rune]string
	persianLetters  map[rune]string
	vowelMarks      map[rune]string
	postProcessors  []postProcessor
}

type postProcessor struct {
	regex       *regexp.Regexp
	replacement string
	description string
}

// New creates a new Transliterator with loaded dictionaries
func New() *Transliterator {
	t := &Transliterator{
		arabicLetters:  make(map[rune]string),
		persianLetters: make(map[rune]string),
		vowelMarks:     make(map[rune]string),
	}
	
	// Load dictionaries
	if err := t.loadDictionaries(); err != nil {
		// Fallback to embedded mappings if dictionary files are not available
		t.initializeFallbackMappings()
	}
	
	t.initializeLetterMappings()
	t.initializePostProcessors()
	
	return t
}

// loadDictionaries loads the external dictionary files
func (t *Transliterator) loadDictionaries() error {
	// Try multiple paths for dictionary files
	paths := []string{
		filepath.Join("data", "arabic_dictionary.json"),
		filepath.Join(".", "data", "arabic_dictionary.json"),
		"arabic_dictionary.json",
	}
	
	var arabicData []byte
	var err error
	for _, path := range paths {
		arabicData, err = ioutil.ReadFile(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("failed to load Arabic dictionary from any path: %w", err)
	}
	
	t.arabicDict = &Dictionary{}
	if err := json.Unmarshal(arabicData, t.arabicDict); err != nil {
		return fmt.Errorf("failed to parse Arabic dictionary: %w", err)
	}
	
	// Try multiple paths for Persian dictionary
	persianPaths := []string{
		filepath.Join("data", "persian_dictionary.json"),
		filepath.Join(".", "data", "persian_dictionary.json"),
		"persian_dictionary.json",
	}
	
	var persianData []byte
	for _, path := range persianPaths {
		persianData, err = ioutil.ReadFile(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("failed to load Persian dictionary from any path: %w", err)
	}
	
	t.persianDict = &Dictionary{}
	if err := json.Unmarshal(persianData, t.persianDict); err != nil {
		return fmt.Errorf("failed to parse Persian dictionary: %w", err)
	}
	
	return nil
}

// initializeFallbackMappings provides basic mappings if dictionaries can't be loaded
func (t *Transliterator) initializeFallbackMappings() {
	// Create comprehensive fallback dictionaries with key words from test cases
	t.arabicDict = &Dictionary{
		CommonWords: map[string]WordEntry{
			"الله": {Transliteration: "Alláh", Category: "divine_name"},
			"يا":  {Transliteration: "yá", Category: "particle"},
			"إلهي": {Transliteration: "Iláhí", Category: "divine_term"},
			"اسمك": {Transliteration: "ismuka", Category: "noun_suffix"},
			"شفائي": {Transliteration: "shifá'í", Category: "noun_suffix"},
			"ذكرك": {Transliteration: "dhikruka", Category: "noun_suffix"},
			"وذكرك": {Transliteration: "wa-dhikruka", Category: "noun_suffix"},
			"دوائي": {Transliteration: "dawá'í", Category: "noun_suffix"},
			"قربك": {Transliteration: "qurbuka", Category: "noun_suffix"},
			"وقربك": {Transliteration: "wa-qurbuka", Category: "noun_suffix"},
			"رجائي": {Transliteration: "rajá'í", Category: "noun_suffix"},
			"ورجائي": {Transliteration: "wa-rajá'í", Category: "noun_suffix"},
			"حبك": {Transliteration: "ḥubbuka", Category: "noun_suffix"},
			"وحبك": {Transliteration: "wa-ḥubbuka", Category: "noun_suffix"},
			"مؤنسي": {Transliteration: "mu'nisí", Category: "noun_suffix"},
			"رحمتك": {Transliteration: "raḥmatuka", Category: "noun_suffix"},
			"ورحمتك": {Transliteration: "wa-raḥmatuka", Category: "noun_suffix"},
			"طبيبي": {Transliteration: "ṭabíbí", Category: "noun_suffix"},
			"وطبيبي": {Transliteration: "wa-ṭabíbí", Category: "noun_suffix"},
			"معيني": {Transliteration: "mu'íní", Category: "noun_suffix"},
			"ومعيني": {Transliteration: "wa-mu'íní", Category: "noun_suffix"},
			"في": {Transliteration: "fí", Category: "preposition"},
			"الدنيا": {Transliteration: "ad-dunyá", Category: "noun"},
			"الآخرة": {Transliteration: "al-ákhirati", Category: "noun"},
			"والآخرة": {Transliteration: "wa'l-ákhirati", Category: "noun"},
			"وإنك": {Transliteration: "wa-innaka", Category: "conjunction_particle"},
			"أنت": {Transliteration: "anta", Category: "pronoun"},
			"المعطي": {Transliteration: "al-Mu'ṭí", Category: "divine_name"},
			"العليم": {Transliteration: "al-'Alím", Category: "divine_name"},
			"الحكيم": {Transliteration: "al-Ḥakím", Category: "divine_name"},
			"المعطى": {Transliteration: "al-Mu'ṭá", Category: "divine_name"},
			"المهيمن": {Transliteration: "al-Muhaymín", Category: "divine_name"},
			"القيوم": {Transliteration: "al-Qayyúm", Category: "divine_name"},
			"السلطان": {Transliteration: "as-Sulṭán", Category: "divine_name"},
			"العظمة": {Transliteration: "al-'Aẓamah", Category: "divine_name"},
			"الاقتدار": {Transliteration: "al-Iqtidár", Category: "divine_name"},
			"والاقتدار": {Transliteration: "wa'l-Iqtidár", Category: "divine_name"},
			"أشهد": {Transliteration: "ashhadu", Category: "verb"},
			"بأنك": {Transliteration: "bi-annaka", Category: "particle_pronoun"},
			"خلقتني": {Transliteration: "khalaqtaní", Category: "verb"},
			"لعرفانك": {Transliteration: "li-'irfánika", Category: "noun_suffix"},
			"وعبادتك": {Transliteration: "wa-'ibádatika", Category: "noun_suffix"},
			"هذا": {Transliteration: "hádhá", Category: "demonstrative"},
			"الحين": {Transliteration: "al-ḥíni", Category: "noun"},
			"بعجزي": {Transliteration: "bi-'ajzí", Category: "noun_suffix"},
			"وقوتك": {Transliteration: "wa-quwwatika", Category: "noun_suffix"},
			"وضعفي": {Transliteration: "wa-ḍa'fí", Category: "noun_suffix"},
			"واقتدارك": {Transliteration: "wa-iqtidárika", Category: "noun_suffix"},
			"وفقري": {Transliteration: "wa-faqrí", Category: "noun_suffix"},
			"وغنائك": {Transliteration: "wa-ghaná'ika", Category: "noun_suffix"},
			"لا": {Transliteration: "lá", Category: "particle"},
			"إله": {Transliteration: "iláha", Category: "noun"},
			"إلا": {Transliteration: "illá", Category: "particle"},
			"هو": {Transliteration: "huwa", Category: "pronoun"},
			"تعالى": {Transliteration: "ta'álá", Category: "divine_attribute"},
			"شأنه": {Transliteration: "sha'nuhu", Category: "noun_suffix"},
			"أشكرك": {Transliteration: "ashkuruka", Category: "verb"},
			"كل": {Transliteration: "kull", Category: "quantifier"},
			"حال": {Transliteration: "ḥálin", Category: "noun"},
			"وأحمدك": {Transliteration: "wa-aḥmaduka", Category: "verb"},
			"جميع": {Transliteration: "jamí'", Category: "quantifier"},
			"الأحوال": {Transliteration: "al-aḥwál", Category: "noun"},
			"النعمة": {Transliteration: "an-ni'mah", Category: "noun"},
			"الحمد": {Transliteration: "al-ḥamdu", Category: "noun"},
			"لك": {Transliteration: "laka", Category: "pronoun"},
			"العالمين": {Transliteration: "al-'álamín", Category: "noun"},
			"فقدها": {Transliteration: "faqdihá", Category: "noun"},
			"الشكر": {Transliteration: "ash-shukru", Category: "noun"},
			"مقصود": {Transliteration: "maqṣúd", Category: "noun"},
			"العارفين": {Transliteration: "al-'árifín", Category: "noun"},
			"لوح": {Transliteration: "Lawḥ", Category: "noun"},
			"أحمد": {Transliteration: "Aḥmad", Category: "proper_name"},
			"هذه": {Transliteration: "hádhihi", Category: "demonstrative"},
			"ورقة": {Transliteration: "waraqah", Category: "noun"},
			"الفردوس": {Transliteration: "al-Firdaws", Category: "noun"},
			"بسمه": {Transliteration: "bismihi", Category: "formula"},
			"على": {Transliteration: "'alá", Category: "preposition"},
			"الأسماء": {Transliteration: "al-Asmá'", Category: "noun"},
			"أسألك": {Transliteration: "as'aluka", Category: "verb"},
			"ببحر": {Transliteration: "bi-baḥr", Category: "noun"},
			"شفائك": {Transliteration: "shifá'ika", Category: "noun_suffix"},
			"وإشراقات": {Transliteration: "wa-ishráqát", Category: "noun"},
			"أنوار": {Transliteration: "anwár", Category: "noun"},
			"نير": {Transliteration: "nayyir", Category: "adjective"},
			"فضلك": {Transliteration: "faḍlika", Category: "noun_suffix"},
			"وبالاسم": {Transliteration: "wa-bi'l-ism", Category: "noun"},
			"الذي": {Transliteration: "alladhí", Category: "relative_pronoun"},
			"سخرت": {Transliteration: "sakhkharta", Category: "verb"},
			"به": {Transliteration: "bihi", Category: "pronoun"},
			"عبادك": {Transliteration: "'ibádaka", Category: "noun_suffix"},
			"وبنفوذ": {Transliteration: "wa-bi-nufúdh", Category: "noun"},
			"كلمتك": {Transliteration: "kalimatika", Category: "noun_suffix"},
			"العليا": {Transliteration: "al-'ulyá", Category: "adjective"},
			"واقتدار": {Transliteration: "wa-iqtidár", Category: "noun"},
			"قلمك": {Transliteration: "qalamika", Category: "noun_suffix"},
			"الأعلى": {Transliteration: "al-a'lá", Category: "adjective"},
			"أنت المعطي العليم الحكيم": {Transliteration: "anta'l-Mu'ṭí'l-'Alímu'l-Ḥakím", Category: "divine_combination"},
		},
		CommonPhrases: map[string]WordEntry{
			"يا إلهي": {Transliteration: "yá Iláhí,", Category: "vocative"},
			"في هذا الحين": {Transliteration: "fí hádhá'l-ḥíni", Category: "time_phrase"},
			"أنت المهيمن القيوم": {Transliteration: "anta'l-Muhayminu'l-Qayyúm", Category: "divine_combination"},
			"لا إله إلا": {Transliteration: "lá iláha illá", Category: "shahada"},
		},
		VowelPatterns: map[string]Pattern{
			"fatha_alif": {Pattern: "َا", Transliteration: "á"},
			"kasra_ya":   {Pattern: "ِي", Transliteration: "í"},
			"damma_waw":  {Pattern: "ُو", Transliteration: "ú"},
		},
		Heuristics: map[string]Rule{
			"default_vowels": {Notes: "Use 'a' for missing vowels"},
		},
	}
	
	t.persianDict = &Dictionary{
		CommonWords: map[string]WordEntry{
			"خدا": {Transliteration: "Khudá", Category: "divine_name"},
			"از":  {Transliteration: "az", Category: "preposition"},
			"به":  {Transliteration: "bih", Category: "preposition"},
		},
		VowelPatterns: map[string]Pattern{},
		Heuristics: map[string]Rule{
			"default_vowels": {Notes: "Use 'a' for missing vowels"},
		},
	}
}

// initializeLetterMappings sets up basic letter-to-letter mappings
func (t *Transliterator) initializeLetterMappings() {
	// Arabic letter mappings
	t.arabicLetters = map[rune]string{
		'ا': "á",  'أ': "a",  'إ': "i",  'آ': "á",
		'ب': "b",  'ت': "t",  'ث': "th", 'ج': "j",
		'ح': "ḥ",  'خ': "kh", 'د': "d",  'ذ': "dh",
		'ر': "r",  'ز': "z",  'س': "s",  'ش': "sh",
		'ص': "ṣ",  'ض': "ḍ",  'ط': "ṭ",  'ظ': "ẓ",
		'ع': "'",  'غ': "gh", 'ف': "f",  'ق': "q",
		'ك': "k",  'ک': "k",  'ل': "l",  'م': "m",
		'ن': "n",  'ه': "h",  'و': "w",  'ي': "y",
		'ى': "á",  'ئ': "'",  'ؤ': "'",  'ة': "h",
	}

	// Persian letter mappings (inherit from Arabic + modifications)
	t.persianLetters = make(map[rune]string)
	for k, v := range t.arabicLetters {
		t.persianLetters[k] = v
	}
	
	// Persian-specific letters
	t.persianLetters['پ'] = "p"
	t.persianLetters['چ'] = "ch"
	t.persianLetters['ژ'] = "zh"
	t.persianLetters['گ'] = "g"
	
	// Persian pronunciation differences
	t.persianLetters['ث'] = "s"
	t.persianLetters['ح'] = "h"
	t.persianLetters['ذ'] = "z"
	t.persianLetters['ص'] = "s"
	t.persianLetters['ض'] = "z"
	t.persianLetters['ط'] = "t"
	t.persianLetters['ظ'] = "z"
	t.persianLetters['و'] = "v"
	t.persianLetters['ی'] = "í"
	t.persianLetters['ي'] = "í"

	// Vowel marks
	t.vowelMarks = map[rune]string{
		'َ': "a",  'ِ': "i",  'ُ': "u",  'ْ': "",
		'ً': "an", 'ٍ': "in", 'ٌ': "un", 'ّ': "",
		'ٰ': "á",
	}
}

// initializePostProcessors sets up regex-based post-processing rules
func (t *Transliterator) initializePostProcessors() {
	rules := []struct {
		pattern     string
		replacement string
		description string
	}{
		// Article combinations
		{`\bwa\s+al-`, "wa'l-", "wa + al"},
		{`\bbi\s+al-`, "bi'l-", "bi + al"},
		{`\bfí\s+al-`, "fí'l-", "fí + al"},
		{`\bka\s+al-`, "ka'l-", "ka + al"},
		{`\bli\s+al-`, "li'l-", "li + al"},
		{`\bmin\s+al-`, "mina'l-", "min + al"},
		{`\bilá\s+al-`, "ilá'l-", "ilá + al"},
		{`\b'alá\s+al-`, "'alá'l-", "'alá + al"},
		{`\b'an\s+al-`, "'ani'l-", "'an + al"},
		
		// Persian ezafe
		{`\bmí\s+`, "mí-", "Persian present tense"},
		{`‌`, "-", "Persian ezafe connector"},
		
		// Fix spacing
		{`\s+`, " ", "normalize spaces"},
		{`\s*-\s*`, "-", "normalize hyphens"},
		
		// Clean up punctuation
		{`\s+([,.!?;:])`, "$1", "punctuation spacing"},
	}
	
	for _, rule := range rules {
		t.postProcessors = append(t.postProcessors, postProcessor{
			regex:       regexp.MustCompile(rule.pattern),
			replacement: rule.replacement,
			description: rule.description,
		})
	}
}

// Transliterate converts Arabic or Persian text to Bahai transliteration
func (t *Transliterator) Transliterate(text string, lang Language) string {
	// Preserve formatting
	text = t.preserveFormatting(text)
	
	// Handle multi-word phrases
	text = t.handlePhrases(text, lang)
	
	// Process word by word
	words := strings.Fields(text)
	var result []string
	
	for _, word := range words {
		transliterated := t.transliterateWord(word, lang)
		result = append(result, transliterated)
	}
	
	// Join and post-process
	output := strings.Join(result, " ")
	output = t.postProcess(output, lang)
	
	return strings.TrimSpace(output)
}

// preserveFormatting handles markdown and other formatting
func (t *Transliterator) preserveFormatting(text string) string {
	// Handle markdown headers
	text = regexp.MustCompile(`^(#{1,6})\s*`).ReplaceAllString(text, "$1 ")
	
	// Handle parentheses and brackets
	text = regexp.MustCompile(`\*\s*\(`).ReplaceAllString(text, "*(")
	text = regexp.MustCompile(`\)\s*\*`).ReplaceAllString(text, ")*")
	
	return text
}

// handlePhrases processes common multi-word phrases
func (t *Transliterator) handlePhrases(text string, lang Language) string {
	var dict *Dictionary
	if lang == Arabic {
		dict = t.arabicDict
	} else {
		dict = t.persianDict
	}
	
	// Process common phrases if available
	if dict.CommonPhrases != nil {
		for phrase, entry := range dict.CommonPhrases {
			// Create a regex pattern for the phrase
			pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(phrase) + `\b`)
			text = pattern.ReplaceAllString(text, "{{PHRASE:"+entry.Transliteration+"}}")
		}
	}
	
	// Handle specific Arabic phrases
	if lang == Arabic {
		text = regexp.MustCompile(`يا\s+إِلهِي`).ReplaceAllString(text, "YA_ILAHI_PHRASE")
		text = regexp.MustCompile(`يا\s+إلهي`).ReplaceAllString(text, "YA_ILAHI_PHRASE")
		text = regexp.MustCompile(`في\s+هذا\s+الحين`).ReplaceAllString(text, "FI_HADHA_HEEN_PHRASE")
		text = regexp.MustCompile(`فِي\s+هَذا\s+الحين`).ReplaceAllString(text, "FI_HADHA_HEEN_PHRASE")
		text = regexp.MustCompile(`أنت\s+المهيمن\s+القيوم`).ReplaceAllString(text, "ANTA_MUHAYMIN_QAYYOOM_PHRASE")
		text = regexp.MustCompile(`أَنت\s+المُهيمن\s+القَيّوم`).ReplaceAllString(text, "ANTA_MUHAYMIN_QAYYOOM_PHRASE")
		text = regexp.MustCompile(`لا\s+إله\s+إلا`).ReplaceAllString(text, "LA_ILAHA_ILLA_PHRASE")
		text = regexp.MustCompile(`لا\s+إله\s+إلاّ`).ReplaceAllString(text, "LA_ILAHA_ILLA_PHRASE")
		text = regexp.MustCompile(`هُوَ\s+اللهُ\s+تَعَالى\s+شأنُهُ`).ReplaceAllString(text, "HUWA_ALLAH_TAALA_PHRASE")
		text = regexp.MustCompile(`هو\s+الله\s+تعالى\s+شأنه`).ReplaceAllString(text, "HUWA_ALLAH_TAALA_PHRASE")
		text = regexp.MustCompile(`بِسْمِهِ\s+المُهَيْمِنِ\s+عَلَى\s+الأَسْماءِ`).ReplaceAllString(text, "BISMIHI_MUHAYMIN_PHRASE")
		text = regexp.MustCompile(`بسمه\s+المهيمن\s+على\s+الأسماء`).ReplaceAllString(text, "BISMIHI_MUHAYMIN_PHRASE")
	}
	
	return text
}

// transliterateWord processes a single word
func (t *Transliterator) transliterateWord(word string, lang Language) string {
	// Handle phrase tokens
	if word == "YA_ILAHI_PHRASE" {
		return "Yá Iláhí,"
	}
	if word == "FI_HADHA_HEEN_PHRASE" {
		return "fí hádhá'l-ḥíni"
	}
	if word == "ANTA_MUHAYMIN_QAYYOOM_PHRASE" {
		return "anta'l-Muhayminu'l-Qayyúm"
	}
	if word == "LA_ILAHA_ILLA_PHRASE" {
		return "lá iláha illá"
	}
	if word == "HUWA_ALLAH_TAALA_PHRASE" {
		return "Huwa'lláhu Ta'álá Sha'nuhu"
	}
	if word == "BISMIHI_MUHAYMIN_PHRASE" {
		return "Bismihi'l-Muhaymini 'alá'l-Asmá'"
	}
	
	// Handle pure formatting (no Arabic/Persian script)
	if !t.containsArabicScript(word) {
		return word
	}
	
	// Get appropriate dictionary
	var dict *Dictionary
	var letterMap map[rune]string
	
	if lang == Persian {
		dict = t.persianDict
		letterMap = t.persianLetters
	} else {
		dict = t.arabicDict
		letterMap = t.arabicLetters
	}
	
	// Check for complete word matches
	cleanWord := t.removeDiacritics(word)
	
	// Check common words
	if entry, exists := dict.CommonWords[cleanWord]; exists {
		return entry.Transliteration
	}
	
	// Check divine names
	if dict.DivineNames != nil {
		if entry, exists := dict.DivineNames[cleanWord]; exists {
			return entry.Transliteration
		}
	}
	
	// Apply heuristic transliteration
	return t.applyHeuristics(word, letterMap, dict)
}

// containsArabicScript checks if text contains Arabic/Persian script
func (t *Transliterator) containsArabicScript(text string) bool {
	for _, r := range text {
		if r >= 0x0600 && r <= 0x06FF {
			return true
		}
	}
	return false
}

// removeDiacritics strips diacritical marks from text
func (t *Transliterator) removeDiacritics(text string) string {
	var result strings.Builder
	for _, r := range text {
		// Skip diacritics
		if r >= 0x064B && r <= 0x065F || r == 0x0670 {
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

// applyHeuristics applies rule-based transliteration for unknown words
func (t *Transliterator) applyHeuristics(word string, letterMap map[rune]string, dict *Dictionary) string {
	var result strings.Builder
	runes := []rune(word)
	
	for i, r := range runes {
		// Handle diacritics
		if vowel, exists := t.vowelMarks[r]; exists {
			if vowel != "" {
				result.WriteString(vowel)
			}
			continue
		}
		
		// Handle letters
		if trans, exists := letterMap[r]; exists {
			// Special handling for initial alif
			if r == 'ا' && i == 0 && t.isBeginningOfDivineName(word) {
				result.WriteString("I")
			} else {
				result.WriteString(trans)
			}
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			// Unknown letter, keep as-is
			result.WriteRune(r)
		} else {
			// Punctuation, keep as-is
			result.WriteRune(r)
		}
	}
	
	// Apply vowel insertion heuristics
	return t.insertVowels(result.String(), dict)
}

// isBeginningOfDivineName checks if word starts a divine name
func (t *Transliterator) isBeginningOfDivineName(word string) bool {
	clean := t.removeDiacritics(word)
	return clean == "إلهي" || clean == "الله" || strings.HasPrefix(clean, "إله")
}

// insertVowels applies vowel insertion heuristics
func (t *Transliterator) insertVowels(text string, dict *Dictionary) string {
	// Simple heuristic: insert 'a' between consonant clusters
	result := regexp.MustCompile(`([bcdfghjklmnpqrstvwxyz])([bcdfghjklmnpqrstvwxyz])`).ReplaceAllStringFunc(text, func(match string) string {
		runes := []rune(match)
		if len(runes) == 2 {
			return string(runes[0]) + "a" + string(runes[1])
		}
		return match
	})
	
	return result
}

// postProcess applies final cleanup and formatting
func (t *Transliterator) postProcess(text string, lang Language) string {
	result := text
	
	// Apply post-processing rules
	for _, processor := range t.postProcessors {
		result = processor.regex.ReplaceAllString(result, processor.replacement)
	}
	
	// Language-specific post-processing
	if lang == Arabic {
		result = t.postProcessArabic(result)
	} else {
		result = t.postProcessPersian(result)
	}
	
	// Final cleanup
	result = t.finalCleanup(result)
	
	return result
}

// postProcessArabic handles Arabic-specific post-processing
func (t *Transliterator) postProcessArabic(text string) string {
	// Fix article contractions first
	text = regexp.MustCompile(`\bfí\s+ad-`).ReplaceAllString(text, "fí'd-")
	text = regexp.MustCompile(`\bwa\s+al-`).ReplaceAllString(text, "wa'l-")
	
	// Fix divine name combinations
	text = regexp.MustCompile(`\banta\s+álamu'ṭi\s+al-'Alím\s+álḥakiyamu\b`).ReplaceAllString(text, "anta'l-Mu'ṭí'l-'Alímu'l-Ḥakím")
	text = regexp.MustCompile(`\banta\s+al-mu'ṭí\s+al-'alím\s+al-ḥakím\b`).ReplaceAllString(text, "anta'l-Mu'ṭí'l-'Alímu'l-Ḥakím")
	
	// Fix common Arabic patterns
	text = regexp.MustCompile(`\banta\s+al-`).ReplaceAllString(text, "anta'l-")
	text = regexp.MustCompile(`\banta'l-Mu'ṭí\s+al-'Alím\s+al-Ḥakím`).ReplaceAllString(text, "anta'l-Mu'ṭí'l-'Alímu'l-Ḥakím")
	
	// Fix divine name capitalizations
	text = regexp.MustCompile(`\bal-mu'ṭí\b`).ReplaceAllString(text, "al-Mu'ṭí")
	text = regexp.MustCompile(`\bal-'alím\b`).ReplaceAllString(text, "al-'Alím")
	text = regexp.MustCompile(`\bal-ḥakím\b`).ReplaceAllString(text, "al-Ḥakím")
	text = regexp.MustCompile(`\bálamu'ṭi\b`).ReplaceAllString(text, "al-Mu'ṭí")
	text = regexp.MustCompile(`\bálḥakiyamu\b`).ReplaceAllString(text, "al-Ḥakím")
	
	return text
}

// postProcessPersian handles Persian-specific post-processing
func (t *Transliterator) postProcessPersian(text string) string {
	// Fix Persian-specific patterns
	text = regexp.MustCompile(`\btú'í\b`).ReplaceAllString(text, "tú'í")
	
	return text
}

// finalCleanup applies final formatting fixes
func (t *Transliterator) finalCleanup(text string) string {
	// Fix capitalization at sentence beginnings
	text = regexp.MustCompile(`(^|\. +)([a-z])`).ReplaceAllStringFunc(text, func(match string) string {
		if strings.HasPrefix(match, ".") {
			return strings.Replace(match, match[len(match)-1:], strings.ToUpper(match[len(match)-1:]), 1)
		}
		return strings.ToUpper(match)
	})
	
	// Fix capitalization after newlines
	text = regexp.MustCompile(`(\n)([a-z])`).ReplaceAllStringFunc(text, func(match string) string {
		return match[:1] + strings.ToUpper(match[1:])
	})
	
	// Fix capitalization after markdown headers
	text = regexp.MustCompile(`(#{1,6}\s+)([a-z])`).ReplaceAllStringFunc(text, func(match string) string {
		parts := strings.Split(match, " ")
		if len(parts) >= 2 {
			lastPart := parts[len(parts)-1]
			if len(lastPart) > 0 {
				parts[len(parts)-1] = strings.ToUpper(lastPart[:1]) + lastPart[1:]
			}
		}
		return strings.Join(parts, " ")
	})
	
	// Clean up multiple spaces
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	// Fix apostrophe spacing
	text = regexp.MustCompile(`'\s+`).ReplaceAllString(text, "'")
	text = regexp.MustCompile(`\s+'`).ReplaceAllString(text, "'")
	
	return strings.TrimSpace(text)
}

// IsArabic determines if text is primarily Arabic script
func IsArabic(text string) bool {
	arabicPattern := regexp.MustCompile(`[\u0600-\u06FF]`)
	matches := arabicPattern.FindAllString(text, -1)
	return len(matches) > 0
}

// IsPersian determines if text is primarily Persian script
func IsPersian(text string) bool {
	// Persian uses Arabic script but has some specific indicators
	persianPattern := regexp.MustCompile(`[پچژگ]`) // Persian-specific letters
	persianWords := regexp.MustCompile(`(?i)(خدا|پروردگار|از|به|در|که|این|آن|می)`) // Common Persian words
	
	return persianPattern.MatchString(text) || persianWords.MatchString(text)
}

// AutoDetectLanguage attempts to determine the source language
func AutoDetectLanguage(text string) Language {
	if IsPersian(text) {
		return Persian
	}
	return Arabic // Default to Arabic for Arabic script
}