-- Book pages for on-platform reading/writing
CREATE TABLE IF NOT EXISTS book_pages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
  page_no INTEGER NOT NULL CHECK (page_no > 0),
  content_html TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(book_id, page_no)
);

