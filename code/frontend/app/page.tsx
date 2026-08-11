"use client";

import { FormEvent, useEffect, useState } from "react";

type Todo = { id: string; title: string; status: "active" | "completed"; createdAt: string; updatedAt: string };
const API = (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080") + "/api/v1";
const SAVE_ERROR = "We could not save that change. Please try again.";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = init?.body ? { "Content-Type": "application/json", ...(init.headers ?? {}) } : init?.headers;
  const res = await fetch(API + path, { ...init, headers });
  if (!res.ok) throw new Error("api_error");
  if (res.status === 204) return undefined as T;
  return res.json();
}

export default function HomePage() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [title, setTitle] = useState("");
  const [validation, setValidation] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  async function load() {
    setLoading(true); setError("");
    try { const res = await request<{ data: Todo[] }>("/todos"); setTodos(res.data); }
    catch { setError("We could not load your tasks. Please try again."); }
    finally { setLoading(false); }
  }
  useEffect(() => { void load(); }, []);

  async function add(e: FormEvent) {
    e.preventDefault(); const trimmed = title.trim();
    if (!trimmed) { setValidation("Task cannot be blank."); return; }
    setValidation(""); setError("");
    try { const todo = await request<Todo>("/todos", { method: "POST", body: JSON.stringify({ title: trimmed }) }); setTodos((items) => [...items, todo]); setTitle(""); }
    catch { setError(SAVE_ERROR); }
  }
  async function toggle(todo: Todo) {
    const next = todo.status === "active" ? "completed" : "active";
    setTodos((items) => items.map((item) => item.id === todo.id ? { ...item, status: next } : item)); setError("");
    try { const saved = await request<Todo>(`/todos/${todo.id}`, { method: "PATCH", body: JSON.stringify({ status: next }) }); setTodos((items) => items.map((item) => item.id === todo.id ? saved : item)); }
    catch { setTodos((items) => items.map((item) => item.id === todo.id ? todo : item)); setError(SAVE_ERROR); }
  }
  async function remove(todo: Todo) {
    setTodos((items) => items.filter((item) => item.id !== todo.id)); setError("");
    try { await request<void>(`/todos/${todo.id}`, { method: "DELETE" }); }
    catch { setTodos((items) => items.some((item) => item.id === todo.id) ? items : [...items, todo]); setError(SAVE_ERROR); }
  }

  return <main className="app-shell"><section className="app-container text-center"><p className="mb-4 inline-flex rounded-full bg-primarySoft px-4 py-2 text-sm font-extrabold text-primary">Simple, saved, and ready</p><h1 className="mx-auto max-w-3xl text-5xl font-black tracking-tight text-[#0B1220] md:text-7xl">Todo List App v2</h1><p className="mx-auto mt-6 max-w-2xl text-lg leading-8 text-muted">A clean blue-and-white todo workspace for adding, completing, and deleting tasks.</p></section><section className="todo-card" aria-labelledby="todo-card-title"><form onSubmit={add} className="mb-6 flex flex-col gap-3 sm:flex-row"><label className="sr-only" htmlFor="todo-title">Todo title</label><input id="todo-title" className="input-field" placeholder="Add a task, e.g. Send weekly update" value={title} onChange={(e) => setTitle(e.target.value)} maxLength={200} /><button className="primary-button" type="submit">Add task</button></form>{validation && <p className="mb-3 font-bold text-red-600">{validation}</p>}{error && <p className="mb-3 font-bold text-red-600">{error}</p>}<h2 id="todo-card-title" className="mb-4 text-2xl font-black text-text">Your tasks</h2>{loading ? <div className="state-panel">Loading tasks…</div> : todos.length === 0 ? <div className="state-panel"><p className="font-bold text-text">No tasks yet</p><p className="mt-2">Add your first task to get started.</p></div> : <ul className="space-y-3">{todos.map((todo) => <li className="todo-row" key={todo.id}><button type="button" onClick={() => toggle(todo)} className="font-bold text-primary">{todo.status === "completed" ? "✓" : "○"}</button><span className={todo.status === "completed" ? "flex-1 text-muted line-through" : "flex-1 text-text"}>{todo.title}</span><button type="button" onClick={() => remove(todo)} className="font-bold text-red-600">Delete</button></li>)}</ul>}</section></main>;
}
