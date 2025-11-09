package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with frontmatter",
			input:    "---\nname: test\npriority: high\n---\nContent here",
			expected: "Content here",
		},
		{
			name:     "without frontmatter",
			input:    "Just plain content",
			expected: "Just plain content",
		},
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripFrontmatter(tt.input)
			if result != tt.expected {
				t.Errorf("stripFrontmatter() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiple spaces",
			input:    "hello    world",
			expected: "hello world",
		},
		{
			name:     "tabs and newlines",
			input:    "hello\t\nworld",
			expected: "hello world",
		},
		{
			name:     "mixed whitespace",
			input:    "foo   \n\n  bar\t\tbaz",
			expected: "foo bar baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeWhitespace() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestRemoveStopWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with stop words",
			input:    "the quick brown fox",
			expected: "quick brown fox",
		},
		{
			name:     "multiple stop words",
			input:    "this is a test of the system",
			expected: "test system",
		},
		{
			name:     "no stop words",
			input:    "Python Django Flask",
			expected: "Python Django Flask",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeStopWords(tt.input)
			if result != tt.expected {
				t.Errorf("removeStopWords() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestPreprocessText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string // Words that should be in output
		excludes []string // Words that should NOT be in output
	}{
		{
			name:     "full preprocessing",
			input:    "---\nname: test\n---\nThis is a Python expert skill for Django and Flask",
			contains: []string{"Python", "expert", "skill", "Django", "Flask"},
			excludes: []string{"This", "for"},
		},
		{
			name:     "whitespace collapse",
			input:    "Machine    learning\n\n\nwith     neural     networks",
			contains: []string{"Machine", "learning", "neural", "networks"},
			excludes: []string{"with"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preprocessText(tt.input)

			for _, word := range tt.contains {
				if !strings.Contains(result, word) {
					t.Errorf("Expected output to contain %q, got %q", word, result)
				}
			}

			for _, word := range tt.excludes {
				if strings.Contains(result, word) {
					t.Errorf("Expected output to NOT contain %q, got %q", word, result)
				}
			}
		})
	}
}

func TestHasValidFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "valid frontmatter with name",
			content:  "---\nname: test-skill\npriority: high\n---\nContent",
			expected: true,
		},
		{
			name:     "no frontmatter",
			content:  "Just content without frontmatter",
			expected: false,
		},
		{
			name:     "frontmatter without name",
			content:  "---\npriority: high\n---\nContent",
			expected: false,
		},
		{
			name:     "empty name field",
			content:  "---\nname: \npriority: high\n---\nContent",
			expected: false,
		},
		{
			name:     "no closing frontmatter",
			content:  "---\nname: test\npriority: high\nContent",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasValidFrontmatter(tt.content)
			if result != tt.expected {
				t.Errorf("hasValidFrontmatter() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsValidSkillFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		content  string
		expected bool
	}{
		{
			name:     "valid .md with frontmatter",
			path:     "/path/to/skill.md",
			content:  "---\nname: test\n---\nContent",
			expected: true,
		},
		{
			name:     "not .md extension",
			path:     "/path/to/skill.txt",
			content:  "---\nname: test\n---\nContent",
			expected: false,
		},
		{
			name:     ".md without frontmatter",
			path:     "/path/to/skill.md",
			content:  "Just content",
			expected: false,
		},
		{
			name:     ".json file",
			path:     "/path/to/settings.json",
			content:  "{}",
			expected: false,
		},
		{
			name:     "uppercase .MD extension",
			path:     "/path/to/skill.MD",
			content:  "---\nname: test\n---\nContent",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidSkillFile(tt.path, tt.content)
			if result != tt.expected {
				t.Errorf("isValidSkillFile() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestExtractMetadata(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		path             string
		expectedName     string
		expectedPriority string
		expectedType     string
	}{
		{
			name:             "with frontmatter",
			content:          "---\nname: test-skill\npriority: high\ntype: skill\n---\nContent",
			path:             "/path/to/file.md",
			expectedName:     "test-skill",
			expectedPriority: "high",
			expectedType:     "skill",
		},
		{
			name:             "without frontmatter - skill path",
			content:          "Just content",
			path:             "/path/skills/foo.md",
			expectedName:     "/path/skills/foo.md",
			expectedPriority: "medium",
			expectedType:     "skill",
		},
		{
			name:             "without frontmatter - agent path",
			content:          "Just content",
			path:             "/path/agents/bar.md",
			expectedName:     "/path/agents/bar.md",
			expectedPriority: "medium",
			expectedType:     "agent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, priority, itemType := extractMetadata(tt.content, tt.path)
			if name != tt.expectedName {
				t.Errorf("name = %v, expected %v", name, tt.expectedName)
			}
			if priority != tt.expectedPriority {
				t.Errorf("priority = %v, expected %v", priority, tt.expectedPriority)
			}
			if itemType != tt.expectedType {
				t.Errorf("type = %v, expected %v", itemType, tt.expectedType)
			}
		})
	}
}

func TestLoadItems(t *testing.T) {
	// Test loading single file
	t.Run("single file", func(t *testing.T) {
		testFile := filepath.Join("testdata", "skills", "foo.md")
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Skip("testdata not available")
		}

		items, err := loadItems(testFile)
		if err != nil {
			t.Fatalf("loadItems() error = %v", err)
		}

		if len(items) != 1 {
			t.Errorf("expected 1 item, got %d", len(items))
		}

		// Just verify we got an item
		if items[0].Name == "" {
			t.Errorf("expected non-empty name")
		}
	})

	// Test loading directory
	t.Run("directory", func(t *testing.T) {
		testDir := filepath.Join("testdata", "skills")
		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			t.Skip("testdata not available")
		}

		items, err := loadItems(testDir)
		if err != nil {
			t.Fatalf("loadItems() error = %v", err)
		}

		if len(items) < 1 {
			t.Errorf("expected at least 1 item, got %d", len(items))
		}

		// Check that items have names
		for _, item := range items {
			if item.Name == "" {
				t.Errorf("item has empty name")
			}
			if item.Content == "" {
				t.Errorf("item has empty content")
			}
		}
	})
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1.0, 0.0, 0.0},
			b:        []float32{1.0, 0.0, 0.0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1.0, 0.0, 0.0},
			b:        []float32{0.0, 1.0, 0.0},
			expected: 0.0,
		},
		{
			name:     "different length",
			a:        []float32{1.0, 0.0},
			b:        []float32{1.0, 0.0, 0.0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			diff := result - tt.expected
			if diff < -0.001 || diff > 0.001 {
				t.Errorf("cosineSimilarity() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHashContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		same    bool
	}{
		{
			name:    "same content produces same hash",
			content: "test content",
			same:    true,
		},
		{
			name:    "different content produces different hash",
			content: "different content",
			same:    false,
		},
	}

	baseHash := hashContent("test content")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashContent(tt.content)

			if tt.same && result != baseHash {
				t.Errorf("expected same hash, got different")
			}
			if !tt.same && result == baseHash {
				t.Errorf("expected different hash, got same")
			}

			// Hash should be 64 characters (SHA256 hex)
			if len(result) != 64 {
				t.Errorf("expected hash length 64, got %d", len(result))
			}
		})
	}
}

func TestCacheFunctions(t *testing.T) {
	// Test embedding cache
	t.Run("embedding cache", func(t *testing.T) {
		content := "test content for caching"
		embedding := []float32{0.1, 0.2, 0.3, 0.4}

		// Save to cache
		err := saveCachedEmbedding(content, embedding)
		if err != nil {
			t.Fatalf("saveCachedEmbedding() error = %v", err)
		}

		// Load from cache
		loaded, found := loadCachedEmbedding(content)
		if !found {
			t.Errorf("expected to find cached embedding")
		}

		if len(loaded) != len(embedding) {
			t.Errorf("expected length %d, got %d", len(embedding), len(loaded))
		}

		for i := range embedding {
			if loaded[i] != embedding[i] {
				t.Errorf("embedding[%d] = %v, expected %v", i, loaded[i], embedding[i])
			}
		}
	})
}
