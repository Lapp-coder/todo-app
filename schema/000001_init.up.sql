CREATE TABLE users
(
    id            SERIAL       NOT NULL UNIQUE,
    name          VARCHAR(30)  NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE todo_lists
(
    id          SERIAL       NOT NULL UNIQUE,
    title       VARCHAR(40)  NOT NULL,
    description VARCHAR(100)
);

CREATE TABLE todo_items
(
    id          SERIAL       NOT NULL UNIQUE,
    list_id     INT          NOT NULL,
    title       VARCHAR(40)  NOT NULL,
    description VARCHAR(100),
    done        BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE TABLE users_lists
(
    user_id INT NOT NULL,
    list_id INT NOT NULL
);
