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
func (r *PostsRepository) FindPostWithStatsById(
	ctx context.Context,
	dbtx repository.DBTX,
	userId string,
	postId string,
) (*model.PostWithUserAndCounts, error) {
	row := dbtx.QueryRowContext(ctx, `
		SELECT 
		  p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
		  u.id, u.username, u.display_name, u.icon_url,
		  COALESCE(like_counts.like_count, 0) AS like_count,
		  COALESCE(repost_counts.repost_count, 0) AS repost_count,
		  COALESCE(comment_counts.comment_count, 0) AS comment_count,
		  CASE WHEN user_likes.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_liked,
		  CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_reposted
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS like_count FROM likes GROUP BY post_id
		) like_counts ON like_counts.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS repost_count FROM reposts GROUP BY post_id
		) repost_counts ON repost_counts.post_id = p.id
		LEFT JOIN (
		  SELECT reply_to, COUNT(*) AS comment_count FROM posts WHERE reply_to IS NOT NULL GROUP BY reply_to
		) comment_counts ON comment_counts.reply_to = p.id
		LEFT JOIN (
		  SELECT post_id, user_id FROM likes WHERE user_id = ?
		) user_likes ON user_likes.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, user_id FROM reposts WHERE user_id = ?
		) user_reposts ON user_reposts.post_id = p.id
		WHERE p.id = ?
	`, userId, userId, postId)

	var p model.PostWithUserAndCounts
	if err := row.Scan(
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
		&p.UserActions.Liked,
		&p.UserActions.Reposted,
	); err != nil {
		return nil, err
	}
	return &p, nil
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

func (r *PostsRepository) FindPostsWithStats(
	ctx context.Context,
	dbtx repository.DBTX,
	userId string,
	limit int,
	offset int,
) (*[]model.PostWithUserAndCounts, error) {
	rows, err := dbtx.QueryContext(ctx, `
		SELECT 
		  p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
		  u.id, u.username, u.display_name, u.icon_url,
		  COALESCE(like_counts.like_count, 0) AS like_count,
		  COALESCE(repost_counts.repost_count, 0) AS repost_count,
		  COALESCE(comment_counts.comment_count, 0) AS comment_count,
		  CASE WHEN user_likes.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_liked,
		  CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_reposted
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS like_count FROM likes GROUP BY post_id
		) like_counts ON like_counts.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS repost_count FROM reposts GROUP BY post_id
		) repost_counts ON repost_counts.post_id = p.id
		LEFT JOIN (
		  SELECT reply_to, COUNT(*) AS comment_count FROM posts WHERE reply_to IS NOT NULL GROUP BY reply_to
		) comment_counts ON comment_counts.reply_to = p.id
		LEFT JOIN (
		  SELECT post_id, user_id FROM likes WHERE user_id = ?
		) user_likes ON user_likes.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, user_id FROM reposts WHERE user_id = ?
		) user_reposts ON user_reposts.post_id = p.id
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`, userId, userId, limit, offset)
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
			&p.UserActions.Liked,    // bool型想定
			&p.UserActions.Reposted, // bool型想定
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}

func (r *PostsRepository) FindRepliesWithStats(
	ctx context.Context,
	dbtx repository.DBTX,
	userId string,
	parentPostId string,
	limit int,
	offset int,
) (*[]model.PostWithUserAndCounts, error) {
	rows, err := dbtx.QueryContext(ctx, `
		SELECT 
		  p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
		  u.id, u.username, u.display_name, u.icon_url,
		  COALESCE(like_counts.like_count, 0) AS like_count,
		  COALESCE(repost_counts.repost_count, 0) AS repost_count,
		  COALESCE(comment_counts.comment_count, 0) AS comment_count,
		  CASE WHEN user_likes.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_liked,
		  CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_reposted
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS like_count FROM likes GROUP BY post_id
		) like_counts ON like_counts.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS repost_count FROM reposts GROUP BY post_id
		) repost_counts ON repost_counts.post_id = p.id
		LEFT JOIN (
		  SELECT reply_to, COUNT(*) AS comment_count FROM posts WHERE reply_to IS NOT NULL GROUP BY reply_to
		) comment_counts ON comment_counts.reply_to = p.id
		LEFT JOIN (
		  SELECT post_id, user_id FROM likes WHERE user_id = ?
		) user_likes ON user_likes.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, user_id FROM reposts WHERE user_id = ?
		) user_reposts ON user_reposts.post_id = p.id
		WHERE p.reply_to = ?
		ORDER BY p.created_at ASC
		LIMIT ? OFFSET ?
	`, userId, userId, parentPostId, limit, offset)
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
			&p.UserActions.Liked,
			&p.UserActions.Reposted,
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}

func (r *PostsRepository) FindPostsByUserIdWithStats(
	ctx context.Context,
	dbtx repository.DBTX,
	targetUserId string, // 投稿者のuserId
	viewerUserId string, // ログイン中のユーザー（リアクション確認用）
	limit int,
	offset int,
) (*[]model.PostWithUserAndCounts, error) {
	rows, err := dbtx.QueryContext(ctx, `
		SELECT 
		  p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
		  u.id, u.username, u.display_name, u.icon_url,
		  COALESCE(like_counts.like_count, 0) AS like_count,
		  COALESCE(repost_counts.repost_count, 0) AS repost_count,
		  COALESCE(comment_counts.comment_count, 0) AS comment_count,
		  CASE WHEN user_likes.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_liked,
		  CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_reposted
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS like_count FROM likes GROUP BY post_id
		) like_counts ON like_counts.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, COUNT(*) AS repost_count FROM reposts GROUP BY post_id
		) repost_counts ON repost_counts.post_id = p.id
		LEFT JOIN (
		  SELECT reply_to, COUNT(*) AS comment_count FROM posts WHERE reply_to IS NOT NULL GROUP BY reply_to
		) comment_counts ON comment_counts.reply_to = p.id
		LEFT JOIN (
		  SELECT post_id, user_id FROM likes WHERE user_id = ?
		) user_likes ON user_likes.post_id = p.id
		LEFT JOIN (
		  SELECT post_id, user_id FROM reposts WHERE user_id = ?
		) user_reposts ON user_reposts.post_id = p.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`, viewerUserId, viewerUserId, targetUserId, limit, offset)
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
			&p.UserActions.Liked,
			&p.UserActions.Reposted,
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}
