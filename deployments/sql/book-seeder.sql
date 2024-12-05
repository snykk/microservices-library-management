INSERT INTO books (id, title, author_id, category_id, stock, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Clean Code', 'a99926f2-100d-41a9-83f5-e35a186ef838', 'a556f178-600d-86f2-e7f6-b88f621f7843', 10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'The Pragmatic Programmer', 'b111b734-200e-42b8-a3c2-f44b287cf839', 'a556f178-600d-86f2-e7f6-b88f621f7843', 8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Design Patterns', 'c223c845-300f-53c9-b4d3-b55c398d6840', 'b667e289-700e-97f3-f8f7-b99f732f7844', 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Refactoring', 'd334d956-400a-64d0-c5e4-b66d409e6841', 'b667e289-700e-97f3-f8f7-b99f732f7844', 12, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'You Don’t Know JS', 'e445e067-500b-75e1-d6f5-b77e510f6842', 'a556f178-600d-86f2-e7f6-b88f621f7843', 20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Introduction to Algorithms', 'd334d956-400a-64d0-c5e4-b66d409e6841', 'c778f390-800f-a8f4-a9f8-b00f843f7845', 7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Grokking Algorithms', 'c223c845-300f-53c9-b4d3-b55c398d6840', 'c778f390-800f-a8f4-a9f8-b00f843f7845', 6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Head First Design Patterns', 'c223c845-300f-53c9-b4d3-b55c398d6840', 'b667e289-700e-97f3-f8f7-b99f732f7844', 15, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Cracking the Coding Interview', 'b111b734-200e-42b8-a3c2-f44b287cf839', 'c778f390-800f-a8f4-a9f8-b00f843f7845', 9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Fluent Python', 'e445e067-500b-75e1-d6f5-b77e510f6842', 'a556f178-600d-86f2-e7f6-b88f621f7843', 14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
