create table money (
  id serial primary key,
  amount numeric(10, 2) not null,
  currency varchar(3) not null,
  type varchar(255) not null,
  created_at timestamp not null default now(),
  name varchar not null default 'Unknown',
  description text default null
);

insert into money (amount, currency, type, name, description) values (200.00, 'USD', 'Deposit', 'Income', 'Cash in hand');