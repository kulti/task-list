CREATE TABLE IF NOT EXISTS task_lists(
    id serial PRIMARY KEY,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS task_list_map(
    task_id int REFERENCES tasks (id) ON UPDATE CASCADE ON DELETE CASCADE,
    list_id int REFERENCES task_lists (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT pkey PRIMARY KEY (task_id, list_id)
);
