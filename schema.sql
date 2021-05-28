create table "user"
(
    id       serial primary key,
    username varchar not null,
    password varchar not null
);

create table snippet
(
    id        serial primary key,
    contents  varchar   not null,
    language  varchar   not null,
    author    int       not null,
    likes     int       not null default 0,
    dislikes  int       not null default 0,
    createdAt timestamp not null,
    constraint fk_author foreign key (author) references "user" (id)
);

create table comment
(
    id        serial primary key,
    contents  varchar   not null,
    snippet   int       not null,
    author    int       not null,
    createdAt timestamp not null,
    constraint fk_author foreign key (author) references "user" (id),
    constraint fk_snippet foreign key (snippet) references snippet (id)
);

create table vote
(
    snippet int not null,
    "user" int not null,
    vote int not null,
    primary key(snippet, "user")
);