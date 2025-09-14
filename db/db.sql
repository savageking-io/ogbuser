DROP TABLE IF EXISTS platforms;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS group_permissions;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS platform_type;
DROP TYPE IF EXISTS permission_domain;

CREATE TYPE platform_type AS ENUM ('steam', 'eos', 'winstore', 'xbox', 'ps', 'web');
CREATE TYPE permission_domain AS ENUM ('own', 'party', 'guild', 'global');

CREATE TABLE users
(
	id         SERIAL PRIMARY KEY,
	username   VARCHAR(50)  NOT NULL UNIQUE,
	password   VARCHAR(255) NOT NULL,
	email      VARCHAR(100) NOT NULL UNIQUE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE platforms
(
	id               SERIAL PRIMARY KEY,
	user_id          INTEGER       NOT NULL REFERENCES users (id),
	platform_name    platform_type NOT NULL,
	platform_user_id VARCHAR(100)  NOT NULL,
	created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	deleted_at       TIMESTAMP WITH TIME ZONE,
	UNIQUE (user_id, platform_name)
);

CREATE TABLE groups
(
	id          SERIAL PRIMARY KEY,
	parent_id   INTEGER REFERENCES groups (id),
	name        VARCHAR(100) NOT NULL,
	description VARCHAR(255),
	is_special  BOOLEAN NOT NULL DEFAULT FALSE,
	created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	deleted_at  TIMESTAMP WITH TIME ZONE,
	UNIQUE (name)
);

CREATE TABLE group_members
(
	id         SERIAL PRIMARY KEY,
	group_id   INTEGER NOT NULL REFERENCES groups (id),
	user_id    INTEGER NOT NULL REFERENCES users (id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE group_permissions
(
	id         SERIAL PRIMARY KEY,
	group_id   INTEGER           NOT NULL REFERENCES groups (id),
	permission VARCHAR(100)      NOT NULL,
	read       BOOLEAN           NOT NULL DEFAULT FALSE,
	write      BOOLEAN           NOT NULL DEFAULT FALSE,
	delete     BOOLEAN           NOT NULL DEFAULT FALSE,
	domain     permission_domain NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE   DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE   DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE user_sessions
(
	id            SERIAL PRIMARY KEY,
	user_id       INTEGER       NOT NULL REFERENCES users (id),
	token         VARCHAR(255)  NOT NULL,
	platform_name platform_type NOT NULL,
	created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	deleted_at    TIMESTAMP WITH TIME ZONE,
	UNIQUE (user_id, token)
);

INSERT INTO users (username, password, email, created_at, updated_at)
VALUES ('root', '$argon2id$v=19$m=65536,t=3,p=2$dmVyeXN0cm9uZ3NhbHQ$2xQImWCDVqmTG0F9ALqoV1RSG2Y98i5Jl3hcXxathms', 'admin@localhost', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       ('jane_smith', '$argon2id$v=19$m=65536,t=3,p=2$dmVyeXN0cm9uZ3NhbHQ$tTF5B137G/sEiXKnTpCHN16j9ZOJ3ri2UPPbnIS875w', 'john.smith@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       ('alice_wonder', '$argon2id$v=19$m=65536,t=3,p=2$dmVyeXN0cm9uZ3NhbHQ$tTF5B137G/sEiXKnTpCHN16j9ZOJ3ri2UPPbnIS875w', 'alice.wonder@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample platforms for users
INSERT INTO platforms (user_id, platform_name, platform_user_id, created_at, updated_at)
VALUES (1, 'steam', 'john_steam_123', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (1, 'xbox', 'john_xbox_456', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (2, 'ps', 'jane_ps_789', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (3, 'steam', 'alice_steam_012', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample groups
INSERT INTO groups (name, parent_id, is_special, created_at, updated_at)
VALUES
	   ('Super Administrators', 1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
	   ('Players', 2, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
	   ('Moderators', LASTVAL(), FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
	   ('Administrators', LASTVAL(), FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample group memberships
INSERT INTO group_members (group_id, user_id, created_at, updated_at)
VALUES (1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample group permissions
INSERT INTO group_permissions (group_id, permission, read, write, delete, domain, created_at, updated_at)
VALUES (1, 'manage_users', true, true, true, 'global', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (1, 'manage_content', true, true, true, 'global', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (2, 'moderate_content', true, true, false, 'global', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (3, 'view_content', true, false, false, 'global', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert sample user sessions
INSERT INTO user_sessions (user_id, token, platform_name, created_at, updated_at)
VALUES (1, 'token_john_123456', 'steam', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (2, 'token_jane_789012', 'ps', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (3, 'token_alice_345678', 'steam', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);