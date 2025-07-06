package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	
	"github.com/LaPingvino/bahai-transliterator"
)

func main() {
	var (
		language = flag.String("lang", "auto", "Language: arabic, persian, or auto")
		file     = flag.String("file", "", "Input file (if not provided, reads from stdin)")
		verbose  = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	trans := transliterator.New()

	var input string
	if *file != "" {
		// Read from file
		content, err := os.ReadFile(*file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		input = string(content)
	} else {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		input = strings.Join(lines, "\n")
	}

	// Determine language
	var lang transliterator.Language
	switch strings.ToLower(*language) {
	case "arabic", "ar":
		lang = transliterator.Arabic
	case "persian", "fa", "farsi":
		lang = transliterator.Persian
	case "auto":
		lang = transliterator.AutoDetectLanguage(input)
		if *verbose {
			langName := "Arabic"
			if lang == transliterator.Persian {
				langName = "Persian"
			}
			fmt.Fprintf(os.Stderr, "Detected language: %s\n", langName)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid language: %s\n", *language)
		os.Exit(1)
	}

	// Transliterate
	result := trans.Transliterate(input, lang)
	
	if *verbose {
		fmt.Fprintf(os.Stderr, "Input length: %d characters\n", len(input))
		fmt.Fprintf(os.Stderr, "Output length: %d characters\n", len(result))
		fmt.Fprintf(os.Stderr, "---\n")
	}

	fmt.Print(result)
}