CREATE TABLE public.orders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    shipping character varying(255) DEFAULT 'dine in'::character varying,
    tax double precision,
    total double precision,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    user_id integer,
    payment_id integer,
    status character varying(255) DEFAULT 'pending'::character varying
);

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payments(id);

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
