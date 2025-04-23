CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,

    active BOOL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    role_id INT REFERENCES roles,
    email TEXT NOT NULL UNIQUE,
    telephone TEXT NOT NULL UNIQUE,

    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    active BOOL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE rooms (
    room_id SERIAL PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    capacity INT NOT NULL,

    active BOOL DEFAULT true,
    created_by BIGSERIAL NOT NULL REFERENCES users,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE bookings (
    booking_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users,
    room_id INT NOT NULL REFERENCES rooms,
    datetime_start TIMESTAMPTZ NOT NULL,
    datetime_end TIMESTAMPTZ NOT NULL CHECK ( datetime_start < datetime_end ),

    active BOOL DEFAULT true,
    created_by BIGSERIAL NOT NULL REFERENCES users,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);


CREATE TABLE scopes
(
    scope_id   BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT '',

    active BOOL DEFAULT true,
    created_by BIGSERIAL NOT NULL REFERENCES users,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE routes
(
    route_id    BIGSERIAL PRIMARY KEY,
    url         TEXT NOT NULL,
    description TEXT DEFAULT '',

    active BOOL DEFAULT true,
    created_by BIGSERIAL NOT NULL REFERENCES users,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE permissions
(
    role_id  BIGINT REFERENCES roles, -- for each role we specify chosen routes
    route_id BIGINT REFERENCES routes, -- one-to-many relationship
    scope_id INT REFERENCES scopes,

    active BOOL DEFAULT true,
    created_by BIGSERIAL NOT NULL REFERENCES users,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);