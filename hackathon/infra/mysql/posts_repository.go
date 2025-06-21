package mysql

import (
	"context"
	"fmt"
	"hackathon/model"
	"hackathon/repository"
)

type PostsRepository struct {
}

var _ repository.PostsRepository = (*PostsRepository)(nil)

func NewPostsRepository() *PostsRepository {
	return &PostsRepository{}
}

// 【変更後】投稿1件取得（ネタバレキーワード対応）
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
         CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_reposted,
         s.keyword
       FROM posts p
       LEFT JOIN spoilers s ON p.id = s.post_id
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
		&p.SpoilerWord, // ネタバレキーワード用のスキャン先を追加
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

// 【変更後】投稿とリポストを時系列で取得（ネタバレキーワード対応）
func (r *PostsRepository) FindPostsWithStats(
	ctx context.Context,
	dbtx repository.DBTX,
	userId string,
	limit int,
	offset int,
) (*[]model.PostWithUserAndCounts, error) {
	rows, err := dbtx.QueryContext(ctx, `
		SELECT
			timeline.id,
			timeline.text,
			timeline.reply_to,
			timeline.repost_ref,
			timeline.media_type,
			timeline.media_url,
			timeline.created_at,
			timeline.author_id,
			timeline.author_username,
			timeline.author_display_name,
			timeline.author_icon_url,
			COALESCE(like_counts.like_count, 0),
			COALESCE(repost_counts.repost_count, 0),
			COALESCE(comment_counts.comment_count, 0),
			CASE WHEN user_likes.user_id IS NOT NULL THEN 1 ELSE 0 END,
			CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END,
			timeline.reposted_by_username,
			s.keyword
		FROM (
			-- 通常の投稿
			SELECT
				p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
				u.id AS author_id, u.username AS author_username, u.display_name AS author_display_name, u.icon_url AS author_icon_url,
				p.created_at AS event_time,
				CAST(NULL AS CHAR(50)) COLLATE utf8mb4_0900_ai_ci AS reposted_by_username
			FROM posts p
			JOIN users u ON p.user_id = u.id
		
			UNION ALL
		
			-- リポストされた投稿
			SELECT
				p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
				u_orig.id AS author_id, u_orig.username AS author_username, u_orig.display_name AS author_display_name, u_orig.icon_url AS author_icon_url,
				r.created_at AS event_time,
				u_reposter.username AS reposted_by_username
			FROM reposts r
			JOIN posts p ON r.post_id = p.id
			JOIN users u_orig ON p.user_id = u_orig.id
			JOIN users u_reposter ON r.user_id = u_reposter.id
		) AS timeline
		LEFT JOIN spoilers s ON timeline.id = s.post_id
		LEFT JOIN (SELECT post_id, COUNT(*) AS like_count FROM likes GROUP BY post_id) AS like_counts ON like_counts.post_id = timeline.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS repost_count FROM reposts GROUP BY post_id) AS repost_counts ON repost_counts.post_id = timeline.id
		LEFT JOIN (SELECT reply_to, COUNT(*) AS comment_count FROM posts WHERE reply_to IS NOT NULL GROUP BY reply_to) AS comment_counts ON comment_counts.reply_to = timeline.id
		LEFT JOIN (SELECT post_id, user_id FROM likes WHERE user_id = ?) AS user_likes ON user_likes.post_id = timeline.id
		LEFT JOIN (SELECT post_id, user_id FROM reposts WHERE user_id = ?) AS user_reposts ON user_reposts.post_id = timeline.id
		ORDER BY timeline.event_time DESC
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
			&p.UserActions.Liked,
			&p.UserActions.Reposted,
			&p.RepostedBy,
			&p.SpoilerWord, // ネタバレキーワード用のスキャン先を追加
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}

// 【変更後】リプライ取得（ネタバレキーワード対応）
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
         CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END AS user_reposted,
         s.keyword
       FROM posts p
       LEFT JOIN spoilers s ON p.id = s.post_id
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
			&p.SpoilerWord, // ネタバレキーワード用のスキャン先を追加
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}

// 【変更後】特定ユーザーの投稿とリポストを時系列で取得（ネタバレキーワード対応）
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
			timeline.id,
			timeline.text,
			timeline.reply_to,
			timeline.repost_ref,
			timeline.media_type,
			timeline.media_url,
			timeline.created_at,
			timeline.author_id,
			timeline.author_username,
			timeline.author_display_name,
			timeline.author_icon_url,
			COALESCE(like_counts.like_count, 0),
			COALESCE(repost_counts.repost_count, 0),
			COALESCE(comment_counts.comment_count, 0),
			CASE WHEN user_likes.user_id IS NOT NULL THEN 1 ELSE 0 END,
			CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END,
			timeline.reposted_by_username,
			s.keyword
		FROM (
			-- targetUserが投稿したもの
			SELECT
				p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
				u.id AS author_id, u.username AS author_username, u.display_name AS author_display_name, u.icon_url AS author_icon_url,
				p.created_at AS event_time,
				CAST(NULL AS CHAR(50)) COLLATE utf8mb4_0900_ai_ci AS reposted_by_username
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.user_id = ? -- targetUserId
		
			UNION ALL
		
			-- targetUserがリポストしたもの
			SELECT
				p.id, p.text, p.reply_to, p.repost_ref, p.media_type, p.media_url, p.created_at,
				u_orig.id AS author_id, u_orig.username AS author_username, u_orig.display_name AS author_display_name, u_orig.icon_url AS author_icon_url,
				r.created_at AS event_time,
				u_reposter.username AS reposted_by_username
			FROM reposts r
			JOIN posts p ON r.post_id = p.id
			JOIN users u_orig ON p.user_id = u_orig.id
			JOIN users u_reposter ON r.user_id = u_reposter.id
			WHERE r.user_id = ? -- targetUserId
		) AS timeline
		LEFT JOIN spoilers s ON timeline.id = s.post_id
		LEFT JOIN (SELECT post_id, COUNT(*) AS like_count FROM likes GROUP BY post_id) AS like_counts ON like_counts.post_id = timeline.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS repost_count FROM reposts GROUP BY post_id) AS repost_counts ON repost_counts.post_id = timeline.id
		LEFT JOIN (SELECT reply_to, COUNT(*) AS comment_count FROM posts WHERE reply_to IS NOT NULL GROUP BY reply_to) AS comment_counts ON comment_counts.reply_to = timeline.id
		LEFT JOIN (SELECT post_id, user_id FROM likes WHERE user_id = ?) AS user_likes ON user_likes.post_id = timeline.id -- viewerUserId
		LEFT JOIN (SELECT post_id, user_id FROM reposts WHERE user_id = ?) AS user_reposts ON user_reposts.post_id = timeline.id -- viewerUserId
		ORDER BY timeline.event_time DESC
		LIMIT ? OFFSET ?
    `, targetUserId, targetUserId, viewerUserId, viewerUserId, limit, offset)
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
			&p.RepostedBy,
			&p.SpoilerWord, // ネタバレキーワード用のスキャン先を追加
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}

func (r *PostsRepository) SearchPostsByKeywordWithStats(
	ctx context.Context,
	dbtx repository.DBTX,
	viewerUserId string,
	keyword string,
	limit int,
	offset int,
) (*[]model.PostWithUserAndCounts, error) {
	rows, err := dbtx.QueryContext(ctx, `
		SELECT
			p.id,
			p.text,
			p.reply_to,
			p.repost_ref,
			p.media_type,
			p.media_url,
			p.created_at,
			u.id AS author_id,
			u.username AS author_username,
			u.display_name AS author_display_name,
			u.icon_url AS author_icon_url,
			COALESCE(like_counts.like_count, 0),
			COALESCE(repost_counts.repost_count, 0),
			COALESCE(comment_counts.comment_count, 0),
			CASE WHEN user_likes.user_id IS NOT NULL THEN 1 ELSE 0 END,
			CASE WHEN user_reposts.user_id IS NOT NULL THEN 1 ELSE 0 END,
			NULL AS reposted_by_username,
			s.keyword
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN spoilers s ON p.id = s.post_id
		LEFT JOIN (SELECT post_id, COUNT(*) AS like_count FROM likes GROUP BY post_id) AS like_counts ON like_counts.post_id = p.id
		LEFT JOIN (SELECT post_id, COUNT(*) AS repost_count FROM reposts GROUP BY post_id) AS repost_counts ON repost_counts.post_id = p.id
		LEFT JOIN (SELECT reply_to, COUNT(*) AS comment_count FROM posts WHERE reply_to IS NOT NULL GROUP BY reply_to) AS comment_counts ON comment_counts.reply_to = p.id
		LEFT JOIN (SELECT post_id, user_id FROM likes WHERE user_id = ?) AS user_likes ON user_likes.post_id = p.id
		LEFT JOIN (SELECT post_id, user_id FROM reposts WHERE user_id = ?) AS user_reposts ON user_reposts.post_id = p.id
		WHERE p.text LIKE ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`, viewerUserId, viewerUserId, "%"+keyword+"%", limit, offset)
	fmt.Println("検索キーワード:", keyword)
	fmt.Println("LIKEクエリ:", "%"+keyword+"%")

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
			&p.RepostedBy,  // NULL固定
			&p.SpoilerWord, // ネタバレ
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return &results, nil
}
