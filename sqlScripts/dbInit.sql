-- room
----------------------------------------------------
DROP TABLE if exists room cascade;
CREATE TABLE room
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(80),
    last_modified TIMESTAMP,
    created TIMESTAMP
)

TABLESPACE pg_default;

-- reservation
----------------------------------------------------
DROP TABLE if exists reservation cascade;
CREATE TABLE reservation
(
    id SERIAL PRIMARY KEY,
    room_id INTEGER NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    expired BOOLEAN DEFAULT false,
    last_modified TIMESTAMP,
    created TIMESTAMP,
    CONSTRAINT room_id FOREIGN KEY (room_id)
        REFERENCES public.room (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;