CREATE TABLE IF NOT EXISTS blog_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    content TEXT
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment_text TEXT
);

CREATE TABLE IF NOT EXISTS blog_posts_comments (
    blog_post_id INTEGER,
    comment_id INTEGER,
    PRIMARY KEY (blog_post_id, comment_id)
    FOREIGN KEY (blog_post_id) REFERENCES blog_posts(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);