--
-- PostgreSQL database dump
--

\restrict nExRQZANPGLVcBQO7Bz9lSeuNb4TALvalGd8B6rEdDJhlSg4HWTrfU7IshjQ778

-- Dumped from database version 18.2
-- Dumped by pg_dump version 18.2

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: chapters; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chapters (
    chapter_id character(15) NOT NULL,
    chapter_name character varying(50),
    story_id character varying(15),
    file_id character(15),
    last_updated bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public.chapters OWNER TO postgres;

--
-- Name: edges; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.edges (
    from_chap character(15) NOT NULL,
    to_chap character(15) NOT NULL
);


ALTER TABLE public.edges OWNER TO postgres;

--
-- Name: stories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stories (
    story_id character(15) NOT NULL,
    username character varying(18),
    story_name character varying(50) NOT NULL,
    visibility smallint DEFAULT 0 NOT NULL,
    description character varying(1000),
    last_updated bigint DEFAULT 0 NOT NULL,
    chapters integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.stories OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    username character varying(18) CONSTRAINT "users_userName_not_null" NOT NULL,
    email text NOT NULL,
    user_id text CONSTRAINT "users_userID_not_null" NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: chapters; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.chapters (chapter_id, chapter_name, story_id, file_id, last_updated) FROM stdin;
RE65GUWMBGE5O75	testchapter	G2GZXOWXRW4VP6Z	MPJ5RCFGA4LV5CN	11
4QN3UOU4Y5ZYLT2	new test chapter	LRJQB4P56TIEGWB	JPUJ3WPPCPINWIP	16
Y3WL5ENJ5RDD3G7	this is the third test chapter	LRJQB4P56TIEGWB	6XWWPGRM4IBI3IC	30
4WDGYBTG4JNGP3B	this is the fourth chapter	LRJQB4P56TIEGWB	PB2WVEHP2CG74OZ	26
CRR67YN63ZW26SH	secondchapter	LRJQB4P56TIEGWB	VWO2OLE5A3S3AT7	13
IOD53JZBLNLPYPO	unlinked chapter	LRJQB4P56TIEGWB	OJPNWDTAQQGUY5I	16
IGP6YCD42KYV5YT	new chapter	A6EJQ6VMWGPRNO2	RX4EDLDKHMY3A4D	11
PDNZE7LQLNRM25V	new second chap for edge making	A6EJQ6VMWGPRNO2	CVSB3YDPR6N2KOL	31
35LXXOFSWUX22HE	<script>alert("hewo")</script>	6QYBTUGYVRYFEDY	ZHA6ZXWDBLHITVB	30
\.


--
-- Data for Name: edges; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.edges (from_chap, to_chap) FROM stdin;
4QN3UOU4Y5ZYLT2	Y3WL5ENJ5RDD3G7
Y3WL5ENJ5RDD3G7	4WDGYBTG4JNGP3B
Y3WL5ENJ5RDD3G7	CRR67YN63ZW26SH
PDNZE7LQLNRM25V	IGP6YCD42KYV5YT
IGP6YCD42KYV5YT	PDNZE7LQLNRM25V
\.


--
-- Data for Name: stories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stories (story_id, username, story_name, visibility, description, last_updated, chapters) FROM stdin;
QY7KU2WBREQZREF	loaf3000	loaf st 1	1	loaf	0	0
C2LPOQZBFQ7YR75	loaf3000	loaf st 2	2	public loaf	0	0
LRJQB4P56TIEGWB	loaf3000	testStory	2	no desc	0	0
G2GZXOWXRW4VP6Z	loaf3000	story name	2	testing	0	0
A6EJQ6VMWGPRNO2	loaf3000	newnewstory	2	no desc	0	0
6QYBTUGYVRYFEDY	loaf3000	<script>alert("hewwo")</script>	2	<script>alert("hewwo")</script>	0	0
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (username, email, user_id) FROM stdin;
loaf3000	siddarth.8202ma@gmail.com	107858227317742008043
tester	bcs_2025084@iiitm.ac.in	111981022634701517894
loaf_tester	siddarthevvs@gmail.com	107824132031111155755
user123	dummy@dummy.com	D-KXUNDKFXVBIRS4S
user1234	dummy@dummy.com	D-OLAZ6DHM6MVSPW4
\.


--
-- Name: users Uniq; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT "Uniq" UNIQUE (username);


--
-- Name: chapters chapters_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chapters
    ADD CONSTRAINT chapters_pk PRIMARY KEY (chapter_id);


--
-- Name: edges edges_fk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.edges
    ADD CONSTRAINT edges_fk PRIMARY KEY (from_chap, to_chap);


--
-- Name: users pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT pkey PRIMARY KEY (user_id);


--
-- Name: stories story_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stories
    ADD CONSTRAINT story_pk PRIMARY KEY (story_id);


--
-- Name: stories unique_story_per_user; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stories
    ADD CONSTRAINT unique_story_per_user UNIQUE (username, story_name);


--
-- Name: fki_stories_chapters_story_id_fk; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_stories_chapters_story_id_fk ON public.chapters USING btree (story_id);


--
-- Name: from_chap_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX from_chap_index ON public.edges USING btree (from_chap) WITH (deduplicate_items='true');


--
-- Name: stories_story_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX stories_story_id_index ON public.stories USING btree (story_id) WITH (deduplicate_items='false');


--
-- Name: stories_username_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX stories_username_index ON public.stories USING btree (username) WITH (deduplicate_items='true');


--
-- Name: users_user_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_user_id_index ON public.users USING btree (user_id) WITH (deduplicate_items='true');


--
-- Name: users_username_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_username_index ON public.users USING btree (username) WITH (deduplicate_items='true');


--
-- Name: edges chapters_edges_from_chapter_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.edges
    ADD CONSTRAINT chapters_edges_from_chapter_fk FOREIGN KEY (from_chap) REFERENCES public.chapters(chapter_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: edges chapters_edges_from_to_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.edges
    ADD CONSTRAINT chapters_edges_from_to_fk FOREIGN KEY (to_chap) REFERENCES public.chapters(chapter_id);


--
-- Name: chapters stories_chapters_story_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chapters
    ADD CONSTRAINT stories_chapters_story_id_fk FOREIGN KEY (story_id) REFERENCES public.stories(story_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: stories users_stories_username_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stories
    ADD CONSTRAINT users_stories_username_fk FOREIGN KEY (username) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict nExRQZANPGLVcBQO7Bz9lSeuNb4TALvalGd8B6rEdDJhlSg4HWTrfU7IshjQ778

