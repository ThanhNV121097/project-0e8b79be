export default function HomePage() {
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
          A clean blue-and-white todo workspace for adding, completing, and deleting tasks. The interactive experience will be wired to the database-backed API in the feature stories.
        </p>
      </section>

      <section className="todo-card" aria-labelledby="todo-card-title">
        <div className="mb-6 flex flex-col gap-3 sm:flex-row">
          <label className="sr-only" htmlFor="todo-title">Todo title</label>
          <input id="todo-title" className="input-field" placeholder="Add a task, e.g. Send weekly update" disabled />
          <button className="primary-button" type="button" disabled>Add task</button>
        </div>
        <h2 id="todo-card-title" className="mb-4 text-2xl font-black text-text">Your tasks</h2>
        <div className="state-panel">
          <p className="font-bold text-text">Scaffold ready</p>
          <p className="mt-2">Todo interactions and persistence are intentionally implemented in later stories.</p>
        </div>
      </section>
    </main>
  );
}
