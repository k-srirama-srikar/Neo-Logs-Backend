-- Create the users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Function to insert a user into the database
CREATE OR REPLACE FUNCTION insert_user(_name VARCHAR, _email VARCHAR, _password TEXT)
RETURNS VOID AS $$
BEGIN
    INSERT INTO users (name, email, password)
    VALUES (_name, _email, _password);
END;
$$ LANGUAGE plpgsql;

-- Function to retrieve a user by email
CREATE OR REPLACE FUNCTION get_user_by_email(_email VARCHAR)
RETURNS TABLE(id INT, name VARCHAR, email VARCHAR, created_at TIMESTAMP) AS $$
BEGIN
    RETURN QUERY SELECT id, name, email, created_at FROM users WHERE email = _email;
END;
$$ LANGUAGE plpgsql;

-- Function to update a user's password
CREATE OR REPLACE FUNCTION update_user_password(_email VARCHAR, _new_password TEXT)
RETURNS VOID AS $$
BEGIN
    UPDATE users
    SET password = _new_password, updated_at = CURRENT_TIMESTAMP
    WHERE email = _email;
END;
$$ LANGUAGE plpgsql;

-- Function to delete a user by email
CREATE OR REPLACE FUNCTION delete_user_by_email(_email VARCHAR)
RETURNS VOID AS $$
BEGIN
    DELETE FROM users WHERE email = _email;
END;
$$ LANGUAGE plpgsql;
