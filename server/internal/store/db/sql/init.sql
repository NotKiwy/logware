create table if not exists users (
    id bigint primary key,
    mode text not null,
    pub bytea not null,
    vault bytea,
    prefs jsonb not null default '{}',
    created timestamptz not null default now()
);

create table if not exists conns (
    id text primary key,
    user_id bigint not null,
    chat_id bigint not null,
    rights jsonb not null default '{}',
    active boolean not null,
    updated timestamptz not null default now()
);

create index if not exists conns_user on conns (user_id, updated desc);

create table if not exists msgs (
    id bigserial primary key,
    user_id bigint not null,
    conn_id text not null,
    chat_id bigint not null,
    msg_id bigint not null,
    ver int not null default 0,
    out boolean not null default false,
    kind text not null,
    body bytea not null,
    media bytea,
    created timestamptz not null default now(),
    unique (conn_id, chat_id, msg_id, ver)
);

create index if not exists msgs_user_created on msgs (user_id, created);

create table if not exists events (
    id text primary key,
    user_id bigint not null,
    kind text not null,
    refs bigint[] not null,
    created timestamptz not null default now()
);

create index if not exists events_user on events (user_id);

create table if not exists mods (
    user_id bigint not null,
    id text not null,
    active boolean not null default false,
    cfg jsonb not null default '{}',
    last timestamptz,
    primary key (user_id, id)
);

create table if not exists ignores (
    user_id bigint not null,
    chat_id bigint not null,
    created timestamptz not null default now(),
    primary key (user_id, chat_id)
);

create table if not exists kws (
    user_id bigint not null,
    hash bytea not null,
    box bytea not null,
    created timestamptz not null default now(),
    primary key (user_id, hash)
);

create table if not exists subs (
    user_id bigint primary key,
    plan text not null,
    until timestamptz not null,
    charge text not null,
    updated timestamptz not null default now()
);

create table if not exists pays (
    charge text primary key,
    user_id bigint not null,
    plan text not null,
    amount int not null,
    until timestamptz not null,
    created timestamptz not null default now()
);

create table if not exists devs (
    user_id bigint primary key,
    nick text not null unique,
    bio text not null default '',
    created timestamptz not null default now()
);

create table if not exists reqs (
    id bigserial primary key,
    user_id bigint not null,
    name text not null,
    about text not null default '',
    ver text not null default '',
    note text not null default '',
    fname text not null,
    file bytea not null,
    status text not null default 'wait',
    reason text not null default '',
    msg bigint not null default 0,
    created timestamptz not null default now(),
    updated timestamptz not null default now()
);

create index if not exists reqs_user on reqs (user_id, created desc);

alter table reqs add column if not exists repo text not null default '';
alter table reqs alter column file drop not null;
alter table reqs alter column fname set default '';

create table if not exists plugs (
    id bigserial primary key,
    user_id bigint not null,
    slug text not null unique,
    name text not null,
    about text not null default '',
    repo text not null default '',
    ver text not null default '',
    pic bytea,
    pv int not null default 0,
    live boolean not null default false,
    created timestamptz not null default now(),
    updated timestamptz not null default now()
);

create index if not exists plugs_user on plugs (user_id, created desc);

alter table reqs add column if not exists plug_id bigint;
alter table reqs add column if not exists kind text not null default 'new';
alter table reqs add column if not exists pic bytea;
alter table reqs add column if not exists done timestamptz;

insert into plugs (user_id, slug, name, about, repo, ver, live, created)
select r.user_id, 'p' || r.id, r.name, r.about, r.repo, r.ver, r.status = 'ok', r.created
from reqs r where r.plug_id is null
on conflict (slug) do nothing;

update reqs r set plug_id = p.id from plugs p where r.plug_id is null and p.slug = 'p' || r.id;

alter table users add column if not exists name text not null default '';
alter table users add column if not exists nick text not null default '';
alter table users add column if not exists seen timestamptz;
alter table users add column if not exists banned boolean not null default false;

alter table pays add column if not exists cur text not null default 'XTR';
alter table subs add column if not exists warned timestamptz;

alter table users add column if not exists prem boolean not null default false;

create table if not exists support (
    user_id bigint primary key,
    thread bigint not null,
    open boolean not null default true,
    updated timestamptz not null default now()
);

create table if not exists kv (
    k text primary key,
    v text not null
);

create table if not exists promos (
    code text primary key,
    plan text not null,
    days integer not null,
    cap integer not null default 0,
    used integer not null default 0,
    ends timestamptz,
    off boolean not null default false,
    created timestamptz not null default now()
);

create table if not exists promo_uses (
    code text not null references promos(code) on delete cascade,
    user_id bigint not null,
    at timestamptz not null default now(),
    primary key (code, user_id)
);

alter table conns add column if not exists name text not null default '';
alter table conns add column if not exists nick text not null default '';
