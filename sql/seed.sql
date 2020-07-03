--
-- PostgreSQL database seed
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET search_path = public, pg_catalog;

--
-- TOC entry 3357 (class 0 OID 63183)
-- Dependencies: 219
-- Data for Name: oauth_clients; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO oauth_clients (id, secret, domain, data) VALUES ('4x45RYbu21vYFHABaPjl', 'cbmDBX2LYiwvlsJntdpa', 'http://127.0.0.1:8080', '{}');

--
-- TOC entry 3336 (class 0 OID 62903)
-- Dependencies: 198
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO users (id, username, password, firstname, middlename, lastname, gender, email) VALUES ('63f284c3-1891-4905-a183-57f621aca134', 'admin', '$2a$14$TH23lPu7kA9QiRqW8SCNJOg182LKQ7okjhCThCN.ICSw9dgmBk2a2', 'admin', NULL, 'admin', 'm', 'admin@admin.com');

