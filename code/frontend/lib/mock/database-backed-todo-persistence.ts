export type TodoDto = {
  id: string;
  title: string;
  completed: boolean;
  createdAt: string;
  updatedAt: string;
};

export type TodoListResponse = {
  data: TodoDto[];
  meta: {
    total: number;
    order: "created_at_asc";
  };
};

export type TodoMutationResponse = {
  data: TodoDto;
};

export type TodoDeleteResponse = {
  data: {
    id: string;
    deleted: true;
  };
};

export type TodoApiError = {
  error: {
    code: "LOAD_FAILED" | "CREATE_FAILED" | "UPDATE_FAILED" | "DELETE_FAILED" | "NOT_FOUND" | "VALIDATION_FAILED";
    message: string;
  };
};

export const persistedTodoResponse: TodoListResponse = {
  data: [
    {
      id: "todo_20260811_001",
      title: "Buy milk",
      completed: false,
      createdAt: "2026-08-11T08:00:00.000Z",
      updatedAt: "2026-08-11T08:00:00.000Z",
    },
    {
      id: "todo_20260811_002",
      title: "Pay rent",
      completed: true,
      createdAt: "2026-08-11T09:15:00.000Z",
      updatedAt: "2026-08-11T09:30:00.000Z",
    },
    {
      id: "todo_20260811_003",
      title: "Send weekly update",
      completed: false,
      createdAt: "2026-08-11T10:20:00.000Z",
      updatedAt: "2026-08-11T10:20:00.000Z",
    },
  ],
  meta: {
    total: 3,
    order: "created_at_asc",
  },
};

export const oneHundredTodoResponse: TodoListResponse = {
  data: Array.from({ length: 100 }, (_, index) => {
    const padded = String(index + 1).padStart(3, "0");
    const timestamp = new Date(Date.UTC(2026, 7, 11, 12, index, 0)).toISOString();
    return {
      id: `todo_20260811_bulk_${padded}`,
      title: `Saved database task ${padded}`,
      completed: index % 3 === 0,
      createdAt: timestamp,
      updatedAt: timestamp,
    };
  }),
  meta: {
    total: 100,
    order: "created_at_asc",
  },
};

export const emptyTodoResponse: TodoListResponse = {
  data: [],
  meta: {
    total: 0,
    order: "created_at_asc",
  },
};

export const loadTodoError: TodoApiError = {
  error: {
    code: "LOAD_FAILED",
    message: "We could not load the saved tasks. Please try again.",
  },
};

export const saveTodoError: TodoApiError = {
  error: {
    code: "CREATE_FAILED",
    message: "We could not save that task yet. Your text is still here, so please try again.",
  },
};

export const missingTodoError: TodoApiError = {
  error: {
    code: "NOT_FOUND",
    message: "That task was already changed elsewhere. The list has been refreshed.",
  },
};
