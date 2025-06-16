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

// 投稿＋like数・リポスト数一覧（JOIN集計）
func (r *PostsRepository) FindAllWithCounts(ctx context.Context, dbtx repository.DBTX, limit, offset int) (*[]model.PostWithUserAndCounts, error) {
	rows, err := dbtx.QueryContext(ctx, `
		SELECT 
		  p.id, p.user_id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
		  u.id, u.username, u.display_name, u.email, u.bio, u.icon_url, u.created_at, u.updated_at,
		  COUNT(DISTINCT l.id) AS like_count,
		  COUNT(DISTINCT r2.id) AS repost_count
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
		LEFT JOIN posts r2 ON r2.repost_ref = p.id
		GROUP BY p.id, u.id, u.username, u.display_name, u.email, u.bio, u.icon_url, u.created_at, u.updated_at
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
			&p.Post.Id, &p.Post.UserId, &p.Post.Text, &p.Post.ReplyTo, &p.Post.RepostRef, &p.Post.MediaType, &p.Post.MediaURL, &p.Post.CreatedAt,
			&p.User.Id, &p.User.UserName, &p.User.IconURL,
			&p.LikeCount, &p.RepostCount,
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}
