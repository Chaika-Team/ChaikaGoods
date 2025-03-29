--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2
-- Dumped by pg_dump version 17.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';

--
-- Name: clone_schema(); Type: FUNCTION; Schema: -; Owner: postgres
--

CREATE OR REPLACE FUNCTION public.clone_schema(source_schema text, dest_schema text) RETURNS void AS
$$

DECLARE
  object text;
  buffer text;
  default_ text;
  column_ text;
BEGIN
  EXECUTE 'CREATE SCHEMA ' || dest_schema ;
 
  -- TODO: Find a way to make this sequence's owner is the correct table.
  FOR object IN
    SELECT sequence_name::text FROM information_schema.SEQUENCES WHERE sequence_schema = source_schema
  LOOP
    EXECUTE 'CREATE SEQUENCE ' || dest_schema || '.' || object;
  END LOOP;
 
  FOR object IN
    SELECT table_name::text FROM information_schema.TABLES WHERE table_schema = source_schema
  LOOP
    buffer := dest_schema || '.' || object;
    EXECUTE 'CREATE TABLE ' || buffer || ' (LIKE ' || source_schema || '.' || object || ' INCLUDING CONSTRAINTS INCLUDING INDEXES INCLUDING DEFAULTS)';
   
    FOR column_, default_ IN
      SELECT column_name::text, replace(column_default::text, source_schema, dest_schema) FROM information_schema.COLUMNS where table_schema = dest_schema AND table_name = object AND column_default LIKE 'nextval(%' || source_schema || '%::regclass)'
    LOOP
      EXECUTE 'ALTER TABLE ' || buffer || ' ALTER COLUMN ' || column_ || ' SET DEFAULT ' || default_;
    END LOOP;
  END LOOP;
 
END;

$$ LANGUAGE plpgsql VOLATILE;

--
-- Name: check_single_dev_version(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.check_single_dev_version() RETURNS trigger
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN
    -- Проверяем, существует ли уже версия с is_dev = TRUE
    IF (NEW.is_dev = TRUE AND EXISTS(SELECT 1 FROM public.version WHERE is_dev = TRUE)) THEN
        RAISE EXCEPTION 'An active development version already exists.';
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.check_single_dev_version() OWNER TO postgres;

--
-- Name: prevent_non_dev_changes(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.prevent_non_dev_changes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Проверяем, является ли версия, к которой привязывается изменение, разработочной
    IF (SELECT is_dev FROM versions WHERE version_id = NEW.version_id) = FALSE THEN
        RAISE EXCEPTION 'Cannot add changes to a non-development version.';
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.prevent_non_dev_changes() OWNER TO postgres;

--
-- Name: set_default_version_id(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.set_default_version_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    current_dev_version INT;
BEGIN
    -- Проверяем, указан ли version_id
    IF NEW.version_id IS NULL THEN
        -- Получаем текущую dev версию
        SELECT version_id INTO current_dev_version FROM version WHERE is_dev = TRUE ORDER BY creation_date DESC LIMIT 1;
        IF current_dev_version IS NULL THEN
            RAISE EXCEPTION 'No active development version available.';
        END IF;
        -- Присваиваем текущую dev версию
        NEW.version_id := current_dev_version;
    ELSE
        -- Проверяем, что указанная версия является dev версией
        IF NOT EXISTS (SELECT 1 FROM version WHERE version_id = NEW.version_id AND is_dev = TRUE) THEN
            RAISE EXCEPTION 'Specified version_id is not a development version.';
        END IF;
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.set_default_version_id() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: changes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.changes (
    change_id integer NOT NULL,
    version_id integer NOT NULL,
    operation smallint NOT NULL,
    new_value jsonb,
    change_timestamp timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    considered boolean DEFAULT false
);


ALTER TABLE public.changes OWNER TO postgres;

--
-- Name: changes_change_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.changes_change_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.changes_change_id_seq OWNER TO postgres;

--
-- Name: changes_change_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.changes_change_id_seq OWNED BY public.changes.change_id;


--
-- Name: package; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.package (
    packageid integer NOT NULL,
    packagename character varying(255) NOT NULL UNIQUE,
    description text
);


ALTER TABLE public.package OWNER TO postgres;

--
-- Name: package_packageid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.package_packageid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.package_packageid_seq OWNER TO postgres;

--
-- Name: package_packageid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.package_packageid_seq OWNED BY public.package.packageid;


--
-- Name: packagecontent; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.packagecontent (
    packagecontentid integer NOT NULL,
    packageid integer NOT NULL,
    productid integer NOT NULL,
    quantity integer NOT NULL
);


ALTER TABLE public.packagecontent OWNER TO postgres;

--
-- Name: packagecontent_packagecontentid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.packagecontent_packagecontentid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.packagecontent_packagecontentid_seq OWNER TO postgres;

--
-- Name: packagecontent_packagecontentid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.packagecontent_packagecontentid_seq OWNED BY public.packagecontent.packagecontentid;


--
-- Name: product; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    price numeric(10,2) NOT NULL,
    imageurl character varying(255),
    sku character varying(100) NOT NULL UNIQUE
);


ALTER TABLE public.product OWNER TO postgres;

--
-- Name: product_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.product_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_id_seq OWNER TO postgres;

--
-- Name: product_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.product_id_seq OWNED BY public.product.id;


--
-- Name: version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.version (
    version_id integer NOT NULL,
    creation_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    is_dev boolean DEFAULT true,
    applied boolean DEFAULT false
);


ALTER TABLE public.version OWNER TO postgres;

--
-- Name: versions_version_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.versions_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.versions_version_id_seq OWNER TO postgres;

--
-- Name: versions_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.versions_version_id_seq OWNED BY public.version.version_id;


--
-- Name: changes change_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.changes ALTER COLUMN change_id SET DEFAULT nextval('public.changes_change_id_seq'::regclass);


--
-- Name: package packageid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.package ALTER COLUMN packageid SET DEFAULT nextval('public.package_packageid_seq'::regclass);


--
-- Name: packagecontent packagecontentid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.packagecontent ALTER COLUMN packagecontentid SET DEFAULT nextval('public.packagecontent_packagecontentid_seq'::regclass);


--
-- Name: product id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product ALTER COLUMN id SET DEFAULT nextval('public.product_id_seq'::regclass);


--
-- Name: version version_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.version ALTER COLUMN version_id SET DEFAULT nextval('public.versions_version_id_seq'::regclass);


--
-- Name: changes changes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.changes
    ADD CONSTRAINT changes_pkey PRIMARY KEY (change_id);


--
-- Name: package package_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.package
    ADD CONSTRAINT package_pkey PRIMARY KEY (packageid);


--
-- Name: packagecontent packagecontent_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.packagecontent
    ADD CONSTRAINT packagecontent_pkey PRIMARY KEY (packagecontentid);


--
-- Name: product product_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_pkey PRIMARY KEY (id);


--
-- Name: version versions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.version
    ADD CONSTRAINT versions_pkey PRIMARY KEY (version_id);


--
-- Name: idx_package_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_package_id ON public.packagecontent USING btree (packageid);


--
-- Name: idx_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_product_id ON public.packagecontent USING btree (productid);


--
-- Name: version trigger_prevent_multiple_dev_versions; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_prevent_multiple_dev_versions BEFORE INSERT ON public.version FOR EACH ROW EXECUTE FUNCTION public.check_single_dev_version();


--
-- Name: changes trigger_set_default_version_id; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_set_default_version_id BEFORE INSERT ON public.changes FOR EACH ROW EXECUTE FUNCTION public.set_default_version_id();


--
-- Name: changes changes_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.changes
    ADD CONSTRAINT changes_version_id_fkey FOREIGN KEY (version_id) REFERENCES public.version(version_id) ON DELETE CASCADE;


--
-- Name: packagecontent packagecontent_packageid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.packagecontent
    ADD CONSTRAINT packagecontent_packageid_fkey FOREIGN KEY (packageid) REFERENCES public.package(packageid);


--
-- Name: packagecontent packagecontent_productid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.packagecontent
    ADD CONSTRAINT packagecontent_productid_fkey FOREIGN KEY (productid) REFERENCES public.product(id);


--
-- Name: TABLE changes; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.changes TO application_user;


--
-- Name: SEQUENCE changes_change_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT USAGE ON SEQUENCE public.changes_change_id_seq TO application_user;


--
-- Name: TABLE package; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.package TO application_user;


--
-- Name: SEQUENCE package_packageid_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT USAGE ON SEQUENCE public.package_packageid_seq TO application_user;


--
-- Name: TABLE packagecontent; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.packagecontent TO application_user;


--
-- Name: SEQUENCE packagecontent_packagecontentid_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT USAGE ON SEQUENCE public.packagecontent_packagecontentid_seq TO application_user;


--
-- Name: TABLE product; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.product TO application_user;


--
-- Name: SEQUENCE product_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT USAGE ON SEQUENCE public.product_id_seq TO application_user;


--
-- Name: TABLE version; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,INSERT,UPDATE ON TABLE public.version TO application_user;


--
-- Name: SEQUENCE versions_version_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT SELECT,USAGE ON SEQUENCE public.versions_version_id_seq TO application_user;


--
-- PostgreSQL database dump complete
--

