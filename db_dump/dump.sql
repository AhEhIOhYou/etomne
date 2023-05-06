--
-- PostgreSQL database dump
--

-- Dumped from database version 14.7 (Ubuntu 14.7-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.7 (Ubuntu 14.7-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;


\connect model3d

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;


CREATE TABLE public.file (
    id integer NOT NULL,
    title character varying NOT NULL,
    owner_id integer NOT NULL,
    url text,
    extension character varying NOT NULL,
    created_at timestamp without time zone
);


ALTER TABLE public.file OWNER TO ahehiohyou;

CREATE SEQUENCE public.file_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.file_id_seq OWNER TO ahehiohyou;

ALTER SEQUENCE public.file_id_seq OWNED BY public.file.id;

CREATE TABLE public.model (
    id integer NOT NULL,
    user_id integer NOT NULL,
    title character varying NOT NULL,
    description text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


ALTER TABLE public.model OWNER TO ahehiohyou;

CREATE TABLE public.model_files (
    model_id integer NOT NULL,
    file_id integer NOT NULL
);


ALTER TABLE public.model_files OWNER TO ahehiohyou;

CREATE SEQUENCE public.model_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.model_id_seq OWNER TO ahehiohyou;

ALTER SEQUENCE public.model_id_seq OWNED BY public.model.id;

CREATE TABLE public."user" (
    id integer NOT NULL,
    name character varying NOT NULL,
    email character varying NOT NULL,
    password text NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    is_admin boolean DEFAULT false NOT NULL
);


ALTER TABLE public."user" OWNER TO ahehiohyou;

CREATE SEQUENCE public.user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_id_seq OWNER TO ahehiohyou;

ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;

ALTER TABLE ONLY public.file ALTER COLUMN id SET DEFAULT nextval('public.file_id_seq'::regclass);

ALTER TABLE ONLY public.model ALTER COLUMN id SET DEFAULT nextval('public.model_id_seq'::regclass);

ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);

COPY public."user" (id, name, email, password, created_at, updated_at, is_admin) FROM stdin;
1	Admin	test@test.ru	$2a$10$Ig6Sl.aG4UOSqgBOQv/zpONNdh/Z82nakzvxK4myuX4y7BKF3dHqO	2023-05-04 19:46:45.251149	2023-05-05 16:30:32.752427	t
\.

SELECT pg_catalog.setval('public.file_id_seq', 1, true);

SELECT pg_catalog.setval('public.model_id_seq', 1, true);

SELECT pg_catalog.setval('public.user_id_seq', 2, true);

ALTER TABLE ONLY public.file
    ADD CONSTRAINT file_id_pk PRIMARY KEY (id);


ALTER TABLE ONLY public.model
    ADD CONSTRAINT model_id_pk PRIMARY KEY (id);

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_email_pk UNIQUE (email);

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_id_pk PRIMARY KEY (id);

ALTER TABLE ONLY public.file
    ADD CONSTRAINT file_user_id_fk FOREIGN KEY (owner_id) REFERENCES public."user"(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.model_files
    ADD CONSTRAINT model_files_file_id_fk FOREIGN KEY (file_id) REFERENCES public.file(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.model_files
    ADD CONSTRAINT model_files_model_id_fk FOREIGN KEY (model_id) REFERENCES public.model(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.model
    ADD CONSTRAINT model_user_id_fk FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;
