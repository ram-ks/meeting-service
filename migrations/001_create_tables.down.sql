DROP INDEX IF EXISTS idx_availability_slot;
DROP INDEX IF EXISTS idx_availability_event;
DROP INDEX IF EXISTS idx_participants_email;
DROP INDEX IF EXISTS idx_participants_event;
DROP INDEX IF EXISTS idx_time_slots_event;
DROP INDEX IF EXISTS idx_events_status;
DROP INDEX IF EXISTS idx_events_organizer;

DROP TABLE IF EXISTS availability;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS time_slots;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;
