package media

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MediaFile represents an uploaded file as stored in the database.
type MediaFile struct {
	Key         string    `json:"key"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

// MediaDriver is the interface for file storage backends.
type MediaDriver interface {
	// Upload saves r under a generated key derived from filename.
	// Returns the key, bytes written, and any error.
	Upload(filename string, r io.Reader, contentType string) (key string, size int64, err error)
	// Delete removes the file identified by key.
	Delete(key string) error
	// PublicURL returns the absolute URL at which key is accessible.
	PublicURL(key string) string
}

// LocalDriver stores files in a local directory and serves them via the backend.
type LocalDriver struct {
	mediaDir string
	siteURL  string
}

func NewLocalDriver(mediaDir, siteURL string) *LocalDriver {
	return &LocalDriver{mediaDir: mediaDir, siteURL: siteURL}
}

var unsafeChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = unsafeChars.ReplaceAllString(name, "_")
	if name == "" || name == "." {
		name = "file"
	}
	return strings.ToLower(name)
}

func (d *LocalDriver) Upload(filename string, r io.Reader, _ string) (string, int64, error) {
	if err := os.MkdirAll(d.mediaDir, 0755); err != nil {
		return "", 0, fmt.Errorf("creating media dir: %w", err)
	}

	safe := sanitizeFilename(filename)
	key := uuid.New().String() + "-" + safe
	dest := filepath.Join(d.mediaDir, key)

	f, err := os.Create(dest)
	if err != nil {
		return "", 0, fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		os.Remove(dest)
		return "", 0, fmt.Errorf("writing file: %w", err)
	}

	return key, n, nil
}

func (d *LocalDriver) Delete(key string) error {
	if strings.Contains(key, "/") || strings.Contains(key, "..") {
		return fmt.Errorf("invalid key")
	}
	path := filepath.Join(d.mediaDir, key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing file: %w", err)
	}
	return nil
}

func (d *LocalDriver) PublicURL(key string) string {
	return strings.TrimRight(d.siteURL, "/") + "/media/" + key
}
