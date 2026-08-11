"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";

import {
  createTodo,
  deleteTodo,
  listTodos,
  TodoDto,
  todoErrorMessage,
  updateTodoStatus,
} from "@/lib/api/todos";

const loadErrorMessage = "We could not load your tasks. Please try again.";

export default function HomePage() {
  const [todos, setTodos] = useState<TodoDto[]>([]);
  const [title, setTitle] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [saveError, setSaveError] = useState<string | null>(null);
  const [isAdding, setIsAdding] = useState(false);
  const [pendingIds, setPendingIds] = useState<Set<string>>(new Set());

  const counts = useMemo(() => {
    const completed = todos.filter((todo) => todo.status === "completed").length;
    return { total: todos.length, active: todos.length - completed, completed };
  }, [todos]);

  async function loadSavedTodos() {
    setIsLoading(true);
    setLoadError(null);
    try {
      const response = await listTodos();
      setTodos(response.data);
    } catch {
      setLoadError(loadErrorMessage);
    } finally {
      setIsLoading(false);
    }
  }

  useEffect(() => {
    void loadSavedTodos();
  }, []);

  async function handleAdd(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedTitle = title.trim();
    if (!trimmedTitle || isAdding) {
      return;
    }

    setIsAdding(true);
    setSaveError(null);
    try {
      const savedTodo = await createTodo(trimmedTitle);
      setTodos((current) => [...current, savedTodo]);
      setTitle("");
    } catch {
      setSaveError(todoErrorMessage);
    } finally {
      setIsAdding(false);
    }
  }

  async function handleToggle(todo: TodoDto) {
    if (pendingIds.has(todo.id)) {
      return;
    }

    const nextStatus = todo.status === "completed" ? "active" : "completed";
    const previousTodo = todo;
    setSaveError(null);
    setPendingIds((current) => new Set(current).add(todo.id));
    setTodos((current) => current.map((item) => (item.id === todo.id ? { ...item, status: nextStatus } : item)));

    try {
      const savedTodo = await updateTodoStatus(todo.id, nextStatus);
      setTodos((current) => current.map((item) => (item.id === todo.id ? savedTodo : item)));
    } catch {
      setTodos((current) => current.map((item) => (item.id === todo.id ? previousTodo : item)));
      setSaveError(todoErrorMessage);
    } finally {
      setPendingIds((current) => {
        const next = new Set(current);
        next.delete(todo.id);
        return next;
      });
    }
  }

  async function handleDelete(todo: TodoDto) {
    if (pendingIds.has(todo.id)) {
      return;
    }

    setSaveError(null);
    setPendingIds((current) => new Set(current).add(todo.id));
    setTodos((current) => current.filter((item) => item.id !== todo.id));

    try {
      await deleteTodo(todo.id);
    } catch {
      setTodos((current) => {
        if (current.some((item) => item.id === todo.id)) {
          return current;
        }
        return [...current, todo].sort((left, right) => {
          const createdDiff = Date.parse(left.createdAt) - Date.parse(right.createdAt);
          return createdDiff || left.id.localeCompare(right.id);
        });
      });
      setSaveError(todoErrorMessage);
    } finally {
      setPendingIds((current) => {
        const next = new Set(current);
        next.delete(todo.id);
        return next;
      });
    }
  }

  return (
    <main className="app-shell">
      <section className="app-container text-center">
        <p className="mb-4 inline-flex rounded-full bg-primarySoft px-4 py-2 text-sm font-extrabold text-primary">
          Simple, saved, and ready
        </p>
        <h1 className="mx-auto max-w-3xl text-5xl font-black tracking-tight text-[#0B1220] md:text-7xl">
          Todo List App v2
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-lg leading-8 text-muted">
          A clean blue-and-white todo workspace for adding, completing, and deleting tasks. Your tasks are saved to the database and stay available after refresh.
        </p>
      </section>

      <section className="todo-card" aria-labelledby="todo-card-title">
        <form className="mb-4 flex flex-col gap-3 sm:flex-row" onSubmit={handleAdd}>
          <label className="sr-only" htmlFor="todo-title">Todo title</label>
          <input
            id="todo-title"
            className="input-field"
            maxLength={200}
            onChange={(event) => setTitle(event.target.value)}
            placeholder="Add a task, e.g. Send weekly update"
            value={title}
          />
          <button className="primary-button disabled:cursor-not-allowed disabled:opacity-60" disabled={isAdding || !title.trim()} type="submit">
            {isAdding ? "Saving..." : "Add task"}
          </button>
        </form>

        {saveError ? (
          <div className="mb-4 rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm font-bold text-red-700" role="alert">
            {saveError}
          </div>
        ) : null}

        <div className="mb-4 flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
          <h2 id="todo-card-title" className="text-2xl font-black text-text">Your tasks</h2>
          <p className="text-sm font-bold text-muted">
            {counts.total} total · {counts.active} active · {counts.completed} completed
          </p>
        </div>

        {isLoading ? (
          <div className="state-panel" aria-live="polite">
            <p className="font-bold text-text">Loading your saved tasks...</p>
            <p className="mt-2">This should only take a moment.</p>
          </div>
        ) : loadError ? (
          <div className="state-panel" role="alert">
            <p className="font-bold text-text">{loadError}</p>
            <button className="primary-button mt-4" type="button" onClick={() => void loadSavedTodos()}>
              Try again
            </button>
          </div>
        ) : todos.length === 0 ? (
          <div className="state-panel">
            <p className="font-bold text-text">No tasks yet</p>
            <p className="mt-2">Add your first task above and it will be saved here.</p>
          </div>
        ) : (
          <ul className="space-y-3" aria-label="Saved tasks">
            {todos.map((todo) => {
              const isPending = pendingIds.has(todo.id);
              const isCompleted = todo.status === "completed";
              return (
                <li className="todo-row" key={todo.id}>
                  <button
                    aria-label={isCompleted ? `Mark ${todo.title} active` : `Mark ${todo.title} complete`}
                    className={`h-7 w-7 rounded-full border-2 ${isCompleted ? "border-emerald-500 bg-emerald-500 text-white" : "border-blue-300 bg-white text-transparent"}`}
                    disabled={isPending}
                    onClick={() => void handleToggle(todo)}
                    type="button"
                  >
                    ✓
                  </button>
                  <span className={`flex-1 font-bold ${isCompleted ? "text-muted line-through" : "text-text"}`}>{todo.title}</span>
                  <button
                    className="rounded-xl bg-red-50 px-3 py-2 text-sm font-black text-red-600 disabled:cursor-not-allowed disabled:opacity-60"
                    disabled={isPending}
                    onClick={() => void handleDelete(todo)}
                    type="button"
                  >
                    Delete
                  </button>
                </li>
              );
            })}
          </ul>
        )}
      </section>
    </main>
  );
}
