CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    telegram_id BIGINT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE group_skills (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    group_skill_id INT REFERENCES group_skills(id)
);

CREATE TABLE user_skills (
    user_id INT REFERENCES users(id),
    skill_id INT REFERENCES skills(id),
    level INT CHECK (level BETWEEN 1 AND 5),
    PRIMARY KEY (user_id, skill_id)
);

CREATE TABLE performance_reviews (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    skill_id INT REFERENCES skills(id),
    question TEXT,
    answer TEXT,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    reviewed_at TIMESTAMPTZ DEFAULT NOW()
);