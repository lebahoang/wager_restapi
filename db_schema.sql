-- public.wagers definition

-- Drop table

DROP TABLE IF EXISTS
	public.wagers,
	public.purchases
CASCADE;

CREATE TABLE public.wagers (
	id bigserial NOT NULL,
	total_wager_value integer NOT NULL,
	odds integer NOT NULL,
	selling_percentage integer NOT NULL,
	selling_price numeric(100, 2) NOT NULL,
	current_selling_price numeric(100, 2) NOT NULL,
	percentage_sold integer,
	amount_sold numeric(100, 2),
	placed_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT wagers_pkey PRIMARY KEY (id)
);

-- public.purchases definition

CREATE TABLE public.purchases (
	id bigserial NOT NULL,
	wager_id integer NOT NULL,
	buying_price numeric(100, 2) NOT NULL,
	bought_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT purchases_pkey PRIMARY KEY (id)
);


-- public.purchases foreign keys

ALTER TABLE public.purchases ADD CONSTRAINT purchases_id_fkey FOREIGN KEY (wager_id) REFERENCES public.wagers(id) ON DELETE CASCADE;