package transliterator

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Language represents the source language
type Language int

const (
	Arabic Language = iota
	Persian
)

// Dictionary represents the structure of our transliteration dictionaries
type Dictionary struct {
	Metadata struct {
		Version     string `json:"version"`
		Description string `json:"description"`
		LastUpdated string `json:"last_updated"`
	} `json:"metadata"`
	CommonWords          map[string]WordEntry   `json:"common_words"`
	DivineNames          map[string]WordEntry   `json:"divine_names"`
	CommonPhrases        map[string]Pattern     `json:"common_phrases"`
	VowelPatterns        map[string]Pattern     `json:"vowel_patterns"`
	ArticleRules         map[string]interface{} `json:"article_rules"`
	EzafeRules           map[string]interface{} `json:"ezafe_rules"`
	Heuristics           map[string]interface{} `json:"heuristics"`
	VerbalPrefixes       map[string]WordEntry   `json:"verbal_prefixes"`
	Suffixes             map[string]interface{} `json:"suffixes"`
	StressPatterns       map[string]interface{} `json:"stress_patterns"`
	MorphologicalPatterns map[string]interface{} `json:"morphological_patterns"`
	ConsonantChanges     map[string]interface{} `json:"consonant_changes"`
}

// WordEntry represents a dictionary entry
type WordEntry struct {
	Transliteration string `json:"transliteration"`
	Category        string `json:"category"`
	Notes           string `json:"notes"`
	Root            string `json:"root"`
	Meaning         string `json:"meaning"`
}

// Pattern represents a transliteration pattern
type Pattern struct {
	Pattern         string `json:"pattern"`
	Transliteration string `json:"transliteration"`
	Notes           string `json:"notes"`
}

// Rule represents a transliteration rule
type Rule struct {
	Pattern         string   `json:"pattern"`
	Transliteration string   `json:"transliteration"`
	Notes           string   `json:"notes"`
	Examples        []string `json:"examples"`
}

// Transliterator represents the new dictionary-first transliterator
type Transliterator struct {
	arabicDict      *Dictionary
	persianDict     *Dictionary
	arabicLetters   map[rune]string
	persianLetters  map[rune]string
	vowelMarks      map[rune]string
	phraseTokens    map[string]string
	minimalRegexes  []minimalRegex
}

// minimalRegex represents essential regex patterns that cannot be handled by dictionary
type minimalRegex struct {
	regex       *regexp.Regexp
	replacement string
	description string
	essential   bool // true if cannot be replaced by dictionary lookup
}

// New creates a new dictionary-first transliterator
func New() (*Transliterator, error) {
	t := &Transliterator{
		phraseTokens: make(map[string]string),
	}

	// Load dictionaries first
	if err := t.loadDictionaries(); err != nil {
		return nil, fmt.Errorf("failed to load dictionaries: %v", err)
	}

	// Initialize minimal letter mappings (fallback only)
	t.initializeLetterMappings()

	// Initialize essential regex patterns only
	t.initializeEssentialPatterns()

	return t, nil
}

// loadDictionaries loads the JSON dictionaries
func (t *Transliterator) loadDictionaries() error {
	// Load Arabic dictionary
	arabicData, err := os.ReadFile("data/arabic_dictionary.json")
	if err != nil {
		return fmt.Errorf("failed to read Arabic dictionary: %v", err)
	}

	t.arabicDict = &Dictionary{}
	if err := json.Unmarshal(arabicData, t.arabicDict); err != nil {
		return fmt.Errorf("failed to parse Arabic dictionary: %v", err)
	}

	// Load Persian dictionary
	persianData, err := os.ReadFile("data/persian_dictionary.json")
	if err != nil {
		return fmt.Errorf("failed to read Persian dictionary: %v", err)
	}

	t.persianDict = &Dictionary{}
	if err := json.Unmarshal(persianData, t.persianDict); err != nil {
		return fmt.Errorf("failed to parse Persian dictionary: %v", err)
	}

	return nil
}

// initializeLetterMappings sets up basic letter mappings as fallback
func (t *Transliterator) initializeLetterMappings() {
	// Arabic letters (minimal fallback set)
	t.arabicLetters = map[rune]string{
		'ا': "á", 'أ': "a", 'إ': "i", 'آ': "á", 'ب': "b", 'ت': "t", 'ث': "th", 'ج': "j", 'ح': "ḥ", 'خ': "kh",
		'د': "d", 'ذ': "dh", 'ر': "r", 'ز': "z", 'س': "s", 'ش': "sh", 'ص': "ṣ",
		'ض': "ḍ", 'ط': "ṭ", 'ظ': "ẓ", 'ع': "'", 'غ': "gh", 'ف': "f", 'ق': "q",
		'ك': "k", 'ک': "k", 'ل': "l", 'م': "m", 'ن': "n", 'ه': "h", 'و': "w", 'ي': "y",
		'ى': "á", 'ة': "h", 'ء': "'", 'ؤ': "u'", 'ئ': "i'",
	}

	// Persian letters (minimal fallback set)
	t.persianLetters = map[rune]string{
		'ا': "á", 'ب': "b", 'پ': "p", 'ت': "t", 'ث': "th", 'ج': "j", 'چ': "ch",
		'ح': "ḥ", 'خ': "kh", 'د': "d", 'ذ': "dh", 'ر': "r", 'ز': "z", 'ژ': "zh",
		'س': "s", 'ش': "sh", 'ص': "ṣ", 'ض': "ḍ", 'ط': "ṭ", 'ظ': "ẓ", 'ع': "'",
		'غ': "gh", 'ف': "f", 'ق': "q", 'ک': "k", 'گ': "g", 'ل': "l", 'م': "m",
		'ن': "n", 'و': "v", 'ه': "h", 'ی': "í", 'ى': "á", 'ة': "h", 'ء': "'",
	}

	// Diacritics
	t.vowelMarks = map[rune]string{
		'َ': "a", 'ِ': "i", 'ُ': "u", 'ً': "an", 'ٍ': "in", 'ٌ': "un",
		'ْ': "", 'ّ': "", 'ٓ': "", 'ٔ': "", 'ٕ': "",
	}
}

// initializeEssentialPatterns sets up only the most essential regex patterns
func (t *Transliterator) initializeEssentialPatterns() {
	essentialPatterns := []struct {
		pattern     string
		replacement string
		description string
		essential   bool
	}{
		// Only essential formatting that cannot be handled by dictionary
		{`\s+`, " ", "normalize spaces", true},
		{`\s*-\s*`, "-", "normalize hyphens", true},
		{`\s+([,.!?;:])`, "$1", "punctuation spacing", true},
		{`'\s+`, "'", "apostrophe spacing", true},
		{`\s+'`, "'", "apostrophe spacing", true},
		
		// Persian ezafe connector (essential structural element)
		{`‌`, "-", "Persian ezafe connector", true},
		
		// Sentence capitalization (essential formatting) - only after periods with spaces
		{`(\. +)([a-z])`, "${1}${2}", "sentence capitalization", true},
		{`(\n)([a-z])`, "${1}${2}", "line capitalization", true},
	}

	for _, pattern := range essentialPatterns {
		if pattern.essential {
			t.minimalRegexes = append(t.minimalRegexes, minimalRegex{
				regex:       regexp.MustCompile(pattern.pattern),
				replacement: pattern.replacement,
				description: pattern.description,
				essential:   pattern.essential,
			})
		}
	}
}

// Transliterate transliterates text using dictionary-first approach
func (t *Transliterator) Transliterate(text string, lang Language) string {
	// Handle phrase-level patterns first
	text = t.handlePhrasesFromDict(text, lang)
	
	// Process word by word with dictionary priority
	words := strings.Fields(text)
	var result []string
	
	for _, word := range words {
		transliterated := t.transliterateWordV2(word, lang)
		result = append(result, transliterated)
	}
	
	// Join and apply minimal post-processing
	output := strings.Join(result, " ")
	output = t.applyEssentialPostProcessing(output, lang)
	
	return strings.TrimSpace(output)
}

// handlePhrasesFromDict handles multi-word phrases using dictionary data
func (t *Transliterator) handlePhrasesFromDict(text string, lang Language) string {
	var dict *Dictionary
	if lang == Persian {
		dict = t.persianDict
	} else {
		dict = t.arabicDict
	}
	
	// Check for common phrases in dictionary
	if dict.CommonPhrases != nil {
		// Sort phrases by length (longest first) to avoid partial matches
		phrases := make([]string, 0, len(dict.CommonPhrases))
		for phrase := range dict.CommonPhrases {
			phrases = append(phrases, phrase)
		}
		sort.Slice(phrases, func(i, j int) bool {
			return len(phrases[i]) > len(phrases[j])
		})
		
		for _, phrase := range phrases {
			if entry, exists := dict.CommonPhrases[phrase]; exists {
				text = strings.ReplaceAll(text, phrase, entry.Transliteration)
			}
		}
	}
	
	return text
}

// transliterateWordV2 uses dictionary-first approach for word transliteration
func (t *Transliterator) transliterateWordV2(word string, lang Language) string {
	// Skip if no Arabic/Persian script
	if !t.containsArabicScript(word) {
		return word
	}
	
	// Get appropriate dictionary
	var dict *Dictionary
	if lang == Persian {
		dict = t.persianDict
	} else {
		dict = t.arabicDict
	}
	
	// Clean word for dictionary lookup
	cleanWord := t.removeDiacritics(word)
	
	// Priority 1: Exact match in common words
	if entry, exists := dict.CommonWords[cleanWord]; exists {
		return entry.Transliteration
	}
	
	// Priority 2: Exact match in divine names
	if dict.DivineNames != nil {
		if entry, exists := dict.DivineNames[cleanWord]; exists {
			return entry.Transliteration
		}
	}
	
	// Priority 3: Compound word analysis using dictionary
	if compound := t.analyzeCompoundWord(cleanWord, dict); compound != "" {
		return compound
	}
	
	// Priority 4: Morphological analysis using dictionary patterns
	if morphological := t.analyzeMorphology(cleanWord, dict); morphological != "" {
		return morphological
	}
	
	// Priority 5: Fallback to heuristic with dictionary guidance
	return t.dictionaryGuidedHeuristic(word, dict, lang)
}

// analyzeCompoundWord attempts to break down compound words using dictionary
func (t *Transliterator) analyzeCompoundWord(word string, dict *Dictionary) string {
	// Try to match prefixes and suffixes from dictionary
	if dict.VerbalPrefixes != nil {
		for prefix := range dict.VerbalPrefixes {
			if strings.HasPrefix(word, prefix) {
				remainder := strings.TrimPrefix(word, prefix)
				if entry, exists := dict.CommonWords[remainder]; exists {
					prefixTrans := dict.VerbalPrefixes[prefix].Transliteration
					return prefixTrans + "-" + entry.Transliteration
				}
			}
		}
	}
	
	// Try compound patterns for Persian
	if dict == t.persianDict {
		return t.analyzePersianCompound(word, dict)
	}
	
	return ""
}

// analyzePersianCompound handles Persian compound words and ezafe constructions
func (t *Transliterator) analyzePersianCompound(word string, dict *Dictionary) string {
	// Look for ezafe patterns
	if strings.Contains(word, "‌") {
		parts := strings.Split(word, "‌")
		var transliterated []string
		
		for _, part := range parts {
			if entry, exists := dict.CommonWords[part]; exists {
				transliterated = append(transliterated, entry.Transliteration)
			} else {
				// Fallback to heuristic for this part
				transliterated = append(transliterated, t.basicHeuristic(part, t.persianLetters))
			}
		}
		
		return strings.Join(transliterated, "-")
	}
	
	return ""
}

// analyzeMorphology attempts morphological analysis using dictionary patterns
func (t *Transliterator) analyzeMorphology(word string, dict *Dictionary) string {
	// Try to identify root + suffix patterns
	if dict.Suffixes != nil {
		// This would need proper implementation based on dictionary structure
		// For now, return empty to use fallback
	}
	
	return ""
}

// dictionaryGuidedHeuristic uses dictionary patterns to guide heuristic transliteration
func (t *Transliterator) dictionaryGuidedHeuristic(word string, dict *Dictionary, lang Language) string {
	var letterMap map[rune]string
	if lang == Persian {
		letterMap = t.persianLetters
	} else {
		letterMap = t.arabicLetters
	}
	
	// Use dictionary heuristics if available
	if dict.Heuristics != nil {
		// Apply dictionary-based vowel insertion patterns
		return t.applyDictionaryHeuristics(word, dict, letterMap)
	}
	
	// Basic letter-by-letter fallback
	return t.basicHeuristic(word, letterMap)
}

// applyDictionaryHeuristics applies heuristic rules from dictionary
func (t *Transliterator) applyDictionaryHeuristics(word string, dict *Dictionary, letterMap map[rune]string) string {
	result := t.basicHeuristic(word, letterMap)
	
	// Apply vowel patterns from dictionary
	if dict.VowelPatterns != nil {
		for pattern, replacement := range dict.VowelPatterns {
			if pattern != "" && replacement.Transliteration != "" {
				result = strings.ReplaceAll(result, pattern, replacement.Transliteration)
			}
		}
	}
	
	return result
}

// basicHeuristic provides basic letter-by-letter transliteration
func (t *Transliterator) basicHeuristic(word string, letterMap map[rune]string) string {
	var result strings.Builder
	runes := []rune(word)
	
	for _, r := range runes {
		// Handle diacritics
		if vowel, exists := t.vowelMarks[r]; exists {
			if vowel != "" {
				result.WriteString(vowel)
			}
			continue
		}
		
		// Handle letters
		if trans, exists := letterMap[r]; exists {
			result.WriteString(trans)
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune(r)
		}
	}
	
	// Apply statistical vowel insertion as final fallback
	return t.insertStatisticalVowels(result.String())
}

// insertStatisticalVowels adds vowels based on common patterns
func (t *Transliterator) insertStatisticalVowels(consonantString string) string {
	if len(consonantString) == 0 {
		return consonantString
	}
	
	// Statistical vowel insertion rules based on common Arabic/Persian patterns
	runes := []rune(consonantString)
	var result strings.Builder
	
	for i, r := range runes {
		result.WriteRune(r)
		
		// Don't insert vowels after vowels or at the end
		if i == len(runes)-1 || isVowel(r) {
			continue
		}
		
		// Check if next character is a vowel
		if i+1 < len(runes) && isVowel(runes[i+1]) {
			continue
		}
		
		// Insert vowel based on statistical patterns
		vowel := t.guessVowel(r, i, runes)
		if vowel != "" {
			result.WriteString(vowel)
		}
	}
	
	return result.String()
}

// guessVowel provides statistical vowel insertion
func (t *Transliterator) guessVowel(currentRune rune, position int, allRunes []rune) string {
	// Don't insert vowels after long vowels like 'á'
	if currentRune == 'á' || currentRune == 'í' || currentRune == 'ú' {
		return ""
	}
	
	// Common vowel patterns based on position and surrounding consonants
	switch currentRune {
	case 'm':
		// 'm' + 'l' often = "malik" pattern
		if position+1 < len(allRunes) && allRunes[position+1] == 'l' {
			return "a"
		}
		return "a"
	case 'l':
		// 'l' + 'k' often = "lik" pattern
		if position+1 < len(allRunes) && allRunes[position+1] == 'k' {
			return "i"
		}
		return "a"
	case 'n', 'r':
		// These consonants often have 'a' after them
		return "a"
	case 'b', 't', 'd', 'g':
		// These often have 'i' in the middle of words
		if position > 0 && position < len(allRunes)-2 {
			return "i"
		}
		return "a"
	case 'k':
		// Don't add vowel after 'k' at the end of words
		if position == len(allRunes)-1 {
			return ""
		}
		return "a"
	case 'f', 's', 'h':
		// These often have short vowels
		return "a"
	case 'q':
		// Emphatic consonants often have 'a'
		return "a"
	case '\'':
		// Pharyngeal sounds often have 'a'
		return "a"
	default:
		// Check for specific multi-character transliterations
		currentStr := string(currentRune)
		switch currentStr {
		case "sh", "gh", "kh", "ṭ", "ḍ", "ṣ", "ẓ", "ḥ":
			return "a"
		}
		
		// Default to 'a' as it's most common, but not at word end
		if position < len(allRunes)-1 {
			return "a"
		}
		return ""
	}
}

// isVowel checks if a character is a vowel
func isVowel(r rune) bool {
	vowels := "aeiouáíúāīū"
	return strings.ContainsRune(vowels, r)
}

// containsArabicScript checks if word contains Arabic/Persian script
func (t *Transliterator) containsArabicScript(word string) bool {
	for _, r := range word {
		if (r >= 0x0600 && r <= 0x06FF) || (r >= 0x0750 && r <= 0x077F) ||
			(r >= 0xFB50 && r <= 0xFDFF) || (r >= 0xFE70 && r <= 0xFEFF) {
			return true
		}
	}
	return false
}

// removeDiacritics removes diacritical marks for dictionary lookup
func (t *Transliterator) removeDiacritics(word string) string {
	var result strings.Builder
	for _, r := range word {
		if _, isDiacritic := t.vowelMarks[r]; !isDiacritic {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// applyEssentialPostProcessing applies only essential post-processing patterns
func (t *Transliterator) applyEssentialPostProcessing(text string, lang Language) string {
	result := text
	
	// Apply essential patterns only
	for _, processor := range t.minimalRegexes {
		if processor.essential {
			if processor.description == "sentence capitalization" || processor.description == "line capitalization" {
				result = processor.regex.ReplaceAllStringFunc(result, func(match string) string {
					parts := strings.Split(match, " ")
					if len(parts) >= 2 {
						lastPart := parts[len(parts)-1]
						if len(lastPart) > 0 {
							parts[len(parts)-1] = strings.ToUpper(lastPart[:1]) + lastPart[1:]
						}
					}
					return strings.Join(parts, " ")
				})
			} else {
				result = processor.regex.ReplaceAllString(result, processor.replacement)
			}
		}
	}
	
	// Language-specific essential processing
	if lang == Persian {
		result = t.applyPersianEzafe(result)
	}
	
	return result
}

// applyPersianEzafe applies Persian ezafe rules from dictionary
func (t *Transliterator) applyPersianEzafe(text string) string {
	// Use ezafe rules from dictionary if available
	if t.persianDict.EzafeRules != nil {
		// Apply ezafe connector rules
		text = strings.ReplaceAll(text, "‌", "-")
	}
	
	return text
}

// IsArabic checks if text is primarily Arabic
func IsArabic(text string) bool {
	arabicCount := 0
	persianCount := 0
	
	for _, r := range text {
		switch r {
		case 'پ', 'چ', 'ژ', 'گ':
			persianCount++
		case 'ض', 'ص', 'ث', 'ق', 'ف', 'غ', 'ع', 'ه', 'خ', 'ح', 'ج', 'د', 'ذ', 'ر', 'ز', 'س', 'ش', 'ت', 'ط', 'ظ', 'ل', 'ن', 'م', 'ك', 'و', 'ي':
			arabicCount++
		}
	}
	
	return arabicCount > persianCount
}

// IsPersian checks if text is primarily Persian
func IsPersian(text string) bool {
	return !IsArabic(text)
}

// AutoDetectLanguage automatically detects the language
func AutoDetectLanguage(text string) Language {
	if IsArabic(text) {
		return Arabic
	}
	return Persian
}