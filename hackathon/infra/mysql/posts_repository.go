package mysql

import (
	"context"
	"hackathon/model"
	"hackathon/repository"
)

type PostsRepository struct {
}

var _ repository.PostsRepository = (*PostsRepository)(nil)

func NewPostsRepository() *PostsRepository {
	return &PostsRepository{}
}

// 投稿1件取得
func (r *PostsRepository) FindPostById(ctx context.Context, dbtx repository.DBTX, id string) (*model.Post, error) {
	row := dbtx.QueryRowContext(ctx, `
		SELECT id, user_id, text, reply_to, repost_ref, media_type, media_url, created_at
		FROM posts
		WHERE id = ?
	`, id)

	var p model.Post
	if err := row.Scan(&p.Id, &p.UserId, &p.Text, &p.ReplyTo, &p.RepostRef, &p.MediaType, &p.MediaURL, &p.CreatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

// ユーザーごとの投稿一覧取得
func (r *PostsRepository) FindPostsByUserId(ctx context.Context, dbtx repository.DBTX, userId string) (*[]model.Post, error) {
	rows, err := dbtx.QueryContext(ctx, `
		SELECT id, user_id, text, reply_to, repost_ref, media_type, media_url, created_at
		FROM posts
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.Id, &p.UserId, &p.Text, &p.ReplyTo, &p.RepostRef, &p.MediaType, &p.MediaURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &posts, nil
}

// 投稿作成
func (r *PostsRepository) InsertPost(ctx context.Context, dbtx repository.DBTX, post model.Post) error {
	_, err := dbtx.ExecContext(ctx, `
		INSERT INTO posts (id, user_id, text, reply_to, repost_ref, media_type, media_url, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, post.Id, post.UserId, post.Text, post.ReplyTo, post.RepostRef, post.MediaType, post.MediaURL, post.CreatedAt)
	return err
}

// 投稿物理削除
func (r *PostsRepository) DeletePost(ctx context.Context, dbtx repository.DBTX, id string) error {
	_, err := dbtx.ExecContext(ctx, `
		DELETE FROM posts
		WHERE id = ?
	`, id)
	return err
}

func (r *PostsRepository) FindAllWithCounts(ctx context.Context, dbtx repository.DBTX, limit, offset int) (*[]model.PostWithUserAndCounts, error) {
	rows, err := dbtx.QueryContext(ctx, `
		SELECT 
		  p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
		  u.id, u.username, u.display_name, u.icon_url,
		  COALESCE(like_counts.like_count, 0) AS like_count,
		  COALESCE(repost_counts.repost_count, 0) AS repost_count,
		  COALESCE(comment_counts.comment_count, 0) AS comment_count
		FROM posts p
		LEFT JOIN user u ON p.user_id = u.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS like_count FROM likes GROUP BY post_id
		) like_counts ON like_counts.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS repost_count FROM reposts GROUP BY post_id
		) repost_counts ON repost_counts.post_id = p.id
		LEFT JOIN (
		  SELECT reply_to, COUNT(*) AS comment_count FROM posts WHERE reply_to IS NOT NULL GROUP BY reply_to
		) comment_counts ON comment_counts.reply_to = p.id
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.PostWithUserAndCounts
	for rows.Next() {
		var p model.PostWithUserAndCounts
		if err := rows.Scan(
			&p.Id,
			&p.Text,
			&p.ReplyTo,
			&p.RepostRef,
			&p.MediaType,
			&p.MediaURL,
			&p.CreatedAt,
			&p.Author.Id,
			&p.Author.Username,
			&p.Author.DisplayName,
			&p.Author.IconURL,
			&p.Stats.Likes,
			&p.Stats.Reposts,
			&p.Stats.Comments,
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}
