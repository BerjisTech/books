-- Book reviews and ratings
CREATE TABLE IF NOT EXISTS book_reviews (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL,
  rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
  title TEXT,
  body TEXT,
  status TEXT NOT NULL DEFAULT 'published',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (book_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_book_reviews_book ON book_reviews(book_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_book_reviews_user ON book_reviews(user_id);

-- Materialized average rating view
CREATE OR REPLACE FUNCTION update_book_rating()
RETURNS trigger AS $$
BEGIN
  UPDATE books SET
    rating = (SELECT ROUND(AVG(rating)::numeric, 1) FROM book_reviews WHERE book_id = COALESCE(NEW.book_id, OLD.book_id) AND status = 'published'),
    review_count = (SELECT COUNT(*) FROM book_reviews WHERE book_id = COALESCE(NEW.book_id, OLD.book_id) AND status = 'published')
  WHERE id = COALESCE(NEW.book_id, OLD.book_id);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add rating columns to books if not present
ALTER TABLE books ADD COLUMN IF NOT EXISTS rating NUMERIC(2,1);
ALTER TABLE books ADD COLUMN IF NOT EXISTS review_count INTEGER NOT NULL DEFAULT 0;

DROP TRIGGER IF EXISTS trg_update_book_rating ON book_reviews;
CREATE TRIGGER trg_update_book_rating
AFTER INSERT OR UPDATE OR DELETE ON book_reviews
FOR EACH ROW EXECUTE FUNCTION update_book_rating();
