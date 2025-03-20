CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    "group" TEXT NOT NULL,
    song TEXT NOT NULL,
    release_date TEXT NOT NULL,
    text TEXT NOT NULL,
    link TEXT NOT NULL
)