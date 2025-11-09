package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	coreauth "github.com/berjistech/berjis-ecosystem/shared/coreauth"
)

type Options struct {
	HS256Secret string
	Env         string
	CoreAPIBase string
	HTTPClient  *http.Client
	Verifier    *coreauth.Verifier
}

const userKey = "userID"

func Middleware(opts Options) fiber.Handler {
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 3 * time.Second}
	}
	return func(c *fiber.Ctx) error {
		authz := c.Get("Authorization")
		var uid string
		tokenStr := ""
		if strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			tokenStr = strings.TrimSpace(authz[7:])
		}
		if tokenStr != "" && strings.TrimSpace(opts.HS256Secret) != "" {
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(opts.HS256Secret), nil
			})
			if err != nil || !token.Valid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "invalid token"})
			}
			if sub, ok := claims["sub"].(string); ok && sub != "" {
				uid = sub
			}
		}
		if uid == "" && opts.Verifier != nil {
			tokenCandidate := tokenStr
			if tokenCandidate == "" {
				tokenCandidate = strings.TrimSpace(c.Cookies("access", ""))
			}
			if tokenCandidate != "" {
				if claims, err := opts.Verifier.Verify(tokenCandidate); err == nil {
					uid = claims.UUID
				} else if errors.Is(err, coreauth.ErrTokenInvalid) || errors.Is(err, coreauth.ErrTokenExpired) || errors.Is(err, coreauth.ErrTokenMissing) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "invalid token"})
				}
			}
		}
		if uid == "" && strings.ToLower(opts.Env) == "development" {
			uid = c.Get("X-User-UUID")
		}
		// If still empty, verify via Core API using cookies/headers
		if uid == "" && strings.TrimSpace(opts.CoreAPIBase) != "" {
			req, _ := http.NewRequest("GET", strings.TrimRight(opts.CoreAPIBase, "/")+"/v1/auth/verify", nil)
			if v := c.Get("Authorization"); v != "" {
				req.Header.Set("Authorization", v)
			}
			if v := c.Get("Cookie"); v != "" {
				req.Header.Set("Cookie", v)
			}
			if resp, err := client.Do(req); err == nil && resp != nil {
				defer resp.Body.Close()
				var raw map[string]any
				if err := json.NewDecoder(resp.Body).Decode(&raw); err == nil {
					if data, _ := raw["data"].(map[string]any); data != nil {
						if valid, ok := data["valid"].(bool); ok && valid {
							if us, ok := data["uid"].(string); ok && us != "" {
								uid = us
							}
							if uid == "" {
								if us, ok := data["uuid"].(string); ok && us != "" {
									uid = us
								}
							}
							if uid == "" {
								if us2, ok := data["userId"].(string); ok && us2 != "" {
									uid = us2
								}
							}
						}
					}
				}
			}
		}
		if uid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "login required"})
		}
		c.Locals(userKey, uid)
		return c.Next()
	}
}

func UserID(c *fiber.Ctx) string {
	if v := c.Locals(userKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
