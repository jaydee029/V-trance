-- +goose Up
CREATE TABLE jobs(
    Job_id uuid PRIMARY KEY ,
    Video_id uuid NOT NULL REFERENCES videos(Video_id),
    Name TEXT NOT NULL,
    Type VARCHAR(11) NOT NULL,
    Options JSONB NOT NULL,
    Status VARCHAR(10) NOT NULL,
    Created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE jobs;
