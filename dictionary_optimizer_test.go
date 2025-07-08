package transliterator

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

// TestCases represents the structure of our JSON test files
type TestCases struct {
	Metadata struct {
		Description string `json:"description"`
		Version     string `json:"version"`
		LastUpdated string `json:"last_updated"`
	} `json:"metadata"`
	TestCases         []TestCase     `json:"test_cases"`
	CommonWordsTest   []WordTest     `json:"common_words_test"`
	HeuristicTestWords []WordTest    `json:"heuristic_test_words"`
	EzafeTestCases    []EzafeTest    `json:"ezafe_test_cases"`
}

type TestCase struct {
	Name     string  `json:"name"`
	Input    string  `json:"input"`
	Expected string  `json:"expected"`
	MinScore float64 `json:"min_score"`
	Category string  `json:"category"`
	Priority string  `json:"priority"`
}

type WordTest struct {
	Word     string `json:"word"`
	Expected string `json:"expected"`
	Category string `json:"category"`
	Notes    string `json:"notes"`
}

type EzafeTest struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
	Notes    string `json:"notes"`
}

// DictionaryOptimizationReport contains suggestions for dictionary improvements
type DictionaryOptimizationReport struct {
	Language              Language
	RedundantWords        []RedundantWord
	SuggestedAdditions    []SuggestedWord
	HeuristicFailures     []HeuristicFailure
	Statistics            OptimizationStats
}

type RedundantWord struct {
	Word              string
	DictionaryResult  string
	HeuristicResult   string
	Confidence        float64
	Reason            string
}

type SuggestedWord struct {
	Word              string
	ExpectedResult    string
	HeuristicResult   string
	Priority          int
	Reason            string
}

type HeuristicFailure struct {
	Word              string
	Expected          string
	HeuristicResult   string
	ErrorType         string
}

type OptimizationStats struct {
	TotalWordsAnalyzed    int
	RedundantWordsFound   int
	SuggestionsGenerated  int
	HeuristicAccuracy     float64
	DictionaryCoverage    float64
}

func TestDictionaryOptimization(t *testing.T) {
	transliterator, err := New()
	if err != nil {
		t.Fatalf("Failed to create transliterator: %v", err)
	}

	// Load test cases from JSON files
	arabicTests, err := loadTestCases("test_cases/arabic_test_cases.json")
	if err != nil {
		t.Logf("Could not load Arabic test cases: %v", err)
		arabicTests = &TestCases{} // Use empty if file doesn't exist
	}

	persianTests, err := loadTestCases("test_cases/persian_test_cases.json")
	if err != nil {
		t.Logf("Could not load Persian test cases: %v", err)
		persianTests = &TestCases{} // Use empty if file doesn't exist
	}

	// Analyze Arabic dictionary
	t.Run("Arabic_Dictionary_Analysis", func(t *testing.T) {
		report := analyzeDictionary(transliterator, Arabic, arabicTests)
		printOptimizationReport(t, report)
	})

	// Analyze Persian dictionary
	t.Run("Persian_Dictionary_Analysis", func(t *testing.T) {
		report := analyzeDictionary(transliterator, Persian, persianTests)
		printOptimizationReport(t, report)
	})
}

func loadTestCases(filename string) (*TestCases, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var testCases TestCases
	err = json.Unmarshal(data, &testCases)
	return &testCases, err
}

func analyzeDictionary(trans *Transliterator, lang Language, tests *TestCases) DictionaryOptimizationReport {
	report := DictionaryOptimizationReport{
		Language: lang,
	}

	var dict *Dictionary
	var letterMap map[rune]string
	
	if lang == Persian {
		dict = trans.persianDict
		letterMap = trans.persianLetters
	} else {
		dict = trans.arabicDict
		letterMap = trans.arabicLetters
	}

	// Collect all words to analyze
	wordsToAnalyze := make(map[string]string) // word -> expected
	
	// Add words from test cases
	for _, wordTest := range tests.CommonWordsTest {
		wordsToAnalyze[wordTest.Word] = wordTest.Expected
	}
	for _, wordTest := range tests.HeuristicTestWords {
		wordsToAnalyze[wordTest.Word] = wordTest.Expected
	}

	// Extract words from dictionary
	for word := range dict.CommonWords {
		if _, exists := wordsToAnalyze[word]; !exists {
			wordsToAnalyze[word] = dict.CommonWords[word].Transliteration
		}
	}

	report.Statistics.TotalWordsAnalyzed = len(wordsToAnalyze)

	// Analyze each word
	correctHeuristics := 0
	for word, expected := range wordsToAnalyze {
		// Get dictionary result
		var dictResult string
		if entry, exists := dict.CommonWords[word]; exists {
			dictResult = entry.Transliteration
		}

		// Get heuristic result (without dictionary lookup)
		heuristicResult := trans.basicHeuristic(word, letterMap)
		heuristicResult = trans.insertStatisticalVowels(heuristicResult)

		// Analyze if dictionary entry is redundant
		if dictResult != "" {
			similarity := calculateSimilarity(dictResult, heuristicResult)
			if similarity > 0.8 { // 80% similarity threshold
				report.RedundantWords = append(report.RedundantWords, RedundantWord{
					Word:             word,
					DictionaryResult: dictResult,
					HeuristicResult:  heuristicResult,
					Confidence:       similarity,
					Reason:          fmt.Sprintf("Heuristic produces %.1f%% similar result", similarity*100),
				})
			}
		}

		// Analyze if word should be added to dictionary
		if dictResult == "" {
			similarity := calculateSimilarity(expected, heuristicResult)
			if similarity < 0.6 { // Less than 60% similarity - heuristic failing
				priority := 1
				if similarity < 0.3 {
					priority = 3 // High priority
				} else if similarity < 0.5 {
					priority = 2 // Medium priority
				}

				report.SuggestedAdditions = append(report.SuggestedAdditions, SuggestedWord{
					Word:            word,
					ExpectedResult:  expected,
					HeuristicResult: heuristicResult,
					Priority:        priority,
					Reason:         fmt.Sprintf("Heuristic only %.1f%% accurate", similarity*100),
				})
			}

			// Track heuristic failures
			if similarity < 0.7 {
				errorType := classifyError(expected, heuristicResult)
				report.HeuristicFailures = append(report.HeuristicFailures, HeuristicFailure{
					Word:            word,
					Expected:        expected,
					HeuristicResult: heuristicResult,
					ErrorType:       errorType,
				})
			} else {
				correctHeuristics++
			}
		} else {
			// Word is in dictionary, count as correct
			correctHeuristics++
		}
	}

	// Calculate statistics
	report.Statistics.RedundantWordsFound = len(report.RedundantWords)
	report.Statistics.SuggestionsGenerated = len(report.SuggestedAdditions)
	report.Statistics.HeuristicAccuracy = float64(correctHeuristics) / float64(len(wordsToAnalyze))
	report.Statistics.DictionaryCoverage = float64(len(dict.CommonWords)) / float64(len(wordsToAnalyze))

	// Sort suggestions by priority
	sort.Slice(report.SuggestedAdditions, func(i, j int) bool {
		return report.SuggestedAdditions[i].Priority > report.SuggestedAdditions[j].Priority
	})

	// Sort redundant words by confidence
	sort.Slice(report.RedundantWords, func(i, j int) bool {
		return report.RedundantWords[i].Confidence > report.RedundantWords[j].Confidence
	})

	return report
}

func calculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	
	// Simple character-based similarity
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)
	
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}
	
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}
	
	// Calculate longest common subsequence ratio
	longer := s1
	shorter := s2
	if len(s2) > len(s1) {
		longer = s2
		shorter = s1
	}
	
	matches := 0
	for i := 0; i < len(shorter); i++ {
		if i < len(longer) && shorter[i] == longer[i] {
			matches++
		}
	}
	
	return float64(matches) / float64(len(longer))
}

func classifyError(expected, actual string) string {
	expected = strings.ToLower(expected)
	actual = strings.ToLower(actual)
	
	if len(actual) < len(expected) {
		return "missing_vowels"
	}
	if len(actual) > len(expected) {
		return "extra_vowels"
	}
	if strings.Contains(expected, "'") && !strings.Contains(actual, "'") {
		return "missing_apostrophes"
	}
	if strings.Contains(expected, "ḥ") && !strings.Contains(actual, "ḥ") {
		return "missing_diacritics"
	}
	return "pattern_mismatch"
}

func printOptimizationReport(t *testing.T, report DictionaryOptimizationReport) {
	langName := "Arabic"
	if report.Language == Persian {
		langName = "Persian"
	}
	
	t.Logf("\n=== %s Dictionary Optimization Report ===", langName)
	
	// Statistics
	t.Logf("\n📊 Statistics:")
	t.Logf("   Words analyzed: %d", report.Statistics.TotalWordsAnalyzed)
	t.Logf("   Dictionary coverage: %.1f%%", report.Statistics.DictionaryCoverage*100)
	t.Logf("   Heuristic accuracy: %.1f%%", report.Statistics.HeuristicAccuracy*100)
	t.Logf("   Redundant entries found: %d", report.Statistics.RedundantWordsFound)
	t.Logf("   Suggestions generated: %d", report.Statistics.SuggestionsGenerated)
	
	// Redundant words (can potentially be removed)
	if len(report.RedundantWords) > 0 {
		t.Logf("\n🗑️  Potentially Redundant Dictionary Entries:")
		for i, word := range report.RedundantWords {
			if i >= 5 { // Limit output
				t.Logf("   ... and %d more", len(report.RedundantWords)-i)
				break
			}
			t.Logf("   '%s': dict='%s' vs heuristic='%s' (%.1f%% similar)",
				word.Word, word.DictionaryResult, word.HeuristicResult, word.Confidence*100)
		}
	}
	
	// Suggested additions (words that need dictionary entries)
	if len(report.SuggestedAdditions) > 0 {
		t.Logf("\n➕ Suggested Dictionary Additions:")
		for i, word := range report.SuggestedAdditions {
			if i >= 10 { // Limit output
				t.Logf("   ... and %d more", len(report.SuggestedAdditions)-i)
				break
			}
			priority := "LOW"
			if word.Priority == 3 {
				priority = "HIGH"
			} else if word.Priority == 2 {
				priority = "MED"
			}
			t.Logf("   [%s] '%s': expected='%s' vs heuristic='%s'",
				priority, word.Word, word.ExpectedResult, word.HeuristicResult)
		}
	}
	
	// Error analysis
	if len(report.HeuristicFailures) > 0 {
		errorTypes := make(map[string]int)
		for _, failure := range report.HeuristicFailures {
			errorTypes[failure.ErrorType]++
		}
		
		t.Logf("\n🔍 Heuristic Error Patterns:")
		for errorType, count := range errorTypes {
			t.Logf("   %s: %d occurrences", errorType, count)
		}
	}
	
	// Recommendations
	t.Logf("\n💡 Recommendations:")
	if report.Statistics.HeuristicAccuracy > 0.8 {
		t.Logf("   ✅ Heuristics performing well (%.1f%% accuracy)", report.Statistics.HeuristicAccuracy*100)
	} else {
		t.Logf("   ⚠️  Heuristics need improvement (%.1f%% accuracy)", report.Statistics.HeuristicAccuracy*100)
	}
	
	if len(report.SuggestedAdditions) > 0 {
		highPriority := 0
		for _, word := range report.SuggestedAdditions {
			if word.Priority >= 2 {
				highPriority++
			}
		}
		t.Logf("   🎯 Add %d high-priority words to dictionary", highPriority)
	}
	
	if len(report.RedundantWords) > 5 {
		t.Logf("   🧹 Consider removing %d redundant dictionary entries", len(report.RedundantWords))
	}
	
	if report.Statistics.DictionaryCoverage < 0.6 {
		t.Logf("   📚 Dictionary coverage low (%.1f%%) - expand vocabulary", report.Statistics.DictionaryCoverage*100)
	}
}

func TestHeuristicQuality(t *testing.T) {
	trans, err := New()
	if err != nil {
		t.Fatalf("Failed to create transliterator: %v", err)
	}

	// Test specific vowel insertion patterns
	testCases := []struct {
		input    string
		expected string
		lang     Language
	}{
		// Arabic patterns
		{"مالك", "malik", Arabic},     // Should get ma-li-k pattern
		{"كتاب", "kitáb", Arabic},     // Should get ki-ta-b pattern  
		{"حكيم", "ḥakím", Arabic},     // Should get ḥa-ki-m pattern
		
		// Persian patterns  
		{"ملکوت", "malakūt", Persian}, // Should handle Persian kaf
		{"کریم", "karīm", Persian},    // Should get ka-rī-m pattern
	}

	t.Logf("\n🧪 Heuristic Quality Analysis:")
	
	for _, tc := range testCases {
		// Get pure heuristic result (no dictionary)
		var letterMap map[rune]string
		if tc.lang == Persian {
			letterMap = trans.persianLetters
		} else {
			letterMap = trans.arabicLetters
		}
		
		result := trans.basicHeuristic(tc.input, letterMap)
		result = trans.insertStatisticalVowels(result)
		
		similarity := calculateSimilarity(tc.expected, result)
		status := "❌"
		if similarity > 0.8 {
			status = "✅"
		} else if similarity > 0.6 {
			status = "⚠️"
		}
		
		t.Logf("   %s '%s': expected='%s' got='%s' (%.1f%% match)",
			status, tc.input, tc.expected, result, similarity*100)
	}
}