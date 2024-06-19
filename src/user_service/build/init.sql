CREATE TABLE UserCredentials (
    id SERIAL PRIMARY KEY,
    userlogin VARCHAR(50) UNIQUE NOT NULL,
    userpassword BYTEA NOT NULL
);

CREATE TABLE UserProfile (
    id INTEGER PRIMARY KEY REFERENCES UserCredentials,
    birthdate DATE,
    email VARCHAR(100),
    first_name VARCHAR(50),
    second_name VARCHAR(50),
    phone_number VARCHAR(20)
);