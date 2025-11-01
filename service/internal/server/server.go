package server

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/berjistech/berjis-ecosystem/books/service/internal/auth"
	"github.com/berjistech/berjis-ecosystem/books/service/internal/billing"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type Options struct {
	AllowedOrigins string
	DB             *sqlx.DB
	Env            string
	AuthHS256      string
	CoreAPIBase    string
}

type Book struct {
	ID              string     `db:"id" json:"id"`
	AuthorID        *string    `db:"author_id" json:"authorId,omitempty"`
	PublisherID     *string    `db:"publisher_id" json:"publisherId,omitempty"`
	Title           string     `db:"title" json:"title"`
	Subtitle        *string    `db:"subtitle" json:"subtitle,omitempty"`
	Description     *string    `db:"description" json:"description,omitempty"`
	ISBN            *string    `db:"isbn" json:"isbn,omitempty"`
	Language        *string    `db:"language" json:"language,omitempty"`
	Categories      []string   `db:"categories" json:"categories,omitempty"`
	CoverURL        *string    `db:"cover_image_url" json:"coverImageUrl,omitempty"`
	EbookURL        *string    `db:"ebook_file_url" json:"ebookFileUrl,omitempty"`
	PriceAmount     *float64   `db:"price_amount" json:"priceAmount,omitempty"`
	PriceCurrency   *string    `db:"price_currency" json:"priceCurrency,omitempty"`
	Shareable       bool       `db:"shareable" json:"shareable"`
	AllowHardcopy   bool       `db:"allow_hardcopy" json:"allowHardcopy"`
	Visibility      string     `db:"visibility" json:"visibility"`
	AuthorName      *string    `db:"author_name" json:"authorName,omitempty"`
	PublisherName   *string    `db:"publisher_name" json:"publisherName,omitempty"`
	Slug            *string    `db:"slug" json:"slug,omitempty"`
	Status          string     `db:"status" json:"status"`
	ReleaseDate     *time.Time `db:"release_date" json:"releaseDate,omitempty"`
	PageCount       *int       `db:"page_count" json:"pageCount,omitempty"`
	PreviewHTML     *string    `db:"preview_html" json:"previewHtml,omitempty"`
	EbookEnabled    bool       `db:"ebook_enabled" json:"ebookEnabled"`
	HardcopyEnabled bool       `db:"hardcopy_enabled" json:"hardcopyEnabled"`
}

type Purchase struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"userId"`
	BookID    string    `db:"book_id" json:"bookId"`
	Kind      string    `db:"kind" json:"kind"`
	Status    string    `db:"status" json:"status"`
	Paid      *float64  `db:"price_paid" json:"pricePaid,omitempty"`
	Currency  *string   `db:"currency" json:"currency,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type ShareLink struct {
	ID        string     `db:"id" json:"id"`
	BookID    string     `db:"book_id" json:"bookId"`
	Token     string     `db:"token" json:"token"`
	MaxUses   *int       `db:"max_uses" json:"maxUses,omitempty"`
	ExpiresAt *time.Time `db:"expires_at" json:"expiresAt,omitempty"`
}

type Club struct {
	ID         string  `db:"id" json:"id"`
	OwnerID    string  `db:"owner_user_id" json:"ownerUserId"`
	Name       string  `db:"name" json:"name"`
	Desc       *string `db:"description" json:"description,omitempty"`
	Visibility string  `db:"visibility" json:"visibility"`
}

type Room struct {
	ID     string `db:"id" json:"id"`
	ClubID string `db:"club_id" json:"clubId"`
	Name   string `db:"name" json:"name"`
}

type Message struct {
	ID      string    `db:"id" json:"id"`
	RoomID  string    `db:"room_id" json:"roomId"`
	UserID  string    `db:"user_id" json:"userId"`
	Content string    `db:"content" json:"content"`
	Created time.Time `db:"created_at" json:"createdAt"`
}

type Meetup struct {
	ID         string     `db:"id" json:"id"`
	Organizer  string     `db:"organizer_user_id" json:"organizerUserId"`
	ClubID     *string    `db:"club_id" json:"clubId,omitempty"`
	BookID     *string    `db:"book_id" json:"bookId,omitempty"`
	Title      string     `db:"title" json:"title"`
	Desc       *string    `db:"description" json:"description,omitempty"`
	Location   *string    `db:"location" json:"location,omitempty"`
	Lat        *float64   `db:"lat" json:"lat,omitempty"`
	Lng        *float64   `db:"lng" json:"lng,omitempty"`
	Start      time.Time  `db:"start_time" json:"startTime"`
	End        *time.Time `db:"end_time" json:"endTime,omitempty"`
	IsPaid     bool       `db:"is_paid" json:"isPaid"`
	Price      *float64   `db:"price_amount" json:"priceAmount,omitempty"`
	Currency   *string    `db:"price_currency" json:"priceCurrency,omitempty"`
	Capacity   *int       `db:"capacity" json:"capacity,omitempty"`
	Visibility string     `db:"visibility" json:"visibility"`
}

type Chapter struct {
	ID        string  `db:"id" json:"id"`
	BookID    string  `db:"book_id" json:"bookId"`
	Number    int     `db:"number" json:"number"`
	Title     string  `db:"title" json:"title"`
	PageStart *int    `db:"page_no_start" json:"pageNoStart,omitempty"`
	Summary   *string `db:"summary" json:"summary,omitempty"`
	AudioURL  *string `db:"audio_url" json:"audioUrl,omitempty"`
}

type BookAsset struct {
	ID        string  `db:"id" json:"id"`
	BookID    string  `db:"book_id" json:"bookId"`
	Kind      string  `db:"kind" json:"kind"`
	Source    string  `db:"source" json:"source"`
	URL       string  `db:"url" json:"url"`
	Checksum  *string `db:"checksum" json:"checksum,omitempty"`
	SizeBytes *int64  `db:"file_size_bytes" json:"fileSizeBytes,omitempty"`
}

type Collaboration struct {
	ID            string     `db:"id" json:"id"`
	PublisherID   string     `db:"publisher_id" json:"publisherId"`
	PublisherName *string    `db:"publisher_name" json:"publisherName,omitempty"`
	AuthorID      string     `db:"author_id" json:"authorId"`
	AuthorName    *string    `db:"author_name" json:"authorName,omitempty"`
	Status        string     `db:"status" json:"status"`
	Notes         *string    `db:"notes" json:"notes,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
	RespondedAt   *time.Time `db:"responded_at" json:"respondedAt,omitempty"`
}

type ClubSession struct {
	ID           string    `db:"id" json:"id"`
	ClubID       string    `db:"club_id" json:"clubId"`
	BookID       *string   `db:"book_id" json:"bookId,omitempty"`
	HostUserID   string    `db:"host_user_id" json:"hostUserId"`
	SessionType  string    `db:"session_type" json:"sessionType"`
	ScheduledAt  time.Time `db:"scheduled_at" json:"scheduledAt"`
	DurationMin  *int      `db:"duration_minutes" json:"durationMinutes,omitempty"`
	Status       string    `db:"status" json:"status"`
	MeetingURL   *string   `db:"meeting_url" json:"meetingUrl,omitempty"`
	RecordingURL *string   `db:"recording_url" json:"recordingUrl,omitempty"`
}

type SessionParticipant struct {
	SessionID string    `db:"session_id" json:"sessionId"`
	UserID    string    `db:"user_id" json:"userId"`
	JoinedAt  time.Time `db:"joined_at" json:"joinedAt"`
	Role      string    `db:"role" json:"role"`
}

type AuthorSummary struct {
	ID           string   `db:"id" json:"id"`
	DisplayName  string   `db:"display_name" json:"displayName"`
	Bio          *string  `db:"bio" json:"bio,omitempty"`
	Genres       []string `db:"genres" json:"genres,omitempty"`
	Experience   *int     `db:"experience_years" json:"experienceYears,omitempty"`
	Status       string   `db:"status" json:"status"`
	ProfileImage *string  `db:"profile_image_url" json:"profileImageUrl,omitempty"`
	BooksCount   int      `db:"books_count" json:"booksCount"`
}

type MessageReaction struct {
	MessageID string    `db:"message_id" json:"messageId"`
	UserID    string    `db:"user_id" json:"userId"`
	Reaction  string    `db:"reaction" json:"reaction"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type MessageReport struct {
	ID         string     `db:"id" json:"id"`
	MessageID  string     `db:"message_id" json:"messageId"`
	ReporterID string     `db:"reporter_user_id" json:"reporterUserId"`
	Reason     *string    `db:"reason" json:"reason,omitempty"`
	Status     string     `db:"status" json:"status"`
	HandledBy  *string    `db:"handled_by_user_id" json:"handledByUserId,omitempty"`
	HandledAt  *time.Time `db:"handled_at" json:"handledAt,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
}

func New(opts Options) *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     opts.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Authorization,Content-Type,Accept",
		AllowCredentials: true,
	}))

	// Health
	app.Get("/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "ok"})
	})

	// Public: catalog search
	app.Get("/v1/public/books", func(c *fiber.Ctx) error {
		if opts.DB == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false})
		}
		q := strings.TrimSpace(c.Query("q"))
		rows := []Book{}
		selectCols := `b.id, b.author_id, b.publisher_id, b.title, b.subtitle, b.description, b.isbn, b.language, b.categories, b.cover_image_url, b.ebook_file_url, b.price_amount, b.price_currency, b.shareable, b.allow_hardcopy, b.visibility, COALESCE(a.display_name,'') AS author_name, COALESCE(p.name,'') AS publisher_name, b.slug, b.status, b.release_date, b.page_count, b.preview_html, b.ebook_enabled, b.hardcopy_enabled`
		if q == "" {
			if err := opts.DB.Select(&rows, `SELECT `+selectCols+` FROM books b LEFT JOIN authors a ON a.id=b.author_id LEFT JOIN publishers p ON p.id=b.publisher_id WHERE b.visibility='public' AND b.status='published' ORDER BY b.created_at DESC LIMIT 50`); err != nil {
				return c.Status(500).JSON(fiber.Map{"success": false})
			}
		} else {
			if err := opts.DB.Select(&rows, `SELECT `+selectCols+` FROM books b LEFT JOIN authors a ON a.id=b.author_id LEFT JOIN publishers p ON p.id=b.publisher_id WHERE b.visibility='public' AND b.status='published' AND (LOWER(b.title) LIKE LOWER($1) OR LOWER(b.description) LIKE LOWER($1)) ORDER BY b.created_at DESC LIMIT 50`, "%"+q+"%"); err != nil {
				return c.Status(500).JSON(fiber.Map{"success": false})
			}
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})

	app.Get("/v1/public/books/:id", func(c *fiber.Ctx) error {
		if opts.DB == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false})
		}
		identifier := strings.TrimSpace(c.Params("id"))
		if identifier == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		selectCols := `b.id, b.author_id, b.publisher_id, b.title, b.subtitle, b.description, b.isbn, b.language, b.categories, b.cover_image_url, b.ebook_file_url, b.price_amount, b.price_currency, b.shareable, b.allow_hardcopy, b.visibility, COALESCE(a.display_name,'') AS author_name, COALESCE(p.name,'') AS publisher_name, b.slug, b.status, b.release_date, b.page_count, b.preview_html, b.ebook_enabled, b.hardcopy_enabled`
		var book Book
		err := opts.DB.QueryRowx(`SELECT `+selectCols+` FROM books b LEFT JOIN authors a ON a.id=b.author_id LEFT JOIN publishers p ON p.id=b.publisher_id WHERE b.status='published' AND b.visibility='public' AND (b.slug=$1 OR b.id=$1)`, identifier).StructScan(&book)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		assets := []BookAsset{}
		if err := opts.DB.Select(&assets, `SELECT id, book_id, kind, source, url, checksum, file_size_bytes FROM book_assets WHERE book_id=$1 ORDER BY created_at ASC`, book.ID); err != nil {
			assets = []BookAsset{}
		}
		chapters := []Chapter{}
		if err := opts.DB.Select(&chapters, `SELECT id, book_id, number, title, page_no_start, summary, audio_url FROM book_chapters WHERE book_id=$1 ORDER BY number ASC`, book.ID); err != nil {
			chapters = []Chapter{}
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{
			"book":     book,
			"assets":   assets,
			"chapters": chapters,
		}})
	})

	app.Get("/v1/public/authors", func(c *fiber.Ctx) error {
		if opts.DB == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false})
		}
		search := strings.TrimSpace(c.Query("q"))
		genre := strings.TrimSpace(c.Query("genre"))
		args := []any{}
		conds := []string{"a.status='active'"}
		nextPlaceholder := func() string { return "$" + strconv.Itoa(len(args)+1) }
		if search != "" {
			ph := nextPlaceholder()
			args = append(args, "%"+strings.ToLower(search)+"%")
			conds = append(conds, "(LOWER(a.display_name) LIKE "+ph+" OR LOWER(COALESCE(a.bio,'')) LIKE "+ph+")")
		}
		if genre != "" {
			ph := nextPlaceholder()
			args = append(args, genre)
			conds = append(conds, ph+" = ANY(a.genres)")
		}
		whereClause := strings.Join(conds, " AND ")
		query := `
            SELECT a.id, a.display_name, a.bio, COALESCE(a.genres, ARRAY[]::text[]) AS genres,
                   a.experience_years, a.status, a.profile_image_url,
                   COALESCE((SELECT COUNT(*) FROM books b WHERE b.author_id=a.id), 0) AS books_count
            FROM authors a
            WHERE ` + whereClause + `
            ORDER BY a.display_name
            LIMIT 100`
		rows := []AuthorSummary{}
		if err := opts.DB.Select(&rows, query, args...); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})

	app.Get("/v1/public/meetups", func(c *fiber.Ctx) error {
		if opts.DB == nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		q := strings.TrimSpace(c.Query("q"))
		rows := []Meetup{}
		if q == "" {
			if err := opts.DB.Select(&rows, `SELECT id, organizer_user_id, club_id, book_id, title, description, location, lat, lng, start_time, end_time, is_paid, price_amount, price_currency, capacity, visibility FROM meetups WHERE start_time >= now() AND visibility='public' ORDER BY start_time ASC LIMIT 100`); err != nil {
				return c.Status(500).JSON(fiber.Map{"success": false})
			}
		} else {
			if err := opts.DB.Select(&rows, `SELECT id, organizer_user_id, club_id, book_id, title, description, location, lat, lng, start_time, end_time, is_paid, price_amount, price_currency, capacity, visibility FROM meetups WHERE start_time >= now() AND visibility='public' AND (LOWER(title) LIKE LOWER($1) OR LOWER(description) LIKE LOWER($1) OR LOWER(location) LIKE LOWER($1)) ORDER BY start_time ASC LIMIT 100`, "%"+q+"%"); err != nil {
				return c.Status(500).JSON(fiber.Map{"success": false})
			}
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})

	// Authenticated routes
	requireAuth := auth.Middleware(auth.Options{HS256Secret: opts.AuthHS256, Env: opts.Env, CoreAPIBase: opts.CoreAPIBase})

	// List my books (as author or publisher)
	app.Get("/v1/me/books", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		rows := []Book{}
		q := `SELECT b.id, b.author_id, b.publisher_id, b.title, b.subtitle, b.description, b.isbn, b.language, b.categories, b.cover_image_url, b.ebook_file_url,
                    b.price_amount, b.price_currency, b.shareable, b.allow_hardcopy, b.visibility,
                    COALESCE(a.display_name,'') AS author_name, COALESCE(p.name,'') AS publisher_name,
                    b.slug, b.status, b.release_date, b.page_count, b.preview_html, b.ebook_enabled, b.hardcopy_enabled
              FROM books b
              LEFT JOIN authors a ON a.id=b.author_id
              LEFT JOIN publishers p ON p.id=b.publisher_id
              WHERE (a.core_user_id=$1) OR (p.core_user_id=$1)
              ORDER BY b.created_at DESC
              LIMIT 200`
		if err := opts.DB.Select(&rows, q, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})

	app.Get("/v1/books/:id", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		selectCols := `b.id, b.author_id, b.publisher_id, b.title, b.subtitle, b.description, b.isbn, b.language, b.categories, b.cover_image_url, b.ebook_file_url, b.price_amount, b.price_currency, b.shareable, b.allow_hardcopy, b.visibility, COALESCE(a.display_name,'') AS author_name, COALESCE(p.name,'') AS publisher_name, b.slug, b.status, b.release_date, b.page_count, b.preview_html, b.ebook_enabled, b.hardcopy_enabled`
		var book Book
		if err := opts.DB.QueryRowx(`SELECT `+selectCols+` FROM books b LEFT JOIN authors a ON a.id=b.author_id LEFT JOIN publishers p ON p.id=b.publisher_id WHERE b.id=$1`, bookID).StructScan(&book); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		assets := []BookAsset{}
		_ = opts.DB.Select(&assets, `SELECT id, book_id, kind, source, url, checksum, file_size_bytes FROM book_assets WHERE book_id=$1 ORDER BY created_at DESC`, bookID)
		chapters := []Chapter{}
		_ = opts.DB.Select(&chapters, `SELECT id, book_id, number, title, page_no_start, summary, audio_url FROM book_chapters WHERE book_id=$1 ORDER BY number ASC`, bookID)
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{
			"book":     book,
			"assets":   assets,
			"chapters": chapters,
		}})
	})

	app.Post("/v1/books", requireAuth, func(c *fiber.Ctx) error {
		var in struct {
			Title       string  `json:"title"`
			Description *string `json:"description"`
			AuthorID    *string `json:"authorId"`
			PublisherID *string `json:"publisherId"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Title) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "invalid body"})
		}
		// exactly one of authorId or publisherId must be provided
		if (in.AuthorID == nil && in.PublisherID == nil) || (in.AuthorID != nil && in.PublisherID != nil) {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "provide authorId or publisherId"})
		}
		uid := auth.UserID(c)
		slug := generateSlug(opts.DB, in.Title)
		// verify ownership of provided profile
		if in.AuthorID != nil {
			var exists int
			if err := opts.DB.QueryRowx(`SELECT 1 FROM authors WHERE id=$1 AND core_user_id=$2`, *in.AuthorID, uid).Scan(&exists); err != nil || exists != 1 {
				return c.Status(403).JSON(fiber.Map{"success": false, "message": "not your author profile"})
			}
			var id string
			if err := opts.DB.Get(&id, `INSERT INTO books (author_id, title, description, visibility, slug, created_by_user_id, updated_by_user_id) VALUES ($1,$2,$3,'private',$4,$5,$5) RETURNING id`, *in.AuthorID, in.Title, in.Description, slug, uid); err != nil {
				return c.Status(500).JSON(fiber.Map{"success": false})
			}
			return c.JSON(fiber.Map{"success": true, "id": id})
		}
		if in.PublisherID != nil {
			var exists int
			if err := opts.DB.QueryRowx(`SELECT 1 FROM publishers WHERE id=$1 AND core_user_id=$2`, *in.PublisherID, uid).Scan(&exists); err != nil || exists != 1 {
				return c.Status(403).JSON(fiber.Map{"success": false, "message": "not your publisher profile"})
			}
			var id string
			if err := opts.DB.Get(&id, `INSERT INTO books (publisher_id, title, description, visibility, slug, created_by_user_id, updated_by_user_id) VALUES ($1,$2,$3,'private',$4,$5,$5) RETURNING id`, *in.PublisherID, in.Title, in.Description, slug, uid); err != nil {
				return c.Status(500).JSON(fiber.Map{"success": false})
			}
			return c.JSON(fiber.Map{"success": true, "id": id})
		}
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "invalid body"})
	})

	app.Patch("/v1/books/:id", requireAuth, func(c *fiber.Ctx) error {
		id := c.Params("id")
		uid := auth.UserID(c)
		// verify ownership (author or publisher profile belonging to user)
		var permitted int
		_ = opts.DB.QueryRowx(`SELECT 1 FROM books b
            LEFT JOIN authors a ON a.id=b.author_id
            LEFT JOIN publishers p ON p.id=b.publisher_id
            WHERE b.id=$1 AND (a.core_user_id=$2 OR p.core_user_id=$2)`, id, uid).Scan(&permitted)
		if permitted != 1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "forbidden"})
		}
		var in struct {
			Title           *string  `json:"title"`
			Subtitle        *string  `json:"subtitle"`
			Description     *string  `json:"description"`
			ISBN            *string  `json:"isbn"`
			Language        *string  `json:"language"`
			CoverImageURL   *string  `json:"coverImageUrl"`
			PriceAmount     *float64 `json:"priceAmount"`
			PriceCurrency   *string  `json:"priceCurrency"`
			Shareable       *bool    `json:"shareable"`
			AllowHardcopy   *bool    `json:"allowHardcopy"`
			EbookEnabled    *bool    `json:"ebookEnabled"`
			HardcopyEnabled *bool    `json:"hardcopyEnabled"`
			PreviewHTML     *string  `json:"previewHtml"`
			ReleaseDate     *string  `json:"releaseDate"`
			Visibility      *string  `json:"visibility"`
			PageCount       *int     `json:"pageCount"`
			Slug            *string  `json:"slug"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		sets := []string{}
		args := []any{}
		idx := 0

		if in.Title != nil {
			val := strings.TrimSpace(*in.Title)
			if val == "" {
				return c.Status(400).JSON(fiber.Map{"success": false, "message": "title required"})
			}
			sets = append(sets, "title=$"+itoa(&idx))
			args = append(args, val)
		}
		if in.Subtitle != nil {
			val := strings.TrimSpace(*in.Subtitle)
			if val == "" {
				sets = append(sets, "subtitle=NULL")
			} else {
				sets = append(sets, "subtitle=$"+itoa(&idx))
				args = append(args, val)
			}
		}
		if in.Description != nil {
			if strings.TrimSpace(*in.Description) == "" {
				sets = append(sets, "description=NULL")
			} else {
				sets = append(sets, "description=$"+itoa(&idx))
				args = append(args, *in.Description)
			}
		}
		if in.ISBN != nil {
			if strings.TrimSpace(*in.ISBN) == "" {
				sets = append(sets, "isbn=NULL")
			} else {
				sets = append(sets, "isbn=$"+itoa(&idx))
				args = append(args, strings.TrimSpace(*in.ISBN))
			}
		}
		if in.Language != nil {
			if strings.TrimSpace(*in.Language) == "" {
				sets = append(sets, "language=NULL")
			} else {
				sets = append(sets, "language=$"+itoa(&idx))
				args = append(args, strings.TrimSpace(*in.Language))
			}
		}
		if in.CoverImageURL != nil {
			if strings.TrimSpace(*in.CoverImageURL) == "" {
				sets = append(sets, "cover_image_url=NULL")
			} else {
				sets = append(sets, "cover_image_url=$"+itoa(&idx))
				args = append(args, strings.TrimSpace(*in.CoverImageURL))
			}
		}
		if in.PriceAmount != nil {
			sets = append(sets, "price_amount=$"+itoa(&idx))
			args = append(args, *in.PriceAmount)
		}
		if in.PriceCurrency != nil {
			if strings.TrimSpace(*in.PriceCurrency) == "" {
				sets = append(sets, "price_currency=NULL")
			} else {
				sets = append(sets, "price_currency=$"+itoa(&idx))
				args = append(args, strings.TrimSpace(*in.PriceCurrency))
			}
		}
		if in.Shareable != nil {
			sets = append(sets, "shareable=$"+itoa(&idx))
			args = append(args, *in.Shareable)
		}
		if in.AllowHardcopy != nil {
			sets = append(sets, "allow_hardcopy=$"+itoa(&idx))
			args = append(args, *in.AllowHardcopy)
		}
		if in.EbookEnabled != nil {
			sets = append(sets, "ebook_enabled=$"+itoa(&idx))
			args = append(args, *in.EbookEnabled)
		}
		if in.HardcopyEnabled != nil {
			sets = append(sets, "hardcopy_enabled=$"+itoa(&idx))
			args = append(args, *in.HardcopyEnabled)
		}
		if in.PreviewHTML != nil {
			if strings.TrimSpace(*in.PreviewHTML) == "" {
				sets = append(sets, "preview_html=NULL")
			} else {
				sets = append(sets, "preview_html=$"+itoa(&idx))
				args = append(args, *in.PreviewHTML)
			}
		}
		if in.ReleaseDate != nil {
			trimmed := strings.TrimSpace(*in.ReleaseDate)
			if trimmed == "" {
				sets = append(sets, "release_date=NULL")
			} else {
				if parsed, err := time.Parse("2006-01-02", trimmed); err == nil {
					sets = append(sets, "release_date=$"+itoa(&idx))
					args = append(args, parsed)
				} else {
					return c.Status(400).JSON(fiber.Map{"success": false, "message": "releaseDate must be YYYY-MM-DD"})
				}
			}
		}
		if in.Visibility != nil {
			vis := strings.ToLower(strings.TrimSpace(*in.Visibility))
			if vis != "public" && vis != "private" {
				return c.Status(400).JSON(fiber.Map{"success": false, "message": "visibility must be public or private"})
			}
			sets = append(sets, "visibility=$"+itoa(&idx))
			args = append(args, vis)
		}
		if in.PageCount != nil {
			if *in.PageCount < 0 {
				return c.Status(400).JSON(fiber.Map{"success": false, "message": "pageCount must be >= 0"})
			}
			sets = append(sets, "page_count=$"+itoa(&idx))
			args = append(args, *in.PageCount)
		}
		if in.Slug != nil {
			desired := slugify(*in.Slug)
			if desired == "" {
				sets = append(sets, "slug=NULL")
			} else {
				if slugTakenByAnother(opts.DB, desired, id) {
					return c.Status(fiber.StatusConflict).JSON(fiber.Map{"success": false, "message": "slug already taken"})
				}
				sets = append(sets, "slug=$"+itoa(&idx))
				args = append(args, desired)
			}
		}

		if len(sets) == 0 {
			return c.JSON(fiber.Map{"success": true})
		}
		sets = append(sets, "updated_at=now()")
		sets = append(sets, "updated_by_user_id=$"+itoa(&idx))
		args = append(args, uid)
		args = append(args, id)
		q := "UPDATE books SET " + strings.Join(sets, ", ") + " WHERE id=$" + itoa(&idx)
		if _, err := opts.DB.Exec(q, args...); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	// Secure ebook URL for reading (requires ownership or completed ebook purchase)
	app.Get("/v1/books/:id/ebook-url", requireAuth, func(c *fiber.Ctx) error {
		id := c.Params("id")
		uid := auth.UserID(c)
		var authorID, publisherID *string
		var ebookURL *string
		if err := opts.DB.QueryRowx(`SELECT author_id, publisher_id, ebook_file_url FROM books WHERE id=$1`, id).Scan(&authorID, &publisherID, &ebookURL); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "not found"})
		}
		if ebookURL == nil || strings.TrimSpace(*ebookURL) == "" {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "no ebook available"})
		}
		allowed := userOwnsBook(opts, uid, id)
		if !allowed {
			var exists int
			_ = opts.DB.QueryRowx(`SELECT 1 FROM purchases WHERE user_id=$1 AND book_id=$2 AND kind='ebook' AND status='completed' LIMIT 1`, uid, id).Scan(&exists)
			if exists == 1 {
				allowed = true
			}
		}
		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "not allowed"})
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"url": *ebookURL}})
	})

	// On-platform reading: list pages (meta)
	app.Get("/v1/books/:id/pages", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !canReadBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		type P struct {
			PageNo  int     `db:"page_no" json:"pageNo"`
			Section string  `db:"section" json:"section"`
			Label   *string `db:"label" json:"label,omitempty"`
			Audio   *string `db:"audio_url" json:"audioUrl,omitempty"`
			Video   *string `db:"video_url" json:"videoUrl,omitempty"`
		}
		rows := []P{}
		if err := opts.DB.Select(&rows, `SELECT page_no, section, label, audio_url, video_url FROM book_pages WHERE book_id=$1 ORDER BY page_no ASC`, bookID); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		total := 0
		if err := opts.DB.Get(&total, `SELECT COUNT(*) FROM book_pages WHERE book_id=$1`, bookID); err != nil {
			total = len(rows)
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"pages": rows, "totalPages": total}})
	})

	// On-platform reading: get page content
	app.Get("/v1/books/:id/pages/:no", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !canReadBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		no := strings.TrimSpace(c.Params("no"))
		var pageNo int
		if v, err := strconv.Atoi(no); err == nil && v > 0 {
			pageNo = v
		} else {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		var html sql.NullString
		var section sql.NullString
		var label sql.NullString
		var audio sql.NullString
		var video sql.NullString
		if err := opts.DB.QueryRowx(`SELECT content_html, section, label, audio_url, video_url FROM book_pages WHERE book_id=$1 AND page_no=$2`, bookID, pageNo).Scan(&html, &section, &label, &audio, &video); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "page not found"})
		}
		total := 0
		_ = opts.DB.Get(&total, `SELECT COUNT(*) FROM book_pages WHERE book_id=$1`, bookID)
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{
			"pageNo":     pageNo,
			"totalPages": total,
			"html":       html.String,
			"section":    section.String,
			"label":      nullableString(label),
			"audioUrl":   nullableString(audio),
			"videoUrl":   nullableString(video),
		}})
	})

	// Authoring: upsert a page (owner only)
	app.Post("/v1/books/:id/pages", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		var in struct {
			PageNo   int     `json:"pageNo"`
			Html     string  `json:"html"`
			Section  *string `json:"section"`
			Label    *string `json:"label"`
			AudioURL *string `json:"audioUrl"`
			VideoURL *string `json:"videoUrl"`
		}
		if err := c.BodyParser(&in); err != nil || in.PageNo <= 0 {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		section := "content"
		if in.Section != nil {
			sec := strings.ToLower(strings.TrimSpace(*in.Section))
			switch sec {
			case "frontmatter", "content", "backmatter":
				section = sec
			default:
				return c.Status(400).JSON(fiber.Map{"success": false, "message": "section must be frontmatter, content, or backmatter"})
			}
		}
		label := strings.TrimSpace(orEmpty(in.Label))
		if label == "" {
			label = ""
		}
		audio := strings.TrimSpace(orEmpty(in.AudioURL))
		video := strings.TrimSpace(orEmpty(in.VideoURL))
		// upsert
		if _, err := opts.DB.Exec(`INSERT INTO book_pages (book_id, page_no, content_html, section, label, audio_url, video_url) VALUES ($1,$2,$3,$4,NULLIF($5,''),NULLIF($6,''),NULLIF($7,''))
            ON CONFLICT (book_id, page_no) DO UPDATE SET content_html=EXCLUDED.content_html, section=EXCLUDED.section, label=EXCLUDED.label, audio_url=EXCLUDED.audio_url, video_url=EXCLUDED.video_url, updated_at=now()`,
			bookID, in.PageNo, in.Html, section, label, audio, video); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	// Chapters: list (auth + can read)
	app.Get("/v1/books/:id/chapters", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !canReadBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		rows := []Chapter{}
		if err := opts.DB.Select(&rows, `SELECT id, book_id, number, title, page_no_start, summary, audio_url FROM book_chapters WHERE book_id=$1 ORDER BY number ASC`, bookID); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})
	// Chapters: upsert (owner only)
	app.Post("/v1/books/:id/chapters", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		var in struct {
			Number      int     `json:"number"`
			Title       string  `json:"title"`
			PageNoStart *int    `json:"pageNoStart"`
			Summary     *string `json:"summary"`
			AudioURL    *string `json:"audioUrl"`
		}
		if err := c.BodyParser(&in); err != nil || in.Number <= 0 || strings.TrimSpace(in.Title) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		summary := strings.TrimSpace(orEmpty(in.Summary))
		audio := strings.TrimSpace(orEmpty(in.AudioURL))
		if _, err := opts.DB.Exec(`INSERT INTO book_chapters (book_id, number, title, page_no_start, summary, audio_url) VALUES ($1,$2,$3,$4,NULLIF($5,''),NULLIF($6,''))
            ON CONFLICT (book_id, number) DO UPDATE SET title=EXCLUDED.title, page_no_start=EXCLUDED.page_no_start, summary=EXCLUDED.summary, audio_url=EXCLUDED.audio_url, updated_at=now()`,
			bookID, in.Number, strings.TrimSpace(in.Title), in.PageNoStart, summary, audio); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Get("/v1/books/:id/assets", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		assets := []BookAsset{}
		if err := opts.DB.Select(&assets, `SELECT id, book_id, kind, source, url, checksum, file_size_bytes FROM book_assets WHERE book_id=$1 ORDER BY created_at DESC`, bookID); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": assets})
	})

	app.Post("/v1/books/:id/assets", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		var in struct {
			Kind          string  `json:"kind"`
			Source        string  `json:"source"`
			URL           string  `json:"url"`
			Checksum      *string `json:"checksum"`
			FileSizeBytes *int64  `json:"fileSizeBytes"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		kind := strings.ToLower(strings.TrimSpace(in.Kind))
		if kind != "ebook" && kind != "audio" && kind != "video" && kind != "supplement" && kind != "hardcopy" {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "invalid asset kind"})
		}
		source := strings.ToLower(strings.TrimSpace(in.Source))
		if source != "upload" && source != "external" {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "invalid source"})
		}
		if strings.TrimSpace(in.URL) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "url required"})
		}
		var id string
		if err := opts.DB.QueryRowx(`INSERT INTO book_assets (book_id, kind, source, url, checksum, file_size_bytes, created_by_user_id)
            VALUES ($1,$2,$3,$4,NULLIF($5,''),$6,$7) RETURNING id`,
			bookID, kind, source, in.URL, orEmpty(in.Checksum), in.FileSizeBytes, uid).Scan(&id); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"id": id}})
	})

	app.Delete("/v1/books/:id/assets/:assetId", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		assetID := c.Params("assetId")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		if _, err := opts.DB.Exec(`DELETE FROM book_assets WHERE id=$1 AND book_id=$2`, assetID, bookID); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/books/:id/publish", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		var in struct {
			ReleaseDate *string `json:"releaseDate"`
		}
		_ = c.BodyParser(&in)
		var release interface{}
		if in.ReleaseDate != nil {
			trimmed := strings.TrimSpace(*in.ReleaseDate)
			if trimmed != "" {
				d, err := time.Parse("2006-01-02", trimmed)
				if err != nil {
					return c.Status(400).JSON(fiber.Map{"success": false, "message": "releaseDate must be YYYY-MM-DD"})
				}
				release = d
			}
		}
		if _, err := opts.DB.Exec(`UPDATE books SET status='published', visibility='public', release_date=COALESCE($2, release_date), updated_at=now(), updated_by_user_id=$3 WHERE id=$1`, bookID, release, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/books/:id/unpublish", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		if _, err := opts.DB.Exec(`UPDATE books SET status='draft', visibility='private', updated_at=now(), updated_by_user_id=$2 WHERE id=$1`, bookID, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	// Profiles: authors/publishers for current user
	app.Get("/v1/me/authors", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		type A struct {
			ID   string `db:"id" json:"id"`
			Name string `db:"display_name" json:"displayName"`
		}
		var rows []A
		if err := opts.DB.Select(&rows, `SELECT id, display_name FROM authors WHERE core_user_id=$1 ORDER BY display_name`, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})
	app.Post("/v1/me/authors", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		var in struct {
			DisplayName string  `json:"displayName"`
			Bio         *string `json:"bio"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.DisplayName) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		var id string
		if err := opts.DB.QueryRowx(`INSERT INTO authors (core_user_id, display_name, bio) VALUES ($1,$2,$3) RETURNING id`, uid, in.DisplayName, in.Bio).Scan(&id); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"id": id}})
	})
	app.Get("/v1/me/publishers", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		type P struct {
			ID   string `db:"id" json:"id"`
			Name string `db:"name" json:"name"`
		}
		var rows []P
		if err := opts.DB.Select(&rows, `SELECT id, name FROM publishers WHERE core_user_id=$1 ORDER BY name`, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})
	app.Post("/v1/me/publishers", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		var in struct {
			Name        string  `json:"name"`
			Description *string `json:"description"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Name) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		var id string
		if err := opts.DB.QueryRowx(`INSERT INTO publishers (core_user_id, name, description) VALUES ($1,$2,$3) RETURNING id`, uid, in.Name, in.Description).Scan(&id); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"id": id}})
	})

	app.Get("/v1/me/collaborations", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		rows := []Collaboration{}
		query := `
            SELECT c.id, c.publisher_id, p.name AS publisher_name, c.author_id, a.display_name AS author_name,
                   c.status, c.notes, c.created_at, c.responded_at
            FROM publisher_author_collaborations c
            LEFT JOIN publishers p ON p.id=c.publisher_id
            LEFT JOIN authors a ON a.id=c.author_id
            WHERE (p.core_user_id=$1 OR a.core_user_id=$1)
            ORDER BY c.created_at DESC
            LIMIT 200`
		if err := opts.DB.Select(&rows, query, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})

	app.Post("/v1/publisher-author-collaborations", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		var in struct {
			PublisherID string  `json:"publisherId"`
			AuthorID    string  `json:"authorId"`
			Notes       *string `json:"notes"`
		}
		if err := c.BodyParser(&in); err != nil ||
			strings.TrimSpace(in.PublisherID) == "" || strings.TrimSpace(in.AuthorID) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		var owns int
		_ = opts.DB.QueryRowx(`SELECT 1 FROM publishers WHERE id=$1 AND core_user_id=$2`, in.PublisherID, uid).Scan(&owns)
		if owns != 1 {
			return c.Status(403).JSON(fiber.Map{"success": false, "message": "not your publisher profile"})
		}
		var authorExists int
		if err := opts.DB.QueryRowx(`SELECT 1 FROM authors WHERE id=$1`, in.AuthorID).Scan(&authorExists); err != nil || authorExists != 1 {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "author not found"})
		}
		notes := strings.TrimSpace(orEmpty(in.Notes))
		var id string
		if err := opts.DB.QueryRowx(`INSERT INTO publisher_author_collaborations (publisher_id, author_id, notes)
                VALUES ($1,$2,NULLIF($3,''))
                ON CONFLICT (publisher_id, author_id) DO UPDATE SET notes=EXCLUDED.notes, status='pending', responded_at=NULL
                RETURNING id`, in.PublisherID, in.AuthorID, notes).Scan(&id); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"id": id}})
	})

	app.Post("/v1/publisher-author-collaborations/:id/respond", requireAuth, func(c *fiber.Ctx) error {
		collabID := c.Params("id")
		uid := auth.UserID(c)
		var in struct {
			Decision string  `json:"decision"`
			Notes    *string `json:"notes"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		decision := strings.ToLower(strings.TrimSpace(in.Decision))
		if decision != "accept" && decision != "decline" && decision != "reject" {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "decision must be accept or decline"})
		}
		// verify user controls author side
		var authorID string
		var authorUserID string
		err := opts.DB.QueryRowx(`SELECT c.author_id, a.core_user_id FROM publisher_author_collaborations c
            LEFT JOIN authors a ON a.id=c.author_id
            WHERE c.id=$1`, collabID).Scan(&authorID, &authorUserID)
		if err != nil || authorID == "" {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		if authorUserID != uid {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		status := "declined"
		if decision == "accept" {
			status = "accepted"
		}
		notes := strings.TrimSpace(orEmpty(in.Notes))
		if _, err := opts.DB.Exec(`UPDATE publisher_author_collaborations
            SET status=$2, responded_at=now(), notes=COALESCE(NULLIF($3,''), notes)
            WHERE id=$1`, collabID, status, notes); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/purchases", requireAuth, func(c *fiber.Ctx) error {
		var in struct {
			BookID   string `json:"bookId"`
			Kind     string `json:"kind"`
			Provider string `json:"provider"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		if in.Kind != "ebook" && in.Kind != "hardcopy" {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "kind must be ebook or hardcopy"})
		}
		uid := auth.UserID(c)
		// Fetch price details
		var priceAmount *float64
		var priceCurrency *string
		if err := opts.DB.QueryRowx(`SELECT price_amount, price_currency FROM books WHERE id=$1`, in.BookID).Scan(&priceAmount, &priceCurrency); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "book not found"})
		}
		var id string
		if err := opts.DB.Get(&id, `INSERT INTO purchases (user_id, book_id, kind, status, price_paid, currency) VALUES ($1,$2,$3,'pending',NULL,COALESCE($4,'')) RETURNING id`, uid, in.BookID, in.Kind, priceCurrency); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		// If price > 0, require payment intent
		amtCents := int64(0)
		curr := ""
		if priceAmount != nil && priceCurrency != nil && *priceAmount > 0 {
			amtCents = int64((*priceAmount) * 100)
			curr = *priceCurrency
		}
		if amtCents > 0 && curr != "" {
			prov := in.Provider
			if prov == "" {
				prov = "mpesa"
			}
			payload := billing.CreateIntentRequest{AmountCents: amtCents, Currency: curr, Description: "Purchase book " + in.BookID, Provider: prov}
			// Build proxy request
			r := new(http.Request)
			r.Header = http.Header{}
			if v := c.Get("Authorization"); v != "" {
				r.Header.Set("Authorization", v)
			}
			if v := c.Get("Cookie"); v != "" {
				r.Header.Set("Cookie", v)
			}
			if v := c.Get("Origin"); v != "" {
				r.Header.Set("Origin", v)
			}
			if resp, status, err := billing.CreatePaymentIntent(opts.CoreAPIBase, r, payload); err == nil && status >= 200 && status < 300 {
				// Persist mapping from purchase to payment intent id (if available)
				if resp != nil && resp.Data != nil {
					if m, ok := resp.Data.(map[string]any); ok {
						if im, ok := m["intent"].(map[string]any); ok {
							if fid, ok := im["id"].(float64); ok {
								intentID := int64(fid)
								_, _ = opts.DB.Exec(`UPDATE purchases SET payment_intent_id=$1 WHERE id=$2`, intentID, id)
							}
						}
					}
				}
				return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{"success": false, "message": "payment required", "data": fiber.Map{"purchaseId": id, "payment": resp.Data}})
			}
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{"success": false, "message": "payment required", "data": fiber.Map{"purchaseId": id}})
		}
		// Free purchase
		paid := 0.0
		if priceAmount != nil {
			paid = *priceAmount
		}
		currency := ""
		if priceCurrency != nil {
			currency = *priceCurrency
		}
		_, _ = opts.DB.Exec(`UPDATE purchases SET status='completed', price_paid=$2, currency=COALESCE(NULLIF($3,''), currency) WHERE id=$1`, id, paid, currency)
		recordRevenueShares(opts.DB, id, in.BookID, paid, currency)
		return c.JSON(fiber.Map{"success": true, "id": id, "status": "completed"})
	})

	// Confirm a purchase based on payment intent status at Core API
	app.Post("/v1/purchases/:id/confirm", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		pid := c.Params("id")
		var intentID *int64
		var cid string
		// Ensure user owns the purchase
		if err := opts.DB.QueryRowx(`SELECT payment_intent_id, book_id FROM purchases WHERE id=$1 AND user_id=$2`, pid, uid).Scan(&intentID, &cid); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "not found"})
		}
		if intentID == nil {
			return c.JSON(fiber.Map{"success": true, "message": "no intent"})
		}
		r := new(http.Request)
		r.Header = http.Header{}
		if v := c.Get("Authorization"); v != "" {
			r.Header.Set("Authorization", v)
		}
		if v := c.Get("Cookie"); v != "" {
			r.Header.Set("Cookie", v)
		}
		if v := c.Get("Origin"); v != "" {
			r.Header.Set("Origin", v)
		}
		resp, status, err := billing.GetPaymentIntent(opts.CoreAPIBase, r, *intentID)
		if err != nil || status < 200 || status >= 300 || resp == nil || resp.Data == nil {
			return c.Status(502).JSON(fiber.Map{"success": false, "message": "billing unavailable"})
		}
		var curStatus string
		if m, ok := resp.Data.(map[string]any); ok {
			if s, ok := m["status"].(string); ok {
				curStatus = s
			}
			// Capture paid amount and currency if needed in future
		}
		if strings.ToLower(curStatus) == "succeeded" || strings.ToLower(curStatus) == "paid" {
			var priceAmount *float64
			var priceCurrency *string
			_ = opts.DB.QueryRowx(`SELECT price_amount, price_currency FROM books WHERE id=$1`, cid).Scan(&priceAmount, &priceCurrency)
			paid := 0.0
			if priceAmount != nil {
				paid = *priceAmount
			}
			currency := ""
			if priceCurrency != nil {
				currency = *priceCurrency
			}
			_, _ = opts.DB.Exec(`UPDATE purchases SET status='completed', price_paid=$2, currency=COALESCE(NULLIF($3,''), currency) WHERE id=$1`, pid, paid, currency)
			recordRevenueShares(opts.DB, pid, cid, paid, currency)
			return c.JSON(fiber.Map{"success": true, "message": "completed"})
		}
		return c.JSON(fiber.Map{"success": true, "message": curStatus})
	})

	// List my purchases
	app.Get("/v1/me/purchases", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		type item struct {
			ID        string    `db:"id" json:"id"`
			BookID    string    `db:"book_id" json:"bookId"`
			Kind      string    `db:"kind" json:"kind"`
			Status    string    `db:"status" json:"status"`
			Currency  *string   `db:"currency" json:"currency,omitempty"`
			PricePaid *float64  `db:"price_paid" json:"pricePaid,omitempty"`
			CreatedAt time.Time `db:"created_at" json:"createdAt"`
			Title     string    `db:"title" json:"title"`
		}
		rows := []item{}
		if err := opts.DB.Select(&rows, `SELECT p.id, p.book_id, p.kind, p.status, p.currency, p.price_paid, p.created_at, COALESCE(b.title,'') AS title FROM purchases p LEFT JOIN books b ON b.id=p.book_id WHERE p.user_id=$1 ORDER BY p.created_at DESC LIMIT 200`, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})

	app.Get("/v1/me/earnings", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		type earning struct {
			PurchaseID    string    `db:"purchase_id" json:"purchaseId"`
			RecipientType string    `db:"recipient_type" json:"recipientType"`
			RecipientID   string    `db:"recipient_id" json:"recipientId"`
			Amount        float64   `db:"amount" json:"amount"`
			Currency      string    `db:"currency" json:"currency"`
			CreatedAt     time.Time `db:"created_at" json:"createdAt"`
		}
		query := `
            SELECT prs.purchase_id, prs.recipient_type, prs.recipient_id, prs.amount, COALESCE(prs.currency,'' ) AS currency, p.created_at
            FROM purchase_revenue_shares prs
            JOIN purchases p ON p.id=prs.purchase_id
            LEFT JOIN authors a ON prs.recipient_type='author' AND a.id=prs.recipient_id
            LEFT JOIN publishers pu ON prs.recipient_type='publisher' AND pu.id=prs.recipient_id
            WHERE (prs.recipient_type='author' AND a.core_user_id=$1)
               OR (prs.recipient_type='publisher' AND pu.core_user_id=$1)
            ORDER BY p.created_at DESC
            LIMIT 500`
		rows := []earning{}
		if err := opts.DB.Select(&rows, query, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		totals := map[string]float64{}
		for _, e := range rows {
			totals[e.Currency] += e.Amount
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"entries": rows, "totals": totals}})
	})

	app.Post("/v1/books/:id/share-links", requireAuth, func(c *fiber.Ctx) error {
		bookID := c.Params("id")
		uid := auth.UserID(c)
		if !userOwnsBook(opts, uid, bookID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		// simple token
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		token := hex.EncodeToString(b)
		var id string
		if err := opts.DB.Get(&id, `INSERT INTO share_links (book_id, created_by_user_id, token) VALUES ($1,$2,$3) RETURNING id`, bookID, auth.UserID(c), token); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "id": id, "token": token})
	})

	app.Get("/v1/clubs", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		clubs := []Club{}
		if err := opts.DB.Select(&clubs, `SELECT id, owner_user_id, name, description, visibility FROM book_clubs WHERE owner_user_id=$1 OR id IN (SELECT club_id FROM club_members WHERE user_id=$1) ORDER BY name`, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": clubs})
	})

	app.Post("/v1/clubs", requireAuth, func(c *fiber.Ctx) error {
		var in struct {
			Name        string  `json:"name"`
			Description *string `json:"description"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Name) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		uid := auth.UserID(c)
		var clubID string
		if err := opts.DB.Get(&clubID, `INSERT INTO book_clubs (owner_user_id, name, description, visibility) VALUES ($1,$2,$3,'public') RETURNING id`, uid, in.Name, in.Description); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		_, _ = opts.DB.Exec(`INSERT INTO club_members (club_id, user_id, role) VALUES ($1,$2,'admin') ON CONFLICT DO NOTHING`, clubID, uid)
		return c.JSON(fiber.Map{"success": true, "id": clubID})
	})

	app.Post("/v1/clubs/:id/join", requireAuth, func(c *fiber.Ctx) error {
		clubID := c.Params("id")
		uid := auth.UserID(c)
		if _, err := opts.DB.Exec(`INSERT INTO club_members (club_id, user_id, role) VALUES ($1,$2,'member') ON CONFLICT DO NOTHING`, clubID, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/clubs/:id/rooms", requireAuth, func(c *fiber.Ctx) error {
		clubID := c.Params("id")
		uid := auth.UserID(c)
		if !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		var in struct {
			Name string `json:"name"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Name) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		var roomID string
		if err := opts.DB.Get(&roomID, `INSERT INTO club_rooms (club_id, name) VALUES ($1,$2) RETURNING id`, clubID, in.Name); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "id": roomID})
	})

	app.Get("/v1/clubs/:id/rooms", requireAuth, func(c *fiber.Ctx) error {
		clubID := c.Params("id")
		uid := auth.UserID(c)
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		rows := []Room{}
		if err := opts.DB.Select(&rows, `SELECT id, club_id, name FROM club_rooms WHERE club_id=$1 ORDER BY name`, clubID); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})

	app.Get("/v1/rooms/:id/messages", requireAuth, func(c *fiber.Ctx) error {
		roomID := c.Params("id")
		uid := auth.UserID(c)
		clubID := roomClubID(opts.DB, roomID)
		if clubID == "" {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		rows := []Message{}
		if err := opts.DB.Select(&rows, `SELECT id, room_id, user_id, content, created_at FROM room_messages WHERE room_id=$1 ORDER BY created_at DESC LIMIT 100`, roomID); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		blocked := blockedSenders(opts.DB, clubID, uid)
		resp := []fiber.Map{}
		for i := len(rows) - 1; i >= 0; i-- {
			msg := rows[i]
			if _, hidden := blocked[msg.UserID]; hidden {
				continue
			}
			reactions := map[string]int{}
			rws, err := opts.DB.Queryx(`SELECT reaction, COUNT(*) FROM room_message_reactions WHERE message_id=$1 GROUP BY reaction`, msg.ID)
			if err == nil {
				for rws.Next() {
					var reaction string
					var count int
					if err := rws.Scan(&reaction, &count); err == nil {
						reactions[reaction] = count
					}
				}
				rws.Close()
			}
			resp = append(resp, fiber.Map{
				"id":        msg.ID,
				"roomId":    msg.RoomID,
				"userId":    msg.UserID,
				"content":   msg.Content,
				"createdAt": msg.Created,
				"reactions": reactions,
			})
		}
		return c.JSON(fiber.Map{"success": true, "data": resp})
	})

	app.Post("/v1/rooms/:id/messages", requireAuth, func(c *fiber.Ctx) error {
		roomID := c.Params("id")
		var in struct {
			Content string `json:"content"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Content) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		uid := auth.UserID(c)
		clubID := roomClubID(opts.DB, roomID)
		if clubID == "" {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		if userIsMuted(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "muted"})
		}
		if _, err := opts.DB.Exec(`INSERT INTO room_messages (room_id, user_id, content) VALUES ($1,$2,$3)`, roomID, uid, in.Content); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/messages/:id/react", requireAuth, func(c *fiber.Ctx) error {
		messageID := c.Params("id")
		uid := auth.UserID(c)
		var in struct {
			Reaction string `json:"reaction"`
			Remove   bool   `json:"remove"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		var roomID string
		if err := opts.DB.QueryRowx(`SELECT room_id FROM room_messages WHERE id=$1`, messageID).Scan(&roomID); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		clubID := roomClubID(opts.DB, roomID)
		if clubID == "" {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		reaction := strings.ToLower(strings.TrimSpace(in.Reaction))
		if in.Remove {
			if reaction == "" {
				return c.Status(400).JSON(fiber.Map{"success": false, "message": "reaction required to remove"})
			}
			_, _ = opts.DB.Exec(`DELETE FROM room_message_reactions WHERE message_id=$1 AND user_id=$2 AND reaction=$3`, messageID, uid, reaction)
			return c.JSON(fiber.Map{"success": true})
		}
		if !isValidReaction(reaction) {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "invalid reaction"})
		}
		if _, err := opts.DB.Exec(`INSERT INTO room_message_reactions (message_id, user_id, reaction) VALUES ($1,$2,$3)
            ON CONFLICT (message_id, user_id, reaction) DO NOTHING`, messageID, uid, reaction); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Delete("/v1/messages/:id/reactions/:reaction", requireAuth, func(c *fiber.Ctx) error {
		messageID := c.Params("id")
		reaction := strings.ToLower(strings.TrimSpace(c.Params("reaction")))
		if !isValidReaction(reaction) {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		uid := auth.UserID(c)
		var roomID string
		if err := opts.DB.QueryRowx(`SELECT room_id FROM room_messages WHERE id=$1`, messageID).Scan(&roomID); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		clubID := roomClubID(opts.DB, roomID)
		if clubID == "" {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		_, _ = opts.DB.Exec(`DELETE FROM room_message_reactions WHERE message_id=$1 AND user_id=$2 AND reaction=$3`, messageID, uid, reaction)
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/messages/:id/report", requireAuth, func(c *fiber.Ctx) error {
		messageID := c.Params("id")
		uid := auth.UserID(c)
		var in struct {
			Reason *string `json:"reason"`
		}
		_ = c.BodyParser(&in)
		var roomID string
		if err := opts.DB.QueryRowx(`SELECT room_id FROM room_messages WHERE id=$1`, messageID).Scan(&roomID); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		clubID := roomClubID(opts.DB, roomID)
		if clubID == "" {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		reason := strings.TrimSpace(orEmpty(in.Reason))
		if _, err := opts.DB.Exec(`INSERT INTO room_message_reports (message_id, reporter_user_id, reason) VALUES ($1,$2,NULLIF($3,''))`, messageID, uid, reason); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/clubs/:id/mute", requireAuth, func(c *fiber.Ctx) error {
		clubID := c.Params("id")
		uid := auth.UserID(c)
		if !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		var in struct {
			TargetUserID    string `json:"targetUserId"`
			DurationMinutes *int   `json:"durationMinutes"`
			Remove          bool   `json:"remove"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.TargetUserID) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		target := strings.TrimSpace(in.TargetUserID)
		if in.Remove {
			_, _ = opts.DB.Exec(`DELETE FROM club_member_mutes WHERE club_id=$1 AND target_user_id=$2`, clubID, target)
			return c.JSON(fiber.Map{"success": true})
		}
		var expires interface{}
		if in.DurationMinutes != nil && *in.DurationMinutes > 0 {
			expires = time.Now().Add(time.Duration(*in.DurationMinutes) * time.Minute)
		} else {
			expires = nil
		}
		if _, err := opts.DB.Exec(`INSERT INTO club_member_mutes (club_id, target_user_id, muted_by_user_id, expires_at)
            VALUES ($1,$2,$3,$4)
            ON CONFLICT (club_id, target_user_id) DO UPDATE SET muted_by_user_id=EXCLUDED.muted_by_user_id, expires_at=EXCLUDED.expires_at, created_at=now()`,
			clubID, target, uid, expires); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/clubs/:id/block", requireAuth, func(c *fiber.Ctx) error {
		clubID := c.Params("id")
		uid := auth.UserID(c)
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		var in struct {
			TargetUserID string `json:"targetUserId"`
			Remove       bool   `json:"remove"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.TargetUserID) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		target := strings.TrimSpace(in.TargetUserID)
		if in.Remove {
			_, _ = opts.DB.Exec(`DELETE FROM club_member_blocks WHERE club_id=$1 AND blocker_user_id=$2 AND blocked_user_id=$3`, clubID, uid, target)
			return c.JSON(fiber.Map{"success": true})
		}
		if _, err := opts.DB.Exec(`INSERT INTO club_member_blocks (club_id, blocker_user_id, blocked_user_id) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`, clubID, uid, target); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/clubs/:id/sessions", requireAuth, func(c *fiber.Ctx) error {
		clubID := c.Params("id")
		uid := auth.UserID(c)
		if !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		var in struct {
			SessionType     string  `json:"sessionType"`
			ScheduledAt     string  `json:"scheduledAt"`
			DurationMinutes *int    `json:"durationMinutes"`
			BookID          *string `json:"bookId"`
			MeetingURL      *string `json:"meetingUrl"`
			RecordingURL    *string `json:"recordingUrl"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		st := strings.ToLower(strings.TrimSpace(in.SessionType))
		if st != "audio" && st != "video" {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "sessionType must be audio or video"})
		}
		scheduledAt, err := time.Parse(time.RFC3339, strings.TrimSpace(in.ScheduledAt))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "scheduledAt must be RFC3339 timestamp"})
		}
		var duration interface{}
		if in.DurationMinutes != nil && *in.DurationMinutes > 0 {
			duration = *in.DurationMinutes
		} else {
			duration = nil
		}
		var id string
		if err := opts.DB.QueryRowx(`INSERT INTO club_sessions (club_id, book_id, host_user_id, session_type, scheduled_at, duration_minutes, status, meeting_url, recording_url)
            VALUES ($1,$2,$3,$4,$5,$6,'scheduled',NULLIF($7,''),NULLIF($8,'')) RETURNING id`,
			clubID, in.BookID, uid, st, scheduledAt, duration, orEmpty(in.MeetingURL), orEmpty(in.RecordingURL)).Scan(&id); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"id": id}})
	})

	app.Get("/v1/clubs/:id/sessions", requireAuth, func(c *fiber.Ctx) error {
		clubID := c.Params("id")
		uid := auth.UserID(c)
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		rows := []ClubSession{}
		if err := opts.DB.Select(&rows, `SELECT id, club_id, book_id, host_user_id, session_type, scheduled_at, duration_minutes, status, meeting_url, recording_url FROM club_sessions WHERE club_id=$1 ORDER BY scheduled_at DESC LIMIT 50`, clubID); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "data": rows})
	})

	app.Post("/v1/sessions/:id/join", requireAuth, func(c *fiber.Ctx) error {
		sessionID := c.Params("id")
		uid := auth.UserID(c)
		var clubID string
		if err := opts.DB.QueryRowx(`SELECT club_id FROM club_sessions WHERE id=$1`, sessionID).Scan(&clubID); err != nil {
			return c.Status(404).JSON(fiber.Map{"success": false})
		}
		if !isClubMember(opts.DB, clubID, uid) && !isClubAdmin(opts.DB, clubID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false})
		}
		if _, err := opts.DB.Exec(`INSERT INTO club_session_participants (session_id, user_id) VALUES ($1,$2) ON CONFLICT (session_id, user_id) DO UPDATE SET joined_at=now()`, sessionID, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	app.Post("/v1/meetups", requireAuth, func(c *fiber.Ctx) error {
		uid := auth.UserID(c)
		var in struct {
			Title       string     `json:"title"`
			Description *string    `json:"description"`
			Location    *string    `json:"location"`
			Start       time.Time  `json:"startTime"`
			End         *time.Time `json:"endTime"`
			IsPaid      bool       `json:"isPaid"`
			Price       *float64   `json:"priceAmount"`
			Currency    *string    `json:"priceCurrency"`
			Capacity    *int       `json:"capacity"`
			ClubID      *string    `json:"clubId"`
			BookID      *string    `json:"bookId"`
		}
		if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Title) == "" {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		var id string
		if err := opts.DB.Get(&id, `INSERT INTO meetups (organizer_user_id, club_id, book_id, title, description, location, start_time, end_time, is_paid, price_amount, price_currency, capacity, visibility) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,'public') RETURNING id`, uid, in.ClubID, in.BookID, in.Title, in.Description, in.Location, in.Start, in.End, in.IsPaid, in.Price, in.Currency, in.Capacity); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "id": id})
	})

	app.Post("/v1/meetups/:id/register", requireAuth, func(c *fiber.Ctx) error {
		meetupID := c.Params("id")
		uid := auth.UserID(c)
		// Check if paid
		var isPaid bool
		var price *float64
		var curr *string
		_ = opts.DB.QueryRowx(`SELECT is_paid, price_amount, price_currency FROM meetups WHERE id=$1`, meetupID).Scan(&isPaid, &price, &curr)
		if isPaid && price != nil && curr != nil && *price > 0 {
			amtCents := int64((*price) * 100)
			payload := billing.CreateIntentRequest{AmountCents: amtCents, Currency: *curr, Description: "Register meetup " + meetupID, Provider: "mpesa"}
			r := new(http.Request)
			r.Header = http.Header{}
			if v := c.Get("Authorization"); v != "" {
				r.Header.Set("Authorization", v)
			}
			if v := c.Get("Cookie"); v != "" {
				r.Header.Set("Cookie", v)
			}
			if v := c.Get("Origin"); v != "" {
				r.Header.Set("Origin", v)
			}
			if resp, status, err := billing.CreatePaymentIntent(opts.CoreAPIBase, r, payload); err == nil && status >= 200 && status < 300 {
				return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{"success": false, "message": "payment required", "data": resp.Data})
			}
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{"success": false, "message": "payment required"})
		}
		if _, err := opts.DB.Exec(`INSERT INTO meetup_registrations (meetup_id, user_id, status) VALUES ($1,$2,'registered') ON CONFLICT (meetup_id, user_id) DO NOTHING`, meetupID, uid); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	return app
}

// itoa increments idx and returns its string representation
func itoa(idx *int) string {
	*idx = *idx + 1
	return strconv.Itoa(*idx)
}

// canReadBook determines if a user can read a given book via in-app pages.
func canReadBook(opts Options, uid string, bookID string) bool {
	if uid == "" || opts.DB == nil {
		return false
	}
	if userOwnsBook(opts, uid, bookID) {
		return true
	}
	// Completed ebook purchase grants reading
	var exists int
	_ = opts.DB.QueryRowx(`SELECT 1 FROM purchases WHERE user_id=$1 AND book_id=$2 AND kind='ebook' AND status='completed' LIMIT 1`, uid, bookID).Scan(&exists)
	return exists == 1
}

// userOwnsBook: true if user is the author or publisher of the book
func userOwnsBook(opts Options, uid string, bookID string) bool {
	var exists int
	_ = opts.DB.QueryRowx(`SELECT 1 FROM books b
        LEFT JOIN authors a ON a.id=b.author_id
        LEFT JOIN publishers p ON p.id=b.publisher_id
        WHERE b.id=$1 AND (a.core_user_id=$2 OR p.core_user_id=$2)`, bookID, uid).Scan(&exists)
	return exists == 1
}

var slugCleanup = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = slugCleanup.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func randomSlugSuffix() string {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b)
}

func slugExists(db *sqlx.DB, slug string) bool {
	if db == nil || slug == "" {
		return false
	}
	var exists bool
	if err := db.QueryRowx(`SELECT EXISTS (SELECT 1 FROM books WHERE slug=$1)`, slug).Scan(&exists); err != nil {
		return false
	}
	return exists
}

func slugTakenByAnother(db *sqlx.DB, slug string, bookID string) bool {
	if db == nil || slug == "" {
		return false
	}
	var found string
	if err := db.QueryRowx(`SELECT id FROM books WHERE slug=$1 LIMIT 1`, slug).Scan(&found); err != nil {
		return false
	}
	return found != "" && found != bookID
}

func generateSlug(db *sqlx.DB, title string) string {
	base := slugify(title)
	slug := base
	if slug == "" {
		slug = randomSlugSuffix()
	}
	attempt := 0
	for slugExists(db, slug) && attempt < 5 {
		attempt++
		if base == "" {
			slug = randomSlugSuffix()
		} else {
			slug = base + "-" + randomSlugSuffix()
		}
	}
	return slug
}

func orEmpty(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}

func nullableString(ns sql.NullString) *string {
	if ns.Valid {
		v := ns.String
		return &v
	}
	return nil
}

func clubRole(db *sqlx.DB, clubID, userID string) string {
	if db == nil || clubID == "" || userID == "" {
		return ""
	}
	var role string
	if err := db.QueryRowx(`SELECT role FROM club_members WHERE club_id=$1 AND user_id=$2`, clubID, userID).Scan(&role); err != nil {
		return ""
	}
	return role
}

func isClubAdmin(db *sqlx.DB, clubID, userID string) bool {
	if db == nil || clubID == "" || userID == "" {
		return false
	}
	role := clubRole(db, clubID, userID)
	if role == "admin" || role == "moderator" {
		return true
	}
	var owner string
	_ = db.QueryRowx(`SELECT owner_user_id FROM book_clubs WHERE id=$1`, clubID).Scan(&owner)
	return owner == userID
}

func isClubMember(db *sqlx.DB, clubID, userID string) bool {
	if db == nil || clubID == "" || userID == "" {
		return false
	}
	var exists int
	_ = db.QueryRowx(`SELECT 1 FROM club_members WHERE club_id=$1 AND user_id=$2`, clubID, userID).Scan(&exists)
	return exists == 1
}

func roomClubID(db *sqlx.DB, roomID string) string {
	if db == nil || roomID == "" {
		return ""
	}
	var clubID string
	_ = db.QueryRowx(`SELECT club_id FROM club_rooms WHERE id=$1`, roomID).Scan(&clubID)
	return clubID
}

func userIsMuted(db *sqlx.DB, clubID, userID string) bool {
	if db == nil || clubID == "" || userID == "" {
		return false
	}
	var expires sql.NullTime
	if err := db.QueryRowx(`SELECT expires_at FROM club_member_mutes WHERE club_id=$1 AND target_user_id=$2`, clubID, userID).Scan(&expires); err != nil {
		return false
	}
	if !expires.Valid {
		return true
	}
	if time.Now().Before(expires.Time) {
		return true
	}
	// Expired mute - cleanup
	_, _ = db.Exec(`DELETE FROM club_member_mutes WHERE club_id=$1 AND target_user_id=$2`, clubID, userID)
	return false
}

func userHasBlocked(db *sqlx.DB, clubID, blockerID, blockedID string) bool {
	if db == nil || clubID == "" || blockerID == "" || blockedID == "" {
		return false
	}
	var exists int
	_ = db.QueryRowx(`SELECT 1 FROM club_member_blocks WHERE club_id=$1 AND blocker_user_id=$2 AND blocked_user_id=$3`, clubID, blockerID, blockedID).Scan(&exists)
	return exists == 1
}

func blockedSenders(db *sqlx.DB, clubID, viewerID string) map[string]struct{} {
	result := map[string]struct{}{}
	if db == nil || clubID == "" || viewerID == "" {
		return result
	}
	rows, err := db.Queryx(`SELECT blocked_user_id FROM club_member_blocks WHERE club_id=$1 AND blocker_user_id=$2`, clubID, viewerID)
	if err != nil {
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err == nil && user != "" {
			result[user] = struct{}{}
		}
	}
	return result
}

func isValidReaction(reaction string) bool {
	switch reaction {
	case "like", "unlike", "upvote", "downvote":
		return true
	default:
		return false
	}
}

func recordRevenueShares(db *sqlx.DB, purchaseID string, bookID string, amount float64, currency string) {
	if db == nil || purchaseID == "" || bookID == "" {
		return
	}
	_, _ = db.Exec(`DELETE FROM purchase_revenue_shares WHERE purchase_id=$1`, purchaseID)
	if amount <= 0 {
		return
	}
	var authorID, publisherID *string
	_ = db.QueryRowx(`SELECT author_id, publisher_id FROM books WHERE id=$1`, bookID).Scan(&authorID, &publisherID)
	type split struct {
		RecipientType string
		RecipientID   string
		Percent       float64
	}
	splits := []split{}
	if authorID != nil && *authorID != "" && publisherID != nil && *publisherID != "" {
		splits = append(splits, split{"author", *authorID, 0.7})
		splits = append(splits, split{"publisher", *publisherID, 0.3})
	} else if authorID != nil && *authorID != "" {
		splits = append(splits, split{"author", *authorID, 1.0})
	} else if publisherID != nil && *publisherID != "" {
		splits = append(splits, split{"publisher", *publisherID, 1.0})
	}
	for _, s := range splits {
		shareAmount := amount * s.Percent
		_, _ = db.Exec(`INSERT INTO purchase_revenue_shares (purchase_id, recipient_type, recipient_id, share_percent, amount, currency) VALUES ($1,$2,$3,$4,$5,$6)`,
			purchaseID, s.RecipientType, s.RecipientID, s.Percent, shareAmount, currency)
	}
}
