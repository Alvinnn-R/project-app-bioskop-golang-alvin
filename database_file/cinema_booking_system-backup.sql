--
-- PostgreSQL database dump
--

\restrict WtVSngrAcOh3GUOHwqhaMe1ppJCyg99n9MnQVPgszQdFFnSiaYYQO7jtccvYhJz

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

-- Started on 2026-01-14 09:26:50

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
-- TOC entry 2 (class 3079 OID 27657)
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- TOC entry 5190 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 237 (class 1259 OID 27591)
-- Name: booking_seats; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.booking_seats (
    id integer NOT NULL,
    booking_id integer NOT NULL,
    seat_id integer NOT NULL,
    price_snapshot numeric(12,2) NOT NULL
);


ALTER TABLE public.booking_seats OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 27590)
-- Name: booking_seats_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.booking_seats_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.booking_seats_id_seq OWNER TO postgres;

--
-- TOC entry 5191 (class 0 OID 0)
-- Dependencies: 236
-- Name: booking_seats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.booking_seats_id_seq OWNED BY public.booking_seats.id;


--
-- TOC entry 235 (class 1259 OID 27567)
-- Name: bookings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bookings (
    id integer NOT NULL,
    user_id integer NOT NULL,
    showtime_id integer NOT NULL,
    status character varying(20) NOT NULL,
    total_amount numeric(12,2) DEFAULT 0,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.bookings OWNER TO postgres;

--
-- TOC entry 234 (class 1259 OID 27566)
-- Name: bookings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.bookings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.bookings_id_seq OWNER TO postgres;

--
-- TOC entry 5192 (class 0 OID 0)
-- Dependencies: 234
-- Name: bookings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.bookings_id_seq OWNED BY public.bookings.id;


--
-- TOC entry 225 (class 1259 OID 27475)
-- Name: cinemas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cinemas (
    id integer NOT NULL,
    name character varying(150) NOT NULL,
    location character varying(255),
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.cinemas OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 27474)
-- Name: cinemas_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cinemas_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cinemas_id_seq OWNER TO postgres;

--
-- TOC entry 5193 (class 0 OID 0)
-- Dependencies: 224
-- Name: cinemas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cinemas_id_seq OWNED BY public.cinemas.id;


--
-- TOC entry 231 (class 1259 OID 27517)
-- Name: movies; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.movies (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    poster_url text,
    genres text[] NOT NULL,
    rating numeric(2,1) DEFAULT 0,
    review_count integer DEFAULT 0,
    release_date date,
    duration_in_minutes integer NOT NULL,
    release_status character varying(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.movies OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 27516)
-- Name: movies_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.movies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.movies_id_seq OWNER TO postgres;

--
-- TOC entry 5194 (class 0 OID 0)
-- Dependencies: 230
-- Name: movies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.movies_id_seq OWNED BY public.movies.id;


--
-- TOC entry 239 (class 1259 OID 27613)
-- Name: payment_methods; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_methods (
    id integer NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.payment_methods OWNER TO postgres;

--
-- TOC entry 238 (class 1259 OID 27612)
-- Name: payment_methods_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.payment_methods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payment_methods_id_seq OWNER TO postgres;

--
-- TOC entry 5195 (class 0 OID 0)
-- Dependencies: 238
-- Name: payment_methods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.payment_methods_id_seq OWNED BY public.payment_methods.id;


--
-- TOC entry 241 (class 1259 OID 27624)
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    id integer NOT NULL,
    booking_id integer NOT NULL,
    payment_method_id integer NOT NULL,
    status character varying(20) NOT NULL,
    payment_details jsonb,
    paid_at timestamp with time zone
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- TOC entry 240 (class 1259 OID 27623)
-- Name: payments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.payments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payments_id_seq OWNER TO postgres;

--
-- TOC entry 5196 (class 0 OID 0)
-- Dependencies: 240
-- Name: payments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.payments_id_seq OWNED BY public.payments.id;


--
-- TOC entry 229 (class 1259 OID 27501)
-- Name: seats; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.seats (
    id integer NOT NULL,
    studio_id integer NOT NULL,
    seat_code character varying(10) NOT NULL
);


ALTER TABLE public.seats OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 27500)
-- Name: seats_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.seats_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seats_id_seq OWNER TO postgres;

--
-- TOC entry 5197 (class 0 OID 0)
-- Dependencies: 228
-- Name: seats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.seats_id_seq OWNED BY public.seats.id;


--
-- TOC entry 223 (class 1259 OID 27455)
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token uuid NOT NULL,
    expired_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 27454)
-- Name: sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sessions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sessions_id_seq OWNER TO postgres;

--
-- TOC entry 5198 (class 0 OID 0)
-- Dependencies: 222
-- Name: sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sessions_id_seq OWNED BY public.sessions.id;


--
-- TOC entry 233 (class 1259 OID 27537)
-- Name: showtimes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.showtimes (
    id integer NOT NULL,
    cinema_id integer NOT NULL,
    studio_id integer NOT NULL,
    movie_id integer NOT NULL,
    show_date date NOT NULL,
    show_time time without time zone NOT NULL,
    price numeric(12,2) NOT NULL
);


ALTER TABLE public.showtimes OWNER TO postgres;

--
-- TOC entry 232 (class 1259 OID 27536)
-- Name: showtimes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.showtimes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.showtimes_id_seq OWNER TO postgres;

--
-- TOC entry 5199 (class 0 OID 0)
-- Dependencies: 232
-- Name: showtimes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.showtimes_id_seq OWNED BY public.showtimes.id;


--
-- TOC entry 227 (class 1259 OID 27485)
-- Name: studios; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.studios (
    id integer NOT NULL,
    cinema_id integer NOT NULL,
    name character varying(100) NOT NULL,
    total_seats integer NOT NULL
);


ALTER TABLE public.studios OWNER TO postgres;

--
-- TOC entry 226 (class 1259 OID 27484)
-- Name: studios_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.studios_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.studios_id_seq OWNER TO postgres;

--
-- TOC entry 5200 (class 0 OID 0)
-- Dependencies: 226
-- Name: studios_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.studios_id_seq OWNED BY public.studios.id;


--
-- TOC entry 221 (class 1259 OID 27435)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    password_hash text NOT NULL,
    is_verified boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 27434)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 5201 (class 0 OID 0)
-- Dependencies: 220
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 4963 (class 2604 OID 27594)
-- Name: booking_seats id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.booking_seats ALTER COLUMN id SET DEFAULT nextval('public.booking_seats_id_seq'::regclass);


--
-- TOC entry 4960 (class 2604 OID 27570)
-- Name: bookings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings ALTER COLUMN id SET DEFAULT nextval('public.bookings_id_seq'::regclass);


--
-- TOC entry 4950 (class 2604 OID 27478)
-- Name: cinemas id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cinemas ALTER COLUMN id SET DEFAULT nextval('public.cinemas_id_seq'::regclass);


--
-- TOC entry 4954 (class 2604 OID 27520)
-- Name: movies id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.movies ALTER COLUMN id SET DEFAULT nextval('public.movies_id_seq'::regclass);


--
-- TOC entry 4964 (class 2604 OID 27616)
-- Name: payment_methods id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_methods ALTER COLUMN id SET DEFAULT nextval('public.payment_methods_id_seq'::regclass);


--
-- TOC entry 4965 (class 2604 OID 27627)
-- Name: payments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments ALTER COLUMN id SET DEFAULT nextval('public.payments_id_seq'::regclass);


--
-- TOC entry 4953 (class 2604 OID 27504)
-- Name: seats id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.seats ALTER COLUMN id SET DEFAULT nextval('public.seats_id_seq'::regclass);


--
-- TOC entry 4948 (class 2604 OID 27458)
-- Name: sessions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions ALTER COLUMN id SET DEFAULT nextval('public.sessions_id_seq'::regclass);


--
-- TOC entry 4959 (class 2604 OID 27540)
-- Name: showtimes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.showtimes ALTER COLUMN id SET DEFAULT nextval('public.showtimes_id_seq'::regclass);


--
-- TOC entry 4952 (class 2604 OID 27488)
-- Name: studios id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.studios ALTER COLUMN id SET DEFAULT nextval('public.studios_id_seq'::regclass);


--
-- TOC entry 4944 (class 2604 OID 27438)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 5180 (class 0 OID 27591)
-- Dependencies: 237
-- Data for Name: booking_seats; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.booking_seats (id, booking_id, seat_id, price_snapshot) FROM stdin;
1	1	1	105000.00
2	1	2	105000.00
\.


--
-- TOC entry 5178 (class 0 OID 27567)
-- Dependencies: 235
-- Data for Name: bookings; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.bookings (id, user_id, showtime_id, status, total_amount, created_at) FROM stdin;
1	1	1	CONFIRMED	210000.00	2026-01-14 09:02:52.365339+07
\.


--
-- TOC entry 5168 (class 0 OID 27475)
-- Dependencies: 225
-- Data for Name: cinemas; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cinemas (id, name, location, created_at) FROM stdin;
1	Cinema XXI	Pakuwon Mall Surabaya	2026-01-14 09:02:14.332884+07
\.


--
-- TOC entry 5174 (class 0 OID 27517)
-- Dependencies: 231
-- Data for Name: movies; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.movies (id, title, poster_url, genres, rating, review_count, release_date, duration_in_minutes, release_status, created_at, updated_at) FROM stdin;
1	Avengers: Infinity War	public/images/avengers_infinity_war.png	{Action,Adventure,Sci-Fi}	4.8	1822	2018-04-27	149	now_playing	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
2	Shang-Chi: Legend of the Ten Rings	public/images/shang_chi.png	{Action,Fantasy}	4.7	1240	2021-09-03	132	now_playing	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
3	Batman v Superman: Dawn of Justice	public/images/batman_vs_superman.png	{Action,Drama}	4.2	980	2016-03-25	151	now_playing	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
4	Guardians of the Galaxy	public/images/guardians_galaxy.png	{Action,Adventure,Sci-Fi}	4.6	1560	2014-08-01	121	now_playing	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
5	Doctor Strange	public/images/doctor_strange.png	{Action,Fantasy}	4.5	1340	2016-11-04	115	now_playing	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
9	Aquaman and the Lost Kingdom	public/images/aquaman_lost_kingdom.png	{Action,Adventure}	4.3	540	2023-12-20	124	coming_soon	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
10	The Marvels	public/images/the_marvels.png	{Action,Sci-Fi}	4.2	430	2023-11-10	105	coming_soon	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
6	Avatar: The Way of Water	public/images/avatar_way_of_water.png	{Adventure,Sci-Fi}	4.6	980	2022-12-16	192	now_playing	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
7	Ant-Man and The Wasp: Quantumania	public/images/antman_quantumania.png	{Action,Comedy,Sci-Fi}	4.1	720	2023-02-17	125	now_playing	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
8	Shazam! Fury of the Gods	public/images/shazam_fury.png	{Action,Comedy}	4.0	610	2023-03-17	130	now_playing	2026-01-14 09:02:26.796911+07	2026-01-14 09:02:26.796911+07
\.


--
-- TOC entry 5182 (class 0 OID 27613)
-- Dependencies: 239
-- Data for Name: payment_methods; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payment_methods (id, name) FROM stdin;
1	QRIS
2	Dana
3	ShopeePay
4	GoPay
5	ATM
6	Visa / Mastercard
\.


--
-- TOC entry 5184 (class 0 OID 27624)
-- Dependencies: 241
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payments (id, booking_id, payment_method_id, status, payment_details, paid_at) FROM stdin;
1	1	1	PAID	{"provider": "QRIS", "transaction_id": "ZP-20240301"}	2026-01-14 09:09:02.869353+07
\.


--
-- TOC entry 5172 (class 0 OID 27501)
-- Dependencies: 229
-- Data for Name: seats; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.seats (id, studio_id, seat_code) FROM stdin;
1	1	A1
2	1	A2
3	1	A3
4	1	A4
5	1	A5
6	1	A6
7	1	A7
8	1	A8
9	1	A9
10	1	A10
11	1	B1
12	1	B2
13	1	B3
14	1	B4
15	1	B5
16	1	B6
17	1	B7
18	1	B8
19	1	B9
20	1	B10
21	1	C1
22	1	C2
23	1	C3
24	1	C4
25	1	C5
26	1	C6
27	1	C7
28	1	C8
29	1	C9
30	1	C10
31	1	D1
32	1	D2
33	1	D3
34	1	D4
35	1	D5
36	1	D6
37	1	D7
38	1	D8
39	1	D9
40	1	D10
41	1	E1
42	1	E2
43	1	E3
44	1	E4
45	1	E5
46	1	E6
47	1	E7
48	1	E8
49	1	E9
50	1	E10
\.


--
-- TOC entry 5166 (class 0 OID 27455)
-- Dependencies: 223
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (id, user_id, token, expired_at, revoked_at, created_at) FROM stdin;
1	1	69975091-66fb-40c9-bf53-f38e62351052	2026-01-15 09:09:12.110301+07	\N	2026-01-14 09:09:12.110301+07
\.


--
-- TOC entry 5176 (class 0 OID 27537)
-- Dependencies: 233
-- Data for Name: showtimes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.showtimes (id, cinema_id, studio_id, movie_id, show_date, show_time, price) FROM stdin;
1	1	1	1	2026-01-14	10:00:00	105000.00
2	1	1	2	2026-01-14	12:30:00	105000.00
3	1	1	3	2026-01-14	15:00:00	105000.00
4	1	1	4	2026-01-14	17:30:00	105000.00
5	1	1	5	2026-01-14	20:00:00	105000.00
6	1	1	6	2026-01-14	22:30:00	105000.00
7	1	1	7	2026-01-14	01:00:00	105000.00
8	1	1	8	2026-01-14	03:30:00	105000.00
9	1	1	9	2026-01-14	06:00:00	105000.00
10	1	1	10	2026-01-14	08:00:00	105000.00
\.


--
-- TOC entry 5170 (class 0 OID 27485)
-- Dependencies: 227
-- Data for Name: studios; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.studios (id, cinema_id, name, total_seats) FROM stdin;
1	1	Studio 1	50
\.


--
-- TOC entry 5164 (class 0 OID 27435)
-- Dependencies: 221
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password_hash, is_verified, created_at, updated_at) FROM stdin;
1	angelina	angelina@mail.com	$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm	t	2026-01-14 09:02:09.141113+07	2026-01-14 09:02:09.141113+07
2	alvin	alvin@mail.com	$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm	t	2026-01-14 09:02:09.141113+07	2026-01-14 09:02:09.141113+07
3	arya	arya@mail.com	$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm	t	2026-01-14 09:02:09.141113+07	2026-01-14 09:02:09.141113+07
4	dwiyanti	dwiyanti@mail.com	$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm	t	2026-01-14 09:02:09.141113+07	2026-01-14 09:02:09.141113+07
\.


--
-- TOC entry 5202 (class 0 OID 0)
-- Dependencies: 236
-- Name: booking_seats_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.booking_seats_id_seq', 4, true);


--
-- TOC entry 5203 (class 0 OID 0)
-- Dependencies: 234
-- Name: bookings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.bookings_id_seq', 3, true);


--
-- TOC entry 5204 (class 0 OID 0)
-- Dependencies: 224
-- Name: cinemas_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.cinemas_id_seq', 1, true);


--
-- TOC entry 5205 (class 0 OID 0)
-- Dependencies: 230
-- Name: movies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.movies_id_seq', 10, true);


--
-- TOC entry 5206 (class 0 OID 0)
-- Dependencies: 238
-- Name: payment_methods_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.payment_methods_id_seq', 6, true);


--
-- TOC entry 5207 (class 0 OID 0)
-- Dependencies: 240
-- Name: payments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.payments_id_seq', 1, true);


--
-- TOC entry 5208 (class 0 OID 0)
-- Dependencies: 228
-- Name: seats_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.seats_id_seq', 50, true);


--
-- TOC entry 5209 (class 0 OID 0)
-- Dependencies: 222
-- Name: sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sessions_id_seq', 1, true);


--
-- TOC entry 5210 (class 0 OID 0)
-- Dependencies: 232
-- Name: showtimes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.showtimes_id_seq', 10, true);


--
-- TOC entry 5211 (class 0 OID 0)
-- Dependencies: 226
-- Name: studios_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.studios_id_seq', 1, true);


--
-- TOC entry 5212 (class 0 OID 0)
-- Dependencies: 220
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 4, true);


--
-- TOC entry 4995 (class 2606 OID 27600)
-- Name: booking_seats booking_seats_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.booking_seats
    ADD CONSTRAINT booking_seats_pkey PRIMARY KEY (id);


--
-- TOC entry 4992 (class 2606 OID 27578)
-- Name: bookings bookings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (id);


--
-- TOC entry 4978 (class 2606 OID 27483)
-- Name: cinemas cinemas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cinemas
    ADD CONSTRAINT cinemas_pkey PRIMARY KEY (id);


--
-- TOC entry 4987 (class 2606 OID 27533)
-- Name: movies movies_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_pkey PRIMARY KEY (id);


--
-- TOC entry 4998 (class 2606 OID 27622)
-- Name: payment_methods payment_methods_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_name_key UNIQUE (name);


--
-- TOC entry 5000 (class 2606 OID 27620)
-- Name: payment_methods payment_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_pkey PRIMARY KEY (id);


--
-- TOC entry 5003 (class 2606 OID 27635)
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- TOC entry 4982 (class 2606 OID 27509)
-- Name: seats seats_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.seats
    ADD CONSTRAINT seats_pkey PRIMARY KEY (id);


--
-- TOC entry 4974 (class 2606 OID 27465)
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- TOC entry 4976 (class 2606 OID 27467)
-- Name: sessions sessions_token_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_token_key UNIQUE (token);


--
-- TOC entry 4990 (class 2606 OID 27549)
-- Name: showtimes showtimes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.showtimes
    ADD CONSTRAINT showtimes_pkey PRIMARY KEY (id);


--
-- TOC entry 4980 (class 2606 OID 27494)
-- Name: studios studios_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.studios
    ADD CONSTRAINT studios_pkey PRIMARY KEY (id);


--
-- TOC entry 4967 (class 2606 OID 27453)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 4969 (class 2606 OID 27449)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4971 (class 2606 OID 27451)
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- TOC entry 4993 (class 1259 OID 27589)
-- Name: idx_bookings_user; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_bookings_user ON public.bookings USING btree (user_id);


--
-- TOC entry 4984 (class 1259 OID 27534)
-- Name: idx_movies_genres; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_movies_genres ON public.movies USING gin (genres);


--
-- TOC entry 4985 (class 1259 OID 27535)
-- Name: idx_movies_release_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_movies_release_status ON public.movies USING btree (release_status);


--
-- TOC entry 5001 (class 1259 OID 27646)
-- Name: idx_payments_booking; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payments_booking ON public.payments USING btree (booking_id);


--
-- TOC entry 4972 (class 1259 OID 27473)
-- Name: idx_sessions_token; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_token ON public.sessions USING btree (token);


--
-- TOC entry 4988 (class 1259 OID 27565)
-- Name: idx_showtimes_cinema_datetime; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_showtimes_cinema_datetime ON public.showtimes USING btree (cinema_id, show_date, show_time);


--
-- TOC entry 4996 (class 1259 OID 27611)
-- Name: ux_booking_seats; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ux_booking_seats ON public.booking_seats USING btree (seat_id, booking_id);


--
-- TOC entry 4983 (class 1259 OID 27515)
-- Name: ux_seats_studio_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ux_seats_studio_code ON public.seats USING btree (studio_id, seat_code);


--
-- TOC entry 5012 (class 2606 OID 27601)
-- Name: booking_seats booking_seats_booking_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.booking_seats
    ADD CONSTRAINT booking_seats_booking_id_fkey FOREIGN KEY (booking_id) REFERENCES public.bookings(id) ON DELETE CASCADE;


--
-- TOC entry 5013 (class 2606 OID 27606)
-- Name: booking_seats booking_seats_seat_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.booking_seats
    ADD CONSTRAINT booking_seats_seat_id_fkey FOREIGN KEY (seat_id) REFERENCES public.seats(id);


--
-- TOC entry 5010 (class 2606 OID 27584)
-- Name: bookings bookings_showtime_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_showtime_id_fkey FOREIGN KEY (showtime_id) REFERENCES public.showtimes(id) ON DELETE CASCADE;


--
-- TOC entry 5011 (class 2606 OID 27579)
-- Name: bookings bookings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5014 (class 2606 OID 27636)
-- Name: payments payments_booking_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_booking_id_fkey FOREIGN KEY (booking_id) REFERENCES public.bookings(id) ON DELETE CASCADE;


--
-- TOC entry 5015 (class 2606 OID 27641)
-- Name: payments payments_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id);


--
-- TOC entry 5006 (class 2606 OID 27510)
-- Name: seats seats_studio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.seats
    ADD CONSTRAINT seats_studio_id_fkey FOREIGN KEY (studio_id) REFERENCES public.studios(id) ON DELETE CASCADE;


--
-- TOC entry 5004 (class 2606 OID 27468)
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5007 (class 2606 OID 27550)
-- Name: showtimes showtimes_cinema_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.showtimes
    ADD CONSTRAINT showtimes_cinema_id_fkey FOREIGN KEY (cinema_id) REFERENCES public.cinemas(id) ON DELETE CASCADE;


--
-- TOC entry 5008 (class 2606 OID 27560)
-- Name: showtimes showtimes_movie_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.showtimes
    ADD CONSTRAINT showtimes_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON DELETE CASCADE;


--
-- TOC entry 5009 (class 2606 OID 27555)
-- Name: showtimes showtimes_studio_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.showtimes
    ADD CONSTRAINT showtimes_studio_id_fkey FOREIGN KEY (studio_id) REFERENCES public.studios(id) ON DELETE CASCADE;


--
-- TOC entry 5005 (class 2606 OID 27495)
-- Name: studios studios_cinema_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.studios
    ADD CONSTRAINT studios_cinema_id_fkey FOREIGN KEY (cinema_id) REFERENCES public.cinemas(id) ON DELETE CASCADE;


-- Completed on 2026-01-14 09:26:50

--
-- PostgreSQL database dump complete
--

\unrestrict WtVSngrAcOh3GUOHwqhaMe1ppJCyg99n9MnQVPgszQdFFnSiaYYQO7jtccvYhJz

