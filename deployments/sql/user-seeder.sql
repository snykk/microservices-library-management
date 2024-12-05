INSERT INTO users (id, email, username, password, verified, role, created_at, updated_at)
VALUES
    -- Admin User
    (uuid_generate_v4(), 'najibfikri13@gmail.com', 'admin', '$2a$10$8/ltlfsE55OXuMsweGkHU.DseW9aiONTkSNCwmk6lW4Ihih4LgZbK', true, 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    
    -- Regular User
    (uuid_generate_v4(), 'najibfikri26@gmail.com', 'regularuser', '$2a$10$8/ltlfsE55OXuMsweGkHU.DseW9aiONTkSNCwmk6lW4Ihih4LgZbK', true, 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
