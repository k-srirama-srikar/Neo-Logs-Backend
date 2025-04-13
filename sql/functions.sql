-- DROP TABLE users;
-- DROP FUNCTION insert_user;
-- DROP FUNCTION get_user_by_email;
-- DROP FUNCTION update_user_password;
-- DROP FUNCTION delete_user_by_email;
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



CREATE TABLE IF NOT EXISTS user_profiles (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    full_name VARCHAR(100),
    bio TEXT,
    profile_picture VARCHAR(255),
    public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS followers (
    follower_id INT REFERENCES users(id) ON DELETE CASCADE,
    following_id INT REFERENCES users(id) ON DELETE CASCADE,
    followed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    CONSTRAINT unique_follow UNIQUE (follower_id, following_id)
);

CREATE OR REPLACE FUNCTION create_user_profile()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_profiles (user_id, full_name, bio, profile_picture, public)
    VALUES (NEW.id, '', '', 'https://neologs.vercel.app/pfp1.jpg', TRUE);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER after_user_insert
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION create_user_profile();

CREATE TABLE IF NOT EXISTS blogs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    tags TEXT[],  -- Array of tags (e.g., ['tech', 'coding'])
    visibility BOOLEAN DEFAULT TRUE, -- TRUE for public, FALSE for private
    status VARCHAR(10) DEFAULT 'draft',  -- 'draft' or 'published'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    blog_id INT REFERENCES blogs(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id INT REFERENCES comments(id) ON DELETE CASCADE, -- NULL if top-level comment
    content TEXT NOT NULL,
    depth INT DEFAULT 0, -- Depth of comment (0 = top-level, 1 = reply, etc.)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- GIN index to improve tag search
CREATE INDEX IF NOT EXISTS idx_blog_tags ON blogs USING GIN(tags);