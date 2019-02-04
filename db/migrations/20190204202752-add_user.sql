-- +migrate Up

CREATE TABLE public.users (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL,
  deleted_at timestamp with time zone,
  username text NOT NULL,
  email text NOT NULL,
  token text,
  token_expiry timestamp with time zone
);

-- +migrate Down

DROP TABLE public.users CASCADE;
