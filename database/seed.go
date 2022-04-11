/**
class that will seed the database with tables and initial data
*/

package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Seed(db *pgxpool.Pool) {

	createTables(db)
	createRoomData(db)
}

// create the room and reservation tables
func createTables(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), `
		DROP TABLE if exists room cascade;
		CREATE TABLE room
		(
			id SERIAL PRIMARY KEY,
			name VARCHAR(80),
			last_modified TIMESTAMP,
			created TIMESTAMP
		)

		TABLESPACE pg_default;
	`)

	if err != nil {
		panic("Error seeding database: " + err.Error())
	}

	_, err = db.Exec(context.Background(), `
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
	`)

	if err != nil {
		panic("Error seeding database: " + err.Error())
	}

	// create triggers for timestamp fields
	_, err = db.Exec(context.Background(), `
		CREATE OR REPLACE FUNCTION trigger_set_last_modified()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.last_modified = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)

	if err != nil {
		panic("Error seeding database: " + err.Error())
	}

	_, err = db.Exec(context.Background(), `
		CREATE OR REPLACE FUNCTION trigger_set_created()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.created = NOW();
			NEW.last_modified = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)

	if err != nil {
		panic("Error seeding database: " + err.Error())
	}

	_, err = db.Exec(context.Background(), `
		CREATE TRIGGER room_insert
		BEFORE INSERT ON room
		FOR EACH ROW
		EXECUTE PROCEDURE trigger_set_created();
		
		CREATE TRIGGER room_update
		BEFORE UPDATE ON room
		FOR EACH ROW
		EXECUTE PROCEDURE trigger_set_last_modified();
	`)

	if err != nil {
		panic("Error seeding database: " + err.Error())
	}

	_, err = db.Exec(context.Background(), `
		CREATE TRIGGER reservation_insert
		BEFORE INSERT ON reservation
		FOR EACH ROW
		EXECUTE PROCEDURE trigger_set_created();
		
		CREATE TRIGGER reservation_update
		BEFORE UPDATE ON reservation
		FOR EACH ROW
		EXECUTE PROCEDURE trigger_set_last_modified();
	`)

	if err != nil {
		panic("Error seeding database: " + err.Error())
	}
}

func createRoomData(db *pgxpool.Pool) {

	rooms := []string{"Meeting Room", "Conference Room", "Lunch Room"}

	for _, name := range rooms {
		_, err := db.Exec(context.Background(), `
			INSERT INTO 
				room (name)
			VALUES 
				($1)
		`, name)

		if err != nil {
			panic("Error reserving room: " + err.Error())
		}
	}
}
