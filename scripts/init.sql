-- Create database schema for URL shortener microservices

-- Create links table
CREATE TABLE IF NOT EXISTS links (
    id VARCHAR(255) PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create stats table
CREATE TABLE IF NOT EXISTS stats (
    id VARCHAR(255) PRIMARY KEY,
    link_id VARCHAR(255) NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    platform INTEGER DEFAULT 0,
    user_agent TEXT,
    ip_address INET,
    referrer TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_stats_link_id ON stats(link_id);
CREATE INDEX IF NOT EXISTS idx_stats_created_at ON stats(created_at);
CREATE INDEX IF NOT EXISTS idx_links_created_at ON links(created_at);

-- Insert some test data
INSERT INTO links (id, original_url) VALUES 
    ('testid1', 'https://example.com/link1'),
    ('testid2', 'https://example.com/link2'),
    ('testid3', 'https://example.com/link3')
ON CONFLICT (id) DO NOTHING;

INSERT INTO stats (id, link_id, platform) VALUES 
    ('stat1', 'testid1', 0),
    ('stat2', 'testid2', 1),
    ('stat3', 'testid3', 2)
ON CONFLICT (id) DO NOTHING;
