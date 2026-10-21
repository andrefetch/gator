package commands

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andrefetch/gator/internal/database"
)

func handlerBrowse(state *State, cmd Command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: browse [limit]")
	}

	limit := int32(2)
	if len(cmd.Args) == 1 {
		parsed, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil || parsed <= 0 {
			return fmt.Errorf("limit must be a positive integer no greater than 2147483647")
		}
		limit = int32(parsed)
	}

	posts, err := state.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("get posts for user: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found from the feeds you follow.")
		return nil
	}

	for _, post := range posts {
		fmt.Printf("%s\n%s\n", post.Title, post.Url)
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", post.PublishedAt.Time.UTC().Format(time.RFC3339))
		}
		if post.Description.Valid {
			fmt.Println(post.Description.String)
		}
		fmt.Println()
	}

	return nil
}
