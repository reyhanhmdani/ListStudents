CREATE TABLE users
(
    id bigint NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin') NOT NULL DEFAULT 'user',
    PRIMARY KEY (id),
    UNIQUE KEY (username),
    UNIQUE KEY (email)
);