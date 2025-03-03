-- Create the database
CREATE DATABASE lendmanager;

-- Connect to the database
\c lendmanager

-- Create friends table
CREATE TABLE IF NOT EXISTS friends (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Create items table
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    friend_id INTEGER REFERENCES friends(id)
);

-- Add some sample data
INSERT INTO friends (name) VALUES ('Courage'), ('Doraemon'), ('Ben');
INSERT INTO items (name, friend_id) VALUES 
('Computer', (SELECT id FROM friends WHERE name = 'Doraemon')),
('Gadget', (SELECT id FROM friends WHERE name = 'Doraemon')),
('Watch', (SELECT id FROM friends WHERE name = 'Ben'));

SELECT f.id AS friend_id, f.name AS friend_name, i.id AS item_id, i.name AS item_name
FROM friends f
LEFT JOIN items i ON f.id = i.friend_id;
