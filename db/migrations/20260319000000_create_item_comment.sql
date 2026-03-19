-- +goose Up
CREATE TABLE todo.item_comments (
    id uuid NOT NULL,
    item_id uuid NOT NULL,
    author_id uuid NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE ONLY todo.item_comments
    ADD CONSTRAINT item_comments_pkey PRIMARY KEY (id);

ALTER TABLE ONLY todo.item_comments
    ADD CONSTRAINT item_comments_item_id_fkey FOREIGN KEY (item_id) REFERENCES todo.items(id) ON DELETE CASCADE;

ALTER TABLE ONLY todo.item_comments
    ADD CONSTRAINT item_comments_author_id_fkey FOREIGN KEY (author_id) REFERENCES auth.users(id) ON DELETE CASCADE;

CREATE INDEX idx_item_comments_item_id ON todo.item_comments (item_id);
CREATE INDEX idx_item_comments_author_id ON todo.item_comments (author_id);

-- +goose Down
DROP TABLE IF EXISTS todo.item_comments;
