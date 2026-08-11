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
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ThanhNV121097/project-0e8b79be/backend/internal/db"
)

const maxBodyBytes int64 = 16 * 1024

type api struct{ database *sql.DB }

type todoDTO struct { ID string `json:"id"`; Title string `json:"title"`; Status string `json:"status"`; CreatedAt string `json:"createdAt"`; UpdatedAt string `json:"updatedAt"` }
type todoListResponse struct { Data []todoDTO `json:"data"`; Meta struct { Total int `json:"total"`; Active int `json:"active"`; Completed int `json:"completed"` } `json:"meta"` }
type fieldError struct { Field string `json:"field"`; Code string `json:"code"`; Message string `json:"message"` }
type errorResponse struct { Error struct { Code string `json:"code"`; Message string `json:"message"`; Details []fieldError `json:"details"`; RequestID string `json:"request_id"` } `json:"error"` }

func main() {
	databaseURL := os.Getenv("DATABASE_URL"); if databaseURL == "" { log.Fatal("DATABASE_URL is required") }
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second); defer cancel()
	database, err := sql.Open("pgx", databaseURL); if err != nil { log.Fatalf("open database: %v", err) }
	defer database.Close()
	if err := database.PingContext(ctx); err != nil { log.Fatalf("ping database: %v", err) }
	if err := db.ApplyMigrations(ctx, database); err != nil { log.Fatalf("apply migrations: %v", err) }
	a := &api{database: database}; mux := http.NewServeMux()
	mux.HandleFunc("/healthz", a.health); mux.HandleFunc("/api/v1/todos", a.todos); mux.HandleFunc("/api/v1/todos/", a.todoByID)
	server := &http.Server{Addr: ":" + listenPort(), Handler: withRequestID(mux), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("backend listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) { log.Fatalf("listen: %v", err) }
}

func (a *api) health(w http.ResponseWriter, r *http.Request) { if r.Method != http.MethodGet { w.WriteHeader(http.StatusMethodNotAllowed); return }; ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second); defer cancel(); if err := a.database.PingContext(ctx); err != nil { writeError(w, r, http.StatusServiceUnavailable, "UNAVAILABLE", "The service is temporarily unavailable.", nil); return }; writeJSON(w, http.StatusOK, map[string]string{"status":"ok"}) }
func (a *api) todos(w http.ResponseWriter, r *http.Request) { switch r.Method { case http.MethodGet: a.listTodos(w, r); case http.MethodPost: a.createTodo(w, r); default: w.WriteHeader(http.StatusMethodNotAllowed) } }
func (a *api) todoByID(w http.ResponseWriter, r *http.Request) { id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/"); if id == "" || strings.Contains(id, "/") || !isUUID(id) { writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The todo identifier is invalid.", nil); return }; switch r.Method { case http.MethodPatch: a.patchTodo(w, r, id); case http.MethodDelete: a.deleteTodo(w, r, id); default: w.WriteHeader(http.StatusMethodNotAllowed) } }

func (a *api) listTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second); defer cancel()
	rows, err := a.database.QueryContext(ctx, `SELECT id::text, title, is_completed, created_at, updated_at FROM todos ORDER BY created_at ASC, id ASC`); if err != nil { writeDBError(w, r, err); return }
	defer rows.Close(); resp := todoListResponse{Data: []todoDTO{}}
	for rows.Next() { var id, title string; var done bool; var created, updated time.Time; if err := rows.Scan(&id, &title, &done, &created, &updated); err != nil { writeDBError(w, r, err); return }; resp.Data = append(resp.Data, mapTodo(id, title, done, created, updated)); if done { resp.Meta.Completed++ } else { resp.Meta.Active++ } }
	if err := rows.Err(); err != nil { writeDBError(w, r, err); return }; resp.Meta.Total = len(resp.Data); writeJSON(w, http.StatusOK, resp)
}

func (a *api) createTodo(w http.ResponseWriter, r *http.Request) {
	if !jsonContent(r) { writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Please send JSON for this request.", nil); return }
	var body struct{ Title any `json:"title"` }; if !decodeBody(w, r, &body) { return }
	title, ok := body.Title.(string); if !ok { writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The todo title must be text.", nil); return }
	title = strings.TrimSpace(title); if title == "" { writeError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Please enter a task before saving.", []fieldError{{"title","BLANK","Title is required."}}); return }; if len([]rune(title)) > 200 { writeError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Title must be 200 characters or fewer.", []fieldError{{"title","TOO_LONG","Title must be 200 characters or fewer."}}); return }
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second); defer cancel(); var id string; var done bool; var created, updated time.Time
	err := a.database.QueryRowContext(ctx, `INSERT INTO todos (title) VALUES ($1) RETURNING id::text, is_completed, created_at, updated_at`, title).Scan(&id, &done, &created, &updated); if err != nil { writeDBError(w, r, err); return }
	w.Header().Set("Location", "/api/v1/todos/"+id); writeJSON(w, http.StatusCreated, mapTodo(id, title, done, created, updated))
}

func (a *api) patchTodo(w http.ResponseWriter, r *http.Request, id string) {
	if !jsonContent(r) { writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Please send JSON for this request.", nil); return }
	var body struct{ Status any `json:"status"` }; if !decodeBody(w, r, &body) { return }
	status, ok := body.Status.(string); if !ok { writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The status must be active or completed.", nil); return }
	if status != "active" && status != "completed" { writeError(w, r, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "The status must be active or completed.", []fieldError{{"status","INVALID","Status must be active or completed."}}); return }
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second); defer cancel(); var outID, title string; var done bool; var created, updated time.Time
	err := a.database.QueryRowContext(ctx, `UPDATE todos SET is_completed = $2, updated_at = now() WHERE id = $1 RETURNING id::text, title, is_completed, created_at, updated_at`, id, status == "completed").Scan(&outID, &title, &done, &created, &updated)
	if errors.Is(err, sql.ErrNoRows) { writeError(w, r, http.StatusNotFound, "NOT_FOUND", "That task no longer exists.", nil); return }; if err != nil { writeDBError(w, r, err); return }
	writeJSON(w, http.StatusOK, mapTodo(outID, title, done, created, updated))
}

func (a *api) deleteTodo(w http.ResponseWriter, r *http.Request, id string) {
	if r.ContentLength > 0 { writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Delete requests must not include a body.", nil); return }
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second); defer cancel(); res, err := a.database.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id); if err != nil { writeDBError(w, r, err); return }
	n, _ := res.RowsAffected(); if n == 0 { writeError(w, r, http.StatusNotFound, "NOT_FOUND", "That task no longer exists.", nil); return }; w.WriteHeader(http.StatusNoContent)
}

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool { r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes); dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields(); if err := dec.Decode(dst); err != nil { writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The request body is not valid JSON.", nil); return false }; if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) { writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "The request body must contain one JSON object.", nil); return false }; return true }
func jsonContent(r *http.Request) bool { ct := r.Header.Get("Content-Type"); return ct != "" && strings.HasPrefix(strings.ToLower(ct), "application/json") }
func mapTodo(id, title string, done bool, created, updated time.Time) todoDTO { status := "active"; if done { status = "completed" }; return todoDTO{ID:id, Title:title, Status:status, CreatedAt:created.UTC().Format(time.RFC3339), UpdatedAt:updated.UTC().Format(time.RFC3339)} }
func writeDBError(w http.ResponseWriter, r *http.Request, err error) { if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) { writeError(w, r, http.StatusServiceUnavailable, "UNAVAILABLE", "The service is temporarily unavailable.", nil); return }; log.Printf("request_id=%s error=%v", requestID(r), err); writeError(w, r, http.StatusInternalServerError, "INTERNAL", "Something went wrong. Please try again.", nil) }
func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string, details []fieldError) { resp := errorResponse{}; resp.Error.Code = code; resp.Error.Message = message; resp.Error.Details = details; resp.Error.RequestID = requestID(r); writeJSON(w, status, resp) }
func writeJSON(w http.ResponseWriter, status int, value any) { w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.WriteHeader(status); if status != http.StatusNoContent { _ = json.NewEncoder(w).Encode(value) } }
func withRequestID(next http.Handler) http.Handler { return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { id := r.Header.Get("X-Request-Id"); if id == "" { id = newRequestID() }; w.Header().Set("X-Request-Id", id); next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id))) }) }
type requestIDKey struct{}
func requestID(r *http.Request) string { if v, ok := r.Context().Value(requestIDKey{}).(string); ok { return v }; return "unknown" }
func newRequestID() string { b := make([]byte, 8); if _, err := rand.Read(b); err != nil { return fmt.Sprintf("%d", time.Now().UnixNano()) }; return hex.EncodeToString(b) }
func isUUID(s string) bool { if len(s) != 36 { return false }; for i, c := range s { if i == 8 || i == 13 || i == 18 || i == 23 { if c != '-' { return false }; continue }; if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') { return false } }; return true }
func listenPort() string { if port := os.Getenv("PORT"); port != "" { return port }; if port := os.Getenv("APP_PORT"); port != "" { return port }; return "8080" }
