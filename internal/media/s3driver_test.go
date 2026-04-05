package media

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockS3Server creates a minimal S3-compatible HTTP server for testing.
// It records the last PUT key and body, and accepts DELETE requests.
type mockS3Server struct {
	lastPutKey  string
	lastPutBody []byte
	lastDelKey  string
}

func (m *mockS3Server) handler() http.Handler {
	mux := http.NewServeMux()

	// PUT /bucket/key
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// path is /<bucket>/<key>
		parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
		if len(parts) < 2 {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		key := parts[1]

		switch r.Method {
		case http.MethodPut:
			m.lastPutKey = key
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "read error", http.StatusInternalServerError)
				return
			}
			m.lastPutBody = body
			w.WriteHeader(http.StatusOK)

		case http.MethodDelete:
			m.lastDelKey = key
			w.WriteHeader(http.StatusNoContent)

		default:
			// Return minimal ListBucketResult for unexpected GETs (SDK probing)
			type ListResult struct {
				XMLName     xml.Name `xml:"ListBucketResult"`
				Name        string   `xml:"Name"`
				IsTruncated bool     `xml:"IsTruncated"`
			}
			w.Header().Set("Content-Type", "application/xml")
			_ = xml.NewEncoder(w).Encode(ListResult{Name: parts[0]})
		}
	})

	return mux
}

func newTestS3Driver(t *testing.T, mock *mockS3Server) (*S3Driver, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(mock.handler())
	bucket := "test-bucket"
	d := NewS3Driver(
		srv.URL,
		bucket,
		"auto",
		"test-access-key",
		"test-secret-key",
		srv.URL+"/"+bucket,
	)
	return d, srv
}

func TestS3Driver_PublicURL(t *testing.T) {
	d := &S3Driver{publicURL: "https://s3.nevaobjects.id/my-bucket"}
	got := d.PublicURL("abc-123-photo.webp")
	want := "https://s3.nevaobjects.id/my-bucket/abc-123-photo.webp"
	if got != want {
		t.Errorf("PublicURL = %q, want %q", got, want)
	}
}

func TestS3Driver_PublicURL_TrailingSlash(t *testing.T) {
	d := NewS3Driver("https://s3.nevaobjects.id", "my-bucket", "auto", "k", "s", "https://s3.nevaobjects.id/my-bucket/")
	got := d.PublicURL("file.jpg")
	if strings.Contains(got, "//file.jpg") {
		t.Errorf("PublicURL has double slash: %q", got)
	}
}

func TestS3Driver_Upload(t *testing.T) {
	mock := &mockS3Server{}
	d, srv := newTestS3Driver(t, mock)
	defer srv.Close()

	content := []byte("hello world")
	key, size, err := d.Upload("photo.jpg", bytes.NewReader(content), "image/jpeg")
	if err != nil {
		t.Fatalf("Upload error: %v", err)
	}
	if size != int64(len(content)) {
		t.Errorf("size = %d, want %d", size, len(content))
	}
	if !strings.HasSuffix(key, "-photo.jpg") {
		t.Errorf("key %q should end with -photo.jpg", key)
	}
	if mock.lastPutKey != key {
		t.Errorf("server received key %q, want %q", mock.lastPutKey, key)
	}
	if !bytes.Equal(mock.lastPutBody, content) {
		t.Errorf("server body mismatch")
	}
}

func TestS3Driver_Delete(t *testing.T) {
	mock := &mockS3Server{}
	d, srv := newTestS3Driver(t, mock)
	defer srv.Close()

	// Upload first so we have a real key
	key, _, err := d.Upload("img.png", bytes.NewReader([]byte("data")), "image/png")
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}

	if err := d.Delete(key); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if mock.lastDelKey != key {
		t.Errorf("server received delete key %q, want %q", mock.lastDelKey, key)
	}
}

func TestS3Driver_Delete_InvalidKey(t *testing.T) {
	mock := &mockS3Server{}
	d, srv := newTestS3Driver(t, mock)
	defer srv.Close()

	cases := []string{"../../etc/passwd", "a/b"}
	for _, k := range cases {
		if err := d.Delete(k); err == nil {
			t.Errorf("Delete(%q) expected error, got nil", k)
		}
	}
}
