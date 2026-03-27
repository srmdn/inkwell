package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Frontmatter is the YAML header written to each post's index.md.
type Frontmatter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	PublishDate string   `yaml:"publishDate"`
	Draft       bool     `yaml:"draft"`
	Tags        []string `yaml:"tags"`
}

// PostFile represents a post file on disk.
type PostFile struct {
	Frontmatter Frontmatter
	Body        string // markdown body (after frontmatter)
}

// Storage handles reading and writing post files under a content directory.
type Storage struct {
	contentDir string
}

func New(contentDir string) *Storage {
	return &Storage{contentDir: contentDir}
}

// PostDir returns the directory for a given slug.
func (s *Storage) PostDir(slug string) string {
	return filepath.Join(s.contentDir, slug)
}

// PostPath returns the full path to a post's index.md.
func (s *Storage) PostPath(slug string) string {
	return filepath.Join(s.contentDir, slug, "index.md")
}

// Read reads and parses a post file from disk.
func (s *Storage) Read(slug string) (*PostFile, error) {
	raw, err := os.ReadFile(s.PostPath(slug))
	if err != nil {
		return nil, fmt.Errorf("reading post %q: %w", slug, err)
	}
	return parse(raw)
}

// Write creates or overwrites a post file on disk.
func (s *Storage) Write(slug string, pf *PostFile) error {
	if err := os.MkdirAll(s.PostDir(slug), 0755); err != nil {
		return fmt.Errorf("creating post dir %q: %w", slug, err)
	}
	content, err := marshal(pf)
	if err != nil {
		return err
	}
	return os.WriteFile(s.PostPath(slug), content, 0644)
}

// Delete removes the post directory and all its contents.
func (s *Storage) Delete(slug string) error {
	if err := os.RemoveAll(s.PostDir(slug)); err != nil {
		return fmt.Errorf("deleting post %q: %w", slug, err)
	}
	return nil
}

// Exists reports whether a post file exists on disk.
func (s *Storage) Exists(slug string) bool {
	_, err := os.Stat(s.PostPath(slug))
	return err == nil
}

func parse(raw []byte) (*PostFile, error) {
	const delim = "---"

	// Trim only leading whitespace so trailing newlines in the body are preserved.
	content := strings.TrimLeft(string(raw), " \t\r\n")

	if !strings.HasPrefix(content, delim) {
		return nil, fmt.Errorf("missing frontmatter delimiter")
	}

	// Skip past the opening --- and its trailing newline.
	rest := content[len(delim):]
	if len(rest) == 0 || rest[0] != '\n' {
		return nil, fmt.Errorf("missing frontmatter delimiter")
	}
	rest = rest[1:]

	// Find the closing ---.
	idx := strings.Index(rest, "\n"+delim)
	if idx == -1 {
		return nil, fmt.Errorf("unclosed frontmatter")
	}

	yamlPart := rest[:idx]

	// Skip \n--- and the blank line marshal writes before the body.
	// marshal writes "---\n\n" so two newlines follow the closing delimiter.
	body := rest[idx+1+len(delim):]
	body = strings.TrimPrefix(body, "\n\n")
	// Fallback for files with only one newline after ---.
	body = strings.TrimPrefix(body, "\n")

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(yamlPart), &fm); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	return &PostFile{Frontmatter: fm, Body: body}, nil
}

func marshal(pf *PostFile) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(pf.Frontmatter); err != nil {
		return nil, fmt.Errorf("encoding frontmatter: %w", err)
	}
	enc.Close()
	buf.WriteString("---\n\n")
	buf.WriteString(pf.Body)
	return buf.Bytes(), nil
}

// FormatPublishDate returns today's date as a publishDate string.
func FormatPublishDate(t time.Time) string {
	return t.Format("2006-01-02")
}
