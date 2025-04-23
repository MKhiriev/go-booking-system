INSERT INTO roles (role_id, name, description)
VALUES (1, 'Super Admin', 'Super Admin: all rights to all Entities'),
       (2, 'Content Manager', 'All rights over Rooms and Bookings'),
       (3, 'HR', 'Human Resources - all rights over Users'),
       (4, 'Event Planner', 'Schedules all bookings. Has all rights over all bookings'),
       (5, 'User', 'Regular user - create booking and change it''s own bookings');

INSERT INTO users (name, role_id, email, telephone, username, password_hash)
VALUES ('Admin', 1, 'admin@booking.app', '+992989991788', 'admin', '65794a68624763694f694a49557a49314e694973496e52356564323137613332613934626134313666383865313631323232373863434936496b705856434a390161e13f3124ae3455747b1a9ed78aa231253ae5c543cd28b9a6605835148299'),  -- Password: `AdminPass`
       ('Content Manager', 2, 'content_manager@booking.app', '+992989991789', 'content_manager', '65794a68624763694f694a49557a49314e694973496e52356564323137613332613934626134313666383865313631323232373863434936496b705856434a39bc004791aa25b3b15386c1ae71c6034d58a0cb4385287f90d1ceb45b9ce6a197'), -- Password: `ContentCop666`
       ('Cadden Jones', 3, 'jonesCadden@booking.app', '+992989991790', 'jones_cadden', '65794a68624763694f694a49557a49314e694973496e52356564323137613332613934626134313666383865313631323232373863434936496b705856434a390844eba7ecb1d12e63ef7cbde4a54cdd800a7579d1666d36ceeb088973a2d97b'), -- Password: `HRFromECorp`
       ('Event Planner', 4, 'events@booking.app', '+992989991791', 'events_admin', '65794a68624763694f694a49557a49314e694973496e52356564323137613332613934626134313666383865313631323232373863434936496b705856434a3966b532b8105c587bc6dd4a3099d0d92ebb5121cfcca6b1612735b3acc529215d'), -- Password: `VerySecurePassword`
       ('Sam Sepiol', 5, 'mr.robot@booking.app', '+992989991792', 'sam.sepiol', '65794a68624763694f694a49557a49314e694973496e52356564323137613332613934626134313666383865313631323232373863434936496b705856434a398c3a9c7c99ea7969eb77cf65438d3e7755e18974af188234576b4d2d90e0d089'); -- Password: `IAmMrRobot`

ALTER TABLE roles
    ADD COLUMN created_by BIGSERIAL NOT NULL REFERENCES users;
UPDATE roles SET created_by = 1
WHERE role_id != 0;

INSERT INTO rooms (room_id, number, capacity, created_by)
VALUES (1, 'Conference room #1', 20, 1),
       (2, 'Conference room #2', 10, 1),
       (3, 'Conference room #3', 5, 1),
       (4, 'Interrogation room', 2, 1),
       (5, 'Sauna', 8, 1);

INSERT INTO bookings (booking_id, user_id, room_id, datetime_start, datetime_end, created_by)
VALUES (1, 3, 1, '2025-04-23 13:00:00.00'::timestamp with time zone, '2025-04-23 14:00:00.00'::timestamp with time zone, 3), -- booking of room 'Conference room #1' by `HR` 13:00-14:00
       (2, 5, 1, '2025-04-23 15:00:00.00'::timestamp with time zone, '2025-04-23 16:00:00.00'::timestamp with time zone, 5), -- booking of room 'Conference room #1' by `Sam Sepiol` 15:00-16:00
       (3, 5, 5, '2025-04-23 10:00:00.00'::timestamp with time zone, '2025-04-23 13:00:00.00'::timestamp with time zone, 5); -- booking of room 'Sauna' by `Sam Sepiol` 10:00-13:00

INSERT INTO scopes (scope_id, name, description)
VALUES (1, 'ALL', 'Read/Update/Delete All records'),
       (2, 'OWNER', 'Read/Update/Delete Only owned records by user');

INSERT INTO routes (route_id, url, description, created_by)
VALUES (1, '/auth/register', 'Register new User. All unathorized users can do that.', 1),
       (2, '/auth/login', 'Log in as a registered User. All unathorized users can do that.', 1),
       (3, '/auth/refresh', 'Refresh access&refresh token.', 1),

       (4, '/user/', 'Get User by id', 1),
       (5, '/user/all', 'Get all Users', 1),
       (6, '/user/update', 'Update User by Id', 1),
       (7, '/user/drop', 'Delete User by Id', 1),

       (8, '/room/', 'Get Room by id', 1),
       (9, '/room/create', 'Create Room', 1),
       (10, '/room/all', 'Get all Rooms', 1),
       (11, '/room/update', 'Update Room by Id', 1),
       (12, '/room/drop', 'Delete Room by Id', 1),

       (13, '/booking/', 'Get Booking by id', 1),
       (14, '/booking/all', 'Get all Bookings', 1),
       (15, '/booking/room', 'Get Bookings by RoomId', 1),
       (16, '/booking/room_time', 'Get Bookings by RoomId and Booking time range', 1),
       (17, '/booking/drop', 'Delete booking by id', 1),
       (18, '/booking/available/room', 'Check if room available for Booking by RoomId and Booking time range', 1),
       (19, '/booking/create', 'Book Room (create booking)', 1),
       (20, '/booking/overlapping', 'Get overlapping Bookings', 1),
       (21, '/booking/update', 'Update Booking by Id', 1);

INSERT INTO permissions (role_id, route_id, scope_id, created_by)
VALUES
--     SUPER ADMIN
    -- AUTH
    (1, 1, 1, 1), -- Register new User
    (1, 2, 1, 1), -- Log in as a registered User
    (1, 3, 1, 1), -- Refresh access&refresh token
    -- USERS
    (1, 4, 1, 1), -- Get User by id
    (1, 5, 1, 1), -- Get all Users
    (1, 6, 1, 1), -- Update User by Id
    (1, 7, 1, 1), -- Delete User by Id
    -- ROOMS
    (1, 8, 1, 1), -- Get Room by id
    (1, 9, 1, 1), -- Create Room
    (1, 10, 1, 1), -- Get all Rooms
    (1, 11, 1, 1), -- Update Room by Id
    (1, 12, 1, 1), -- Delete Room by Id
    -- BOOKINGS
    (1, 13, 1, 1), -- Get Booking by id
    (1, 14, 1, 1), -- Get all Bookings
    (1, 15, 1, 1), -- Get Bookings by RoomId
    (1, 16, 1, 1), -- Get Bookings by RoomId and Booking time range
    (1, 17, 1, 1), -- Delete booking by id
    (1, 18, 1, 1), -- Check if room available for Booking by RoomId and Booking time range
    (1, 19, 1, 1), -- Book Room (create booking)
    (1, 20, 1, 1), -- Get overlapping Bookings
    (1, 21, 1, 1), -- Update Booking by Id

--      CONTENT MANAGER
    -- AUTH
    (2, 3, 1, 1), -- Refresh access&refresh token
    -- USERS
    (2, 4, 1, 1), -- Get User by id
    (2, 5, 1, 1), -- Get all Users
    (2, 6, 2, 1), -- Update User by Id (only himself)
    -- ROOMS
    (2, 8, 1, 1), -- Get Room by id
    (2, 9, 1, 1), -- Create Room
    (2, 10, 1, 1), -- Get all Rooms
    (2, 11, 1, 1), -- Update Room by Id
    (2, 12, 1, 1), -- Delete Room by Id
    -- BOOKINGS
    (2, 13, 1, 1), -- Get Booking by id
    (2, 14, 1, 1), -- Get all Bookings
    (2, 15, 1, 1), -- Get Bookings by RoomId
    (2, 16, 1, 1), -- Get Bookings by RoomId and Booking time range
    (2, 17, 1, 1), -- Delete booking by id
    (2, 18, 1, 1), -- Check if room available for Booking by RoomId and Booking time range
    (2, 19, 1, 1), -- Book Room (create booking)
    (2, 20, 1, 1), -- Get overlapping Bookings
    (2, 21, 1, 1), -- Update Booking by Id

--      HR - Human Resources
    -- AUTH
    (3, 1, 1, 1), -- Register new User
    (3, 3, 1, 1), -- Refresh access&refresh token
    -- USERS
    (3, 4, 1, 1), -- Get User by id
    (3, 5, 1, 1), -- Get all Users
    (3, 6, 1, 1), -- Update User by Id
    (3, 7, 1, 1), -- Delete User by Id
    -- ROOMS
    (3, 8, 1, 1), -- get room by id
    (3, 10, 1, 1), -- get all rooms
    -- BOOKINGS
    (3, 13, 1, 1), -- Get Booking by id
    (3, 15, 1, 1), -- Get Bookings by RoomId
    (3, 16, 1, 1), -- Get Bookings by RoomId and Booking time range
    (3, 17, 2, 1), -- Delete booking by id (only created by himself)
    (3, 18, 1, 1), -- Check if room available for Booking by RoomId and Booking time range
    (3, 19, 1, 1), -- Book Room (create booking)
    (3, 20, 1, 1), -- Get overlapping Bookings
    (3, 21, 2, 1), -- Update Booking by Id (only created by himself)

--     EVENT PLANNER
    -- AUTH
    (4, 3, 1, 1), -- Refresh access&refresh token
    -- USERS
    (4, 4, 1, 1), -- Get User by id
    (3, 5, 1, 1), -- Get all Users
    (4, 6, 2, 1), -- Update User by Id (only created by himself)
    -- ROOMS
    (4, 8, 1, 1), -- Get Room by id
    (4, 10, 1, 1), -- Get all Rooms
    -- BOOKINGS
    (4, 13, 1, 1), -- Get Booking by id
    (4, 14, 1, 1), -- Get all Bookings
    (4, 15, 1, 1), -- Get Bookings by RoomId
    (4, 16, 1, 1), -- Get Bookings by RoomId and Booking time range
    (4, 17, 1, 1), -- Delete booking by id
    (4, 18, 1, 1), -- Check if room available for Booking by RoomId and Booking time range
    (4, 19, 1, 1), -- Book Room (create booking)
    (4, 20, 1, 1), -- Get overlapping Bookings
    (4, 21, 1, 1), -- Update Booking by Id

--     USER
    -- AUTH
    (5, 3, 1, 1), -- Refresh access&refresh token
    -- USERS
    (5, 4, 1, 1), -- Get User by id
    (3, 5, 1, 1), -- Get all Users
    (5, 6, 2, 1), -- Update User by Id
    -- ROOMS
    (5, 8, 1, 1), -- Get Room by id
    (5, 10, 1, 1), -- Get all Rooms
    -- BOOKINGS
    (5, 13, 1, 1), -- Get Booking by id
    (5, 15, 1, 1), -- Get Bookings by RoomId
    (5, 16, 1, 1), -- Get Bookings by RoomId and Booking time range
    (5, 17, 2, 1), -- Delete booking by id  (only created by himself)
    (5, 18, 1, 1), -- Check if room available for Booking by RoomId and Booking time range
    (5, 19, 1, 1), -- Book Room (create booking)
    (5, 20, 1, 1), -- Get overlapping Bookings
    (5, 21, 2, 1); -- Update Booking by Id (only created by himself)
