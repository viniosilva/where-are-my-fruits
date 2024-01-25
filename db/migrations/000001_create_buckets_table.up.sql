CREATE TABLE buckets (
    id bigint NOT NULL AUTO_INCREMENT,
    created_at datetime NOT NULL,
    deleted_at datetime,

    name varchar(128) NOT NULL,
    capacity int,

    PRIMARY KEY (ID)
);
