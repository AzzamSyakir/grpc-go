create table IF NOT EXISTS users
(
    id uuid not null primary key,
    name text not null unique ,
    password TEXT not null  ,
    email TEXT not null unique ,
    Balance int not null ,
    created_at timestamptz not null ,
    updated_at timestamptz not null ,
    deleted_at timestamptz
);

create table IF NOT EXISTS categories
(
    id uuid not null primary key,
    name text not null unique ,
    created_at timestamptz not null ,
    updated_at timestamptz not null ,
    deleted_at timestamptz
);


create table IF NOT EXISTS products
(
    id uuid not null primary key,
    Sku text not null  ,
    Name text not null unique ,
    stock int not null ,
    price int not null ,
    category_id uuid not null ,
    created_at timestamptz not null ,
    updated_at timestamptz not null ,
    deleted_at timestamptz,
    foreign key(category_id) references categories(id)
);

