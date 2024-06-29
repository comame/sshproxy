create database sshproxy;
use sshproxy;

create table authorized_key (
    authorized_key_id int unsigned auto_increment primary key,
    authenticated_user_id varchar(64) not null,
    username varchar(32) not null,
    options json not null,
    pubkey text not null
);
