CREATE TABLE IF NOT EXISTS goods
(
    Id        UInt32,
    ProjectID UInt32,
    Name      String,
    Description Nullable(String),
    Priority  UInt32,
    Removed   Boolean,
    EventTime DateTime
);
