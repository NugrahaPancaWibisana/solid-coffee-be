CREATE TABLE public.dt_order (
    id integer NOT NULL,
    order_id uuid,
    qty integer,
    subtotal double precision,
    menu_id integer,
    product_size_id integer,
    product_type_id integer
);

CREATE SEQUENCE public.dt_order_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.dt_order_id_seq OWNED BY public.dt_order.id;

ALTER TABLE ONLY public.dt_order ALTER COLUMN id SET DEFAULT nextval('public.dt_order_id_seq'::regclass);

ALTER TABLE ONLY public.dt_order
    ADD CONSTRAINT dt_order_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.dt_order
    ADD CONSTRAINT dt_order_menu_id_fkey FOREIGN KEY (menu_id) REFERENCES public.menus(id);

ALTER TABLE ONLY public.dt_order
    ADD CONSTRAINT dt_order_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);

ALTER TABLE ONLY public.dt_order
    ADD CONSTRAINT dt_order_product_size_id_fkey FOREIGN KEY (product_size_id) REFERENCES public.product_size(id);

ALTER TABLE ONLY public.dt_order
    ADD CONSTRAINT dt_order_product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES public.product_type(id);
