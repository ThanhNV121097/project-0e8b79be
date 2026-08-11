export type TodoStatus = "active" | "completed";

export type TodoDto = {
  id: string;
  title: string;
  status: TodoStatus;
  createdAt: string;
  updatedAt: string;
};

export type TodoListResponse = {
  data: TodoDto[];
  meta: { total: number; active: number; completed: number };
};

const API_BASE = (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080").replace(/\/$/, "");
const SAVE_MESSAGE = "We could not save that change. Please try again.";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}/api/v1${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) },
  });
  if (!response.ok) throw new Error(SAVE_MESSAGE);
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}

export function listTodos() {
  return request<TodoListResponse>("/todos", { method: "GET" });
}

export function createTodo(title: string) {
  return request<TodoDto>("/todos", { method: "POST", body: JSON.stringify({ title }) });
}

export function updateTodoStatus(id: string, status: TodoStatus) {
  return request<TodoDto>(`/todos/${id}`, { method: "PATCH", body: JSON.stringify({ status }) });
}

export function deleteTodo(id: string) {
  return request<void>(`/todos/${id}`, { method: "DELETE", headers: {} });
}

export const todoErrorMessage = SAVE_MESSAGE;
