USE clickdb;

CREATE TABLE IF NOT EXISTS events (
    eventID Int64 DEFAULT unique_id(),
    eventType String,
    userID Int64,
    eventTime DateTime,
    payload String
) ENGINE = MergeTree
ORDER BY (eventID, eventTime);
