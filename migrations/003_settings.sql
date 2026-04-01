CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL DEFAULT ''
);

INSERT OR IGNORE INTO settings (key, value) VALUES
    ('site_name',        'My FolioCMS Site'),
    ('site_description', ''),
    ('social_github',    ''),
    ('social_twitter',   ''),
    ('social_linkedin',  '');
