--
-- PostgreSQL database dump
--

-- Dumped from database version 12.7
-- Dumped by pg_dump version 12.7

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

--
-- Name: ipa_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ipa_users (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    uid text NOT NULL,
    reminded_at timestamp without time zone NOT NULL,
    token text NOT NULL,
    notified_at timestamp without time zone,
    notifications_sent integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.ipa_users OWNER TO postgres;

--
-- Name: requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.requests (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    email text NOT NULL,
    name text NOT NULL,
    last_name text NOT NULL,
    email_verification boolean DEFAULT false NOT NULL,
    expiry_date timestamp without time zone NOT NULL,
    token text NOT NULL,
    comment_field text
);


ALTER TABLE public.requests OWNER TO postgres;

--
-- Name: resets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resets (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    expiry_date timestamp without time zone NOT NULL,
    email text NOT NULL,
    token text NOT NULL,
    login character varying(255) NOT NULL
);


ALTER TABLE public.resets OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: ipa_users ipa_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ipa_users
    ADD CONSTRAINT ipa_users_pkey PRIMARY KEY (id);


--
-- Name: resets resets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resets
    ADD CONSTRAINT resets_pkey PRIMARY KEY (id);


--
-- Name: requests users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.requests
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

