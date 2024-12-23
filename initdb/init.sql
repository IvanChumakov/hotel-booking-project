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

CREATE TABLE public.users (
      id uuid NOT NULL,
      "role" varchar NOT NULL,
      login varchar NOT NULL,
      "password" varchar NULL,
      CONSTRAINT users_pk PRIMARY KEY (id)
);

ALTER TABLE public.hotels ADD owner_id uuid NULL;

ALTER TABLE public.hotels ADD CONSTRAINT hotels_users_fk FOREIGN KEY (owner_id) REFERENCES public.users(id);
