-- CREATE DATABASE dev;

-- USE dev;

CREATE TABLE accounts (
  account_id SERIAL PRIMARY KEY,
  balance DECIMAL(10,2)
);
