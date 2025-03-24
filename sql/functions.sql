-- DROP TABLE users;
DROP FUNCTION insert_user;
DROP FUNCTION get_user_by_email;
DROP FUNCTION update_user_password;
DROP FUNCTION delete_user_by_email;
-- Create the users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ALTER TABLE users ADD CONSTRAINT unique_name UNIQUE (name);

-- Function to insert a user into the database
CREATE OR REPLACE FUNCTION insert_user(_name VARCHAR, _email VARCHAR, _password TEXT)
RETURNS VOID AS $$
BEGIN
    INSERT INTO users (name, email, password)
    VALUES (_name, _email, _password);
EXCEPTION
    WHEN unique_violation THEN
        RAISE EXCEPTION 'Username or email already exists' USING ERRCODE = '23505';
END;
$$ LANGUAGE plpgsql;

-- Function to retrieve a user by email
CREATE OR REPLACE FUNCTION get_user_by_email(_identifier VARCHAR)
RETURNS TABLE(id INT, name VARCHAR, email VARCHAR, password TEXT) AS $$
BEGIN
    RETURN QUERY 
    SELECT u.id, u.name, u.email, u.password 
    FROM users u 
    WHERE u.email = _identifier OR u.name=_identifier;
END;
$$ LANGUAGE plpgsql;

-- Function to update a user's password
CREATE OR REPLACE FUNCTION update_user_password(_identifier VARCHAR, _new_password TEXT)
RETURNS VOID AS $$
BEGIN
    UPDATE users u
    SET password = _new_password, updated_at = CURRENT_TIMESTAMP
    WHERE u.email = _identifier OR u.name = _identifier;
END;
$$ LANGUAGE plpgsql;

-- Function to delete a user by email
CREATE OR REPLACE FUNCTION delete_user_by_email(_identifier VARCHAR)
RETURNS VOID AS $$
BEGIN
    DELETE FROM users u WHERE u.email = _identifier OR u.name = _identifier;
END;
$$ LANGUAGE plpgsql;
