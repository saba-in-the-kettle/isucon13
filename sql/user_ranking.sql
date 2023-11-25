create table IF NOT EXISTS user_ranking as
select u.id, u.name, count(*) as reaction_count from users u
     INNER JOIN livestreams l ON l.user_id = u.id
     INNER JOIN reactions r ON r.livestream_id = l.id
group by u.id;
alter table user_ranking add tips_count bigint default 0 not null;
