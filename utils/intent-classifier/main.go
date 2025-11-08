package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/hybridgroup/yzma/pkg/download"
	"github.com/hybridgroup/yzma/pkg/llama"
)

// Item represents a file with its content for matching
type Item struct {
	Name     string
	Path     string
	Content  string
	Priority string
	Type     string // "skill" or "agent"
}

// Match represents a matched item with its similarity score
type Match struct {
	Name       string
	Path       string
	Similarity float32
	Priority   string
	Type       string // "skill" or "agent"
}

// Common English stop words (lightweight list)
var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
	"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
	"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
	"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
	"this": true, "these": true, "those": true, "or": true, "but": true, "can": true,
	"have": true, "do": true, "does": true, "did": true, "doing": true,
}

// preprocessText compacts text for embedding by removing noise and redundancy
func preprocessText(text string) string {
	// 1. Strip YAML frontmatter (already parsed, just noise for embedding)
	text = stripFrontmatter(text)

	// 2. Normalize whitespace - collapse multiple spaces/newlines
	text = normalizeWhitespace(text)

	// 3. Remove stop words to reduce dimensionality
	text = removeStopWords(text)

	// 4. Final whitespace cleanup
	text = strings.TrimSpace(text)

	return text
}

// stripFrontmatter removes YAML frontmatter from text
func stripFrontmatter(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "---") {
		// Find closing ---
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				// Return everything after frontmatter
				return strings.Join(lines[i+1:], "\n")
			}
		}
	}
	return text
}

// normalizeWhitespace collapses multiple spaces, tabs, and newlines
func normalizeWhitespace(text string) string {
	// Replace multiple whitespace with single space
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}

// removeStopWords removes common English stop words
func removeStopWords(text string) string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	var filtered []string
	for _, word := range words {
		lower := strings.ToLower(word)
		if !stopWords[lower] && len(word) > 1 { // Keep words > 1 char
			filtered = append(filtered, word)
		}
	}

	return strings.Join(filtered, " ")
}

func main() {
	// Define flags
	prompt := flag.String("prompt", "", "User prompt to match against (required)")
	embed := flag.String("embed", "", "File or directory path to search and match (required)")
	threshold := flag.Float64("threshold", 0.4, "Similarity threshold (0.0-1.0)")
	embeddingModel := flag.String("embedding-model", "https://huggingface.co/second-state/All-MiniLM-L6-v2-Embedding-GGUF/resolve/main/all-MiniLM-L6-v2-Q5_K_M.gguf", "Embedding model URL or path")
	libPath := flag.String("lib", "", "llama.cpp library path (auto-detect if empty)")
	processor := flag.String("processor", "cpu", "Processor type: cpu, cuda, vulkan, metal")
	outputType := flag.String("output-type", "auto", "Output type: auto, skills, or agents (auto-detects from directory structure)")
	llamaLogLevel := flag.Int("llama-log-level", 0, "Llama.cpp log level (0=disabled, 1=error, 2=warn, 3=info)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Required flags:")
		fmt.Fprintln(os.Stderr, "  -prompt string")
		fmt.Fprintln(os.Stderr, "        User prompt to match against")
		fmt.Fprintln(os.Stderr, "  -embed string")
		fmt.Fprintln(os.Stderr, "        File or directory to embed and match")
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		fmt.Fprintln(os.Stderr, "  -threshold float")
		fmt.Fprintln(os.Stderr, "        Similarity threshold (default: 0.4)")
		fmt.Fprintln(os.Stderr, "  -output-type string")
		fmt.Fprintln(os.Stderr, "        Output type: auto, skills, or agents (default: auto)")
		fmt.Fprintln(os.Stderr, "  -embedding-model string")
		fmt.Fprintln(os.Stderr, "        Embedding model URL or local path")
		fmt.Fprintln(os.Stderr, "        (default: all-MiniLM-L6-v2)")
		fmt.Fprintln(os.Stderr, "  -lib string")
		fmt.Fprintln(os.Stderr, "        llama.cpp library path (default: auto-download)")
		fmt.Fprintln(os.Stderr, "  -processor string")
		fmt.Fprintln(os.Stderr, "        Processor type: cpu, cuda, vulkan, metal (default: cpu)")
	}

	flag.Parse()

	// Validate required flags
	if *prompt == "" || *embed == "" {
		fmt.Fprintln(os.Stderr, "Error: -prompt and -embed are required")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	// Resolve embedding model to GGUF file path
	embeddingModelPath := resolveModel(*embeddingModel, "embedding")

	// Auto-download llama.cpp if not found
	if *libPath == "" {
		*libPath = ensureLlamaLib(*processor)
	}

	// Load llama.cpp library
	if err := llama.Load(*libPath); err != nil {
		if strings.Contains(err.Error(), "libffi") {
			fmt.Fprintf(os.Stderr, "❌ Missing libffi dependency\n")
			fmt.Fprintln(os.Stderr, "\nInstall libffi for your system:")
			fmt.Fprintln(os.Stderr, "  • Ubuntu/Debian: sudo apt install libffi8")
			fmt.Fprintln(os.Stderr, "  • Fedora/RHEL:   sudo dnf install libffi")
			fmt.Fprintln(os.Stderr, "  • Arch Linux:    sudo pacman -S libffi")
			fmt.Fprintln(os.Stderr, "  • macOS:         brew install libffi")
			fmt.Fprintln(os.Stderr, "  • Nix:           nix profile install nixpkgs#libffi")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Failed to load llama.cpp library: %v\n", err)
		fmt.Fprintln(os.Stderr, "Hint: Ensure llama.cpp shared library is available")
		fmt.Fprintln(os.Stderr, "      You can specify it with --lib /path/to/libllama.so")
		os.Exit(1)
	}

	// Initialize llama.cpp
	llama.Init()
	defer llama.BackendFree()

	// Set llama.cpp log level (0 = silent)
	if *llamaLogLevel == 0 {
		llama.LogSet(llama.LogSilent())
	}

	// Load backends from the library path
	llama.GGMLBackendLoadAllFromPath(*libPath)

	// Load embedding model
	model := llama.ModelLoadFromFile(embeddingModelPath, llama.ModelDefaultParams())
	if model == 0 {
		fmt.Fprintf(os.Stderr, "Failed to load embedding model from %s\n", embeddingModelPath)
		os.Exit(1)
	}
	defer llama.ModelFree(model)

	// Create context for embeddings
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 512
	ctxParams.NBatch = 512
	ctxParams.Embeddings = 1 // Enable embeddings mode

	lctx := llama.InitFromModel(model, ctxParams)
	if lctx == 0 {
		fmt.Fprintf(os.Stderr, "Failed to create context from model\n")
		os.Exit(1)
	}
	defer llama.Free(lctx)

	// Load items from file or directory
	items, err := loadItems(*embed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load items: %v\n", err)
		os.Exit(1)
	}

	// Compute prompt embedding (preprocess first)
	processedPrompt := preprocessText(strings.ToLower(*prompt))
	promptEmbed, err := getEmbedding(model, lctx, processedPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to embed prompt: %v\n", err)
		os.Exit(1)
	}

	// Embedding similarity mode - match items
	matches := matchItems(model, lctx, promptEmbed, items, float32(*threshold))

	// Output results
	if len(matches) > 0 {
		outputWithTemplate(matches, *outputType)
	}
}

// matchWithLLMAndTemplate uses LLM to analyze files and generate template output
func getEmbedding(model llama.Model, lctx llama.Context, text string) ([]float32, error) {
	// Tokenize
	vocab := llama.ModelGetVocab(model)
	count := llama.Tokenize(vocab, text, nil, true, true)
	if count <= 0 {
		return nil, fmt.Errorf("tokenization returned no tokens")
	}
	tokens := make([]llama.Token, count)
	llama.Tokenize(vocab, text, tokens, true, true)

	// Encode (use Encode for embedding models like BERT)
	batch := llama.BatchGetOne(tokens)
	if llama.Encode(lctx, batch) != 0 {
		return nil, fmt.Errorf("encode failed")
	}

	// Get embeddings
	nEmbd := llama.ModelNEmbd(model)
	vec := llama.GetEmbeddingsSeq(lctx, 0, nEmbd)

	// Normalize
	var sum float64
	for _, v := range vec {
		sum += float64(v * v)
	}
	sum = math.Sqrt(sum)
	norm := float32(1.0 / sum)

	normalized := make([]float32, len(vec))
	for i, v := range vec {
		normalized[i] = v * norm
	}

	return normalized, nil
}

// loadItems reads file or directory and loads content
func loadItems(path string) ([]Item, error) {
	var items []Item

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Single file
	if !info.IsDir() {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		name, priority, itemType := extractMetadata(string(content), path)
		items = append(items, Item{
			Name:     name,
			Path:     path,
			Content:  string(content),
			Priority: priority,
			Type:     itemType,
		})
		return items, nil
	}

	// Directory - recursively walk and read all files
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", p, err)
			return nil // Continue walking
		}

		// Extract name, priority, and type from frontmatter
		name, priority, itemType := extractMetadata(string(content), p)

		items = append(items, Item{
			Name:     name,
			Path:     p,
			Content:  string(content),
			Priority: priority,
			Type:     itemType,
		})

		return nil
	})

	return items, err
}

// extractMetadata extracts name, priority, and type from frontmatter and path
func extractMetadata(content string, path string) (string, string, string) {
	name := ""
	priority := "medium" // default priority
	itemType := ""       // will be detected from path or frontmatter

	// Try to parse frontmatter
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "---") {
		// Look for name, priority, and type in frontmatter
		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "---" {
				break // End of frontmatter
			}

			// Look for name: value
			if strings.HasPrefix(line, "name:") {
				nameValue := strings.TrimSpace(strings.TrimPrefix(line, "name:"))
				nameValue = strings.Trim(nameValue, "\"'") // Remove quotes if present
				if nameValue != "" {
					name = nameValue
				}
			}

			// Look for priority: value
			if strings.HasPrefix(line, "priority:") {
				priorityValue := strings.TrimSpace(strings.TrimPrefix(line, "priority:"))
				priorityValue = strings.Trim(priorityValue, "\"'") // Remove quotes if present
				if priorityValue != "" {
					priority = priorityValue
				}
			}

			// Look for type: value
			if strings.HasPrefix(line, "type:") {
				typeValue := strings.TrimSpace(strings.TrimPrefix(line, "type:"))
				typeValue = strings.Trim(typeValue, "\"'") // Remove quotes if present
				if typeValue != "" {
					itemType = typeValue
				}
			}
		}
	}

	// If no name found in frontmatter, fallback to absolute path
	if name == "" {
		absPath, err := filepath.Abs(path)
		if err == nil {
			name = absPath
		} else {
			name = path
		}
	}

	// Auto-detect type from path if not specified in frontmatter
	if itemType == "" {
		if strings.Contains(path, "/agents/") || strings.Contains(path, "\\agents\\") {
			itemType = "agent"
		} else if strings.Contains(path, "/skills/") || strings.Contains(path, "\\skills\\") {
			itemType = "skill"
		} else {
			// Default to skill if can't determine
			itemType = "skill"
		}
	}

	return name, priority, itemType
}

// matchItems computes similarity between prompt and item contents
func matchItems(model llama.Model, lctx llama.Context, promptEmbed []float32, items []Item, threshold float32) []Match {
	var matches []Match

	// Get model's max context size and be conservative (leave room for special tokens)
	maxTokens := int(llama.NCtx(lctx)) - 10

	// Get context params for recreation
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 512
	ctxParams.NBatch = 512
	ctxParams.Embeddings = 1

	for _, item := range items {
		// Preprocess text before embedding (removes stop words, whitespace, frontmatter)
		itemText := preprocessText(strings.ToLower(item.Content))

		// Try to load from cache first
		itemEmbed, cached := loadCachedEmbedding(itemText)

		if !cached {
			// Not in cache - generate embedding
			// Conservative char-based truncation (estimate ~3 chars per token to be safe)
			maxChars := maxTokens * 3
			if len(itemText) > maxChars {
				itemText = itemText[:maxChars]
			}

			// Create fresh context for each file to avoid state accumulation
			itemCtx := llama.InitFromModel(model, ctxParams)
			if itemCtx == 0 {
				fmt.Fprintf(os.Stderr, "Warning: failed to create context for %s\n", item.Name)
				continue
			}

			var err error
			itemEmbed, err = getEmbedding(model, itemCtx, itemText)
			llama.Free(itemCtx)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to embed %s: %v\n", item.Name, err)
				continue
			}

			// Save to cache for next time
			if err := saveCachedEmbedding(itemText, itemEmbed); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to cache embedding for %s: %v\n", item.Name, err)
			}
		}

		// Compute cosine similarity
		similarity := cosineSimilarity(promptEmbed, itemEmbed)

		if similarity >= threshold {
			matches = append(matches, Match{
				Name:       item.Name,
				Path:       item.Path,
				Similarity: similarity,
				Priority:   item.Priority,
				Type:       item.Type,
			})
		}
	}

	return matches
}

// cosineSimilarity computes cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float32
	for i := range a {
		dotProduct += a[i] * b[i]
	}

	// Vectors are already normalized, so dot product = cosine similarity
	return dotProduct
}

// hashContent returns SHA256 hash of content
func hashContent(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])
}

// getCacheFile returns path to cache file for given hash
func getCacheFile(hash, cacheType string) string {
	cacheDir := filepath.Join(getCacheDir(), cacheType)
	os.MkdirAll(cacheDir, 0755)
	return filepath.Join(cacheDir, hash+".cache")
}

// loadCachedEmbedding loads embedding from cache if exists
func loadCachedEmbedding(content string) ([]float32, bool) {
	hash := hashContent(content)
	cacheFile := getCacheFile(hash, "embeddings")

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, false
	}

	// Decode float32 array
	if len(data)%4 != 0 {
		return nil, false
	}

	embedding := make([]float32, len(data)/4)
	for i := range embedding {
		bits := uint32(data[i*4]) | uint32(data[i*4+1])<<8 | uint32(data[i*4+2])<<16 | uint32(data[i*4+3])<<24
		embedding[i] = math.Float32frombits(bits)
	}

	return embedding, true
}

// saveCachedEmbedding saves embedding to cache
func saveCachedEmbedding(content string, embedding []float32) error {
	hash := hashContent(content)
	cacheFile := getCacheFile(hash, "embeddings")

	// Encode float32 array
	data := make([]byte, len(embedding)*4)
	for i, v := range embedding {
		bits := math.Float32bits(v)
		data[i*4] = byte(bits)
		data[i*4+1] = byte(bits >> 8)
		data[i*4+2] = byte(bits >> 16)
		data[i*4+3] = byte(bits >> 24)
	}

	return os.WriteFile(cacheFile, data, 0644)
}

// outputWithTemplate renders matches grouped by type and priority
func outputWithTemplate(matches []Match, outputType string) {
	// Separate matches by type and priority
	skillsByPriority := map[string][]string{
		"critical": {},
		"high":     {},
		"medium":   {},
		"low":      {},
	}
	agentsByPriority := map[string][]string{
		"critical": {},
		"high":     {},
		"medium":   {},
		"low":      {},
	}

	for _, match := range matches {
		priority := strings.ToLower(match.Priority)
		if priority != "critical" && priority != "high" && priority != "medium" && priority != "low" {
			priority = "medium" // default
		}

		if match.Type == "agent" {
			agentsByPriority[priority] = append(agentsByPriority[priority], match.Name)
		} else {
			skillsByPriority[priority] = append(skillsByPriority[priority], match.Name)
		}
	}

	// Count totals
	hasSkills := len(skillsByPriority["critical"])+len(skillsByPriority["high"])+len(skillsByPriority["medium"])+len(skillsByPriority["low"]) > 0
	hasAgents := len(agentsByPriority["critical"])+len(agentsByPriority["high"])+len(agentsByPriority["medium"])+len(agentsByPriority["low"]) > 0

	var output strings.Builder
	output.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if hasSkills && hasAgents {
		output.WriteString("🎯 SKILLS & AGENTS ACTIVATION CHECK\n")
	} else if hasAgents {
		output.WriteString("🤖 AGENTS ACTIVATION CHECK\n")
	} else {
		output.WriteString("🎯 SKILLS ACTIVATION CHECK\n")
	}

	output.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	output.WriteString("\n")

	// Output skills section
	if hasSkills {
		outputSection(&output, "SKILLS", "", skillsByPriority)
	}

	// Output agents section
	if hasAgents {
		if hasSkills {
			output.WriteString("\n") // Extra spacing between sections
		}
		outputSection(&output, "AGENTS", "@", agentsByPriority)
	}

	// Build action text
	var actionParts []string
	if hasSkills {
		actionParts = append(actionParts, "Use Skill tool")
	}
	if hasAgents {
		var agentList []string
		for _, priority := range []string{"critical", "high", "medium", "low"} {
			for _, agent := range agentsByPriority[priority] {
				agentList = append(agentList, "@"+agent)
			}
		}
		if len(agentList) > 0 {
			actionParts = append(actionParts, "Use "+strings.Join(agentList, ", "))
		}
	}

	if len(actionParts) > 0 {
		output.WriteString("ACTION: " + strings.Join(actionParts, " and ") + "\n")
	}

	output.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	fmt.Print(output.String())
}

// outputSection outputs a single section (skills or agents) grouped by priority
func outputSection(output *strings.Builder, label string, prefix string, itemsByPriority map[string][]string) {
	if len(itemsByPriority["critical"]) > 0 {
		output.WriteString("⚠️  CRITICAL " + label + " (REQUIRED):\n")
		for _, item := range itemsByPriority["critical"] {
			output.WriteString("  → " + prefix + item + "\n")
		}
		output.WriteString("\n")
	}

	if len(itemsByPriority["high"]) > 0 {
		output.WriteString("📚 RECOMMENDED " + label + ":\n")
		for _, item := range itemsByPriority["high"] {
			output.WriteString("  → " + prefix + item + "\n")
		}
		output.WriteString("\n")
	}

	if len(itemsByPriority["medium"]) > 0 {
		output.WriteString("💡 SUGGESTED " + label + ":\n")
		for _, item := range itemsByPriority["medium"] {
			output.WriteString("  → " + prefix + item + "\n")
		}
		output.WriteString("\n")
	}

	if len(itemsByPriority["low"]) > 0 {
		output.WriteString("📌 OPTIONAL " + label + ":\n")
		for _, item := range itemsByPriority["low"] {
			output.WriteString("  → " + prefix + item + "\n")
		}
		output.WriteString("\n")
	}
}

// resolveModel resolves a model URL or path to a local GGUF file path
func resolveModel(modelSpec string, modelType string) string {
	// If it's already a local file path, return it
	if _, err := os.Stat(modelSpec); err == nil {
		return modelSpec
	}

	// Must be a URL - download it
	if !strings.HasPrefix(modelSpec, "http://") && !strings.HasPrefix(modelSpec, "https://") {
		fmt.Fprintf(os.Stderr, "❌ Model must be either a local path or a URL: %s\n", modelSpec)
		os.Exit(1)
	}

	// Extract filename from URL
	parts := strings.Split(modelSpec, "/")
	filename := parts[len(parts)-1]

	// Remove query parameters if present
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}

	// Build cache path
	cacheDir := filepath.Join(getCacheDir(), "models", modelType)
	os.MkdirAll(cacheDir, 0755)
	modelPath := filepath.Join(cacheDir, filename)

	// If model already cached, return path
	if _, err := os.Stat(modelPath); err == nil {
		return modelPath
	}

	// Download model using HTTP client
	fmt.Printf("📥 Downloading %s model...\n", modelType)
	fmt.Printf("   From: %s\n", modelSpec)

	if err := downloadFile(modelSpec, modelPath); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to download model: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Model downloaded successfully")
	return modelPath
}

// downloadFile downloads a file from a URL to a local path
func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// ensureLlamaLib ensures llama.cpp library is available
func ensureLlamaLib(processor string) string {
	// 1. Check if already exists in current directory
	libName := download.LibraryName(runtime.GOOS)
	if _, err := os.Stat(libName); err == nil {
		return "."
	}

	// 2. Check cache directory
	cacheDir := getCacheDir()
	os.MkdirAll(cacheDir, 0755)

	libPath := filepath.Join(cacheDir, libName)
	if _, err := os.Stat(libPath); err == nil {
		return cacheDir
	}

	// 3. Download llama.cpp
	fmt.Println("📥 Downloading llama.cpp library (first time setup)...")

	version, err := download.LlamaLatestVersion()
	if err != nil {
		fmt.Println("⚠️  Could not get latest version, using default...")
		version = "b6795"
	}

	fmt.Printf("📦 Installing llama.cpp version %s (%s)...\n", version, processor)
	if err := download.Get(runtime.GOOS, processor, version, cacheDir); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to download llama.cpp: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ llama.cpp library installed successfully")
	return cacheDir
}

// getCacheDir returns a cross-platform cache directory
func getCacheDir() string {
	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			return filepath.Join(localAppData, "intent-classifier")
		}
		return filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "intent-classifier")
	}

	// Unix-like systems (Linux, macOS, BSD)
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
		return filepath.Join(xdgCache, "intent-classifier")
	}

	return filepath.Join(os.Getenv("HOME"), ".cache", "intent-classifier")
}
