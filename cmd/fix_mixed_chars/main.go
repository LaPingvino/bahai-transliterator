package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/LaPingvino/bahai-transliterator"
)

type Config struct {
	DatabasePath string
	DryRun       bool
	BatchSize    int
	Language     string
}

type FixRecord struct {
	Version         string
	SourceID        string
	Language        string
	OriginalText    string
	CleanedText     string
	HasMixedChars   bool
	ArabicChars     []string
}

func main() {
	var config Config
	
	flag.StringVar(&config.DatabasePath, "db", "", "Path to the bahaiwritings database directory")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be fixed without making changes")
	flag.IntVar(&config.BatchSize, "batch-size", 20, "Number of records to process in each batch")
	flag.StringVar(&config.Language, "lang", "both", "Language to fix: 'fa-translit', 'ar-translit', or 'both'")
	flag.Parse()

	if config.DatabasePath == "" {
		fmt.Println("Usage: fix_mixed_chars -db /path/to/bahaiwritings")
		fmt.Println("  -db string       Path to the bahaiwritings database directory")
		fmt.Println("  -dry-run         Show what would be fixed without making changes")
		fmt.Println("  -batch-size int  Number of records to process in each batch (default 20)")
		fmt.Println("  -lang string     Language to fix: 'fa-translit', 'ar-translit', or 'both' (default 'both')")
		os.Exit(1)
	}

	if err := fixMixedCharacters(config); err != nil {
		log.Fatalf("Error fixing mixed characters: %v", err)
	}
}

func fixMixedCharacters(config Config) error {
	fmt.Printf("Analyzing transliterations for mixed characters in database: %s\n", config.DatabasePath)
	fmt.Printf("Language filter: %s\n", config.Language)
	fmt.Printf("Dry run: %t\n", config.DryRun)

	// Get all transliteration records
	records, err := getTransliterationRecords(config.DatabasePath, config.Language)
	if err != nil {
		return fmt.Errorf("failed to get records: %v", err)
	}

	fmt.Printf("Found %d transliteration records to analyze\n", len(records))

	// Analyze and fix mixed characters
	var problemRecords []FixRecord
	for _, record := range records {
		fixRecord := analyzeMixedCharacters(record)
		if fixRecord.HasMixedChars {
			problemRecords = append(problemRecords, fixRecord)
		}
	}

	fmt.Printf("\nFound %d records with mixed characters\n", len(problemRecords))

	if len(problemRecords) == 0 {
		fmt.Println("No mixed character issues found!")
		return nil
	}

	// Show examples of problems found
	fmt.Println("\nExamples of mixed character issues:")
	for i, record := range problemRecords {
		if i >= 5 { // Show only first 5 examples
			break
		}
		fmt.Printf("\nRecord %s (%s):\n", record.SourceID, record.Language)
		fmt.Printf("  Arabic chars found: %v\n", record.ArabicChars)
		fmt.Printf("  Original: %s\n", truncateString(record.OriginalText, 100))
		fmt.Printf("  Cleaned:  %s\n", truncateString(record.CleanedText, 100))
	}

	if config.DryRun {
		fmt.Printf("\nDry run complete. %d records would be updated.\n", len(problemRecords))
		return nil
	}

	// Initialize transliterator for re-processing
	t, err := transliterator.New()
	if err != nil {
		return fmt.Errorf("failed to initialize transliterator: %v", err)
	}

	// Update records with properly cleaned transliterations
	fmt.Printf("\nUpdating %d records...\n", len(problemRecords))
	updatedCount := 0

	for i, record := range problemRecords {
		if i%config.BatchSize == 0 {
			fmt.Printf("Processing batch %d-%d...\n", i+1, min(i+config.BatchSize, len(problemRecords)))
		}

		// Get the original text and re-transliterate it properly
		originalText, err := getOriginalText(config.DatabasePath, record.SourceID, record.Language)
		if err != nil {
			fmt.Printf("  Warning: Could not get original text for %s: %v\n", record.SourceID, err)
			continue
		}

		var lang transliterator.Language
		if strings.HasPrefix(record.Language, "fa") {
			lang = transliterator.Persian
		} else {
			lang = transliterator.Arabic
		}

		// Re-transliterate with fixed algorithm
		newTranslit := t.Transliterate(originalText, lang)
		
		// Clean any remaining mixed characters as fallback
		cleanedTranslit := cleanMixedCharacters(newTranslit)

		fmt.Printf("  Updating %s\n", record.SourceID)
		if err := updateRecord(config.DatabasePath, record.Version, cleanedTranslit); err != nil {
			return fmt.Errorf("failed to update record %s: %v", record.Version, err)
		}
		updatedCount++
	}

	fmt.Printf("\nSuccessfully updated %d records\n", updatedCount)

	if !config.DryRun {
		fmt.Println("\nCommitting changes to database...")
		if err := commitChanges(config.DatabasePath); err != nil {
			return fmt.Errorf("failed to commit changes: %v", err)
		}
	}

	return nil
}

func getTransliterationRecords(dbPath, langFilter string) ([]FixRecord, error) {
	var whereClause string
	if langFilter == "both" {
		whereClause = "WHERE language IN ('fa-translit', 'ar-translit')"
	} else if langFilter == "fa-translit" || langFilter == "ar-translit" {
		whereClause = fmt.Sprintf("WHERE language = '%s'", langFilter)
	} else {
		return nil, fmt.Errorf("invalid language filter: %s", langFilter)
	}

	query := fmt.Sprintf(`SELECT version, source_id, language, text FROM writings %s ORDER BY source_id`, whereClause)

	cmd := exec.Command("dolt", "sql", "-q", query, "-r", "csv")
	cmd.Dir = dbPath
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute dolt sql: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(string(output)))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV output: %v", err)
	}

	var records []FixRecord
	// Skip header row
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < 4 {
			continue
		}
		
		record := FixRecord{
			Version:      rows[i][0],
			SourceID:     rows[i][1],
			Language:     rows[i][2],
			OriginalText: rows[i][3],
		}
		records = append(records, record)
	}

	return records, nil
}

func analyzeMixedCharacters(record FixRecord) FixRecord {
	fixRecord := record
	fixRecord.CleanedText = cleanMixedCharacters(record.OriginalText)
	fixRecord.HasMixedChars = record.OriginalText != fixRecord.CleanedText
	
	if fixRecord.HasMixedChars {
		fixRecord.ArabicChars = findArabicCharacters(record.OriginalText)
	}
	
	return fixRecord
}

func findArabicCharacters(text string) []string {
	var arabicChars []string
	seen := make(map[string]bool)
	
	for _, r := range text {
		if isArabicScript(r) {
			char := string(r)
			if !seen[char] {
				arabicChars = append(arabicChars, char)
				seen[char] = true
			}
		}
	}
	
	return arabicChars
}

func isArabicScript(r rune) bool {
	// Check if character is in Arabic or Arabic supplement blocks
	return (r >= 0x0600 && r <= 0x06FF) || // Arabic
		   (r >= 0x0750 && r <= 0x077F) || // Arabic Supplement
		   (r >= 0x08A0 && r <= 0x08FF) || // Arabic Extended-A
		   (r >= 0xFB50 && r <= 0xFDFF) || // Arabic Presentation Forms-A
		   (r >= 0xFE70 && r <= 0xFEFF)    // Arabic Presentation Forms-B
}

func cleanMixedCharacters(text string) string {
	// Define mappings for common Arabic characters that might appear in transliterations
	arabicToLatin := map[rune]string{
		'ی': "i",     // Persian ye
		'ا': "a",     // Alif
		'ع': "'",     // Ain
		'ح': "h",     // Ha
		'خ': "kh",    // Kha
		'د': "d",     // Dal
		'ذ': "dh",    // Dhal
		'ر': "r",     // Ra
		'ز': "z",     // Zay
		'س': "s",     // Sin
		'ش': "sh",    // Shin
		'ص': "s",     // Sad
		'ض': "d",     // Dad
		'ط': "t",     // Ta
		'ظ': "z",     // Za
		'غ': "gh",    // Ghain
		'ف': "f",     // Fa
		'ق': "q",     // Qaf
		'ک': "k",     // Kaf
		'گ': "g",     // Gaf (Persian)
		'ل': "l",     // Lam
		'م': "m",     // Meem
		'ن': "n",     // Noon
		'ه': "h",     // He
		'و': "w",     // Waw
		'ء': "'",     // Hamza
		'ؤ': "u'",    // Waw with hamza
		'ئ': "i'",    // Ya with hamza
		'ة': "h",     // Ta marbuta
		'آ': "a",     // Alif with madda
		'أ': "a",     // Alif with hamza above
		'إ': "i",     // Alif with hamza below
		'ژ': "zh",    // Zhe (Persian)
		'چ': "ch",    // Che (Persian)
		'پ': "p",     // Pe (Persian)
		'ڤ': "v",     // Ve (Persian)
		'َ': "a",     // Fatha
		'ِ': "i",     // Kasra
		'ُ': "u",     // Damma
		'ً': "an",    // Tanween fath
		'ٍ': "in",    // Tanween kasr
		'ٌ': "un",    // Tanween damm
		'ْ': "",      // Sukun
		'ّ': "",      // Shadda
		'ٰ': "a",     // Alif khanjariya
	}

	var result strings.Builder
	for _, r := range text {
		if isArabicScript(r) {
			if replacement, exists := arabicToLatin[r]; exists {
				result.WriteString(replacement)
			}
			// If no mapping exists, skip the character
		} else {
			result.WriteRune(r)
		}
	}

	// Clean up multiple spaces and other artifacts
	cleaned := result.String()
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)
	
	return cleaned
}

func getOriginalText(dbPath, sourceID, translitLang string) (string, error) {
	var originalLang string
	if strings.HasPrefix(translitLang, "fa") {
		originalLang = "fa"
	} else {
		originalLang = "ar"
	}

	query := fmt.Sprintf(`SELECT text FROM writings WHERE source_id = '%s' AND language = '%s' LIMIT 1`, sourceID, originalLang)

	cmd := exec.Command("dolt", "sql", "-q", query, "-r", "csv")
	cmd.Dir = dbPath
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute dolt sql: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(string(output)))
	rows, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to parse CSV output: %v", err)
	}

	if len(rows) < 2 || len(rows[1]) < 1 {
		return "", fmt.Errorf("no original text found for source_id %s", sourceID)
	}

	return rows[1][0], nil
}

func updateRecord(dbPath, version, newText string) error {
	// Escape single quotes in the text for SQL
	escapedText := strings.ReplaceAll(newText, "'", "''")
	
	query := fmt.Sprintf(`UPDATE writings SET text = '%s' WHERE version = '%s'`, escapedText, version)

	cmd := exec.Command("dolt", "sql", "-q", query)
	cmd.Dir = dbPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update record: %v, output: %s", err, string(output))
	}

	return nil
}

func commitChanges(dbPath string) error {
	commands := [][]string{
		{"dolt", "add", "."},
		{"dolt", "commit", "-m", "Fix mixed Arabic characters in transliterations"},
		{"dolt", "push"},
	}

	for _, cmdArgs := range commands {
		fmt.Printf("Executing: %s\n", strings.Join(cmdArgs, " "))
		
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = dbPath
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to execute %s: %v, output: %s", strings.Join(cmdArgs, " "), err, string(output))
		}
		
		if len(output) > 0 {
			fmt.Printf("Output: %s\n", string(output))
		}
	}

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}