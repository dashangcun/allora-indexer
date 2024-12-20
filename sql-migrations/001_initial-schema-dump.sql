--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4
-- Dumped by pg_dump version 16.3

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
-- Name: networklossbundlevaluetype; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.networklossbundlevaluetype AS ENUM (
    'InfererValues',
    'ForecasterValues',
    'OneOutInfererValues',
    'OneInForecasterValues',
    'OneOutForecasterValues'
);


--
-- Name: reputervaluetype; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.reputervaluetype AS ENUM (
    'InfererValues',
    'ForecasterValues',
    'OneOutInfererValues',
    'OneInForecasterValues',
    'OneOutForecasterValues'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.addresses (
    id integer NOT NULL,
    pub_key character varying(255) DEFAULT NULL::character varying,
    type character varying(255) DEFAULT NULL::character varying,
    memo character varying(255) DEFAULT NULL::character varying,
    address character varying(255) DEFAULT NULL::character varying
);


--
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- Name: block_info; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.block_info (
    block_hash character varying(255),
    block_total_parts integer,
    block_part_set_header_hash character varying(255),
    block_version character varying(255),
    chain_id character varying(255),
    height bigint NOT NULL,
    block_time timestamp without time zone,
    last_block_hash character varying(255),
    last_block_total_parts integer,
    last_block_part_set_header_hash character varying(255),
    last_commit_hash character varying(255),
    data_hash character varying(255),
    validators_hash character varying(255),
    next_validators_hash character varying(255),
    consensus_hash character varying(255),
    app_hash character varying(255),
    last_results_hash character varying(255),
    evidence_hash character varying(255),
    proposer_address character varying(255)
);


--
-- Name: bundle_values; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bundle_values (
    bundle_id integer,
    reputer_value_type public.reputervaluetype,
    value character varying(255),
    worker character varying(255)
);


--
-- Name: consensus_params; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.consensus_params (
    id integer NOT NULL,
    max_bytes character varying(255),
    max_gas character varying(255),
    max_age_duration character varying(255),
    max_age_num_blocks character varying(255),
    evidence_max_bytes character varying(255),
    pub_key_types text
);


--
-- Name: consensus_params_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.consensus_params_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: consensus_params_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.consensus_params_id_seq OWNED BY public.consensus_params.id;


--
-- Name: ecosystem_token_mint; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ecosystem_token_mint (
    id integer NOT NULL,
    height_tx bigint,
    block_height bigint,
    token_amount numeric(72,18)
);


--
-- Name: ecosystem_token_mint_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ecosystem_token_mint_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ecosystem_token_mint_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ecosystem_token_mint_id_seq OWNED BY public.ecosystem_token_mint.id;


--
-- Name: ema_scores; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ema_scores (
    id integer NOT NULL,
    height_tx bigint,
    height bigint,
    topic_id integer,
    type character varying(255),
    address character varying(255),
    score numeric(72,18),
    is_active boolean
);


--
-- Name: ema_scores_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ema_scores_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ema_scores_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ema_scores_id_seq OWNED BY public.ema_scores.id;


--
-- Name: events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.events (
    id integer NOT NULL,
    height_tx bigint,
    height bigint,
    type character varying(255),
    sender character varying(255),
    data jsonb,
    hash numeric
);


--
-- Name: events_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.events_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.events_id_seq OWNED BY public.events.id;


--
-- Name: forecast_values; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.forecast_values (
    forecast_id integer,
    value character varying(255),
    inferer character varying(255)
);


--
-- Name: forecaster_network_regret; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.forecaster_network_regret (
    id integer NOT NULL,
    height_tx bigint,
    block_height bigint,
    topic_id bigint,
    addresses text[],
    regrets numeric(72,18)[]
);


--
-- Name: forecaster_network_regret_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.forecaster_network_regret_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: forecaster_network_regret_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.forecaster_network_regret_id_seq OWNED BY public.forecaster_network_regret.id;


--
-- Name: forecasts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.forecasts (
    id integer NOT NULL,
    message_height integer,
    message_id integer,
    nonce_block_height integer,
    topic_id integer,
    block_height integer,
    forecaster character varying(255),
    extra_data character varying(255)
);


--
-- Name: forecasts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.forecasts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: forecasts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.forecasts_id_seq OWNED BY public.forecasts.id;


--
-- Name: inferences; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.inferences (
    id integer NOT NULL,
    message_height integer,
    message_id integer,
    nonce_block_height integer,
    topic_id integer,
    block_height integer,
    inferer character varying(255),
    value text,
    extra_data text,
    proof text
);


--
-- Name: inferences_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.inferences_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: inferences_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.inferences_id_seq OWNED BY public.inferences.id;


--
-- Name: inferer_network_regret; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.inferer_network_regret (
    id integer NOT NULL,
    height_tx bigint,
    block_height bigint,
    topic_id bigint,
    addresses text[],
    regrets numeric(72,18)[]
);


--
-- Name: inferer_network_regret_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.inferer_network_regret_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: inferer_network_regret_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.inferer_network_regret_id_seq OWNED BY public.inferer_network_regret.id;


--
-- Name: last_commit_values; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.last_commit_values (
    id integer NOT NULL,
    height_tx bigint,
    height bigint,
    topic_id integer,
    is_worker boolean
);


--
-- Name: last_commit_values_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.last_commit_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: last_commit_values_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.last_commit_values_id_seq OWNED BY public.last_commit_values.id;


--
-- Name: listening_coefficients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.listening_coefficients (
    id integer NOT NULL,
    actor_type character varying(255),
    topic_id bigint,
    block_height bigint,
    addresses text[],
    coefficients numeric(72,18)[]
);


--
-- Name: listening_coefficients_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.listening_coefficients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: listening_coefficients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.listening_coefficients_id_seq OWNED BY public.listening_coefficients.id;


--
-- Name: messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.messages (
    id integer NOT NULL,
    height bigint,
    type character varying(255),
    sender character varying(255),
    data jsonb,
    hash numeric,
    result jsonb,
    tx_hash character varying(255)
);


--
-- Name: messages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.messages_id_seq OWNED BY public.messages.id;


--
-- Name: naive_inferer_network_regret; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.naive_inferer_network_regret (
    id integer NOT NULL,
    height_tx bigint,
    block_height bigint,
    topic_id bigint,
    addresses text[],
    regrets numeric(72,18)[]
);


--
-- Name: naive_inferer_network_regret_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.naive_inferer_network_regret_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: naive_inferer_network_regret_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.naive_inferer_network_regret_id_seq OWNED BY public.naive_inferer_network_regret.id;


--
-- Name: networkloss_bundle_values; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.networkloss_bundle_values (
    bundle_id integer,
    reputer_value_type public.networklossbundlevaluetype,
    value character varying(255),
    worker character varying(255)
);


--
-- Name: networklosses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.networklosses (
    id integer NOT NULL,
    height_tx bigint,
    height bigint,
    topic_id integer,
    naive_value character varying(255),
    combined_value character varying(255)
);


--
-- Name: networklosses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.networklosses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: networklosses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.networklosses_id_seq OWNED BY public.networklosses.id;


--
-- Name: query_results; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.query_results (
    id integer NOT NULL,
    key text NOT NULL,
    query_type text NOT NULL,
    metadata jsonb,
    value jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: query_results_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.query_results_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: query_results_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.query_results_id_seq OWNED BY public.query_results.id;


--
-- Name: reputer_bundles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reputer_bundles (
    id integer NOT NULL,
    reputer_payload_id integer,
    pubkey character varying(255),
    signature character varying(255),
    reputer character varying(255),
    topic_id integer,
    extra_data character varying(255),
    naive_value character varying(255),
    combined_value character varying(255),
    reputer_request_worker_nonce integer,
    reputer_request_reputer_nonce integer
);


--
-- Name: reputer_bundles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.reputer_bundles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: reputer_bundles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.reputer_bundles_id_seq OWNED BY public.reputer_bundles.id;


--
-- Name: reputer_payload; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reputer_payload (
    id integer NOT NULL,
    message_height integer,
    message_id integer,
    sender character varying(255),
    worker_nonce_block_height integer,
    reputer_nonce_block_height integer,
    topic_id integer
);


--
-- Name: reputer_payload_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.reputer_payload_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: reputer_payload_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.reputer_payload_id_seq OWNED BY public.reputer_payload.id;


--
-- Name: reputer_stakes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reputer_stakes (
    id integer NOT NULL,
    type text NOT NULL,
    topic_id integer NOT NULL,
    sender text NOT NULL,
    amount numeric(72,18),
    reputer_address text,
    delegator_address text,
    height integer NOT NULL
);


--
-- Name: reputer_stakes_2; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reputer_stakes_2 (
    id integer,
    type text,
    topic_id integer,
    sender text,
    amount numeric(72,18),
    reputer_address text,
    delegator_address text,
    height integer
);


--
-- Name: reputer_stakes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.reputer_stakes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: reputer_stakes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.reputer_stakes_id_seq OWNED BY public.reputer_stakes.id;


--
-- Name: reputer_total_staked_per_topic; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reputer_total_staked_per_topic (
    id integer NOT NULL,
    reputer_address character varying(255),
    topic_id integer NOT NULL,
    total_stake numeric(72,18),
    height integer,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: reputer_total_staked_per_topic_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.reputer_total_staked_per_topic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: reputer_total_staked_per_topic_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.reputer_total_staked_per_topic_id_seq OWNED BY public.reputer_total_staked_per_topic.id;


--
-- Name: research_metrics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.research_metrics (
    id integer NOT NULL,
    topic_id integer NOT NULL,
    epoch integer NOT NULL,
    address character varying(255),
    metric_value double precision,
    metric_name character varying(255) NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: research_metrics_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.research_metrics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: research_metrics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.research_metrics_id_seq OWNED BY public.research_metrics.id;


--
-- Name: reward_current_block_emission; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reward_current_block_emission (
    id integer NOT NULL,
    height_tx bigint,
    block_height bigint,
    token_amount numeric(72,18)
);


--
-- Name: reward_current_block_emission_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.reward_current_block_emission_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: reward_current_block_emission_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.reward_current_block_emission_id_seq OWNED BY public.reward_current_block_emission.id;


--
-- Name: rewards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.rewards (
    id integer NOT NULL,
    height_tx bigint,
    height bigint,
    topic_id integer,
    type character varying(255),
    address character varying(255),
    value numeric(72,18)
);


--
-- Name: rewards_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.rewards_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: rewards_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.rewards_id_seq OWNED BY public.rewards.id;


--
-- Name: scores; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.scores (
    id integer NOT NULL,
    height_tx bigint,
    height bigint,
    topic_id integer,
    type character varying(255),
    address character varying(255),
    value numeric(72,18)
);


--
-- Name: scores_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.scores_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: scores_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.scores_id_seq OWNED BY public.scores.id;


--
-- Name: tokenomics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tokenomics (
    id integer NOT NULL,
    height_tx bigint,
    staked_amount numeric(72,18),
    circulating_supply numeric(72,18),
    emissions_amount numeric(72,18),
    ecosystem_mint_amount numeric(72,18)
);


--
-- Name: tokenomics_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tokenomics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tokenomics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tokenomics_id_seq OWNED BY public.tokenomics.id;


--
-- Name: topic_forecasting_scores; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.topic_forecasting_scores (
    id integer NOT NULL,
    height_tx bigint,
    topic_id integer,
    score character varying(255)
);


--
-- Name: topic_forecasting_scores_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.topic_forecasting_scores_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: topic_forecasting_scores_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.topic_forecasting_scores_id_seq OWNED BY public.topic_forecasting_scores.id;


--
-- Name: topic_initial_regret; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.topic_initial_regret (
    id integer NOT NULL,
    topic_id bigint,
    height_tx bigint,
    block_height bigint,
    regret numeric(72,18)
);


--
-- Name: topic_initial_regret_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.topic_initial_regret_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: topic_initial_regret_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.topic_initial_regret_id_seq OWNED BY public.topic_initial_regret.id;


--
-- Name: topic_rewards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.topic_rewards (
    id integer NOT NULL,
    height_tx bigint,
    topic_id integer,
    reward character varying(255)
);


--
-- Name: topic_rewards_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.topic_rewards_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: topic_rewards_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.topic_rewards_id_seq OWNED BY public.topic_rewards.id;


--
-- Name: topics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.topics (
    id integer NOT NULL,
    creator character varying(255),
    metadata character varying(255),
    loss_logic character varying(255),
    loss_method character varying(255),
    inference_logic character varying(255),
    inference_method character varying(255),
    epoch_length character varying(255),
    ground_truth_lag character varying(255),
    default_arg character varying(255),
    pnorm character varying(255),
    alpha_regret character varying(255),
    preward_reputer character varying(255),
    preward_inference character varying(255),
    preward_forecast character varying(255),
    f_tolerance character varying(255),
    allow_negative boolean,
    message_height integer,
    message_id integer
);


--
-- Name: transfers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.transfers (
    id integer NOT NULL,
    message_height integer,
    message_id integer,
    from_address character varying(255),
    topic_id integer,
    to_address character varying(255) DEFAULT NULL::character varying,
    amount character varying(255),
    denom character varying(255)
);


--
-- Name: transfers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.transfers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: transfers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.transfers_id_seq OWNED BY public.transfers.id;


--
-- Name: validator_commission; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.validator_commission (
    id integer NOT NULL,
    height_tx bigint,
    validator character varying(255),
    amount numeric(72,18)
);


--
-- Name: validator_commission_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.validator_commission_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: validator_commission_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.validator_commission_id_seq OWNED BY public.validator_commission.id;


--
-- Name: validator_rewards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.validator_rewards (
    id integer NOT NULL,
    height_tx bigint,
    validator character varying(255),
    amount numeric(72,18)
);


--
-- Name: validator_rewards_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.validator_rewards_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: validator_rewards_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.validator_rewards_id_seq OWNED BY public.validator_rewards.id;


--
-- Name: validator_withdraw_commission; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.validator_withdraw_commission (
    id integer NOT NULL,
    height_tx bigint,
    validator character varying(255),
    amount numeric(72,18)
);


--
-- Name: validator_withdraw_commission_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.validator_withdraw_commission_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: validator_withdraw_commission_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.validator_withdraw_commission_id_seq OWNED BY public.validator_withdraw_commission.id;


--
-- Name: worker_registrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.worker_registrations (
    message_height integer,
    message_id integer,
    topic_id integer,
    sender character varying(255),
    owner character varying(255),
    worker_libp2pkey character varying(255),
    is_reputer boolean
);


--
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- Name: consensus_params id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.consensus_params ALTER COLUMN id SET DEFAULT nextval('public.consensus_params_id_seq'::regclass);


--
-- Name: ecosystem_token_mint id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ecosystem_token_mint ALTER COLUMN id SET DEFAULT nextval('public.ecosystem_token_mint_id_seq'::regclass);


--
-- Name: ema_scores id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ema_scores ALTER COLUMN id SET DEFAULT nextval('public.ema_scores_id_seq'::regclass);


--
-- Name: events id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events ALTER COLUMN id SET DEFAULT nextval('public.events_id_seq'::regclass);


--
-- Name: forecaster_network_regret id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.forecaster_network_regret ALTER COLUMN id SET DEFAULT nextval('public.forecaster_network_regret_id_seq'::regclass);


--
-- Name: forecasts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.forecasts ALTER COLUMN id SET DEFAULT nextval('public.forecasts_id_seq'::regclass);


--
-- Name: inferences id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inferences ALTER COLUMN id SET DEFAULT nextval('public.inferences_id_seq'::regclass);


--
-- Name: inferer_network_regret id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inferer_network_regret ALTER COLUMN id SET DEFAULT nextval('public.inferer_network_regret_id_seq'::regclass);


--
-- Name: last_commit_values id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.last_commit_values ALTER COLUMN id SET DEFAULT nextval('public.last_commit_values_id_seq'::regclass);


--
-- Name: listening_coefficients id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.listening_coefficients ALTER COLUMN id SET DEFAULT nextval('public.listening_coefficients_id_seq'::regclass);


--
-- Name: messages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages ALTER COLUMN id SET DEFAULT nextval('public.messages_id_seq'::regclass);


--
-- Name: naive_inferer_network_regret id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.naive_inferer_network_regret ALTER COLUMN id SET DEFAULT nextval('public.naive_inferer_network_regret_id_seq'::regclass);


--
-- Name: networklosses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.networklosses ALTER COLUMN id SET DEFAULT nextval('public.networklosses_id_seq'::regclass);


--
-- Name: query_results id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.query_results ALTER COLUMN id SET DEFAULT nextval('public.query_results_id_seq'::regclass);


--
-- Name: reputer_bundles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reputer_bundles ALTER COLUMN id SET DEFAULT nextval('public.reputer_bundles_id_seq'::regclass);


--
-- Name: reputer_payload id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reputer_payload ALTER COLUMN id SET DEFAULT nextval('public.reputer_payload_id_seq'::regclass);


--
-- Name: reputer_stakes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reputer_stakes ALTER COLUMN id SET DEFAULT nextval('public.reputer_stakes_id_seq'::regclass);


--
-- Name: reputer_total_staked_per_topic id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reputer_total_staked_per_topic ALTER COLUMN id SET DEFAULT nextval('public.reputer_total_staked_per_topic_id_seq'::regclass);


--
-- Name: research_metrics id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.research_metrics ALTER COLUMN id SET DEFAULT nextval('public.research_metrics_id_seq'::regclass);


--
-- Name: reward_current_block_emission id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reward_current_block_emission ALTER COLUMN id SET DEFAULT nextval('public.reward_current_block_emission_id_seq'::regclass);


--
-- Name: rewards id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.rewards ALTER COLUMN id SET DEFAULT nextval('public.rewards_id_seq'::regclass);


--
-- Name: scores id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.scores ALTER COLUMN id SET DEFAULT nextval('public.scores_id_seq'::regclass);


--
-- Name: tokenomics id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tokenomics ALTER COLUMN id SET DEFAULT nextval('public.tokenomics_id_seq'::regclass);


--
-- Name: topic_forecasting_scores id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topic_forecasting_scores ALTER COLUMN id SET DEFAULT nextval('public.topic_forecasting_scores_id_seq'::regclass);


--
-- Name: topic_initial_regret id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topic_initial_regret ALTER COLUMN id SET DEFAULT nextval('public.topic_initial_regret_id_seq'::regclass);


--
-- Name: topic_rewards id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topic_rewards ALTER COLUMN id SET DEFAULT nextval('public.topic_rewards_id_seq'::regclass);


--
-- Name: transfers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfers ALTER COLUMN id SET DEFAULT nextval('public.transfers_id_seq'::regclass);


--
-- Name: validator_commission id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.validator_commission ALTER COLUMN id SET DEFAULT nextval('public.validator_commission_id_seq'::regclass);


--
-- Name: validator_rewards id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.validator_rewards ALTER COLUMN id SET DEFAULT nextval('public.validator_rewards_id_seq'::regclass);


--
-- Name: validator_withdraw_commission id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.validator_withdraw_commission ALTER COLUMN id SET DEFAULT nextval('public.validator_withdraw_commission_id_seq'::regclass);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: block_info block_info_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.block_info
    ADD CONSTRAINT block_info_pkey PRIMARY KEY (height);


--
-- Name: consensus_params consensus_params_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.consensus_params
    ADD CONSTRAINT consensus_params_pkey PRIMARY KEY (id);


--
-- Name: ecosystem_token_mint ecosystem_token_mint_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ecosystem_token_mint
    ADD CONSTRAINT ecosystem_token_mint_pkey PRIMARY KEY (id);


--
-- Name: ema_scores ema_scores_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ema_scores
    ADD CONSTRAINT ema_scores_pkey PRIMARY KEY (id);


--
-- Name: events events_height_data; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_height_data UNIQUE (height, hash, type);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: forecaster_network_regret forecaster_network_regret_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.forecaster_network_regret
    ADD CONSTRAINT forecaster_network_regret_pkey PRIMARY KEY (id);


--
-- Name: forecasts forecasts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.forecasts
    ADD CONSTRAINT forecasts_pkey PRIMARY KEY (id);


--
-- Name: inferences inferences_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inferences
    ADD CONSTRAINT inferences_pkey PRIMARY KEY (id);


--
-- Name: inferer_network_regret inferer_network_regret_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inferer_network_regret
    ADD CONSTRAINT inferer_network_regret_pkey PRIMARY KEY (id);


--
-- Name: last_commit_values last_commit_values_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.last_commit_values
    ADD CONSTRAINT last_commit_values_pkey PRIMARY KEY (id);


--
-- Name: listening_coefficients listening_coefficients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.listening_coefficients
    ADD CONSTRAINT listening_coefficients_pkey PRIMARY KEY (id);


--
-- Name: messages messages_height_data; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_height_data UNIQUE (height, hash);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: naive_inferer_network_regret naive_inferer_network_regret_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.naive_inferer_network_regret
    ADD CONSTRAINT naive_inferer_network_regret_pkey PRIMARY KEY (id);


--
-- Name: networklosses networklosses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.networklosses
    ADD CONSTRAINT networklosses_pkey PRIMARY KEY (id);


--
-- Name: query_results query_results_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.query_results
    ADD CONSTRAINT query_results_pkey PRIMARY KEY (id);


--
-- Name: reputer_bundles reputer_bundles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reputer_bundles
    ADD CONSTRAINT reputer_bundles_pkey PRIMARY KEY (id);


--
-- Name: reputer_payload reputer_payload_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reputer_payload
    ADD CONSTRAINT reputer_payload_pkey PRIMARY KEY (id);


--
-- Name: reputer_stakes reputer_stakes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reputer_stakes
    ADD CONSTRAINT reputer_stakes_pkey PRIMARY KEY (id);


--
-- Name: reputer_total_staked_per_topic reputer_total_staked_per_topic_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reputer_total_staked_per_topic
    ADD CONSTRAINT reputer_total_staked_per_topic_pkey PRIMARY KEY (id);


--
-- Name: research_metrics research_metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.research_metrics
    ADD CONSTRAINT research_metrics_pkey PRIMARY KEY (id);


--
-- Name: reward_current_block_emission reward_current_block_emission_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reward_current_block_emission
    ADD CONSTRAINT reward_current_block_emission_pkey PRIMARY KEY (id);


--
-- Name: rewards rewards_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.rewards
    ADD CONSTRAINT rewards_pkey PRIMARY KEY (id);


--
-- Name: scores scores_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.scores
    ADD CONSTRAINT scores_pkey PRIMARY KEY (id);


--
-- Name: tokenomics tokenomics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tokenomics
    ADD CONSTRAINT tokenomics_pkey PRIMARY KEY (id);


--
-- Name: topic_forecasting_scores topic_forecasting_scores_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topic_forecasting_scores
    ADD CONSTRAINT topic_forecasting_scores_pkey PRIMARY KEY (id);


--
-- Name: topic_initial_regret topic_initial_regret_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topic_initial_regret
    ADD CONSTRAINT topic_initial_regret_pkey PRIMARY KEY (id);


--
-- Name: topic_rewards topic_rewards_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topic_rewards
    ADD CONSTRAINT topic_rewards_pkey PRIMARY KEY (id);


--
-- Name: topics topics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT topics_pkey PRIMARY KEY (id);


--
-- Name: transfers transfers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transfers
    ADD CONSTRAINT transfers_pkey PRIMARY KEY (id);


--
-- Name: last_commit_values unique_actor_last_commit_entry; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.last_commit_values
    ADD CONSTRAINT unique_actor_last_commit_entry UNIQUE (topic_id, is_worker);


--
-- Name: ema_scores unique_ema_score_entry; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ema_scores
    ADD CONSTRAINT unique_ema_score_entry UNIQUE (topic_id, type, address, height);


--
-- Name: networklosses unique_networkloss_entry; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.networklosses
    ADD CONSTRAINT unique_networkloss_entry UNIQUE (height_tx, height, topic_id);


--
-- Name: rewards unique_reward_entry; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.rewards
    ADD CONSTRAINT unique_reward_entry UNIQUE (height, topic_id, type, address);


--
-- Name: scores unique_score_entry; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.scores
    ADD CONSTRAINT unique_score_entry UNIQUE (height, topic_id, type, address);


--
-- Name: topic_forecasting_scores unique_topic_forecasting_scores_entry; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topic_forecasting_scores
    ADD CONSTRAINT unique_topic_forecasting_scores_entry UNIQUE (topic_id, height_tx);


--
-- Name: topic_rewards unique_topic_rewards_entry; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topic_rewards
    ADD CONSTRAINT unique_topic_rewards_entry UNIQUE (topic_id, height_tx);


--
-- Name: validator_commission validator_commission_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.validator_commission
    ADD CONSTRAINT validator_commission_pkey PRIMARY KEY (id);


--
-- Name: validator_rewards validator_rewards_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.validator_rewards
    ADD CONSTRAINT validator_rewards_pkey PRIMARY KEY (id);


--
-- Name: validator_withdraw_commission validator_withdraw_commission_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.validator_withdraw_commission
    ADD CONSTRAINT validator_withdraw_commission_pkey PRIMARY KEY (id);


--
-- Name: idx_actor_last_commit_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_actor_last_commit_topic_id ON public.last_commit_values USING btree (topic_id);


--
-- Name: idx_ecosystem_token_mint_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ecosystem_token_mint_block_height ON public.ecosystem_token_mint USING btree (block_height);


--
-- Name: idx_ema_scores_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ema_scores_height ON public.ema_scores USING btree (height);


--
-- Name: idx_emascores_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_emascores_topic_id ON public.ema_scores USING btree (topic_id);


--
-- Name: idx_events_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_events_height ON public.events USING btree (height);


--
-- Name: idx_forecaster_network_regret_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_forecaster_network_regret_block_height ON public.forecaster_network_regret USING btree (block_height);


--
-- Name: idx_forecasts_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_forecasts_block_height ON public.forecasts USING btree (block_height);


--
-- Name: idx_forecasts_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_forecasts_topic_id ON public.forecasts USING btree (topic_id);


--
-- Name: idx_inferences_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_inferences_block_height ON public.inferences USING btree (block_height);


--
-- Name: idx_inferences_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_inferences_topic_id ON public.inferences USING btree (topic_id);


--
-- Name: idx_inferences_topic_inferer_block; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_inferences_topic_inferer_block ON public.inferences USING btree (topic_id, inferer, block_height DESC);


--
-- Name: idx_inferer_network_regret_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_inferer_network_regret_block_height ON public.inferer_network_regret USING btree (block_height);


--
-- Name: idx_inferer_network_regret_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_inferer_network_regret_topic_id ON public.inferer_network_regret USING btree (topic_id);


--
-- Name: idx_last_commit_values_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_last_commit_values_height ON public.last_commit_values USING btree (height);


--
-- Name: idx_listening_coefficients_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_listening_coefficients_block_height ON public.listening_coefficients USING btree (block_height);


--
-- Name: idx_listening_coefficients_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_listening_coefficients_topic_id ON public.listening_coefficients USING btree (topic_id);


--
-- Name: idx_naive_inferer_network_regret_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_naive_inferer_network_regret_block_height ON public.naive_inferer_network_regret USING btree (block_height);


--
-- Name: idx_naive_inferer_network_regret_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_naive_inferer_network_regret_topic_id ON public.naive_inferer_network_regret USING btree (topic_id);


--
-- Name: idx_networklosses_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_networklosses_height ON public.networklosses USING btree (height);


--
-- Name: idx_networklosses_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_networklosses_topic_id ON public.networklosses USING btree (topic_id);


--
-- Name: idx_reputer_bundles_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reputer_bundles_topic_id ON public.reputer_bundles USING btree (topic_id);


--
-- Name: idx_reputer_payload_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reputer_payload_topic_id ON public.reputer_payload USING btree (topic_id);


--
-- Name: idx_reputer_stakes_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reputer_stakes_topic_id ON public.reputer_stakes USING btree (topic_id);


--
-- Name: idx_reputer_topic; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reputer_topic ON public.reputer_total_staked_per_topic USING btree (reputer_address, topic_id);


--
-- Name: idx_reward_current_block_emission_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reward_current_block_emission_block_height ON public.reward_current_block_emission USING btree (block_height);


--
-- Name: idx_rewards_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_rewards_height ON public.rewards USING btree (height);


--
-- Name: idx_rewards_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_rewards_topic_id ON public.rewards USING btree (topic_id);


--
-- Name: idx_scores_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_scores_height ON public.scores USING btree (height);


--
-- Name: idx_scores_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_scores_topic_id ON public.scores USING btree (topic_id);


--
-- Name: idx_tokenomics_height_tx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tokenomics_height_tx ON public.tokenomics USING btree (height_tx);


--
-- Name: idx_topic_forecasting_scores_height_tx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_topic_forecasting_scores_height_tx ON public.topic_forecasting_scores USING btree (height_tx);


--
-- Name: idx_topic_forecasting_scores_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_topic_forecasting_scores_topic_id ON public.topic_forecasting_scores USING btree (topic_id);


--
-- Name: idx_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_topic_id ON public.reputer_total_staked_per_topic USING btree (topic_id);


--
-- Name: idx_topic_initial_regret_block_height; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_topic_initial_regret_block_height ON public.topic_initial_regret USING btree (block_height);


--
-- Name: idx_topic_initial_regret_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_topic_initial_regret_topic_id ON public.topic_initial_regret USING btree (topic_id);


--
-- Name: idx_topic_reward_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_topic_reward_topic_id ON public.topic_rewards USING btree (topic_id);


--
-- Name: idx_topic_rewards_height_tx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_topic_rewards_height_tx ON public.topic_rewards USING btree (height_tx);


--
-- Name: idx_transfers_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_transfers_topic_id ON public.transfers USING btree (topic_id);


--
-- Name: idx_validator_commission_height_tx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_validator_commission_height_tx ON public.validator_commission USING btree (height_tx);


--
-- Name: idx_worker_registrations_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_worker_registrations_topic_id ON public.worker_registrations USING btree (topic_id);


--
-- Name: messages_result_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX messages_result_code ON public.messages USING btree (((result ->> 'code'::text)));


--
-- Name: query_results_key_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX query_results_key_idx ON public.query_results USING btree (key);


--
-- Name: query_results_query_type_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX query_results_query_type_idx ON public.query_results USING btree (query_type);


--
-- Name: query_results_query_type_metadata_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX query_results_query_type_metadata_idx ON public.query_results USING btree (query_type, metadata);


--
-- Name: research_metrics_epoch_topic_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX research_metrics_epoch_topic_idx ON public.research_metrics USING btree (epoch, topic_id);


--
-- Name: research_metrics_topic_metric_epoch_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX research_metrics_topic_metric_epoch_idx ON public.research_metrics USING btree (topic_id, metric_name, epoch);


--
-- PostgreSQL database dump complete
--

