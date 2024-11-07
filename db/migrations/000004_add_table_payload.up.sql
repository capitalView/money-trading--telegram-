create table payload (
    id serial primary key,
    data jsonb not null,
    transaction_id INT UNIQUE NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES money (id) ON UPDATE CASCADE ON DELETE CASCADE
);