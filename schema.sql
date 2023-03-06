create table users
(
    id         bigint primary key,
    username   text   not null,
    email      text   not null,
    password   text        default '',
    discord_id bigint not null,
    avatar_url text        default null,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create unique index users_email_unique on users (email);
create unique index users_username_unique on users (username) where username is not null;

create table users_refresh_tokens
(
    id         bigint primary key,
    user_id    bigint      not null,
    token      text        not null,
    expired_at timestamptz not null,
    issued_at  timestamptz not null,
    created_at timestamptz default now(),
    foreign key (user_id) references users (id)
);

create table home
(
    id             bigint primary key,
    name           varchar(255),
    user_id        bigint references users (id),
    background_url varchar(255),
    created_at     timestamptz default now(),
    updated_at     timestamptz default now()
);

create unique index home_name_user_id_unique on home (name, user_id);


create table home_members
(
    id         bigint primary key,
    home_id    bigint references home (id),
    user_id    bigint references users (id),
    permission bigint      default 0,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create unique index home_members_home_id_user_id_unique on home_members (home_id, user_id);

create table devices
(
    id                  bigint primary key,
    name                text not null,
    icon                text not null,
    model               text        default '',
    kind                smallint,
    home_id             bigint references home (id),
    user_id             bigint references users (id),
    token               text not null,
    ip_address          text,
    location            text,
    latest_heartbeat_at timestamptz default null,
    created_at          timestamptz default now(),
    updated_at          timestamptz default now()
);

create unique index devices_home_id_name_unique on devices (home_id, name);



create table device_ir_remotes
(
    id         bigint primary key,
    device_id  bigint references devices (id),
    home_id    bigint references home (id) default null,
    created_at timestamptz                 default now(),
    updated_at timestamptz                 default now()
);

-- 1:1 devices:device_ir_remotes
create unique index device_ir_remotes_device_id_unique on device_ir_remotes (device_id);

create table device_ir_remote_virtual_keys
(
    id         bigint primary key,
    remote_id  bigint references device_ir_remotes (id),
    name       text not null,
    kind       text not null,
    icon       text not null,
    properties jsonb       default '{}'::jsonb,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create table device_ir_remote_commands
(
    id         bigint primary key, -- command id
    remote_id  bigint references device_ir_remotes (id),
    virtual_id bigint references device_ir_remote_virtual_keys (id),
    name       text    not null,
    code       text    not null,
    frequency  integer not null,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);