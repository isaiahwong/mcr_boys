-- CREATE DATABASE dev;

-- USE dev;

CREATE TABLE accounts (
  account_id SERIAL PRIMARY KEY,
  balance DECIMAL(10,2),
  owner_name VARCHAR(50)
);

CREATE TABLE transactions (
  transaction_id SERIAL PRIMARY KEY,
  amount DECIMAL(10,2),
  source_account_id INT,
  destination_account_id INT,
  timestamp TIMESTAMP
);

CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,
  username VARCHAR(50),
  password VARCHAR(50)
);

CREATE TABLE permissions (
  permission_id SERIAL PRIMARY KEY,
  user_id INT,
  account_id INT,
  permission_type VARCHAR(10)
);
