-- Enhance books metadata and collaboration schema

-- Books metadata extensions
ALTER TABLE books ADD COLUMN IF NOT EXISTS slug TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS idx_books_slug ON books(slug);
ALTER TABLE books ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'draft';
ALTER TABLE books ADD COLUMN IF NOT EXISTS release_date DATE;
ALTER TABLE books ADD COLUMN IF NOT EXISTS page_count INTEGER;
ALTER TABLE books ADD COLUMN IF NOT EXISTS preview_html TEXT;
ALTER TABLE books ADD COLUMN IF NOT EXISTS ebook_enabled BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE books ADD COLUMN IF NOT EXISTS hardcopy_enabled BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE books ADD COLUMN IF NOT EXISTS created_by_user_id TEXT;
ALTER TABLE books ADD COLUMN IF NOT EXISTS updated_by_user_id TEXT;

-- Author profile enrichment
ALTER TABLE authors ADD COLUMN IF NOT EXISTS genres TEXT[];
ALTER TABLE authors ADD COLUMN IF NOT EXISTS availability TEXT;
ALTER TABLE authors ADD COLUMN IF NOT EXISTS experience_years INTEGER;
ALTER TABLE authors ADD COLUMN IF NOT EXISTS hourly_rate NUMERIC;
ALTER TABLE authors ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';
ALTER TABLE authors ADD COLUMN IF NOT EXISTS profile_image_url TEXT;

-- Publisher profile enrichment
ALTER TABLE publishers ADD COLUMN IF NOT EXISTS focus_genres TEXT[];
ALTER TABLE publishers ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';
ALTER TABLE publishers ADD COLUMN IF NOT EXISTS logo_url TEXT;

-- Book pages structure enhancements
ALTER TABLE book_pages ADD COLUMN IF NOT EXISTS section TEXT NOT NULL DEFAULT 'content';
ALTER TABLE book_pages ADD COLUMN IF NOT EXISTS label TEXT;
ALTER TABLE book_pages ADD COLUMN IF NOT EXISTS audio_url TEXT;
ALTER TABLE book_pages ADD COLUMN IF NOT EXISTS video_url TEXT;

-- Book chapters enrichment
ALTER TABLE book_chapters ADD COLUMN IF NOT EXISTS summary TEXT;
ALTER TABLE book_chapters ADD COLUMN IF NOT EXISTS audio_url TEXT;

-- Collaboration between publishers and authors
CREATE TABLE IF NOT EXISTS publisher_author_collaborations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  publisher_id UUID NOT NULL REFERENCES publishers(id) ON DELETE CASCADE,
  author_id UUID NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'pending',
  notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  responded_at TIMESTAMPTZ,
  UNIQUE(publisher_id, author_id)
);
CREATE INDEX IF NOT EXISTS idx_collaborations_author ON publisher_author_collaborations(author_id);
CREATE INDEX IF NOT EXISTS idx_collaborations_publisher ON publisher_author_collaborations(publisher_id);

-- Digital and physical assets linked to books
CREATE TABLE IF NOT EXISTS book_assets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
  kind TEXT NOT NULL CHECK (kind IN ('ebook','audio','video','supplement','hardcopy')),
  source TEXT NOT NULL CHECK (source IN ('upload','external')),
  url TEXT NOT NULL,
  checksum TEXT,
  file_size_bytes BIGINT,
  created_by_user_id TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_book_assets_book ON book_assets(book_id);
CREATE INDEX IF NOT EXISTS idx_book_assets_kind ON book_assets(kind);

-- Purchases revenue sharing for payouts
CREATE TABLE IF NOT EXISTS purchase_revenue_shares (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  purchase_id UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
  recipient_type TEXT NOT NULL CHECK (recipient_type IN ('author','publisher')),
  recipient_id UUID NOT NULL,
  share_percent NUMERIC NOT NULL,
  amount NUMERIC,
  currency TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_purchase_shares_purchase ON purchase_revenue_shares(purchase_id);

-- Club live sessions for group reading (audio/video)
CREATE TABLE IF NOT EXISTS club_sessions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  club_id UUID NOT NULL REFERENCES book_clubs(id) ON DELETE CASCADE,
  book_id UUID REFERENCES books(id) ON DELETE SET NULL,
  host_user_id TEXT NOT NULL,
  session_type TEXT NOT NULL CHECK (session_type IN ('audio','video')),
  scheduled_at TIMESTAMPTZ NOT NULL,
  duration_minutes INTEGER,
  status TEXT NOT NULL DEFAULT 'scheduled',
  meeting_url TEXT,
  recording_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_club_sessions_club ON club_sessions(club_id, scheduled_at);

CREATE TABLE IF NOT EXISTS club_session_participants (
  session_id UUID NOT NULL REFERENCES club_sessions(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  role TEXT NOT NULL DEFAULT 'participant',
  PRIMARY KEY (session_id, user_id)
);

-- Discussion interactions
CREATE TABLE IF NOT EXISTS room_message_reactions (
  message_id UUID NOT NULL REFERENCES room_messages(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL,
  reaction TEXT NOT NULL CHECK (reaction IN ('like','unlike','upvote','downvote')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY(message_id, user_id, reaction)
);

CREATE TABLE IF NOT EXISTS room_message_reports (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  message_id UUID NOT NULL REFERENCES room_messages(id) ON DELETE CASCADE,
  reporter_user_id TEXT NOT NULL,
  reason TEXT,
  status TEXT NOT NULL DEFAULT 'pending',
  handled_by_user_id TEXT,
  handled_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_room_reports_message ON room_message_reports(message_id);

CREATE TABLE IF NOT EXISTS club_member_mutes (
  club_id UUID NOT NULL REFERENCES book_clubs(id) ON DELETE CASCADE,
  target_user_id TEXT NOT NULL,
  muted_by_user_id TEXT NOT NULL,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (club_id, target_user_id)
);

CREATE TABLE IF NOT EXISTS club_member_blocks (
  club_id UUID NOT NULL REFERENCES book_clubs(id) ON DELETE CASCADE,
  blocker_user_id TEXT NOT NULL,
  blocked_user_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (club_id, blocker_user_id, blocked_user_id)
);
