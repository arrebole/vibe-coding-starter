import React, { useEffect, useMemo, useState } from 'react'
import { createRoot } from 'react-dom/client'
import { Check, Loader2, Plus, Trash2 } from 'lucide-react'
import './styles.css'

type ApiResponse<T> = {
  code: number
  message: string
  data?: T
}

type Todo = {
  id: number
  title: string
  completed: boolean
  created_at: string
  updated_at: string
}

const apiBaseURL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

function App() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [title, setTitle] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  const completedCount = useMemo(() => todos.filter((todo) => todo.completed).length, [todos])

  async function loadTodos() {
    setLoading(true)
    setError('')
    try {
      const response = await request<Todo[]>('/api/v1/todos')
      setTodos(response)
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载 Todo 失败')
    } finally {
      setLoading(false)
    }
  }

  async function createTodo(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const value = title.trim()
    if (!value) return

    setSaving(true)
    setError('')
    try {
      const todo = await request<Todo>('/api/v1/todos', {
        method: 'POST',
        body: JSON.stringify({ title: value }),
      })
      setTodos((items) => [todo, ...items])
      setTitle('')
    } catch (err) {
      setError(err instanceof Error ? err.message : '创建 Todo 失败')
    } finally {
      setSaving(false)
    }
  }

  async function toggleTodo(todo: Todo) {
    setError('')
    try {
      const updated = await request<Todo>(`/api/v1/todos/${todo.id}`, {
        method: 'PATCH',
        body: JSON.stringify({ completed: !todo.completed }),
      })
      setTodos((items) => items.map((item) => (item.id === updated.id ? updated : item)))
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新 Todo 失败')
    }
  }

  async function deleteTodo(id: number) {
    setError('')
    try {
      await request<{ deleted: boolean }>(`/api/v1/todos/${id}`, { method: 'DELETE' })
      setTodos((items) => items.filter((item) => item.id !== id))
    } catch (err) {
      setError(err instanceof Error ? err.message : '删除 Todo 失败')
    }
  }

  useEffect(() => {
    void loadTodos()
  }, [])

  return (
    <main className="min-h-screen bg-zinc-50 text-zinc-950">
      <section className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
        <header className="flex flex-col gap-3 border-b border-zinc-200 pb-5 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h1 className="text-2xl font-semibold tracking-normal">Vibe Coding Start</h1>
            <p className="mt-2 text-sm text-zinc-600">Go + Gin + GORM + Redis + React 模板</p>
          </div>
          <div className="flex gap-3 text-sm text-zinc-600">
            <span>{todos.length} 个任务</span>
            <span>{completedCount} 个完成</span>
          </div>
        </header>

        <form className="flex gap-2" onSubmit={createTodo}>
          <input
            className="h-11 flex-1 rounded-md border border-zinc-300 bg-white px-3 text-sm outline-none transition focus:border-zinc-900"
            maxLength={200}
            placeholder="添加 Todo"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
          />
          <button
            className="inline-flex h-11 items-center gap-2 rounded-md bg-zinc-950 px-4 text-sm font-medium text-white transition hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-60"
            disabled={saving || !title.trim()}
            type="submit"
          >
            {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4" />}
            新增
          </button>
        </form>

        {error && <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{error}</div>}

        <div className="overflow-hidden rounded-md border border-zinc-200 bg-white">
          {loading ? (
            <div className="flex h-48 items-center justify-center text-sm text-zinc-500">
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              正在加载
            </div>
          ) : todos.length === 0 ? (
            <div className="flex h-48 items-center justify-center text-sm text-zinc-500">暂无 Todo</div>
          ) : (
            <ul className="divide-y divide-zinc-200">
              {todos.map((todo) => (
                <li className="grid grid-cols-[auto_1fr_auto] items-center gap-3 px-4 py-3" key={todo.id}>
                  <button
                    aria-label={todo.completed ? '标记为未完成' : '标记为完成'}
                    className={`flex h-8 w-8 items-center justify-center rounded-md border transition ${
                      todo.completed ? 'border-emerald-600 bg-emerald-600 text-white' : 'border-zinc-300 text-transparent hover:border-zinc-500'
                    }`}
                    onClick={() => void toggleTodo(todo)}
                    type="button"
                  >
                    <Check className="h-4 w-4" />
                  </button>
                  <span className={`min-w-0 text-sm ${todo.completed ? 'text-zinc-400 line-through' : 'text-zinc-900'}`}>{todo.title}</span>
                  <button
                    aria-label="删除 Todo"
                    className="flex h-8 w-8 items-center justify-center rounded-md text-zinc-500 transition hover:bg-zinc-100 hover:text-red-600"
                    onClick={() => void deleteTodo(todo.id)}
                    type="button"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </section>
    </main>
  )
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${apiBaseURL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
    ...init,
  })

  const body = (await response.json()) as ApiResponse<T>
  if (!response.ok || body.code !== 0) {
    throw new Error(body.message || '请求失败')
  }
  return body.data as T
}

createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
