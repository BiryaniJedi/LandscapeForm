--
-- PostgreSQL database dump
--

\restrict pU3Rt1s4lQOlSM5eLAy1Mcu8pliWL2XiGnJMGEPmb5PK7s1SZnlSjlhjBF42QZ2

-- Dumped from database version 16.11 (Debian 16.11-1.pgdg13+1)
-- Dumped by pg_dump version 16.11 (Homebrew)

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

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: enforce_lawn_form(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.enforce_lawn_form() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  IF (SELECT form_type FROM forms WHERE id = NEW.form_id) <> 'lawn' THEN
    RAISE EXCEPTION 'Form % is not a lawn form', NEW.form_id;
  END IF;
  RETURN NEW;
END;
$$;


ALTER FUNCTION public.enforce_lawn_form() OWNER TO postgres;

--
-- Name: enforce_shrub_form(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.enforce_shrub_form() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  IF (SELECT form_type FROM forms WHERE id = NEW.form_id) <> 'shrub' THEN
    RAISE EXCEPTION 'Form % is not a shrub form', NEW.form_id;
  END IF;
  RETURN NEW;
END;
$$;


ALTER FUNCTION public.enforce_shrub_form() OWNER TO postgres;

--
-- Name: prevent_form_type_change(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.prevent_form_type_change() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  IF OLD.form_type IS DISTINCT FROM NEW.form_type THEN
    RAISE EXCEPTION 'form_type cannot be changed once set';
  END IF;
  RETURN NEW;
END;
$$;


ALTER FUNCTION public.prevent_form_type_change() OWNER TO postgres;

--
-- Name: set_updated_at(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;


ALTER FUNCTION public.set_updated_at() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: chemicals; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chemicals (
    id smallint NOT NULL,
    category text NOT NULL,
    brand_name text NOT NULL,
    chemical_name text NOT NULL,
    epa_reg_no text NOT NULL,
    recipe text NOT NULL,
    unit text NOT NULL,
    CONSTRAINT chemicals_category_check CHECK ((category = ANY (ARRAY['lawn'::text, 'shrub'::text])))
);


ALTER TABLE public.chemicals OWNER TO postgres;

--
-- Name: chemicals_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.chemicals_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.chemicals_id_seq OWNER TO postgres;

--
-- Name: chemicals_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.chemicals_id_seq OWNED BY public.chemicals.id;


--
-- Name: forms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.forms (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_by uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    form_type text NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    street_number text NOT NULL,
    street_name text NOT NULL,
    town text NOT NULL,
    zip_code text NOT NULL,
    home_phone text NOT NULL,
    other_phone text NOT NULL,
    call_before boolean DEFAULT false NOT NULL,
    is_holiday boolean DEFAULT false NOT NULL,
    CONSTRAINT forms_form_type_check CHECK ((form_type = ANY (ARRAY['shrub'::text, 'lawn'::text]))),
    CONSTRAINT forms_zip_code_check CHECK ((zip_code ~ '^\d{5}(-\d{4})?$'::text))
);


ALTER TABLE public.forms OWNER TO postgres;

--
-- Name: lawn_forms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.lawn_forms (
    form_id uuid NOT NULL,
    lawn_area_sq_ft integer NOT NULL,
    fert_only boolean DEFAULT false NOT NULL
);


ALTER TABLE public.lawn_forms OWNER TO postgres;

--
-- Name: notes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notes (
    id smallint NOT NULL,
    form_id uuid NOT NULL,
    note text NOT NULL
);


ALTER TABLE public.notes OWNER TO postgres;

--
-- Name: notes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notes_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.notes_id_seq OWNER TO postgres;

--
-- Name: notes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notes_id_seq OWNED BY public.notes.id;


--
-- Name: pesticide_applications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pesticide_applications (
    id smallint NOT NULL,
    form_id uuid NOT NULL,
    chem_used smallint NOT NULL,
    app_timestamp timestamp with time zone NOT NULL,
    rate text NOT NULL,
    amount_applied numeric(10,2) NOT NULL,
    location_code character varying(2) NOT NULL
);


ALTER TABLE public.pesticide_applications OWNER TO postgres;

--
-- Name: pesticide_applications_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pesticide_applications_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pesticide_applications_id_seq OWNER TO postgres;

--
-- Name: pesticide_applications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pesticide_applications_id_seq OWNED BY public.pesticide_applications.id;


--
-- Name: shrub_forms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.shrub_forms (
    form_id uuid NOT NULL,
    flea_only boolean DEFAULT false NOT NULL
);


ALTER TABLE public.shrub_forms OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    pending boolean DEFAULT true NOT NULL,
    role text DEFAULT 'employee'::text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    date_of_birth date DEFAULT '2000-01-01'::date NOT NULL,
    username text NOT NULL,
    password_hash text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: chemicals id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chemicals ALTER COLUMN id SET DEFAULT nextval('public.chemicals_id_seq'::regclass);


--
-- Name: notes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notes ALTER COLUMN id SET DEFAULT nextval('public.notes_id_seq'::regclass);


--
-- Name: pesticide_applications id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesticide_applications ALTER COLUMN id SET DEFAULT nextval('public.pesticide_applications_id_seq'::regclass);


--
-- Data for Name: chemicals; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.chemicals (id, category, brand_name, chemical_name, epa_reg_no, recipe, unit) FROM stdin;
1	lawn	This	That	was210	this over that	oz
2	shrub	THATEW	THIAHWISWDH	123	rji	lb
3	lawn	breas	lkadlkn	1234	epoadn	opkj
4	shrub	Those	THIasd	123	tomwd	lbs
5	lawn	THISAIDHIHAWOIHd	OIWHADOIHAOSIDHOWIHO	12301928309a09d	jahkhdlaHjhlksJdlkjh	OZ
6	shrub	TAhis	THiasd	1230983	kjnwkqj	jnjn
7	shrub	4	4	4	4	4
8	shrub	5	5	5	5	5
9	shrub	6	6	6	6	6
10	shrub	7	7	7	7	7
\.


--
-- Data for Name: forms; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.forms (id, created_by, created_at, form_type, updated_at, first_name, last_name, street_number, street_name, town, zip_code, home_phone, other_phone, call_before, is_holiday) FROM stdin;
b9f761e1-24d7-48ba-bf34-74a8a55367f6	6db48f55-9c24-4b21-903a-0f7a3525f15f	2026-01-29 02:45:17.97095+00	lawn	2026-01-29 02:45:17.97095+00	test	form	123	this st	this town	00000	1231231231	1231231234	f	f
8d41859c-ac52-4f0a-9e5e-63ae72be2b92	6db48f55-9c24-4b21-903a-0f7a3525f15f	2026-01-29 02:46:08.22316+00	lawn	2026-01-29 02:46:08.22316+00	Sanay	Dap	123	thiasd	that town	00291	1234567891	1231231231	f	f
5f8674d0-e9be-4d7f-9dae-893bbb470e7a	6db48f55-9c24-4b21-903a-0f7a3525f15f	2026-01-29 22:14:49.843585+00	shrub	2026-01-29 22:14:49.843585+00	e	Test	123	This street	Guacamoleville	07045	1231231234	1231231234	f	t
\.


--
-- Data for Name: lawn_forms; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.lawn_forms (form_id, lawn_area_sq_ft, fert_only) FROM stdin;
b9f761e1-24d7-48ba-bf34-74a8a55367f6	10000	t
8d41859c-ac52-4f0a-9e5e-63ae72be2b92	100	t
\.


--
-- Data for Name: notes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notes (id, form_id, note) FROM stdin;
\.


--
-- Data for Name: pesticide_applications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pesticide_applications (id, form_id, chem_used, app_timestamp, rate, amount_applied, location_code) FROM stdin;
6	b9f761e1-24d7-48ba-bf34-74a8a55367f6	1	2025-12-27 17:00:00+00	2 oz / 1000 sqft	1.00	FL
7	b9f761e1-24d7-48ba-bf34-74a8a55367f6	3	2025-12-26 01:00:00+00	3 oz / 1000 sqft	4.00	FL
8	8d41859c-ac52-4f0a-9e5e-63ae72be2b92	3	2025-11-11 17:08:00+00	3 oz / 1000 sq ft	1093.00	FA
11	5f8674d0-e9be-4d7f-9dae-893bbb470e7a	4	2031-01-01 05:01:00+00	4 units/1000 sq ft	8.20	2B
12	5f8674d0-e9be-4d7f-9dae-893bbb470e7a	2	2031-01-02 17:48:00+00	5 units/1020 sq ft	10.00	2A
\.


--
-- Data for Name: shrub_forms; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.shrub_forms (form_id, flea_only) FROM stdin;
5f8674d0-e9be-4d7f-9dae-893bbb470e7a	t
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, created_at, updated_at, pending, role, first_name, last_name, date_of_birth, username, password_hash) FROM stdin;
6db48f55-9c24-4b21-903a-0f7a3525f15f	2026-01-26 09:37:17.719802+00	2026-01-26 09:37:37.299851+00	f	admin	sanay	sanay	2004-06-27	sanay	$2a$10$SnTYsci1EWBc9vAcQKpyfeMsWtSlXm3u2SXTywxu8EDK6jvtCio3q
b671cb45-dce0-4cfa-9706-3b78f76dccfb	2026-01-26 10:12:06.807507+00	2026-01-26 10:12:21.419322+00	f	employee	r	r	2001-01-01	r	$2a$10$7d89nG/owmvQozWK34xvvOfF3xb08TSVgrkJr8xONkVjxl/KlhvTG
a6798769-66f4-429e-90fa-5cfa546a7acf	2026-01-26 09:54:08.50615+00	2026-01-27 11:58:26.165132+00	f	employee	ee	ee	1911-11-11	e	$2a$10$2QuyYYU651nuvMoB.WMAKevomt0daCKwbJq6qVERyyQ3AcPMUxeZe
97d9ea80-cf1b-48ed-a597-510e8502d85c	2026-01-27 11:53:14.880779+00	2026-01-30 00:28:21.799028+00	f	employee	wasd	wasd	2002-02-22	erer	$2a$10$YTGmJjnsgXFczqV/RN84ZuacZLbvYvWWEz3WdSu01jfyNguo1XOES
ce90240a-58cf-4f29-8e60-c84b20e11456	2026-02-05 00:36:39.932309+00	2026-02-05 00:37:56.120615+00	f	employee	test	this	2000-06-06	testUserName	$2a$10$ouQDOihE4yZA/Saa4UetdOhEekyk3Ko81PeRMvlJiP6hxxg0eHuaO
\.


--
-- Name: chemicals_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.chemicals_id_seq', 10, true);


--
-- Name: notes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notes_id_seq', 1, false);


--
-- Name: pesticide_applications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pesticide_applications_id_seq', 12, true);


--
-- Name: chemicals chemicals_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chemicals
    ADD CONSTRAINT chemicals_pkey PRIMARY KEY (id);


--
-- Name: forms forms_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forms
    ADD CONSTRAINT forms_pkey PRIMARY KEY (id);


--
-- Name: lawn_forms lawn_forms_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.lawn_forms
    ADD CONSTRAINT lawn_forms_pkey PRIMARY KEY (form_id);


--
-- Name: notes notes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notes
    ADD CONSTRAINT notes_pkey PRIMARY KEY (id, form_id);


--
-- Name: pesticide_applications pesticide_applications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesticide_applications
    ADD CONSTRAINT pesticide_applications_pkey PRIMARY KEY (id);


--
-- Name: shrub_forms shrub_forms_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shrub_forms
    ADD CONSTRAINT shrub_forms_pkey PRIMARY KEY (form_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_chemicals_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_chemicals_id ON public.chemicals USING btree (id);


--
-- Name: idx_forms_home_phone; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_home_phone ON public.forms USING btree (home_phone);


--
-- Name: idx_forms_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_name ON public.forms USING btree (first_name, last_name);


--
-- Name: idx_forms_name_lower; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_name_lower ON public.forms USING btree (lower(first_name), lower(last_name));


--
-- Name: idx_forms_street_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_street_name ON public.forms USING btree (street_name);


--
-- Name: idx_forms_street_name_lower; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_street_name_lower ON public.forms USING btree (lower(street_name));


--
-- Name: idx_forms_street_number; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_street_number ON public.forms USING btree (street_number);


--
-- Name: idx_forms_town; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_town ON public.forms USING btree (town);


--
-- Name: idx_forms_town_lower; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_town_lower ON public.forms USING btree (lower(town));


--
-- Name: idx_forms_user_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_user_created_at ON public.forms USING btree (created_by, created_at DESC);


--
-- Name: idx_forms_zip_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_forms_zip_code ON public.forms USING btree (zip_code);


--
-- Name: idx_pesticide_applications_app_timestamp; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pesticide_applications_app_timestamp ON public.pesticide_applications USING btree (app_timestamp);


--
-- Name: forms trg_forms_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trg_forms_updated BEFORE UPDATE ON public.forms FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: lawn_forms trg_lawn_type_check; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trg_lawn_type_check BEFORE INSERT OR UPDATE ON public.lawn_forms FOR EACH ROW EXECUTE FUNCTION public.enforce_lawn_form();


--
-- Name: forms trg_prevent_form_type_change; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trg_prevent_form_type_change BEFORE UPDATE ON public.forms FOR EACH ROW EXECUTE FUNCTION public.prevent_form_type_change();


--
-- Name: shrub_forms trg_shrub_type_check; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trg_shrub_type_check BEFORE INSERT OR UPDATE ON public.shrub_forms FOR EACH ROW EXECUTE FUNCTION public.enforce_shrub_form();


--
-- Name: users trg_users_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trg_users_updated BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: forms forms_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forms
    ADD CONSTRAINT forms_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: lawn_forms lawn_forms_form_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.lawn_forms
    ADD CONSTRAINT lawn_forms_form_id_fkey FOREIGN KEY (form_id) REFERENCES public.forms(id) ON DELETE CASCADE;


--
-- Name: notes notes_form_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notes
    ADD CONSTRAINT notes_form_id_fkey FOREIGN KEY (form_id) REFERENCES public.forms(id) ON DELETE CASCADE;


--
-- Name: pesticide_applications pesticide_applications_chem_used_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesticide_applications
    ADD CONSTRAINT pesticide_applications_chem_used_fkey FOREIGN KEY (chem_used) REFERENCES public.chemicals(id);


--
-- Name: pesticide_applications pesticide_applications_form_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesticide_applications
    ADD CONSTRAINT pesticide_applications_form_id_fkey FOREIGN KEY (form_id) REFERENCES public.forms(id) ON DELETE CASCADE;


--
-- Name: shrub_forms shrub_forms_form_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shrub_forms
    ADD CONSTRAINT shrub_forms_form_id_fkey FOREIGN KEY (form_id) REFERENCES public.forms(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict pU3Rt1s4lQOlSM5eLAy1Mcu8pliWL2XiGnJMGEPmb5PK7s1SZnlSjlhjBF42QZ2

