"use client";

import { FormEvent, useMemo, useState } from "react";
import {
  createMockTodo,
  getMockTodoList,
  MAX_TODO_TITLE_LENGTH,
  mockTodoError,
  toTodoListResponse,
  type TodoDto,
} from "@/lib/mock/single-page-todo-experience";
import styles from "./SinglePageTodoExperience.module.css";

type ViewState = "default" | "loading" | "empty" | "error";

export default function SinglePageTodoExperience() {
  const initialResponse = useMemo(() => getMockTodoList(), []);
  const [todos, setTodos] = useState<TodoDto[]>(initialResponse.data);
  const [title, setTitle] = useState("");
  const [validationMessage, setValidationMessage] = useState("");
  const [saveMessage, setSaveMessage] = useState("");
  const [viewState, setViewState] = useState<ViewState>("default");

  const response = toTodoListResponse(todos);
  const visibleTodos = viewState === "empty" ? [] : response.data;

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedTitle = title.trim();
    setSaveMessage("");

    if (!trimmedTitle) {
      setValidationMessage("The task cannot be blank.");
      return;
    }

    if (trimmedTitle.length > MAX_TODO_TITLE_LENGTH) {
      setValidationMessage(`Tasks must be ${MAX_TODO_TITLE_LENGTH} characters or fewer.`);
      return;
    }

    setValidationMessage("");

    if (viewState === "error") {
      setSaveMessage(mockTodoError.error.message);
      return;
    }

    setTodos((currentTodos) => [createMockTodo(trimmedTitle), ...currentTodos]);
    setTitle("");
    setViewState("default");
  }

  function toggleTodo(todoId: string) {
    setSaveMessage("");

    if (viewState === "error") {
      setSaveMessage(mockTodoError.error.message);
      return;
    }

    setTodos((currentTodos) =>
      currentTodos.map((todo) =>
        todo.id === todoId
          ? {
              ...todo,
              status: todo.status === "completed" ? "active" : "completed",
              updatedAt: new Date().toISOString(),
            }
          : todo,
      ),
    );
  }

  function deleteTodo(todoId: string) {
    setSaveMessage("");

    if (viewState === "error") {
      setSaveMessage(mockTodoError.error.message);
      return;
    }

    setTodos((currentTodos) => {
      const nextTodos = currentTodos.filter((todo) => todo.id !== todoId);
      if (nextTodos.length === 0) {
        setViewState("empty");
      }
      return nextTodos;
    });
  }

  function resetTodos() {
    const resetResponse = getMockTodoList();
    setTodos(resetResponse.data);
    setViewState("default");
    setSaveMessage("");
    setValidationMessage("");
  }

  return (
    <main className={styles.shell}>
      <section className={styles.hero} aria-labelledby="todo-page-title">
        <p className={styles.eyebrow}>Simple, saved, and ready</p>
        <h1 id="todo-page-title">Todo List App v2</h1>
        <p>A clean blue-and-white workspace for adding tasks, checking them off, and keeping the day moving without extra navigation.</p>
      </section>

      <section className={styles.card} aria-labelledby="todo-card-title">
        <div className={styles.cardHeader}>
          <div>
            <h2 id="todo-card-title">Today</h2>
            <p>Interactive todo preview</p>
          </div>
          <span className={styles.pill} aria-label={`${response.meta.active} active and ${response.meta.completed} completed tasks`}>
            {response.meta.active} active · {response.meta.completed} done
          </span>
        </div>

        <form className={styles.toolbar} onSubmit={handleSubmit} noValidate>
          <div className={styles.field}>
            <label className={styles.srOnly} htmlFor="todo-title">Task name</label>
            <input id="todo-title" value={title} onChange={(event) => setTitle(event.target.value)} className={styles.input} maxLength={MAX_TODO_TITLE_LENGTH + 1} placeholder="Add a task, e.g. Send weekly update" aria-describedby={validationMessage ? "todo-title-error" : undefined} />
            {validationMessage ? <p className={styles.validation} id="todo-title-error" role="alert">{validationMessage}</p> : null}
          </div>
          <button className={styles.primaryButton} type="submit">Add</button>
        </form>

        <div className={styles.stateControls} aria-label="Preview todo states">
          <button type="button" onClick={resetTodos}>Default</button>
          <button type="button" onClick={() => setViewState("loading")}>Loading</button>
          <button type="button" onClick={() => setViewState("empty")}>Empty</button>
          <button type="button" onClick={() => setViewState("error")}>Error</button>
        </div>

        {saveMessage ? <p className={styles.notice} role="alert">{saveMessage}</p> : null}

        <div className={styles.summary} aria-live="polite">
          <span>{response.meta.total} total</span>
          <span>{response.meta.active} active</span>
          <span>{response.meta.completed} completed</span>
        </div>

        {viewState === "loading" ? (
          <div className={styles.statePanel} aria-live="polite">
            <div className={styles.loadingDots} aria-hidden="true"><span /><span /><span /></div>
            <h3>Loading tasks</h3>
            <p>Your todos are being prepared.</p>
          </div>
        ) : visibleTodos.length === 0 ? (
          <div className={styles.statePanel} aria-live="polite">
            <div className={styles.emptyIcon} aria-hidden="true">＋</div>
            <h3>No tasks yet</h3>
            <p>Add your first task to start a simple focused list.</p>
          </div>
        ) : (
          <ul className={styles.todoList} aria-live="polite">
            {visibleTodos.map((todo) => (
              <li className={styles.todoRow} data-completed={todo.status === "completed"} key={todo.id}>
                <button className={styles.checkButton} type="button" aria-label={todo.status === "completed" ? `Mark ${todo.title} incomplete` : `Mark ${todo.title} complete`} onClick={() => toggleTodo(todo.id)}>{todo.status === "completed" ? "✓" : ""}</button>
                <span className={styles.todoTitle}>{todo.title}</span>
                <button className={styles.deleteButton} type="button" aria-label={`Delete ${todo.title}`} onClick={() => deleteTodo(todo.id)}>×</button>
              </li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
}
