-- init.sql
CREATE TABLE IF NOT EXISTS data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO data (name) VALUES ('Sofia');
INSERT INTO data (name) VALUES ('Laura');

