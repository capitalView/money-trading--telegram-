alter table money add column updated_at timestamp not null default now();

update money set updated_at = created_at where id > 1;