package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ThanhNV121097/project-0e8b79be/backend/internal/db"
)

type todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type app struct {
	db      *sql.DB
	limiter *limiter
}

type apiError struct {
	Error errBody `json:"error"`
}

type errBody struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   []errDetail `json:"details"`
	RequestID string      `json:"request_id"`
}

type errDetail struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type contextKey string

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

const requestIDKey contextKey = "request_id"

var uuidRE = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	database, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if err := database.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}
	if err := db.ApplyMigrations(ctx, database); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	a := &app{db: database, limiter: newLimiter()}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", a.health)
	mux.HandleFunc("/api/v1/todos", a.todos)
	mux.HandleFunc("/api/v1/todos/", a.todoByID)

	server := &http.Server{
		Addr:              ":" + listenPort(),
		Handler:           a.withRequestLog(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v", err)
	}
}

func (a *app) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := a.db.PingContext(ctx); err != nil {
		writeErr(w, r, http.StatusServiceUnavailable, "UNAVAILABLE", "Service is unavailable.", nil)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *app) todos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/todos" {
		writeErr(w, r, http.StatusNotFound, "NOT_FOUND", "Todo was not found.", nil)
		return
	}
	switch r.Method {
	case http.MethodGet:
		a.list(w, r)
	case http.MethodPost:
		a.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *app) todoByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
	if strings.Contains(id, "/") || !uuidRE.MatchString(id) {
		writeErr(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request is invalid.", nil)
		return
	}
	switch r.Method {
	case http.MethodPatch:
		a.patch(w, r, id)
	case http.MethodDelete:
		a.delete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *app) list(w http.ResponseWriter, r *http.Request) {
	if !a.allowRead(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	rows, err := a.db.QueryContext(ctx, `SELECT id::text,title,CASE WHEN is_completed THEN 'completed' ELSE 'active' END,created_at,updated_at FROM todos ORDER BY created_at ASC,id ASC`)
	if err != nil {
		dbErr(w, r, err)
		return
	}
	defer rows.Close()

	data := []todo{}
	active := 0
	completed := 0
	for rows.Next() {
		t, err := scanTodo(rows)
		if err != nil {
			dbErr(w, r, err)
			return
		}
		if t.Status == "completed" {
			completed++
		} else {
			active++
		}
		data = append(data, t)
	}
	if err := rows.Err(); err != nil {
		dbErr(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": data, "meta": map[string]int{"total": len(data), "active": active, "completed": completed}})
}

func (a *app) create(w http.ResponseWriter, r *http.Request) {
	if !a.allowWriteJSON(w, r) {
		return
	}
	var payload struct {
		Title any `json:"title"`
	}
	if !decode(w, r, &payload) {
		return
	}
	titleValue, ok := payload.Title.(string)
	if !ok {
		writeErr(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request is invalid.", nil)
		return
	}
	title := strings.TrimSpace(titleValue)
	if title == "" {
		writeErr(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Todo title is invalid.", []errDetail{{Field: "title", Code: "BLANK", Message: "Title cannot be blank."}})
		return
	}
	if len([]rune(title)) > 200 {
		writeErr(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Todo title is invalid.", []errDetail{{Field: "title", Code: "TOO_LONG", Message: "Title must be 200 characters or fewer."}})
		return
	}

	t, err := a.queryOne(r.Context(), `INSERT INTO todos (title) VALUES ($1) RETURNING id::text,title,CASE WHEN is_completed THEN 'completed' ELSE 'active' END,created_at,updated_at`, title)
	if err != nil {
		dbErr(w, r, err)
		return
	}
	w.Header().Set("Location", "/api/v1/todos/"+t.ID)
	writeJSON(w, http.StatusCreated, t)
}
func (a *app) patch(w http.ResponseWriter, r *http.Request, id string) {
	if !a.allowWriteJSON(w, r) {
		return
	}
	var payload struct {
		Status any `json:"status"`
	}
	if !decode(w, r, &payload) {
		return
	}
	statusValue, ok := payload.Status.(string)
	if !ok {
		writeErr(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request is invalid.", nil)
		return
	}
	if statusValue != "active" && statusValue != "completed" {
		writeErr(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Todo status is invalid.", []errDetail{{Field: "status", Code: "INVALID", Message: "Status must be active or completed."}})
		return
	}

	t, err := a.queryOne(r.Context(), `UPDATE todos SET is_completed=$2,updated_at=now() WHERE id=$1 RETURNING id::text,title,CASE WHEN is_completed THEN 'completed' ELSE 'active' END,created_at,updated_at`, id, statusValue == "completed")
	if errors.Is(err, sql.ErrNoRows) {
		writeErr(w, r, http.StatusNotFound, "NOT_FOUND", "Todo was not found.", nil)
		return
	}
	if err != nil {
		dbErr(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (a *app) delete(w http.ResponseWriter, r *http.Request, id string) {
	if !a.allowWrite(w, r) {
		return
	}
	if r.TransferEncoding != nil || r.ContentLength > 0 {
		writeErr(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request is invalid.", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	res, err := a.db.ExecContext(ctx, `DELETE FROM todos WHERE id=$1`, id)
	if err != nil {
		dbErr(w, r, err)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		writeErr(w, r, http.StatusNotFound, "NOT_FOUND", "Todo was not found.", nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *app) queryOne(ctx context.Context, query string, args ...any) (todo, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return scanTodo(a.db.QueryRowContext(queryCtx, query, args...))
}

func (a *app) allowRead(w http.ResponseWriter, r *http.Request) bool {
	return a.allow(w, r, 120)
}

func (a *app) allowWrite(w http.ResponseWriter, r *http.Request) bool {
	return a.allow(w, r, 60)
}

func (a *app) allowWriteJSON(w http.ResponseWriter, r *http.Request) bool {
	if !a.allowWrite(w, r) {
		return false
	}
	contentType := r.Header.Get("Content-Type")
	if contentType == "" || !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		writeErr(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request is invalid.", nil)
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	return true
}

func (a *app) allow(w http.ResponseWriter, r *http.Request, maxHits int) bool {
	if !a.limiter.allow(clientSource(r), maxHits) {
		w.Header().Set("Retry-After", "60")
		writeErr(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests. Please try again shortly.", nil)
		return false
	}
	return true
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(v); err != nil {
		writeErr(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request is invalid.", nil)
		return false
	}
	if dec.Decode(&struct{}{}) != io.EOF {
		writeErr(w, r, http.StatusBadRequest, "BAD_REQUEST", "Request is invalid.", nil)
		return false
	}
	return true
}

func scanTodo(row interface{ Scan(...any) error }) (todo, error) {
	var t todo
	var createdAt time.Time
	var updatedAt time.Time
	err := row.Scan(&t.ID, &t.Title, &t.Status, &createdAt, &updatedAt)
	if err == nil {
		t.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		t.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
	}
	return t, err
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, r *http.Request, status int, code string, message string, details []errDetail) {
	requestID := getRequestID(r)
	w.Header().Set("X-Request-Id", requestID)
	writeJSON(w, status, apiError{Error: errBody{Code: code, Message: message, Details: details, RequestID: requestID}})
}

func dbErr(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		writeErr(w, r, http.StatusServiceUnavailable, "UNAVAILABLE", "Service is unavailable.", nil)
		return
	}
	writeErr(w, r, http.StatusInternalServerError, "INTERNAL", "Something went wrong.", nil)
}

func getRequestID(r *http.Request) string {
	if requestID, ok := r.Context().Value(requestIDKey).(string); ok && requestID != "" {
		return requestID
	}
	return newRequestID()
}

func newRequestID() string {
	return hex.EncodeToString([]byte(time.Now().UTC().Format("20060102150405.000000000")))
}

func clientSource(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (a *app) withRequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = newRequestID()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-Id", requestID)

		wrapped := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(wrapped, r)
		sourceHash := sha256.Sum256([]byte(clientSource(r)))
		log.Printf("request_id=%s timestamp=%s method=%s path=%s status=%d duration_ms=%d response_bytes=%d client_source_hash=%s", requestID, start.UTC().Format(time.RFC3339), r.Method, r.URL.Path, wrapped.status, time.Since(start).Milliseconds(), wrapped.bytes, hex.EncodeToString(sourceHash[:8]))
	})
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(body)
	w.bytes += n
	return n, err
}

type limiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time
}

func newLimiter() *limiter {
	return &limiter{hits: map[string][]time.Time{}}
}

func (l *limiter) allow(key string, maxHits int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-time.Minute)
	hits := l.hits[key]
	kept := 0
	for _, hit := range hits {
		if hit.After(cutoff) {
			hits[kept] = hit
			kept++
		}
	}
	hits = hits[:kept]
	if len(hits) >= maxHits {
		l.hits[key] = hits
		return false
	}
	l.hits[key] = append(hits, now)
	return true
}

func listenPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		return port
	}
	return "8080"
}
