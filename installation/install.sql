-- auto-generated definition
create table history
(
	id       integer
		constraint history_pk
			primary key autoincrement,
	task_key varchar not null,
	duration integer not null,
	added    varchar not null
);