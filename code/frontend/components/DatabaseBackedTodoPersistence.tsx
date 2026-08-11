"use client";

import { FormEvent, useMemo, useState } from "react";
import {
  emptyTodoResponse,
  loadTodoError,
  persistedTodoResponse,
  saveTodoError,
  type TodoDto,
} from "../lib/mock/database-backed-todo-persistence";
import styles from "./DatabaseBackedTodoPersistence.module.css";

type ViewMode = "default" | "loading" | "empty" | "error";

const initialTodos = [...persistedTodoResponse.data].sort((a, b) => a.createdAt.localeCompare(b.createdAt));

export function DatabaseBackedTodoPersistence() {
  const [mode, setMode] = useState<ViewMode>("default");
  const [todos, setTodos] = useState<TodoDto[]>(initialTodos);
  const [input, setInput] = useState("");
  const [validation, setValidation] = useState("");
  const [notice, setNotice] = useState("Saved tasks are shown in oldest-first order.");
  const [failNextSave, setFailNextSave] = useState(false);

  const summary = useMemo(() => {
    const completed = todos.filter((todo) => todo.completed).length;
    return `${completed} completed of ${todos.length} saved tasks`;
  }, [todos]);

  function loadSavedTodos(nextMode: ViewMode = "default") {
    setValidation("");
    if (nextMode === "loading") {
      setMode("loading");
      setNotice("Loading saved tasks from the shared list.");
      return;
    }
    if (nextMode === "empty") {
      setTodos(emptyTodoResponse.data);
      setMode("empty");
      setNotice("The shared list is empty. Add the first task when ready.");
      return;
    }
    if (nextMode === "error") {
      setMode("error");
      setNotice(loadTodoError.error.message);
      return;
    }
    setTodos(initialTodos);
    setMode("default");
    setNotice("Saved tasks loaded from the shared list.");
  }

  function addTodo(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const title = input.trim();
    if (!title) {
      setValidation("Enter a task before adding it.");
      return;
    }
    if (failNextSave) {
      setNotice(saveTodoError.error.message);
      setFailNextSave(false);
      return;
    }
    const now = new Date().toISOString();
    const todo: TodoDto = {
      id: `todo_mock_${Date.now()}`,
      title,
      completed: false,
      createdAt: now,
      updatedAt: now,
    };
    setTodos((current) => [...current, todo]);
    setInput("");
    setValidation("");
    setMode("default");
    setNotice("Task saved to the shared list.");
  }

  function toggleTodo(id: string) {
    if (failNextSave) {
      setNotice("We could not save that change. The task is back to its last saved state.");
      setFailNextSave(false);
      return;
    }
    setTodos((current) => current.map((todo) => todo.id === id ? { ...todo, completed: !todo.completed, updatedAt: new Date().toISOString() } : todo));
    setNotice("Task status saved.");
  }

  function deleteTodo(id: string) {
    if (failNextSave) {
      setNotice("We could not delete that task. It is still in the saved list.");
      setFailNextSave(false);
      return;
    }
    setTodos((current) => {
      const next = current.filter((todo) => todo.id !== id);
      if (next.length === 0) setMode("empty");
      return next;
    });
    setNotice("Task deleted from the shared list.");
  }

  return (
    <section className={styles.card} aria-labelledby="persisted-todo-title">
      <div className={styles.top}>
        <div>
          <h2 id="persisted-todo-title" className={styles.title}>Today</h2>
          <p className={styles.subtitle}>Database-backed todo preview</p>
        </div>
        <span className={styles.pill}>Shared saved list</span>
      </div>

      <form className={styles.toolbar} onSubmit={addTodo} noValidate>
        <div className={styles.field}>
          <label className={styles.sr} htmlFor="persisted-todo-input">Task name</label>
          <input
            id="persisted-todo-input"
            className={styles.input}
            value={input}
            maxLength={80}
            onChange={(event) => setInput(event.target.value)}
            placeholder="Add a task, e.g. Send weekly update"
            aria-describedby={validation ? "persisted-input-error" : undefined}
          />
          {validation ? <p id="persisted-input-error" className={styles.validation}>{validation}</p> : null}
        </div>
        <button className={styles.primary} type="submit">Add</button>
      </form>

      <div className={styles.notice} role="status">{notice}</div>

      <div className={styles.controls} aria-label="Preview states">
        <button className={styles.ghost} type="button" onClick={() => loadSavedTodos("default")}>Saved</button>
        <button className={styles.ghost} type="button" onClick={() => loadSavedTodos("loading")}>Loading</button>
        <button className={styles.ghost} type="button" onClick={() => loadSavedTodos("empty")}>Empty</button>
        <button className={styles.ghost} type="button" onClick={() => loadSavedTodos("error")}>Error</button>
        <button className={styles.ghost} type="button" onClick={() => setFailNextSave(true)}>Fail next save</button>
      </div>

      {mode === "loading" ? <LoadingState /> : null}
      {mode === "error" ? <ErrorState onRetry={() => loadSavedTodos("default")} /> : null}
      {mode === "empty" ? <EmptyState /> : null}
      {mode === "default" ? (
        <>
          <p className={styles.summary}>{summary}</p>
          <ul className={styles.list} aria-live="polite">
            {todos.map((todo) => (
              <li className={styles.row} key={todo.id}>
                <button className={`${styles.check} ${todo.completed ? styles.checked : ""}`} type="button" onClick={() => toggleTodo(todo.id)} aria-label={todo.completed ? "Mark incomplete" : "Mark complete"}>{todo.completed ? "✓" : ""}</button>
                <span className={todo.completed ? styles.completedText : styles.todoText}>{todo.title}</span>
                <button className={styles.delete} type="button" onClick={() => deleteTodo(todo.id)} aria-label={`Delete ${todo.title}`}>×</button>
              </li>
            ))}
          </ul>
        </>
      ) : null}
    </section>
  );
}

function LoadingState() {
  return <div className={styles.state}><div className={styles.dots}><span /><span /><span /></div><h3>Loading saved tasks</h3><p>Please wait while the shared list is loaded.</p></div>;
}

function EmptyState() {
  return <div className={styles.state}><div className={styles.emptyIcon}>＋</div><h3>No saved tasks yet</h3><p>Add the first task to start the shared list.</p></div>;
}

function ErrorState({ onRetry }: { onRetry: () => void }) {
  return <div className={styles.state}><div className={styles.errorIcon}>!</div><h3>Tasks could not load</h3><p>Something got in the way. Try again without losing the page.</p><button className={styles.primary} type="button" onClick={onRetry}>Try again</button></div>;
}
