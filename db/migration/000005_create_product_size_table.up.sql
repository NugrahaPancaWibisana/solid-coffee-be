CREATE TABLE public.product_size (
    id integer NOT NULL,
    name character varying(255),
    price integer,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.product_size_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.product_size_id_seq OWNED BY public.product_size.id;

ALTER TABLE ONLY public.product_size ALTER COLUMN id SET DEFAULT nextval('public.product_size_id_seq'::regclass);

ALTER TABLE ONLY public.product_size
    ADD CONSTRAINT product_size_pkey PRIMARY KEY (id);
