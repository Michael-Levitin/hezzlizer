DROP TABLE IF EXISTS goods;
DROP TABLE IF EXISTS projects;

CREATE TABLE projects
(
    id         serial primary key,
    name       varchar(30) not null,
    created_at timestamp   not null
);

comment on column projects.id is 'id записи';
comment on column projects.name is 'название';
comment on column projects.created_at is 'дата и время';

CREATE TABLE goods
(
    id          serial,
    project_id  integer,
    name        varchar(60) not null,
    description varchar(60),
    priority    integer,
    removed     bool,
    created_at  timestamp
);

ALTER TABLE goods
    ADD CONSTRAINT pk_composite PRIMARY KEY (id, project_id);
ALTER TABLE goods
    ADD CONSTRAINT project_id
        foreign key (project_id) references projects (id) match full;

CREATE INDEX goods_name ON goods(name);

comment on column goods.id is 'id записи';
comment on column goods.project_id is 'id компании';
comment on column goods.name is 'название';
comment on column goods.description is 'описание';
comment on column goods.priority is 'приоритет';
comment on column goods.removed is 'статус удаления';
comment on column goods.created_at is 'дата и время';