-- +migrate Up

CREATE TABLE public.images (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL,
  deleted_at timestamp with time zone,
  user_id uuid NOT NULL
  href text NOT NULL,
  shortlink text
);

ALTER TABLE ONLY public.images
      ADD CONSTRAINT images_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);

-- +migrate Down

DROP TABLE public.images CASCADE;
