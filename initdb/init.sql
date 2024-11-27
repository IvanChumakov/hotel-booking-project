CREATE TABLE hotels (
        id UUID PRIMARY KEY,
        name VARCHAR(100) NOT NULL
);

CREATE TABLE rooms (
       id UUID PRIMARY KEY,
       price INTEGER NOT NULL,
       hotel_id UUID REFERENCES hotels(id),
       room_number INTEGER NOT NULL
);

CREATE TABLE bookings (
      id UUID PRIMARY KEY,
      hotel_name VARCHAR(100) NOT NULL,
      room_number INTEGER NOT NULL,
      from_date DATE NOT NULL,
      to_date DATE NOT NULL
);
