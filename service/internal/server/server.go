package server

import (
    "crypto/rand"
    "encoding/hex"
    "net/http"
    "strings"
    "time"

    "github.com/jmoiron/sqlx"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/berjistech/berjis-ecosystem/books/service/internal/auth"
    "github.com/berjistech/berjis-ecosystem/books/service/internal/billing"
    "strconv"
    "database/sql"
)

type Options struct {
    AllowedOrigins string
    DB             *sqlx.DB
    Env            string
    AuthHS256      string
    CoreAPIBase    string
}

type Book struct {
    ID            string   `db:"id" json:"id"`
    AuthorID      *string  `db:"author_id" json:"authorId,omitempty"`
    PublisherID   *string  `db:"publisher_id" json:"publisherId,omitempty"`
    Title         string   `db:"title" json:"title"`
    Subtitle      *string  `db:"subtitle" json:"subtitle,omitempty"`
    Description   *string  `db:"description" json:"description,omitempty"`
    ISBN          *string  `db:"isbn" json:"isbn,omitempty"`
    Language      *string  `db:"language" json:"language,omitempty"`
    Categories    []string `db:"categories" json:"categories,omitempty"`
    CoverURL      *string  `db:"cover_image_url" json:"coverImageUrl,omitempty"`
    EbookURL      *string  `db:"ebook_file_url" json:"ebookFileUrl,omitempty"`
    PriceAmount   *float64 `db:"price_amount" json:"priceAmount,omitempty"`
    PriceCurrency *string  `db:"price_currency" json:"priceCurrency,omitempty"`
    Shareable     bool     `db:"shareable" json:"shareable"`
    AllowHardcopy bool     `db:"allow_hardcopy" json:"allowHardcopy"`
    Visibility    string   `db:"visibility" json:"visibility"`
}

type Purchase struct {
    ID        string  `db:"id" json:"id"`
    UserID    string  `db:"user_id" json:"userId"`
    BookID    string  `db:"book_id" json:"bookId"`
    Kind      string  `db:"kind" json:"kind"`
    Status    string  `db:"status" json:"status"`
    Paid      *float64 `db:"price_paid" json:"pricePaid,omitempty"`
    Currency  *string `db:"currency" json:"currency,omitempty"`
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
    ID        string    `db:"id" json:"id"`
    OwnerID   string    `db:"owner_user_id" json:"ownerUserId"`
    Name      string    `db:"name" json:"name"`
    Desc      *string   `db:"description" json:"description,omitempty"`
    Visibility string   `db:"visibility" json:"visibility"`
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
    ID        string     `db:"id" json:"id"`
    Organizer string     `db:"organizer_user_id" json:"organizerUserId"`
    ClubID    *string    `db:"club_id" json:"clubId,omitempty"`
    BookID    *string    `db:"book_id" json:"bookId,omitempty"`
    Title     string     `db:"title" json:"title"`
    Desc      *string    `db:"description" json:"description,omitempty"`
    Location  *string    `db:"location" json:"location,omitempty"`
    Lat       *float64   `db:"lat" json:"lat,omitempty"`
    Lng       *float64   `db:"lng" json:"lng,omitempty"`
    Start     time.Time  `db:"start_time" json:"startTime"`
    End       *time.Time `db:"end_time" json:"endTime,omitempty"`
    IsPaid    bool       `db:"is_paid" json:"isPaid"`
    Price     *float64   `db:"price_amount" json:"priceAmount,omitempty"`
    Currency  *string    `db:"price_currency" json:"priceCurrency,omitempty"`
    Capacity  *int       `db:"capacity" json:"capacity,omitempty"`
    Visibility string    `db:"visibility" json:"visibility"`
}

type Chapter struct {
    ID           string  `db:"id" json:"id"`
    BookID       string  `db:"book_id" json:"bookId"`
    Number       int     `db:"number" json:"number"`
    Title        string  `db:"title" json:"title"`
    PageStart    *int    `db:"page_no_start" json:"pageNoStart,omitempty"`
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
        if opts.DB == nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false}) }
        q := strings.TrimSpace(c.Query("q"))
        rows := []Book{}
        if q == "" {
            if err := opts.DB.Select(&rows, `SELECT id, author_id, publisher_id, title, subtitle, description, isbn, language, categories, cover_image_url, price_amount, price_currency, shareable, allow_hardcopy, visibility FROM books WHERE visibility='public' ORDER BY created_at DESC LIMIT 50`); err != nil {
                return c.Status(500).JSON(fiber.Map{"success": false})
            }
        } else {
            if err := opts.DB.Select(&rows, `SELECT id, author_id, publisher_id, title, subtitle, description, isbn, language, categories, cover_image_url, price_amount, price_currency, shareable, allow_hardcopy, visibility FROM books WHERE visibility='public' AND (LOWER(title) LIKE LOWER($1) OR LOWER(description) LIKE LOWER($1)) ORDER BY created_at DESC LIMIT 50`, "%"+q+"%"); err != nil {
                return c.Status(500).JSON(fiber.Map{"success": false})
            }
        }
        return c.JSON(fiber.Map{"success": true, "data": rows})
    })

    app.Get("/v1/public/meetups", func(c *fiber.Ctx) error {
        if opts.DB == nil { return c.Status(500).JSON(fiber.Map{"success": false}) }
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
        q := `SELECT id, author_id, publisher_id, title, subtitle, description, isbn, language, categories, cover_image_url, price_amount, price_currency, shareable, allow_hardcopy, visibility
              FROM books
              WHERE (author_id IN (SELECT id FROM authors WHERE core_user_id=$1))
                 OR (publisher_id IN (SELECT id FROM publishers WHERE core_user_id=$1))
              ORDER BY created_at DESC
              LIMIT 200`
        if err := opts.DB.Select(&rows, q, uid); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true, "data": rows})
    })

    app.Post("/v1/books", requireAuth, func(c *fiber.Ctx) error {
        var in struct {
            Title string `json:"title"`
            Description *string `json:"description"`
            AuthorID *string `json:"authorId"`
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
        // verify ownership of provided profile
        if in.AuthorID != nil {
            var exists int
            if err := opts.DB.QueryRowx(`SELECT 1 FROM authors WHERE id=$1 AND core_user_id=$2`, *in.AuthorID, uid).Scan(&exists); err != nil || exists != 1 {
                return c.Status(403).JSON(fiber.Map{"success": false, "message": "not your author profile"})
            }
            var id string
            if err := opts.DB.Get(&id, `INSERT INTO books (author_id, title, description, visibility) VALUES ($1,$2,$3,'public') RETURNING id`, *in.AuthorID, in.Title, in.Description); err != nil {
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
            if err := opts.DB.Get(&id, `INSERT INTO books (publisher_id, title, description, visibility) VALUES ($1,$2,$3,'public') RETURNING id`, *in.PublisherID, in.Title, in.Description); err != nil {
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
        if permitted != 1 { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "forbidden"}) }
        var in struct {
            PriceAmount   *float64 `json:"priceAmount"`
            PriceCurrency *string  `json:"priceCurrency"`
            Shareable     *bool    `json:"shareable"`
            AllowHardcopy *bool    `json:"allowHardcopy"`
        }
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"success": false}) }
        sets := []string{}
        args := []any{}
        idx := 1
        if in.PriceAmount != nil { sets = append(sets, "price_amount=$"+itoa(&idx)); args = append(args, *in.PriceAmount) }
        if in.PriceCurrency != nil { sets = append(sets, "price_currency=$"+itoa(&idx)); args = append(args, *in.PriceCurrency) }
        if in.Shareable != nil { sets = append(sets, "shareable=$"+itoa(&idx)); args = append(args, *in.Shareable) }
        if in.AllowHardcopy != nil { sets = append(sets, "allow_hardcopy=$"+itoa(&idx)); args = append(args, *in.AllowHardcopy) }
        if len(sets) == 0 { return c.JSON(fiber.Map{"success": true}) }
        args = append(args, id)
        q := "UPDATE books SET " + strings.Join(sets, ", ") + ", updated_at=now() WHERE id=$" + itoa(&idx)
        if _, err := opts.DB.Exec(q, args...); err != nil { return c.Status(500).JSON(fiber.Map{"success": false}) }
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
            if exists == 1 { allowed = true }
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
        if !canReadBook(opts, uid, bookID) { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false}) }
        type P struct { PageNo int `db:"page_no" json:"pageNo"` }
        rows := []P{}
        if err := opts.DB.Select(&rows, `SELECT page_no FROM book_pages WHERE book_id=$1 ORDER BY page_no ASC`, bookID); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        total := 0
        if err := opts.DB.Get(&total, `SELECT COUNT(*) FROM book_pages WHERE book_id=$1`, bookID); err != nil { total = len(rows) }
        return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"pages": rows, "totalPages": total}})
    })

    // On-platform reading: get page content
    app.Get("/v1/books/:id/pages/:no", requireAuth, func(c *fiber.Ctx) error {
        bookID := c.Params("id")
        uid := auth.UserID(c)
        if !canReadBook(opts, uid, bookID) { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false}) }
        no := strings.TrimSpace(c.Params("no"))
        var pageNo int
        if v, err := strconv.Atoi(no); err == nil && v > 0 { pageNo = v } else { return c.Status(400).JSON(fiber.Map{"success": false}) }
        var html sql.NullString
        if err := opts.DB.QueryRowx(`SELECT content_html FROM book_pages WHERE book_id=$1 AND page_no=$2`, bookID, pageNo).Scan(&html); err != nil {
            return c.Status(404).JSON(fiber.Map{"success": false, "message": "page not found"})
        }
        total := 0
        _ = opts.DB.Get(&total, `SELECT COUNT(*) FROM book_pages WHERE book_id=$1`, bookID)
        return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"pageNo": pageNo, "totalPages": total, "html": html.String}})
    })

    // Authoring: upsert a page (owner only)
    app.Post("/v1/books/:id/pages", requireAuth, func(c *fiber.Ctx) error {
        bookID := c.Params("id")
        uid := auth.UserID(c)
        // verify ownership by author
        var owner string
        _ = opts.DB.Get(&owner, `SELECT COALESCE(author_id,'') FROM books WHERE id=$1`, bookID)
        if owner != uid { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false}) }
        var in struct { PageNo int `json:"pageNo"`; Html string `json:"html"` }
        if err := c.BodyParser(&in); err != nil || in.PageNo <= 0 { return c.Status(400).JSON(fiber.Map{"success": false}) }
        // upsert
        if _, err := opts.DB.Exec(`INSERT INTO book_pages (book_id, page_no, content_html) VALUES ($1,$2,$3)
            ON CONFLICT (book_id, page_no) DO UPDATE SET content_html=EXCLUDED.content_html, updated_at=now()`, bookID, in.PageNo, in.Html); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true})
    })

    // Chapters: list (auth + can read)
    app.Get("/v1/books/:id/chapters", requireAuth, func(c *fiber.Ctx) error {
        bookID := c.Params("id")
        uid := auth.UserID(c)
        if !canReadBook(opts, uid, bookID) { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false}) }
        rows := []Chapter{}
        if err := opts.DB.Select(&rows, `SELECT id, book_id, number, title, page_no_start FROM book_chapters WHERE book_id=$1 ORDER BY number ASC`, bookID); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true, "data": rows})
    })
    // Chapters: upsert (owner only)
    app.Post("/v1/books/:id/chapters", requireAuth, func(c *fiber.Ctx) error {
        bookID := c.Params("id")
        uid := auth.UserID(c)
        if !userOwnsBook(opts, uid, bookID) { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false}) }
        var in struct { Number int `json:"number"`; Title string `json:"title"`; PageNoStart *int `json:"pageNoStart"` }
        if err := c.BodyParser(&in); err != nil || in.Number <= 0 || strings.TrimSpace(in.Title) == "" { return c.Status(400).JSON(fiber.Map{"success": false}) }
        if _, err := opts.DB.Exec(`INSERT INTO book_chapters (book_id, number, title, page_no_start) VALUES ($1,$2,$3,$4)
            ON CONFLICT (book_id, number) DO UPDATE SET title=EXCLUDED.title, page_no_start=EXCLUDED.page_no_start, updated_at=now()`, bookID, in.Number, in.Title, in.PageNoStart); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true})
    })

    // Profiles: authors/publishers for current user
    app.Get("/v1/me/authors", requireAuth, func(c *fiber.Ctx) error {
        uid := auth.UserID(c)
        type A struct { ID string `db:"id" json:"id"`; Name string `db:"display_name" json:"displayName"` }
        var rows []A
        if err := opts.DB.Select(&rows, `SELECT id, display_name FROM authors WHERE core_user_id=$1 ORDER BY display_name`, uid); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true, "data": rows})
    })
    app.Post("/v1/me/authors", requireAuth, func(c *fiber.Ctx) error {
        uid := auth.UserID(c)
        var in struct{ DisplayName string `json:"displayName"`; Bio *string `json:"bio"` }
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
        type P struct { ID string `db:"id" json:"id"`; Name string `db:"name" json:"name"` }
        var rows []P
        if err := opts.DB.Select(&rows, `SELECT id, name FROM publishers WHERE core_user_id=$1 ORDER BY name`, uid); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true, "data": rows})
    })
    app.Post("/v1/me/publishers", requireAuth, func(c *fiber.Ctx) error {
        uid := auth.UserID(c)
        var in struct{ Name string `json:"name"`; Description *string `json:"description"` }
        if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Name) == "" { return c.Status(400).JSON(fiber.Map{"success": false}) }
        var id string
        if err := opts.DB.QueryRowx(`INSERT INTO publishers (core_user_id, name, description) VALUES ($1,$2,$3) RETURNING id`, uid, in.Name, in.Description).Scan(&id); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"id": id}})
    })

    app.Post("/v1/purchases", requireAuth, func(c *fiber.Ctx) error {
        var in struct { BookID string `json:"bookId"`; Kind string `json:"kind"`; Provider string `json:"provider"` }
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"success": false}) }
        if in.Kind != "ebook" && in.Kind != "hardcopy" { return c.Status(400).JSON(fiber.Map{"success": false, "message": "kind must be ebook or hardcopy"}) }
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
        if priceAmount != nil && priceCurrency != nil {
            if *priceAmount > 0 {
                amtCents = int64((*priceAmount) * 100)
                curr = *priceCurrency
            }
        }
        if amtCents > 0 && curr != "" {
            prov := in.Provider
            if prov == "" { prov = "mpesa" }
            payload := billing.CreateIntentRequest{AmountCents: amtCents, Currency: curr, Description: "Purchase book " + in.BookID, Provider: prov}
            // Build proxy request
            r := new(http.Request); r.Header = http.Header{}
            if v := c.Get("Authorization"); v != "" { r.Header.Set("Authorization", v) }
            if v := c.Get("Cookie"); v != "" { r.Header.Set("Cookie", v) }
            if v := c.Get("Origin"); v != "" { r.Header.Set("Origin", v) }
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
        if intentID == nil { return c.JSON(fiber.Map{"success": true, "message": "no intent"}) }
        r := new(http.Request); r.Header = http.Header{}
        if v := c.Get("Authorization"); v != "" { r.Header.Set("Authorization", v) }
        if v := c.Get("Cookie"); v != "" { r.Header.Set("Cookie", v) }
        if v := c.Get("Origin"); v != "" { r.Header.Set("Origin", v) }
        resp, status, err := billing.GetPaymentIntent(opts.CoreAPIBase, r, *intentID)
        if err != nil || status < 200 || status >= 300 || resp == nil || resp.Data == nil {
            return c.Status(502).JSON(fiber.Map{"success": false, "message": "billing unavailable"})
        }
        var curStatus string
        if m, ok := resp.Data.(map[string]any); ok {
            if s, ok := m["status"].(string); ok { curStatus = s }
            // Capture paid amount and currency if needed in future
        }
        if strings.ToLower(curStatus) == "succeeded" || strings.ToLower(curStatus) == "paid" {
            _, _ = opts.DB.Exec(`UPDATE purchases SET status='completed' WHERE id=$1`, pid)
            return c.JSON(fiber.Map{"success": true, "message": "completed"})
        }
        return c.JSON(fiber.Map{"success": true, "message": curStatus})
    })

    // List my purchases
    app.Get("/v1/me/purchases", requireAuth, func(c *fiber.Ctx) error {
        uid := auth.UserID(c)
        type item struct {
            ID string `db:"id" json:"id"`
            BookID string `db:"book_id" json:"bookId"`
            Kind string `db:"kind" json:"kind"`
            Status string `db:"status" json:"status"`
            Currency *string `db:"currency" json:"currency,omitempty"`
            PricePaid *float64 `db:"price_paid" json:"pricePaid,omitempty"`
            CreatedAt time.Time `db:"created_at" json:"createdAt"`
            Title string `db:"title" json:"title"`
        }
        rows := []item{}
        if err := opts.DB.Select(&rows, `SELECT p.id, p.book_id, p.kind, p.status, p.currency, p.price_paid, p.created_at, COALESCE(b.title,'') AS title FROM purchases p LEFT JOIN books b ON b.id=p.book_id WHERE p.user_id=$1 ORDER BY p.created_at DESC LIMIT 200`, uid); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true, "data": rows})
    })

    app.Post("/v1/books/:id/share-links", requireAuth, func(c *fiber.Ctx) error {
        bookID := c.Params("id")
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

    app.Post("/v1/clubs", requireAuth, func(c *fiber.Ctx) error {
        var in struct { Name string `json:"name"`; Description *string `json:"description"` }
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
        var in struct { Name string `json:"name"` }
        if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Name) == "" { return c.Status(400).JSON(fiber.Map{"success": false}) }
        var roomID string
        if err := opts.DB.Get(&roomID, `INSERT INTO club_rooms (club_id, name) VALUES ($1,$2) RETURNING id`, clubID, in.Name); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true, "id": roomID})
    })

    app.Get("/v1/clubs/:id/rooms", requireAuth, func(c *fiber.Ctx) error {
        clubID := c.Params("id")
        rows := []Room{}
        if err := opts.DB.Select(&rows, `SELECT id, club_id, name FROM club_rooms WHERE club_id=$1 ORDER BY name`, clubID); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true, "data": rows})
    })

    app.Get("/v1/rooms/:id/messages", requireAuth, func(c *fiber.Ctx) error {
        roomID := c.Params("id")
        rows := []Message{}
        if err := opts.DB.Select(&rows, `SELECT id, room_id, user_id, content, created_at FROM room_messages WHERE room_id=$1 ORDER BY created_at DESC LIMIT 100`, roomID); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        // reverse to chronological
        for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 { rows[i], rows[j] = rows[j], rows[i] }
        return c.JSON(fiber.Map{"success": true, "data": rows})
    })

    app.Post("/v1/rooms/:id/messages", requireAuth, func(c *fiber.Ctx) error {
        roomID := c.Params("id")
        var in struct { Content string `json:"content"` }
        if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Content) == "" { return c.Status(400).JSON(fiber.Map{"success": false}) }
        uid := auth.UserID(c)
        if _, err := opts.DB.Exec(`INSERT INTO room_messages (room_id, user_id, content) VALUES ($1,$2,$3)`, roomID, uid, in.Content); err != nil {
            return c.Status(500).JSON(fiber.Map{"success": false})
        }
        return c.JSON(fiber.Map{"success": true})
    })

    app.Post("/v1/meetups", requireAuth, func(c *fiber.Ctx) error {
        uid := auth.UserID(c)
        var in struct {
            Title string `json:"title"`
            Description *string `json:"description"`
            Location *string `json:"location"`
            Start time.Time `json:"startTime"`
            End *time.Time `json:"endTime"`
            IsPaid bool `json:"isPaid"`
            Price *float64 `json:"priceAmount"`
            Currency *string `json:"priceCurrency"`
            Capacity *int `json:"capacity"`
            ClubID *string `json:"clubId"`
            BookID *string `json:"bookId"`
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
            r := new(http.Request); r.Header = http.Header{}
            if v := c.Get("Authorization"); v != "" { r.Header.Set("Authorization", v) }
            if v := c.Get("Cookie"); v != "" { r.Header.Set("Cookie", v) }
            if v := c.Get("Origin"); v != "" { r.Header.Set("Origin", v) }
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
    if uid == "" || opts.DB == nil { return false }
    if userOwnsBook(opts, uid, bookID) { return true }
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
