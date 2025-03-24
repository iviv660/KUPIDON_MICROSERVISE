CREATE TABLE likes (
    id SERIAL PRIMARY KEY,       -- Уникальный идентификатор для каждой записи
    from_user_id BIGINT NOT NULL, -- ID пользователя, который поставил лайк
    to_user_id BIGINT NOT NULL
);
