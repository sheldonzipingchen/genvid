-- Seed data for avatars
INSERT INTO avatars (name, display_name, gender, age_range, style, language, is_premium, sort_order) VALUES
('emma_casual', 'Emma', 'female', '20s', 'casual', ARRAY['en', 'es'], false, 1),
('james_pro', 'James', 'male', '30s', 'professional', ARRAY['en'], false, 2),
('sofia_energy', 'Sofia', 'female', '20s', 'energetic', ARRAY['en', 'pt'], true, 3),
('li_friendly', 'Li', 'male', '30s', 'friendly', ARRAY['en', 'zh'], false, 4),
('maria_elegant', 'Maria', 'female', '40s', 'elegant', ARRAY['en', 'es', 'pt'], true, 5),
('alex_trendy', 'Alex', 'male', '20s', 'trendy', ARRAY['en'], false, 6),
('yuki_casual', 'Yuki', 'female', '20s', 'casual', ARRAY['en', 'ja'], true, 7),
('david_pro', 'David', 'male', '40s', 'professional', ARRAY['en'], false, 8),
('priya_friendly', 'Priya', 'female', '30s', 'friendly', ARRAY['en', 'hi'], false, 9),
('marco_trendy', 'Marco', 'male', '20s', 'trendy', ARRAY['en', 'es', 'it'], true, 10);

-- Seed data for script templates
INSERT INTO script_templates (name, category, template_text, language, is_premium) VALUES
('Product Review', 'product_review', 
'I''ve been using {product_name} for {time_period} now, and I have to say... {review_content}. If you''re looking for {benefit}, this is definitely worth checking out! Link in bio!',
'en', false),

('Unboxing Experience', 'unboxing',
'Hey everyone! I just got my {product_name} delivered today. Let''s unbox this together! First impressions... {first_impression}. Stay tuned for my full review!',
'en', false),

('Before & After', 'before_after',
'So I''ve been using {product_name} for {time_period}. Here''s what it looked like before... And here''s after {time_period}. The difference is {result}! Absolutely recommend this!',
'en', true),

('Comparison', 'comparison',
'Today I''m comparing {product_name} with {competitor}. Let''s break it down: price, quality, and overall value. Here''s my honest take... Spoiler: {product_name} wins!',
'en', true),

('Quick Tutorial', 'tutorial',
'Quick tip! Here''s how to use {product_name} for the best results. Step 1: {step1}. Step 2: {step2}. Step 3: {step3}. That''s it! Save this for later!',
'en', false),

('Testimonial', 'testimonial',
'I was skeptical about {product_name} at first, but after {time_period}, I''m a believer! {result_benefit}. If you''re on the fence, just try it. You won''t regret it!',
'en', true),

('Story Hook', 'storytelling',
'So this happened... I was {situation} and then I discovered {product_name}. {story_content}. Now I can''t imagine life without it!',
'en', true),

('Problem Solution', 'product_review',
'Are you struggling with {problem}? I was too, until I found {product_name}. {solution_content}. Game changer! Link in bio to get yours!',
'en', false);
