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
  meta: {
    total: number;
    active: number;
    completed: number;
  };
};

export type TodoMutationError = {
  error: {
    code: "SAVE_FAILED" | "VALIDATION_FAILED";
    message: string;
  };
};

export const MAX_TODO_TITLE_LENGTH = 200;

const seedTodos: TodoDto[] = [
  {
    id: "todo_001",
    title: "Plan the week",
    status: "active",
    createdAt: "2026-08-11T09:00:00.000Z",
    updatedAt: "2026-08-11T09:00:00.000Z",
  },
  {
    id: "todo_002",
    title: "Send project update",
    status: "completed",
    createdAt: "2026-08-11T08:30:00.000Z",
    updatedAt: "2026-08-11T09:15:00.000Z",
  },
];

export function getMockTodoList(): TodoListResponse {
  return toTodoListResponse(seedTodos);
}

export function toTodoListResponse(todos: TodoDto[]): TodoListResponse {
  return {
    data: todos,
    meta: {
      total: todos.length,
      active: todos.filter((todo) => todo.status === "active").length,
      completed: todos.filter((todo) => todo.status === "completed").length,
    },
  };
}

export function createMockTodo(title: string): TodoDto {
  const timestamp = new Date().toISOString();

  return {
    id: `todo_${crypto.randomUUID()}`,
    title,
    status: "active",
    createdAt: timestamp,
    updatedAt: timestamp,
  };
}

export const mockTodoError: TodoMutationError = {
  error: {
    code: "SAVE_FAILED",
    message: "We could not save that change. Please try again.",
  },
};
