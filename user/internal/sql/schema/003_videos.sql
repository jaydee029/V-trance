-- +goose Up
CREATE TABLE videos(
    User_id uuid NOT NULL,
    Video_id uuid PRIMARY KEY,
    Name TEXT NOT NULL,
    Type VARCHAR(25) NOT NULL,
    Resolution INT NOT NULL,
    Video_url TEXT, 
    Stream_url TEXT,
    Created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE videos;