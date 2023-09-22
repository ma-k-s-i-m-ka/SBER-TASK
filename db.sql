DROP TABLE IF EXISTS Task;

CREATE TABLE IF NOT EXISTS Task (
 id              serial       primary key,
 title           text         not null,
 description     text         not null,
 date            timestamptz  not null,
 status          bool         not null
);
