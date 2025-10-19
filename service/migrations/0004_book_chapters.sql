-- Chapters/sections structure for books
CREATE TABLE IF NOT EXISTS book_chapters (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
  number INTEGER NOT NULL CHECK (number > 0),
  title TEXT NOT NULL,
  page_no_start INTEGER CHECK (page_no_start > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(book_id, number)
);

