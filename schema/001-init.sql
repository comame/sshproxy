create database sshproxy;
use sshproxy;

create table authorized_key (
    authorized_key_id int unsigned auto_increment primary key,
    username varchar(32) not null,
    options json not null,
    pubkey text not null
);
