CREATE TABLE public.menus (
    id integer NOT NULL,
    discount double precision,
    type character varying(50),
    product_id integer,
    stock integer,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.menus_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.menus_id_seq OWNED BY public.menus.id;

ALTER TABLE ONLY public.menus ALTER COLUMN id SET DEFAULT nextval('public.menus_id_seq'::regclass);

ALTER TABLE ONLY public.menus
    ADD CONSTRAINT menus_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.menus
    ADD CONSTRAINT menus_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);
