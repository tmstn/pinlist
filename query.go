package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tmstn/pinboard"
	"github.com/urfave/cli/v2"
)

func query(cCtx *cli.Context) error {
	username := cCtx.String("username")

	passphrase, err := readPassphrase(cCtx)
	if err != nil {
		return err
	}

	pl, err := readAuth(username, passphrase)
	if err != nil {
		return err
	}

	params := &pinboard.PostsAllOptions{
		Start:   0,
		Results: 250,
	}

	tags := cCtx.StringSlice("tag")
	if len(tags) > 0 {
		params.Tag = tags
	}

	page := cCtx.Int("page")
	if page > 0 {
		params.Results = page
	}

	posts := []*pinboard.Post{}
	more := true
	count := 0
	remove := cCtx.Bool("remove")
	max := cCtx.Int("max")
	url := cCtx.String("url")

	for more {
		p, err := pl.query(params)
		if err != nil {
			return err
		}

		for _, post := range p {
			if strings.TrimSpace(url) == "" ||
				strings.Contains(post.Href.String(), url) {

				posts = append(posts, post)
				count += 1
				if max > 0 && count == max {
					break
				}
			}
		}

		if len(p) < params.Results || (max > 0 && count == max) {
			break
		}

		params.Start += params.Results
		time.Sleep(time.Millisecond * 100)
	}

	for _, post := range posts {
		fmt.Println(post.Href.String())
	}

	if remove {
		for _, post := range posts {
			err := pl.delete(post)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
