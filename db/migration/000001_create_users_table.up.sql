CREATE TABLE IF NOT EXISTS public.users (
    id integer NOT NULL,
    fullname character varying(255),
    email character varying(255),
    password character varying(255),
    photo character varying(255),
    phone character(13),
    address text,
    role character varying(20),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone,
    lastlogin_at timestamp without time zone
);

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
