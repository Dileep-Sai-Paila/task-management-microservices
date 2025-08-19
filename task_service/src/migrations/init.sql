CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY, -- id of the task
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    user_id INT, --  id of the user this task is assigned to
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now())
);