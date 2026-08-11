package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ThanhNV121097/project-0e8b79be/backend/internal/db"
)

const (
	maxBodyBytes  int64 = 16 * 1024
	readLimit           = 120
	writeLimit          = 60
	rateWindow          = time.Minute
)

type api struct {
	database *sql.DB
	limiter  *rateLimiter
}

type todoDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type todoListResponse struct {
	Data []todoDTO `json:"data"`
	Meta struct {
		Total     int `json:"total"`
		Active    int `json:"active"`
		Completed int `json:"completed"`
	} `json:"meta"`
}

type fieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error struct {
		Code      string       `json:"code"`
		Message   string       `json:"message"`
		Details   []fieldError `json:"details"`
		RequestID string       `json:"request_id"`
	} `json:"error"`
}

type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]rateBucket
	now     func() time.Time
}

type rateBucket struct {
	windowStart time.Time
	reads       int
	writes      int
}

type requestIDKey struct{}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

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

	a := &api{database: database, limiter: newRateLimiter()}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", a.health)
	mux.HandleFunc("/api/v1/todos", a.todos)
	mux.HandleFunc("/api/v1/todos/", a.todoByID)

	server := &http.Server{
		Addr:              ":" + listenPort(),
		Handler:           withRequestID(withRequestLogging(mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v", err)
	}
}

func (a *api) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := a.database.PingContext(ctx); err != nil {
		writeError(w, r, http.StatusServiceUnavailable, "UNAVAILABLE", "The service is temporarily unavailable.", nil)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *api) todos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !a.allow(w, r, true) {
			return
		}
		a.listTodos(w, r)
	case http.MethodPost:
		if !a.allow(w, r, false) {
			return
		}
		a.createTodo(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *api) todoByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
	if id == "" || strings.Contains(id, "/") || !isUUID(id) {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The todo identifier is invalid.", nil)
		return
	}

	switch r.Method {
	case http.MethodPatch:
		if !a.allow(w, r, false) {
			return
		}
		a.patchTodo(w, r, id)
	case http.MethodDelete:
		if !a.allow(w, r, false) {
			return
		}
		a.deleteTodo(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *api) listTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := a.database.QueryContext(ctx, `SELECT id::text, title, is_completed, created_at, updated_at FROM todos ORDER BY created_at ASC, id ASC`)
	if err != nil {
		writeDBError(w, r, err)
		return
	}
	defer rows.Close()

	resp := todoListResponse{Data: []todoDTO{}}
	for rows.Next() {
		var id, title string
		var done bool
		var created, updated time.Time
		if err := rows.Scan(&id, &title, &done, &created, &updated); err != nil {
			writeDBError(w, r, err)
			return
		}
		resp.Data = append(resp.Data, mapTodo(id, title, done, created, updated))
		if done {
			resp.Meta.Completed++
		} else {
			resp.Meta.Active++
		}
	}
	if err := rows.Err(); err != nil {
		writeDBError(w, r, err)
		return
	}
	resp.Meta.Total = len(resp.Data)
	writeJSON(w, http.StatusOK, resp)
}

func (a *api) createTodo(w http.ResponseWriter, r *http.Request) {
	if !jsonContent(r) {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Please send JSON for this request.", nil)
		return
	}

	var body struct {
		Title any `json:"title"`
	}
	if !decodeBody(w, r, &body) {
		return
	}

	title, ok := body.Title.(string)
	if !ok {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The todo title must be text.", nil)
		return
	}

	title = strings.TrimSpace(title)
	if title == "" {
		writeError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Please enter a task before saving.", []fieldError{{Field: "title", Code: "BLANK", Message: "Title is required."}})
		return
	}
	if len([]rune(title)) > 200 {
		writeError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Title must be 200 characters or fewer.", []fieldError{{Field: "title", Code: "TOO_LONG", Message: "Title must be 200 characters or fewer."}})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var id string
	var done bool
	var created, updated time.Time
	err := a.database.QueryRowContext(ctx, `INSERT INTO todos (title) VALUES ($1) RETURNING id::text, is_completed, created_at, updated_at`, title).Scan(&id, &done, &created, &updated)
	if err != nil {
		writeDBError(w, r, err)
		return
	}

	w.Header().Set("Location", "/api/v1/todos/"+id)
	writeJSON(w, http.StatusCreated, mapTodo(id, title, done, created, updated))
}

func (a *api) patchTodo(w http.ResponseWriter, r *http.Request, id string) {
	if !jsonContent(r) {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Please send JSON for this request.", nil)
		return
	}

	var body struct {
		Status any `json:"status"`
	}
	if !decodeBody(w, r, &body) {
		return
	}

	status, ok := body.Status.(string)
	if !ok {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The status must be active or completed.", nil)
		return
	}
	if status != "active" && status != "completed" {
		writeError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "The status must be active or completed.", []fieldError{{Field: "status", Code: "INVALID", Message: "Status must be active or completed."}})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var outID, title string
	var done bool
	var created, updated time.Time
	err := a.database.QueryRowContext(ctx, `UPDATE todos SET is_completed = $2, updated_at = now() WHERE id = $1 RETURNING id::text, title, is_completed, created_at, updated_at`, id, status == "completed").Scan(&outID, &title, &done, &created, &updated)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, r, http.StatusNotFound, "NOT_FOUND", "That task no longer exists.", nil)
		return
	}
	if err != nil {
		writeDBError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, mapTodo(outID, title, done, created, updated))
}

func (a *api) deleteTodo(w http.ResponseWriter, r *http.Request, id string) {
	if hasRequestBody(r) {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Delete requests must not include a body.", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := a.database.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		writeDBError(w, r, err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeError(w, r, http.StatusNotFound, "NOT_FOUND", "That task no longer exists.", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *api) allow(w http.ResponseWriter, r *http.Request, read bool) bool {
	if a.limiter.allow(clientSource(r), read) {
		return true
	}
	w.Header().Set("Retry-After", "60")
	writeError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests. Please wait a moment and try again.", nil)
	return false
}

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The request body is not valid JSON.", nil)
		return false
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The request body must contain one JSON object.", nil)
		return false
	}

	return true
}

func hasRequestBody(r *http.Request) bool {
	if r.Body == nil || r.Body == http.NoBody {
		return false
	}
	if r.ContentLength > 0 || r.ContentLength == -1 {
		return true
	}
	return false
}

func jsonContent(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	return mediaType == "application/json"
}

func mapTodo(id, title string, done bool, created, updated time.Time) todoDTO {
	status := "active"
	if done {
		status = "completed"
	}
	return todoDTO{
		ID:        id,
		Title:     title,
		Status:    status,
		CreatedAt: created.UTC().Format(time.RFC3339),
		UpdatedAt: updated.UTC().Format(time.RFC3339),
	}
}

func writeDBError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) || isUnavailableError(err) {
		writeError(w, r, http.StatusServiceUnavailable, "UNAVAILABLE", "The service is temporarily unavailable.", nil)
		return
	}

	log.Printf("request_id=%s error=%v", requestID(r), err)
	writeError(w, r, http.StatusInternalServerError, "INTERNAL", "Something went wrong. Please try again.", nil)
}

func isUnavailableError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "connection refused") || strings.Contains(message, "connection reset") || strings.Contains(message, "broken pipe") || strings.Contains(message, "server closed") || strings.Contains(message, "failed to connect")
}

func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string, details []fieldError) {
	resp := errorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	resp.Error.Details = details
	resp.Error.RequestID = requestID(r)
	writeJSON(w, status, resp)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if status != http.StatusNoContent {
		_ = json.NewEncoder(w).Encode(value)
	}
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id)))
	})
}

func withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		log.Printf("request_id=%s timestamp=%s method=%s path_template=%s status=%d duration_ms=%d response_bytes=%d client_source_hash=%s", requestID(r), started.UTC().Format(time.RFC3339), r.Method, pathTemplate(r), recorder.status, time.Since(started).Milliseconds(), recorder.bytes, hashString(clientSource(r)))
	})
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	n, err := r.ResponseWriter.Write(body)
	r.bytes += n
	return n, err
}

func requestID(r *http.Request) string {
	if v, ok := r.Context().Value(requestIDKey{}).(string); ok {
		return v
	}
	return "unknown"
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{buckets: map[string]rateBucket{}, now: time.Now}
}

func (l *rateLimiter) allow(source string, read bool) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	bucket := l.buckets[source]
	if bucket.windowStart.IsZero() || now.Sub(bucket.windowStart) >= rateWindow {
		bucket = rateBucket{windowStart: now}
	}

	allowed := false
	if read {
		bucket.reads++
		allowed = bucket.reads <= readLimit
	} else {
		bucket.writes++
		allowed = bucket.writes <= writeLimit
	}
	l.buckets[source] = bucket

	for key, old := range l.buckets {
		if now.Sub(old.windowStart) > 2*rateWindow {
			delete(l.buckets, key)
		}
	}

	return allowed
}

func clientSource(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || host == "" {
		return r.RemoteAddr
	}
	return host
}

func hashString(value string) string {
	b := []byte(value)
	var hash uint32 = 2166136261
	for _, c := range b {
		hash ^= uint32(c)
		hash *= 16777619
	}
	return fmt.Sprintf("%08x", hash)
}

func pathTemplate(r *http.Request) string {
	if strings.HasPrefix(r.URL.Path, "/api/v1/todos/") {
		return "/api/v1/todos/{todo_id}"
	}
	return r.URL.Path
}

func isUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
			continue
		}
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
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
