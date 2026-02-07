package validate

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"prompt-ci/internal/suite"
)

var citationRegex = regexp.MustCompile(`\[doc:([a-z0-9_\-]+)#(c[1-9][0-9]*)\]`)

// ValidateGroundingCitations validates that content has proper citations
func ValidateGroundingCitations(content string, s *suite.Suite) (bool, []string) {
	var failures []string

	// Build doc index
	docIndex := suite.BuildDocIndex(s)

	// Extract all citations
	citations := extractCitations(content)

	// Validate each citation references a valid doc/chunk
	for _, cit := range citations {
		if chunks, exists := docIndex[cit.DocID]; !exists {
			failures = append(failures, fmt.Sprintf("citation [doc:%s#%s] references non-existent doc '%s'", cit.DocID, cit.ChunkID, cit.DocID))
		} else if !chunks[cit.ChunkID] {
			failures = append(failures, fmt.Sprintf("citation [doc:%s#%s] references non-existent chunk '%s'", cit.DocID, cit.ChunkID, cit.ChunkID))
		}
	}

	// Check citation requirement for long sentences
	sentences := splitIntoSentences(content)
	for _, sentence := range sentences {
		// Skip if sentence is only citations or headings
		if isOnlyCitations(sentence) || isHeading(sentence) {
			continue
		}

		// Count non-whitespace tokens
		tokens := countTokens(sentence)
		if tokens >= 6 {
			// Check if sentence contains at least one citation
			if !containsCitation(sentence) {
				failures = append(failures, fmt.Sprintf("sentence with %d tokens lacks citation: '%s'", tokens, truncate(sentence, 50)))
			}
		}
	}

	if len(failures) > 0 {
		return false, failures
	}

	return true, nil
}

// Citation represents an extracted citation
type Citation struct {
	DocID   string
	ChunkID string
	Full    string
}

// extractCitations extracts all citations from content
func extractCitations(content string) []Citation {
	var citations []Citation
	matches := citationRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		citations = append(citations, Citation{
			DocID:   match[1],
			ChunkID: match[2],
			Full:    match[0],
		})
	}
	return citations
}

// splitIntoSentences splits content into sentences
func splitIntoSentences(content string) []string {
	// Split by sentence-ending punctuation followed by whitespace
	re := regexp.MustCompile(`[.!?]\s+`)
	parts := re.Split(content, -1)

	var sentences []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			sentences = append(sentences, part)
		}
	}
	return sentences
}

// countTokens counts non-whitespace tokens in a string
func countTokens(s string) int {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	return len(fields)
}

// containsCitation checks if a string contains a citation
func containsCitation(s string) bool {
	return citationRegex.MatchString(s)
}

// isOnlyCitations checks if a string is only citations
func isOnlyCitations(s string) bool {
	// Remove all citations and see if anything meaningful remains
	stripped := citationRegex.ReplaceAllString(s, "")
	stripped = strings.TrimSpace(stripped)
	// Allow punctuation and whitespace only
	for _, r := range stripped {
		if !unicode.IsPunct(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// isHeading checks if a string looks like a heading
func isHeading(s string) bool {
	s = strings.TrimSpace(s)
	// Headings typically start with # or are very short all-caps
	if strings.HasPrefix(s, "#") {
		return true
	}
	if len(s) < 50 && strings.ToUpper(s) == s && !strings.Contains(s, ".") {
		return true
	}
	return false
}

// truncate truncates a string to maxLen characters
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
