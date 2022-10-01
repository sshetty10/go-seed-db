create table if not exists public.trainers (
  id uuid not null default uuid_generate_v4(),
  name varchar(255) not null,
  city varchar(255) not null,
  age int not null,
  license_id varchar(255) not null,
  primary key(id)
);

insert into public.trainers (name,city,age,license_id) values 
   ('foo','pune',25,'VA-123456'),
   ('bar','mumbai',35,'MD-675436'),
   ('caz','delhi',47,'DE-324589'),
   ('cli','bangalore','DC-578347');