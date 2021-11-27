CREATE TABLE users
(
    id            SERIAL       NOT NULL  PRIMARY KEY UNIQUE,
    name          VARCHAR(30)  NOT NULL,
    email         VARCHAR(255) NOT NULL              UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE todo_lists
(
    id              SERIAL                    NOT NULL UNIQUE,
    user_id         INT REFERENCES users (id) NOT NULL,
    title           VARCHAR(40)               NOT NULL,
    description     VARCHAR(100)              NOT NULL DEFAULT '', 
    completion_date TIMESTAMP                 NOT NULL DEFAULT NOW() 
);

CREATE TABLE todo_items
(
    id              SERIAL                         NOT NULL UNIQUE,
    list_id         INT REFERENCES todo_lists (id) NOT NULL,
    title           VARCHAR(40)                    NOT NULL,
    description     VARCHAR(100)                   NOT NULL DEFAULT '',
    completion_date TIMESTAMP                      NOT NULL DEFAULT NOW(),
    done            BOOLEAN                        NOT NULL DEFAULT FALSE
);
