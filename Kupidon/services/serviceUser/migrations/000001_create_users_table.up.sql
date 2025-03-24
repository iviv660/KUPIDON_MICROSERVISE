CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  telegram_id BIGINT NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  age INT NOT NULL,
  city VARCHAR(100),
  gender VARCHAR(20),
  description TEXT,
  photo TEXT
);
