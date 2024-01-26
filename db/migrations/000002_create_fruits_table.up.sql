CREATE TABLE fruits (
    id bigint NOT NULL AUTO_INCREMENT,
    created_at datetime NOT NULL,
    deleted_at datetime,
    
    bucket_fk bigint,

    name varchar(128) NOT NULL,
    price decimal(8,2) NOT NULL,
    expires_at datetime NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (bucket_fk) REFERENCES buckets(id)
);
