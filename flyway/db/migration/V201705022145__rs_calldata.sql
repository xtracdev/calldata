create table if not exists callrecord (
    logtime timestamp not null,
    txnid varchar(100) constraint firstkey primary key,
    error bool not null,
    host varchar(100) not null,
    category varchar(100),
    name varchar(100) not null,
    sub varchar(50),
    aud varchar(100),
    duration integer not null
);

create table if not exists svccall (
    logtime timestamp not null,
    txnid varchar(100) constraint svccall_pk primary key,
    error bool not null,
    name varchar(100) not null,
    endpoint varchar(100) not null,
    duration integer not null
);