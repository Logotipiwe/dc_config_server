create table namespaces
(
    id   varchar(255) not null
        primary key,
    name varchar(255) not null
);

create table services
(
    id   varchar(255) not null
        primary key,
    name varchar(255) not null
);

create table config_entries
(
    id        varchar(255)         not null
        primary key,
    service   varchar(255)         null,
    namespace varchar(255)         null,
    is_active tinyint(1) default 1 not null,
    name      varchar(255)         not null,
    value     varchar(255)         not null,
    constraint config_entries_namespaces_id_fk
        foreign key (namespace) references namespaces (id),
    constraint config_entries_services_id_fk
        foreign key (service) references services (id)
);

