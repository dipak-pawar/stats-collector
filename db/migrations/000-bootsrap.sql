-- This table sets the foundation for all future database migrations.
create table version (
  id bigserial primary key,
  updated_at timestamp with time zone not null default current_timestamp,
  version int not null unique CONSTRAINT positive_version CHECK (version >= 0)
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE nodes (
  id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
  name VARCHAR (50),
  CONSTRAINT name_unique UNIQUE (name)
);

CREATE TABLE node_metrics (
   timestamp bigint,
   cpu_usage float,
   memory_usage float,
   node_id uuid not null references nodes(id)
);