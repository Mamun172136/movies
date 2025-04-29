-- Database initialization script
-- Creates exact structure and data from your PostgreSQL dump

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

-- Create tables
CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    genre VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(512),
    release_date DATE,
    runtime INTEGER,
    mpaa_rating VARCHAR(10),
    description TEXT,
    image VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE movies_genres (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    genre_id INTEGER REFERENCES genres(id) ON DELETE CASCADE,
    UNIQUE(movie_id, genre_id)
);

-- Insert data in correct order to satisfy foreign keys

-- Genres first (no dependencies)
INSERT INTO genres (id, genre, created_at, updated_at) OVERRIDING SYSTEM VALUE VALUES
(1, 'Comedy', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(2, 'Sci-Fi', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(3, 'Horror', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(4, 'Romance', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(5, 'Action', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(6, 'Thriller', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(7, 'Drama', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(8, 'Mystery', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(9, 'Crime', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(10, 'Animation', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(11, 'Adventure', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(12, 'Fantasy', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(13, 'Superhero', '2022-09-23 00:00:00', '2022-09-23 00:00:00');

-- Movies next (no dependencies)
INSERT INTO movies (id, title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at) OVERRIDING SYSTEM VALUE VALUES
(1, 'Highlander', '1986-03-07', 116, 'R', 'He fought his first battle on the Scottish Highlands in 1536. He will fight his greatest battle on the streets of New York City in 1986. His name is Connor MacLeod. He is immortal.', '/8Z8dptJEypuLoOQro1WugD855YE.jpg', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(2, 'Raiders of the Lost Ark', '1981-06-12', 115, 'PG-13', 'Archaeology professor Indiana Jones ventures to seize a biblical artefact known as the Ark of the Covenant. While doing so, he puts up a fight against Renee and a troop of Nazis.', '/ceG9VzoRAVGwivFU403Wc3AHRys.jpg', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(3, 'The Godfather', '1972-03-24', 175, '18A', 'The aging patriarch of an organized crime dynasty in postwar New York City transfers control of his clandestine empire to his reluctant youngest son.', '/3bhkrj58Vtu7enYsRolD1fZdja1.jpg', '2022-09-23 00:00:00', '2022-09-23 00:00:00');

-- Users (no dependencies)
INSERT INTO users (id, first_name, last_name, email, password, created_at, updated_at) OVERRIDING SYSTEM VALUE VALUES
(1, 'Admin', 'User', 'admin@example.com', '$2a$14$wVsaPvJnJJsomWArouWCtusem6S/.Gauq/GjOIEHpyh2DAMmso1wy', '2022-09-23 00:00:00', '2022-09-23 00:00:00'),
(2, 'Admin', 'User', 'admin1@example.com', '$2a$14$9x6O0D2eAhnj0Ta5z1CAHOeldNLzwZuyaNA49jKBT1izH9qrpJOQ6', '2022-09-23 00:00:00', '2022-09-23 00:00:00');

-- Movie-genre relationships last (depends on movies and genres)
INSERT INTO movies_genres (id, movie_id, genre_id) OVERRIDING SYSTEM VALUE VALUES
(1, 1, 5),
(2, 1, 12),
(3, 2, 5),
(4, 2, 11),
(5, 3, 9),
(6, 3, 7);

-- Reset sequences to proper values
SELECT setval('genres_id_seq', (SELECT MAX(id) FROM genres));
SELECT setval('movies_id_seq', (SELECT MAX(id) FROM movies));
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('movies_genres_id_seq', (SELECT MAX(id) FROM movies_genres));