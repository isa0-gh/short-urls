CREATE TABLE users (
	id uuid PRIMARY KEY DEFAULT uuidv7(),
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE links (
	id uuid PRIMARY KEY DEFAULT uuidv7(),
	owner_id UUID NOT NULL REFERENCES users(id),
	short_code TEXT NOT NULL UNIQUE,
	destination_url TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE telemetry (
	id uuid PRIMARY KEY DEFAULT uuidv7(),
	link_id UUID NOT NULL REFERENCES links(id),
	referrer_url TEXT,
	user_agent TEXT DEFAULT 'unknown',
	ip_address TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);