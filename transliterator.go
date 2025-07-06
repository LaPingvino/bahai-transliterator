package transliterator

import (
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

// Transliterator handles Arabic and Persian to Bahai transliteration
type Transliterator struct {
	// Arabic letter mappings
	arabicLetters map[rune]string
	// Persian letter mappings
	persianLetters map[rune]string
	// Vowel diacritics
	vowelMarks map[rune]string
	// Common word patterns
	arabicWords map[string]string
	persianWords map[string]string
	// Regex patterns for post-processing
	patterns []patternRule
}

type patternRule struct {
	regex       *regexp.Regexp
	replacement string
	description string
}

// New creates a new Transliterator with predefined rules
func New() *Transliterator {
	t := &Transliterator{
		arabicLetters:  make(map[rune]string),
		persianLetters: make(map[rune]string),
		vowelMarks:     make(map[rune]string),
		arabicWords:    make(map[string]string),
		persianWords:   make(map[string]string),
	}
	t.initializeMappings()
	return t
}

func (t *Transliterator) initializeMappings() {
	// Arabic letter mappings
	t.arabicLetters = map[rune]string{
		'ا': "á",  // alif
		'أ': "a",  // hamza above alif
		'إ': "i",  // hamza below alif
		'آ': "á",  // madda alif
		'ب': "b",  // ba
		'ت': "t",  // ta
		'ث': "th", // tha
		'ج': "j",  // jim
		'ح': "ḥ",  // ha (emphatic)
		'خ': "kh", // kha
		'د': "d",  // dal
		'ذ': "dh", // dhal
		'ر': "r",  // ra
		'ز': "z",  // zayn
		'س': "s",  // sin
		'ش': "sh", // shin
		'ص': "ṣ",  // sad (emphatic)
		'ض': "ḍ",  // dad (emphatic)
		'ط': "ṭ",  // ta (emphatic)
		'ظ': "ẓ",  // za (emphatic)
		'ع': "'",  // ayn
		'غ': "gh", // ghayn
		'ف': "f",  // fa
		'ق': "q",  // qaf
		'ك': "k",  // kaf
		'ک': "k",  // kaf (Persian form)
		'ل': "l",  // lam
		'م': "m",  // mim
		'ن': "n",  // nun
		'ه': "h",  // ha
		'و': "w",  // waw
		'ي': "y",  // ya
		'ى': "á",  // alif maqsura
		'ئ': "'",  // hamza on ya
		'ؤ': "'",  // hamza on waw
		'ة': "h",  // ta marbuta
	}

	// Persian letter mappings (inherits Arabic + specific changes)
	t.persianLetters = make(map[rune]string)
	for k, v := range t.arabicLetters {
		t.persianLetters[k] = v
	}
	
	// Persian-specific letters
	t.persianLetters['پ'] = "p"  // pe
	t.persianLetters['چ'] = "ch" // che
	t.persianLetters['ژ'] = "zh" // zhe
	t.persianLetters['گ'] = "g"  // gaf
	
	// Persian pronunciation differences
	t.persianLetters['ث'] = "s"  // se (not th)
	t.persianLetters['ح'] = "h"  // he (not emphatic)
	t.persianLetters['ذ'] = "z"  // zal (not dh)
	t.persianLetters['ص'] = "s"  // sad (not emphatic)
	t.persianLetters['ض'] = "z"  // zad (not emphatic)
	t.persianLetters['ط'] = "t"  // ta (not emphatic)
	t.persianLetters['ظ'] = "z"  // za (not emphatic)
	t.persianLetters['و'] = "v"  // vav (not w)
	t.persianLetters['ی'] = "í"  // ye
	t.persianLetters['ي'] = "í"  // ye

	// Vowel marks
	t.vowelMarks = map[rune]string{
		'َ': "a",  // fatha
		'ِ': "i",  // kasra
		'ُ': "u",  // damma
		'ْ': "",   // sukun
		'ً': "an", // tanwin fath
		'ٍ': "in", // tanwin kasr
		'ٌ': "un", // tanwin damm
		'ّ': "",   // shadda (handled specially)
		'ٰ': "á",  // alif khanjariya
	}

	// Arabic common words
	t.arabicWords = map[string]string{
		"الله":       "Alláh",
		"إله":       "iláh",
		"إلهي":      "Iláhí",
		"يا":        "yá",
		"أشهد":      "ashhadu",
		"بأنك":      "bi-annaka",
		"بأنه":      "bi-annahu",
		"وأنت":      "wa-anta",
		"إنك":       "innaka",
		"إنه":       "innahu",
		"أنت":       "anta",
		"هو":        "huwa",
		"المهيمن":    "al-Muhaymín",
		"القيوم":     "al-Qayyúm",
		"العليم":     "al-'Alím",
		"الحكيم":     "al-Ḥakím",
		"الغفور":     "al-Ghafúr",
		"الكريم":     "al-Karím",
		"الرحمن":     "ar-Raḥmán",
		"الرحيم":     "ar-Raḥím",
		"المقتدر":    "al-Muqtadir",
		"القدير":     "al-Qadír",
		"العزيز":     "al-'Azíz",
		"السلطان":    "as-Sulṭán",
		"بسمه":       "bismihi",
		"الأسماء":    "al-Asmá'",
		"العظمة":     "al-'Aẓamah",
		"الاقتدار":   "al-Iqtidár",
		"الفردوس":   "al-Firdaws",
		"البقاء":     "al-Baqá'",
		"المخلصين":   "al-Mukhlliṣín",
		"الموحدين":   "al-Muwaḥḥidín",
		"العالمين":   "al-'Álamín",
		"العارفين":   "al-'Árifín",
		"المعطي":     "al-Mu'ṭí",
		"لا":        "lá",
		"إلا":       "illá",
		"من":        "min",
		"إلى":       "ilá",
		"على":       "'alá",
		"في":        "fí",
		"عن":        "'an",
		"مع":        "ma'a",
		"كل":        "kull",
		"جميع":      "jamí'",
		"أحمد":      "Aḥmad",
		"لوح":       "Lawḥ",
		"شفاء":      "shifá'",
		"دواء":      "dawá'",
		"رجاء":      "rajá'",
		"حب":        "ḥubb",
		"رحمة":      "raḥmah",
		"طبيب":      "ṭabíb",
		"معين":      "mu'ín",
		"الدنيا":    "ad-dunyá",
		"الآخرة":    "al-ákhirah",
		"اسم":       "ism",
		"اسمك":      "ismuka",
		"شفائي":     "shifá'í",
		"ذكرك":      "dhikruka",
		"دوائي":     "dawá'í",
		"قربك":      "qurbuka",
		"رجائي":     "rajá'í",
		"حبك":       "ḥubbuka",
		"مؤنسي":     "mu'nisí",
		"رحمتك":     "raḥmatuka",
		"طبيبي":     "ṭabíbí",
		"معيني":     "mu'íní",
		"وإنك":      "wa-innaka",
		"يا إلهي":   "yá Iláhí,",
	}

	// Persian common words
	t.persianWords = map[string]string{
		"خدا":       "Khudá",
		"خداوند":    "Khudávand",
		"پروردگار":  "Parvardigár",
		"از":        "az",
		"به":        "bih",
		"در":        "dar",
		"تا":        "tá",
		"با":        "bá",
		"بر":        "bar",
		"که":        "kih",
		"را":        "rá",
		"و":         "va",
		"یا":        "yá",
		"ای":        "ay",
		"تو":        "tú",
		"من":        "man",
		"او":        "ú",
		"این":       "ín",
		"آن":        "án",
		"می":        "mí",
		"خواهد":     "kháhad",
		"است":       "ast",
		"بود":       "búd",
		"نموده":     "namúdih",
		"کرده":      "kardih",
		"فرموده":    "farmúdih",
		"بگو":       "Bigú",
		"گواهی":     "guvāhī",
		"شهادت":     "shahādaat",
		"یکتا":      "yaktá",
		"وحدانیت":   "vaḥdāniyyat",
		"فردانیت":   "fardāniyyat",
		"مالک":      "mālik",
		"ملکوت":     "malakūt",
		"سلطان":     "sulṭān",
		"غیب":       "ghayb",
		"شهود":      "shuhūd",
		"مسکین":     "miskīn",
		"بحر":       "baḥr",
		"غنا":       "ghaná",
		"کریم":      "karīm",
		"رحیم":      "raḥīm",
		"بخشنده":    "bakhshindih",
		"توانا":     "tavāná",
		"دانا":      "dāná",
		"بینا":      "bīná",
		"جان":       "ján",
		"روان":      "ruvān",
		"لسان":      "lisān",
		"واحد":      "vāḥid",
		"فقیر":      "faqīr",
		"سائل":      "sā'il",
		"کنیز":      "kanīz",
		"محبوب":     "maḥbūb",
		"سید":       "sayyid",
		"سند":       "sanad",
		"مقصود":     "maqṣūd",
		"ایران":     "Īrān",
	}

	// Post-processing patterns
	t.patterns = []patternRule{
		// Fix article combinations
		{regexp.MustCompile(`\bwa\s+al-`), "wa'l-", "wa + al"},
		{regexp.MustCompile(`\bwa\s+a([tṭdḍrzsṣshnjl])-`), "wa'$1-", "wa + sun letters"},
		{regexp.MustCompile(`\bbi\s+al-`), "bi'l-", "bi + al"},
		{regexp.MustCompile(`\bfī\s+al-`), "fī'l-", "fī + al"},
		{regexp.MustCompile(`\bfī\s+a([tṭdḍrzsṣshnjl])-`), "fī'$1-", "fī + sun letters"},
		{regexp.MustCompile(`\bka\s+al-`), "ka'l-", "ka + al"},
		{regexp.MustCompile(`\bli\s+al-`), "li'l-", "li + al"},
		{regexp.MustCompile(`\bmin\s+al-`), "mina'l-", "min + al"},
		{regexp.MustCompile(`\bilá\s+al-`), "ilá'l-", "ilá + al"},
		{regexp.MustCompile(`\b'alá\s+al-`), "'alá'l-", "'alá + al"},
		{regexp.MustCompile(`\b'an\s+al-`), "'ani'l-", "'an + al"},
		
		// Fix definite article with sun letters
		{regexp.MustCompile(`\bal-([tṭdḍrzsṣshnjl])`), "a$1-$1", "sun letters"},
		
		// Fix wa- connections
		{regexp.MustCompile(`\bwa([aáiuūíē])`), "wa-$1", "wa + vowel"},
		{regexp.MustCompile(`\bwa([bcdfghjklmnpqrstvwxyz])`), "wa-$1", "wa + consonant"},
		
		// Fix common phrases
		{regexp.MustCompile(`\blá\s+iláha\s+illá\b`), "lá iláha illá", "no god but"},
		{regexp.MustCompile(`\bAllāh\b`), "Alláh", "Allah"},
		
		// Fix Persian ezafe
		{regexp.MustCompile(`\b(\w+)-i\s+(\w+)`), "$1-i $2", "ezafe connection"},
		
		// Clean up spacing
		{regexp.MustCompile(`\s+`), " ", "normalize spaces"},
		{regexp.MustCompile(`\s*-\s*`), "-", "normalize hyphens"},
		
		// Fix capitalization
		{regexp.MustCompile(`\b([a-z])`), "${1}", "word boundaries"},
	}
}

// Transliterate converts Arabic or Persian text to Bahai transliteration
func (t *Transliterator) Transliterate(text string, lang Language) string {
	// Preserve formatting markers
	text = t.preserveFormatting(text)
	
	// Split into words and process each
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

func (t *Transliterator) preserveFormatting(text string) string {
	// Handle markdown headers
	text = regexp.MustCompile(`^(#{1,6})\s*`).ReplaceAllString(text, "$1 ")
	
	// Handle parentheses and brackets
	text = regexp.MustCompile(`\*\s*\(`).ReplaceAllString(text, "*(")
	text = regexp.MustCompile(`\)\s*\*`).ReplaceAllString(text, ")*")
	
	// Handle quotes
	text = regexp.MustCompile(`"\s*"`).ReplaceAllString(text, `""`)
	
	return text
}

func (t *Transliterator) transliterateWord(word string, lang Language) string {
	// Handle pure formatting
	if !t.containsArabicScript(word) {
		return word
	}
	
	// Check for complete word matches first
	var wordMap map[string]string
	var letterMap map[rune]string
	
	if lang == Persian {
		wordMap = t.persianWords
		letterMap = t.persianLetters
	} else {
		wordMap = t.arabicWords
		letterMap = t.arabicLetters
	}
	
	// Clean word for lookup (remove diacritics)
	cleanWord := t.removeDiacritics(word)
	if translation, exists := wordMap[cleanWord]; exists {
		return translation
	}
	
	// Transliterate letter by letter with context
	return t.transliterateLetter(word, letterMap)
}

func (t *Transliterator) containsArabicScript(text string) bool {
	for _, r := range text {
		if r >= 0x0600 && r <= 0x06FF {
			return true
		}
	}
	return false
}

func (t *Transliterator) removeDiacritics(text string) string {
	var result strings.Builder
	for _, r := range text {
		// Skip diacritics
		if r >= 0x064B && r <= 0x065F {
			continue
		}
		if r == 0x0670 { // alif khanjariya
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

func (t *Transliterator) transliterateLetter(word string, letterMap map[rune]string) string {
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
			// Special handling for alif at beginning
			if r == 'ا' && i == 0 {
				// Check if it should be 'I' for divine names
				if t.isBeginningOfDivineName(word) {
					result.WriteString("I")
				} else {
					result.WriteString(trans)
				}
			} else {
				result.WriteString(trans)
			}
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			// Unknown Arabic letter, keep as is
			result.WriteRune(r)
		} else {
			// Punctuation, keep as is
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

func (t *Transliterator) isBeginningOfDivineName(word string) bool {
	clean := t.removeDiacritics(word)
	return clean == "إلهي" || clean == "الله" || strings.HasPrefix(clean, "إله")
}

func (t *Transliterator) postProcess(text string, lang Language) string {
	result := text
	
	// Apply pattern rules
	for _, pattern := range t.patterns {
		result = pattern.regex.ReplaceAllString(result, pattern.replacement)
	}
	
	// Language-specific post-processing
	if lang == Persian {
		result = t.postProcessPersian(result)
	} else {
		result = t.postProcessArabic(result)
	}
	
	// Final cleanup
	result = t.finalCleanup(result)
	
	return result
}

func (t *Transliterator) postProcessPersian(text string) string {
	// Handle Persian-specific patterns
	text = regexp.MustCompile(`\bmí\s+`).ReplaceAllString(text, "mí-")
	text = regexp.MustCompile(`\bkhváhad\s+`).ReplaceAllString(text, "kháhad ")
	
	return text
}

func (t *Transliterator) postProcessArabic(text string) string {
	// Handle Arabic-specific patterns
	text = regexp.MustCompile(`\bwa\s+`).ReplaceAllString(text, "wa-")
	text = regexp.MustCompile(`\bbi\s+`).ReplaceAllString(text, "bi-")
	
	// Fix specific divine name capitalizations
	text = regexp.MustCompile(`\bal-mu'ṭí\b`).ReplaceAllString(text, "al-Mu'ṭí")
	text = regexp.MustCompile(`\bal-'alím\b`).ReplaceAllString(text, "al-'Alím")
	text = regexp.MustCompile(`\bal-ḥakím\b`).ReplaceAllString(text, "al-Ḥakím")
	
	// Fix article + divine name combinations
	text = regexp.MustCompile(`\banta\s+al-`).ReplaceAllString(text, "anta'l-")
	
	// Fix double vowels
	text = regexp.MustCompile(`aa+`).ReplaceAllString(text, "a")
	text = regexp.MustCompile(`áá+`).ReplaceAllString(text, "á")
	text = regexp.MustCompile(`ii+`).ReplaceAllString(text, "i")
	text = regexp.MustCompile(`uu+`).ReplaceAllString(text, "u")
	text = regexp.MustCompile(`ūū+`).ReplaceAllString(text, "ū")
	text = regexp.MustCompile(`íí+`).ReplaceAllString(text, "í")
	text = regexp.MustCompile(`ēē+`).ReplaceAllString(text, "ē")
	
	// Fix 'iy suffix to 'í
	text = regexp.MustCompile(`([aáiuūíē])'?iy\b`).ReplaceAllString(text, "$1'í")
	
	return text
}

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