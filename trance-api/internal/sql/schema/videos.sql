-- +goose Up
CREATE TABLE videos(
    user_id uuid NOT NULL,
    Video_id uuid PRIMARY KEY,
    Name TEXT NOT NULL,
    Type VARCHAR(25) NOT NULL,
    Height INT NOT NULL,
    Width INT NOT NULL,
    video_url TEXT, 
    stream_url TEXT
);

-- +goose Down
DROP TABLE videos;