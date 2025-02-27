INSERT INTO books (id, title, author_id, category_id, stock, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'Clean Code', 'a99926f2-100d-41a9-83f5-e35a186ef838', 'a556f178-600d-46f2-87f6-b88f621f7843', 10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'The Pragmatic Programmer', 'b111b734-200e-42b8-93c2-f44b287cf839', 'a556f178-600d-46f2-87f6-b88f621f7843', 8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Design Patterns', 'c223c845-300f-43c9-84d3-b55c398d6840', 'b667e289-700e-47f3-98f7-b99f732f7844', 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Refactoring', 'd334d956-400a-44d0-95e4-b66d409e6841', 'b667e289-700e-47f3-98f7-b99f732f7844', 12, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'You Don’t Know JS', 'e445e067-500b-45e1-86f5-b77e510f6842', 'a556f178-600d-46f2-87f6-b88f621f7843', 20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Introduction to Algorithms', 'd334d956-400a-44d0-95e4-b66d409e6841', 'c778f390-800f-48f4-89f8-b00f843f7845', 7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Grokking Algorithms', 'c223c845-300f-43c9-84d3-b55c398d6840', 'c778f390-800f-48f4-89f8-b00f843f7845', 6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Head First Design Patterns', 'c223c845-300f-43c9-84d3-b55c398d6840', 'b667e289-700e-47f3-98f7-b99f732f7844', 15, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Cracking the Coding Interview', 'b111b734-200e-42b8-93c2-f44b287cf839', 'c778f390-800f-48f4-89f8-b00f843f7845', 9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Fluent Python', 'e445e067-500b-45e1-86f5-b77e510f6842', 'a556f178-600d-46f2-87f6-b88f621f7843', 14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
