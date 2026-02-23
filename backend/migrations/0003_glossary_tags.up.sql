CREATE TABLE IF NOT EXISTS glossary_tags (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS glossary_term_tags (
  term_id TEXT NOT NULL REFERENCES glossary_terms(id) ON DELETE CASCADE,
  tag_id  TEXT NOT NULL REFERENCES glossary_tags(id) ON DELETE CASCADE,
  PRIMARY KEY (term_id, tag_id)
);
