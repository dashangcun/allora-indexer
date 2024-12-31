CREATE TABLE IF NOT EXISTS public.topic_initial_ema_scores (
    id SERIAL PRIMARY KEY,
    actor_type VARCHAR(255),
    topic_id INT,
    block_height BIGINT,
    score NUMERIC(72,18)
);

CREATE INDEX IF NOT EXISTS idx_topic_initial_ema_scores_topic_id ON public.topic_initial_ema_scores (topic_id);
CREATE INDEX IF NOT EXISTS idx_topic_initial_ema_scores_block_height ON public.topic_initial_ema_scores (block_height);
ALTER TABLE ONLY public.topic_initial_ema_scores ADD CONSTRAINT unique_topic_initial_ema_score_entry UNIQUE (topic_id, actor_type, block_height);
