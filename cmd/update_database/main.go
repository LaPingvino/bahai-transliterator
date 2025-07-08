package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/LaPingvino/bahai-transliterator"
)

type Config struct {
	DatabasePath string
	DryRun       bool
	BatchSize    int
	Language     string
}

func main() {
	var config Config
	
	flag.StringVar(&config.DatabasePath, "db", "", "Path to the bahaiwritings database directory")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be updated without making changes")
	flag.IntVar(&config.BatchSize, "batch-size", 10, "Number of records to process in each batch")
	flag.StringVar(&config.Language, "lang", "both", "Language to update: 'fa', 'ar', or 'both'")
	flag.Parse()

	if config.DatabasePath == "" {
		fmt.Println("Usage: update_database -db /path/to/bahaiwritings")
		fmt.Println("  -db string       Path to the bahaiwritings database directory")
		fmt.Println("  -dry-run         Show what would be updated without making changes")
		fmt.Println("  -batch-size int  Number of records to process in each batch (default 10)")
		fmt.Println("  -lang string     Language to update: 'fa', 'ar', or 'both' (default 'both')")
		os.Exit(1)
	}

	if err := updateDatabase(config); err != nil {
		log.Fatalf("Error updating database: %v", err)
	}
}

func updateDatabase(config Config) error {
	// Initialize transliterator
	t, err := transliterator.New()
	if err != nil {
		return fmt.Errorf("failed to initialize transliterator: %v", err)
	}

	// Connect to database using dolt
	dbPath := filepath.Join(config.DatabasePath, ".dolt")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database path %s does not appear to be a dolt repository", config.DatabasePath)
	}

	// Use dolt sql commands for database operations
	fmt.Printf("Updating transliterations in database: %s\n", config.DatabasePath)
	fmt.Printf("Language filter: %s\n", config.Language)
	fmt.Printf("Dry run: %t\n", config.DryRun)
	fmt.Printf("Batch size: %d\n", config.BatchSize)

	// Process Persian if requested
	if config.Language == "fa" || config.Language == "both" {
		if err := updateLanguage(t, config, "fa", "fa-translit"); err != nil {
			return fmt.Errorf("failed to update Persian: %v", err)
		}
	}

	// Process Arabic if requested
	if config.Language == "ar" || config.Language == "both" {
		if err := updateLanguage(t, config, "ar", "ar-translit"); err != nil {
			return fmt.Errorf("failed to update Arabic: %v", err)
		}
	}

	if !config.DryRun {
		fmt.Println("\nCommitting changes to database...")
		if err := commitChanges(config.DatabasePath); err != nil {
			return fmt.Errorf("failed to commit changes: %v", err)
		}
	}

	return nil
}

func updateLanguage(t *transliterator.Transliterator, config Config, sourceLang, targetLang string) error {
	fmt.Printf("\n=== Processing %s -> %s ===\n", sourceLang, targetLang)

	// Get records to process
	records, err := getRecordsToUpdate(config.DatabasePath, sourceLang, targetLang)
	if err != nil {
		return fmt.Errorf("failed to get records: %v", err)
	}

	fmt.Printf("Found %d records to process\n", len(records))

	var lang transliterator.Language
	if sourceLang == "fa" {
		lang = transliterator.Persian
	} else {
		lang = transliterator.Arabic
	}

	updatedCount := 0
	unchangedCount := 0

	for i, record := range records {
		if i%config.BatchSize == 0 {
			fmt.Printf("Processing batch %d-%d...\n", i+1, min(i+config.BatchSize, len(records)))
		}

		// Transliterate the text
		newTranslit := t.Transliterate(record.Text, lang)

		// Check if it's different from current
		if newTranslit != record.CurrentTranslit {
			fmt.Printf("  Updating %s (source_id: %s)\n", record.Name, record.SourceID)
			if !config.DryRun {
				if err := updateRecord(config.DatabasePath, record.Version, newTranslit); err != nil {
					return fmt.Errorf("failed to update record %s: %v", record.Version, err)
				}
			}
			updatedCount++
		} else {
			unchangedCount++
		}
	}

	fmt.Printf("\nSummary for %s:\n", sourceLang)
	fmt.Printf("  Updated: %d records\n", updatedCount)
	fmt.Printf("  Unchanged: %d records\n", unchangedCount)
	fmt.Printf("  Total: %d records\n", len(records))

	return nil
}

type Record struct {
	Version         string
	SourceID        string
	Name            string
	Text            string
	CurrentTranslit string
}

func getRecordsToUpdate(dbPath, sourceLang, targetLang string) ([]Record, error) {
	// Use dolt sql to get records
	query := fmt.Sprintf(`SELECT w2.version, w1.source_id, COALESCE(w1.name, '') as name, w1.text, w2.text as current_translit FROM writings w1 JOIN writings w2 ON w1.source = w2.source AND w1.source_id = w2.source_id WHERE w1.language = '%s' AND w2.language = '%s' ORDER BY w1.source_id`, sourceLang, targetLang)

	cmd := exec.Command("dolt", "sql", "-q", query, "-r", "csv")
	cmd.Dir = dbPath
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute dolt sql: %v", err)
	}

	// Parse CSV output
	reader := csv.NewReader(strings.NewReader(string(output)))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV output: %v", err)
	}

	var records []Record
	// Skip header row
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < 5 {
			continue // Skip malformed rows
		}
		
		record := Record{
			Version:         rows[i][0],
			SourceID:        rows[i][1],
			Name:            rows[i][2],
			Text:            rows[i][3],
			CurrentTranslit: rows[i][4],
		}
		records = append(records, record)
	}

	return records, nil
}

func updateRecord(dbPath, version, newTranslit string) error {
	// Escape single quotes in the text for SQL
	escapedText := strings.ReplaceAll(newTranslit, "'", "''")
	
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
	// Use dolt commands to add, commit, and push changes
	commands := [][]string{
		{"dolt", "add", "."},
		{"dolt", "commit", "-m", "Update transliterations with improved dictionary-based transliterator"},
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}