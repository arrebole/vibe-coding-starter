CREATE TABLE IF NOT EXISTS public.todos (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

COMMENT ON TABLE public.todos IS 'Todo 任务表';
COMMENT ON COLUMN public.todos.id IS '主键 ID';
COMMENT ON COLUMN public.todos.title IS 'Todo 标题';
COMMENT ON COLUMN public.todos.completed IS 'Todo 是否已完成';
COMMENT ON COLUMN public.todos.created_at IS '创建时间';
COMMENT ON COLUMN public.todos.updated_at IS '更新时间';
COMMENT ON COLUMN public.todos.deleted_at IS '软删除时间，NULL 表示未删除';

CREATE INDEX IF NOT EXISTS idx_todos_title ON public.todos (title);
CREATE INDEX IF NOT EXISTS idx_todos_deleted_at ON public.todos (deleted_at);
