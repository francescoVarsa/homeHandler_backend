-- DROP SEQUENCE public.foodslist_id_seq;

CREATE SEQUENCE public.foodslist_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE public.nutrition_plans_id_seq;

CREATE SEQUENCE public.nutrition_plans_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE public.users_id_seq;

CREATE SEQUENCE public.users_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id serial4 NOT NULL,
	"name" text NULL,
	last_name text NULL,
	"password" text NULL,
	email text NULL,
	created_at date NULL,
	updated_at date NULL,
	reset_request_date text NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);


-- public.nutrition_plans definition

-- Drop table

-- DROP TABLE public.nutrition_plans;

CREATE TABLE public.nutrition_plans (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	plan_name text NULL,
	created_at date NULL,
	updated_at date NULL,
	CONSTRAINT nutrition_plans_pk PRIMARY KEY (id),
	CONSTRAINT nutrition_plans_fk FOREIGN KEY (user_id) REFERENCES public.users(id)
);


-- public.foodslist definition

-- Drop table

-- DROP TABLE public.foodslist;

CREATE TABLE public.foodslist (
	id serial4 NOT NULL,
	food_name text NULL,
	plan_id int4 NULL,
	meal_type text NULL,
	day_of_the_week text NULL,
	CONSTRAINT foodslist_pk PRIMARY KEY (id),
	CONSTRAINT foodslist_fk FOREIGN KEY (plan_id) REFERENCES public.nutrition_plans(id)
);
