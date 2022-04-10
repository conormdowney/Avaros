CREATE OR REPLACE FUNCTION trigger_set_last_modified()
RETURNS TRIGGER AS $$
BEGIN
	NEW.last_modified = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Sets created date AND last modified date
CREATE OR REPLACE FUNCTION trigger_set_created()
RETURNS TRIGGER AS $$
BEGIN
	NEW.created = NOW();
	NEW.last_modified = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- room
-- ===========================================
CREATE TRIGGER room_insert
BEFORE INSERT ON room
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_created();

CREATE TRIGGER room_update
BEFORE UPDATE ON room
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_last_modified();

-- reservation
-- ===========================================
CREATE TRIGGER reservation_insert
BEFORE INSERT ON reservation
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_created();

CREATE TRIGGER reservation_update
BEFORE UPDATE ON reservation
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_last_modified();
