package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

type todo struct { ID string `json:"id"`; Title string `json:"title"`; Status string `json:"status"`; CreatedAt string `json:"createdAt"`; UpdatedAt string `json:"updatedAt"` }
type app struct{ db *sql.DB; limiter *limiter }
type apiError struct{ Error errBody `json:"error"` }
type errBody struct { Code string `json:"code"`; Message string `json:"message"`; Details []errDetail `json:"details"`; RequestID string `json:"request_id"` }
type errDetail struct { Field string `json:"field"`; Code string `json:"code"`; Message string `json:"message"` }
var uuidRE = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func main() { databaseURL := os.Getenv("DATABASE_URL"); if databaseURL == "" { log.Fatal("DATABASE_URL is required") }; ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second); defer cancel(); database, err := sql.Open("pgx", databaseURL); if err != nil { log.Fatalf("open database: %v", err) }; defer database.Close(); if err := database.PingContext(ctx); err != nil { log.Fatalf("ping database: %v", err) }; if err := db.ApplyMigrations(ctx, database); err != nil { log.Fatalf("apply migrations: %v", err) }; a := &app{db: database, limiter: newLimiter()}; mux := http.NewServeMux(); mux.HandleFunc("/healthz", a.health); mux.HandleFunc("/api/v1/todos", a.todos); mux.HandleFunc("/api/v1/todos/", a.todoByID); server := &http.Server{Addr: ":" + listenPort(), Handler: a.withLog(mux), ReadHeaderTimeout: 5*time.Second}; log.Printf("backend listening on %s", server.Addr); if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) { log.Fatalf("listen: %v", err) } }
func (a *app) health(w http.ResponseWriter, r *http.Request) { if r.Method != http.MethodGet { w.WriteHeader(405); return }; ctx,c:=context.WithTimeout(r.Context(),3*time.Second); defer c(); if err:=a.db.PingContext(ctx); err!=nil { writeErr(w,r,503,"UNAVAILABLE","Service is unavailable.",nil); return }; writeJSON(w,200,map[string]string{"status":"ok"}) }
func (a *app) todos(w http.ResponseWriter, r *http.Request) { if r.URL.Path!="/api/v1/todos" { writeErr(w,r,404,"NOT_FOUND","Todo was not found.",nil); return }; switch r.Method { case http.MethodGet: a.list(w,r); case http.MethodPost: a.create(w,r); default: w.WriteHeader(405) } }
func (a *app) todoByID(w http.ResponseWriter, r *http.Request) { id:=strings.TrimPrefix(r.URL.Path,"/api/v1/todos/"); if strings.Contains(id,"/")||!uuidRE.MatchString(id){ writeErr(w,r,400,"BAD_REQUEST","Request is invalid.",nil); return }; switch r.Method { case http.MethodPatch: a.patch(w,r,id); case http.MethodDelete: a.delete(w,r,id); default: w.WriteHeader(405) } }
func (a *app) list(w http.ResponseWriter, r *http.Request) { if !a.limiter.allow(client(r),120){ w.Header().Set("Retry-After","60"); writeErr(w,r,429,"RATE_LIMITED","Too many requests. Please try again shortly.",nil); return }; ctx,c:=context.WithTimeout(r.Context(),3*time.Second); defer c(); rows,err:=a.db.QueryContext(ctx,`SELECT id::text,title,CASE WHEN is_completed THEN 'completed' ELSE 'active' END,created_at,updated_at FROM todos ORDER BY created_at ASC,id ASC`); if err!=nil{ dbErr(w,r,err); return }; defer rows.Close(); data:=[]todo{}; active:=0; comp:=0; for rows.Next(){ var t todo; var cr,up time.Time; if err:=rows.Scan(&t.ID,&t.Title,&t.Status,&cr,&up); err!=nil{ dbErr(w,r,err); return }; t.CreatedAt=cr.UTC().Format(time.RFC3339); t.UpdatedAt=up.UTC().Format(time.RFC3339); if t.Status=="completed"{comp++}else{active++}; data=append(data,t) }; if err:=rows.Err(); err!=nil{dbErr(w,r,err);return}; writeJSON(w,200,map[string]any{"data":data,"meta":map[string]int{"total":len(data),"active":active,"completed":comp}}) }
func (a *app) create(w http.ResponseWriter, r *http.Request) { if !writeOK(a,w,r){return}; var p struct{Title any `json:"title"`}; if !decode(w,r,&p){return}; s,ok:=p.Title.(string); if !ok{ writeErr(w,r,400,"BAD_REQUEST","Request is invalid.",nil); return }; title:=strings.TrimSpace(s); if title==""{ writeErr(w,r,422,"VALIDATION_FAILED","Todo title is invalid.",[]errDetail{{"title","BLANK","Title cannot be blank."}}); return }; if len([]rune(title))>200{ writeErr(w,r,422,"VALIDATION_FAILED","Todo title is invalid.",[]errDetail{{"title","TOO_LONG","Title must be 200 characters or fewer."}}); return }; t,err:=a.queryOne(r.Context(),`INSERT INTO todos (title) VALUES ($1) RETURNING id::text,title,CASE WHEN is_completed THEN 'completed' ELSE 'active' END,created_at,updated_at`,title); if err!=nil{dbErr(w,r,err);return}; w.Header().Set("Location","/api/v1/todos/"+t.ID); writeJSON(w,201,t) }
func (a *app) patch(w http.ResponseWriter, r *http.Request, id string) { if !writeOK(a,w,r){return}; var p struct{Status any `json:"status"`}; if !decode(w,r,&p){return}; s,ok:=p.Status.(string); if !ok{ writeErr(w,r,400,"BAD_REQUEST","Request is invalid.",nil); return }; if s!="active"&&s!="completed"{ writeErr(w,r,422,"VALIDATION_FAILED","Todo status is invalid.",[]errDetail{{"status","INVALID","Status must be active or completed."}}); return }; t,err:=a.queryOne(r.Context(),`UPDATE todos SET is_completed=$2,updated_at=now() WHERE id=$1 RETURNING id::text,title,CASE WHEN is_completed THEN 'completed' ELSE 'active' END,created_at,updated_at`,id,s=="completed"); if errors.Is(err,sql.ErrNoRows){ writeErr(w,r,404,"NOT_FOUND","Todo was not found.",nil); return }; if err!=nil{dbErr(w,r,err);return}; writeJSON(w,200,t) }
func (a *app) delete(w http.ResponseWriter, r *http.Request, id string) { if !writeOK(a,w,r){return}; if r.ContentLength>0{ writeErr(w,r,400,"BAD_REQUEST","Request is invalid.",nil); return }; ctx,c:=context.WithTimeout(r.Context(),3*time.Second); defer c(); res,err:=a.db.ExecContext(ctx,`DELETE FROM todos WHERE id=$1`,id); if err!=nil{dbErr(w,r,err);return}; n,_:=res.RowsAffected(); if n==0{ writeErr(w,r,404,"NOT_FOUND","Todo was not found.",nil); return }; w.WriteHeader(204) }
func (a *app) queryOne(ctx context.Context,q string,args ...any)(todo,error){ c,cn:=context.WithTimeout(ctx,3*time.Second); defer cn(); var t todo; var cr,up time.Time; err:=a.db.QueryRowContext(c,q,args...).Scan(&t.ID,&t.Title,&t.Status,&cr,&up); if err==nil{t.CreatedAt=cr.UTC().Format(time.RFC3339); t.UpdatedAt=up.UTC().Format(time.RFC3339)}; return t,err }
func writeOK(a *app,w http.ResponseWriter,r *http.Request)bool{ if !a.limiter.allow(client(r),60){w.Header().Set("Retry-After","60"); writeErr(w,r,429,"RATE_LIMITED","Too many requests. Please try again shortly.",nil); return false}; if ct:=r.Header.Get("Content-Type"); ct!=""&&!strings.HasPrefix(ct,"application/json"){writeErr(w,r,400,"BAD_REQUEST","Request is invalid.",nil); return false}; r.Body=http.MaxBytesReader(w,r.Body,16*1024); return true }
func decode(w http.ResponseWriter,r *http.Request,v any)bool{ dec:=json.NewDecoder(r.Body); if err:=dec.Decode(v); err!=nil{writeErr(w,r,400,"BAD_REQUEST","Request is invalid.",nil); return false}; return true }
func writeJSON(w http.ResponseWriter,code int,v any){ w.Header().Set("Content-Type","application/json; charset=utf-8"); w.WriteHeader(code); _=json.NewEncoder(w).Encode(v) }
func writeErr(w http.ResponseWriter,r *http.Request,status int,code,msg string,d []errDetail){ rid:=requestID(r); w.Header().Set("X-Request-Id",rid); writeJSON(w,status,apiError{errBody{code,msg,d,rid}}) }
func dbErr(w http.ResponseWriter,r *http.Request,err error){ if errors.Is(err,context.DeadlineExceeded){writeErr(w,r,503,"UNAVAILABLE","Service is unavailable.",nil)}else{writeErr(w,r,500,"INTERNAL","Something went wrong.",nil)} }
func requestID(r *http.Request)string{ if v:=r.Header.Get("X-Request-Id"); v!=""{return v}; return fmt.Sprintf("%d",time.Now().UnixNano()) }
func client(r *http.Request)string{ h,_,e:=net.SplitHostPort(r.RemoteAddr); if e!=nil{return r.RemoteAddr}; return h }
func (a *app) withLog(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){start:=time.Now(); next.ServeHTTP(w,r); sum:=sha256.Sum256([]byte(client(r))); log.Printf("request_id=%s method=%s path=%s duration_ms=%d client_source_hash=%s",requestID(r),r.Method,r.URL.Path,time.Since(start).Milliseconds(),hex.EncodeToString(sum[:8]))})}
type limiter struct{ mu sync.Mutex; hits map[string][]time.Time }; func newLimiter()*limiter{return &limiter{hits:map[string][]time.Time{}}}; func (l *limiter)allow(k string,max int)bool{l.mu.Lock(); defer l.mu.Unlock(); now:=time.Now(); cutoff:=now.Add(-time.Minute); xs:=l.hits[k]; j:=0; for _,t:=range xs{if t.After(cutoff){xs[j]=t;j++}}; xs=xs[:j]; if len(xs)>=max{l.hits[k]=xs; return false}; l.hits[k]=append(xs,now); return true}
func listenPort() string { if port := os.Getenv("PORT"); port != "" { return port }; if port := os.Getenv("APP_PORT"); port != "" { return port }; return "8080" }
