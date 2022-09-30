create table if not exists public.trainers (
  id uuid not null default uuid_generate_v4(),
  name varchar(255) not null,
  city varchar(255) not null,
  age int not null,
  primary key(id)
);

insert into public.trainers (name,city,age) values 
   ('foo','pune',25),
   ('bar','mumbai',35),
   ('caz','delhi',47),
   ('cli','bangalore',57);