package demo

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/srmdn/foliocms/internal/config"
	"github.com/srmdn/foliocms/internal/db"
	"github.com/srmdn/foliocms/internal/storage"
)

// Apply resets the demo instance to a clean state: deletes all posts and
// media, then re-creates the demo user and seed content.
// Called on startup when DEMO_MODE=true and by the reset endpoint.
func Apply(database *db.DB, cfg *config.Config) error {
	store := storage.New(cfg.ContentDir)

	// Delete all post directories from disk.
	rows, err := database.Query(`SELECT slug FROM posts`)
	if err != nil {
		return fmt.Errorf("querying posts: %w", err)
	}
	var slugs []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			rows.Close()
			return fmt.Errorf("scanning slug: %w", err)
		}
		slugs = append(slugs, slug)
	}
	rows.Close()

	for _, slug := range slugs {
		_ = store.Delete(slug)
	}

	if _, err := database.Exec(`DELETE FROM posts`); err != nil {
		return fmt.Errorf("clearing posts: %w", err)
	}

	// Clear media files and DB records.
	if err := os.RemoveAll(cfg.MediaDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("clearing media dir: %w", err)
	}
	if err := os.MkdirAll(cfg.MediaDir, 0755); err != nil {
		return fmt.Errorf("recreating media dir: %w", err)
	}
	if _, err := database.Exec(`DELETE FROM media`); err != nil {
		return fmt.Errorf("clearing media records: %w", err)
	}

	// Upsert demo user with known credentials.
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.DemoPasswd), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing demo password: %w", err)
	}
	if _, err := database.Exec(
		`INSERT INTO users (email, passwd_hash) VALUES (?, ?)
		 ON CONFLICT(email) DO UPDATE SET passwd_hash = excluded.passwd_hash`,
		cfg.DemoEmail, string(hash),
	); err != nil {
		return fmt.Errorf("upserting demo user: %w", err)
	}

	if err := seedPosts(database, store); err != nil {
		return fmt.Errorf("seeding posts: %w", err)
	}

	if err := seedSettings(database); err != nil {
		return fmt.Errorf("seeding settings: %w", err)
	}

	return nil
}

func seedPosts(database *db.DB, store *storage.Storage) error {
	type seedPost struct {
		slug        string
		title       string
		description string
		tags        string
		draft       int
		publishDate string
		body        string
	}

	posts := []seedPost{
		{
			slug:        "getting-started-with-foliocms",
			title:       "Getting Started with FolioCMS",
			description: "A tour of FolioCMS: write posts, manage media, and configure your site.",
			tags:        "tutorial,cms",
			draft:       0,
			publishDate: "2026-04-01",
			body: `FolioCMS is a self-hosted CMS for developers who want a clean writing experience without managing a complex stack.

## What's included

- **Post editor**: A WYSIWYG Markdown editor with full toolbar support: headings, lists, code blocks, tables, and links.
- **Media library**: Upload and manage images directly from the dashboard. Supports local storage and S3-compatible providers.
- **Site settings**: Manage your site name, description, and social links from the dashboard.

## Getting started

1. Download the latest release from the [releases page](https://github.com/srmdn/foliocms/releases).
2. Copy ` + "`" + `.env.example` + "`" + ` to ` + "`" + `.env` + "`" + ` and fill in your configuration.
3. Run ` + "`" + `./folio --setup` + "`" + ` to create your admin account.
4. Start the server: ` + "`" + `./folio` + "`" + `

Your admin dashboard will be available at ` + "`" + `http://localhost:8090/admin` + "`" + `.

## Themes

FolioCMS ships with an Astro SSR theme. After writing and publishing posts, click **Rebuild Site** in Settings to regenerate the frontend.
`,
		},
		{
			slug:        "try-editing-this-post",
			title:       "Try Editing This Post",
			description: "A draft post for you to experiment with the FolioCMS editor.",
			tags:        "demo",
			draft:       1,
			publishDate: "2026-04-06",
			body: `This is a draft post for you to experiment with.

**Here's what to try:**

- Change the title at the top
- Edit this content in the WYSIWYG editor
- Upload a hero image
- Press **Cmd+S** (or **Ctrl+S**) to save

When you're done exploring, click **Reset Demo** in the sidebar to restore everything to its original state.
`,
		},
	}

	for _, p := range posts {
		if err := store.Write(p.slug, &storage.PostFile{
			Frontmatter: storage.Frontmatter{
				Title:       p.title,
				Description: p.description,
				PublishDate: p.publishDate,
				Draft:       p.draft == 1,
				Tags:        splitTags(p.tags),
			},
			Body: p.body,
		}); err != nil {
			return fmt.Errorf("writing post %q: %w", p.slug, err)
		}

		if _, err := database.Exec(
			`INSERT INTO posts (slug, title, description, tags, draft, publish_date)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			p.slug, p.title, p.description, p.tags, p.draft, p.publishDate,
		); err != nil {
			return fmt.Errorf("inserting post %q: %w", p.slug, err)
		}
	}

	return nil
}

func seedSettings(database *db.DB) error {
	settings := map[string]string{
		"site_name":        "Demo Site",
		"site_description": "A demo instance of FolioCMS. Explore, edit, and reset.",
		"social_github":    "https://github.com/srmdn/foliocms",
	}

	for key, value := range settings {
		if _, err := database.Exec(
			`UPDATE settings SET value = ? WHERE key = ?`, value, key,
		); err != nil {
			return fmt.Errorf("updating setting %q: %w", key, err)
		}
	}

	return nil
}

func splitTags(tags string) []string {
	if tags == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i <= len(tags); i++ {
		if i == len(tags) || tags[i] == ',' {
			t := tags[start:i]
			if t != "" {
				result = append(result, t)
			}
			start = i + 1
		}
	}
	return result
}
