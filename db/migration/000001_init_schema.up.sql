CREATE TABLE accounts (
    id integer PRIMARY KEY auto_increment,
    `owner` varchar(255) NOT NULL,
    balance bigint NOT NULL,
    currency varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp
);

CREATE TABLE entries (
    id integer auto_increment PRIMARY KEY,
    account_id integer NOT NULL ,
    amount bigint NOT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp
);

CREATE TABLE transfers (
    id integer auto_increment PRIMARY KEY,
    from_account_id integer NOT NULL ,
    to_account_id integer NOT NULL ,
    amount bigint NOT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp
);

ALTER TABLE entries ADD FOREIGN KEY (account_id) REFERENCES accounts (id);

ALTER TABLE transfers ADD FOREIGN KEY (from_account_id) REFERENCES accounts (id);

ALTER TABLE transfers ADD FOREIGN KEY (to_account_id) REFERENCES accounts (id);

CREATE INDEX index_entries_account_id ON entries (account_id);

CREATE INDEX index_transfers_from_account_id ON transfers (from_account_id);

CREATE INDEX index_transfers_to_account_id ON transfers (to_account_id);

CREATE INDEX index_transfers_from_to_account_id ON transfers (from_account_id, to_account_id);