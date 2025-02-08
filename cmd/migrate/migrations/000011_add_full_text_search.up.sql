DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes WHERE tablename = 'posts' AND indexname = 'posts_search_idx'
    ) THEN
        CREATE INDEX posts_search_idx 
        ON posts 
        USING GIN (to_tsvector('english', title || ' ' || content));
    END IF;
END $$;
