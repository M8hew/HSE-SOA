CREATE TABLE UserCredentials (
    userlogin VARCHAR(50) PRIMARY KEY,
    userpassword BYTEA NOT NULL
);

CREATE TABLE UserProfile (
    userlogin VARCHAR(50) PRIMARY KEY,
    birthdate DATE,
    email VARCHAR(100),
    first_name VARCHAR(50),
    second_name VARCHAR(50),
    phone_number VARCHAR(20),
    FOREIGN KEY (userlogin) REFERENCES UserCredentials(userlogin)
);