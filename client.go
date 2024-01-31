package main

import (
	"github.com/tmstn/pinboard"
)

type user struct {
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`
}

type client struct {
	user user
	pb   *pinboard.Client
}

func (c *client) authenticate() error {
	_, err := c.pb.User.Secret()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) query(params *pinboard.PostsAllOptions) ([]*pinboard.Post, error) {
	posts, err := c.pb.Posts.All(params)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (c *client) delete(post *pinboard.Post) error {
	err := c.pb.Posts.Delete(post.Href.String())
	if err != nil {
		return err
	}

	return nil
}

func newClient(user user) *client {
	return &client{
		user: user,
		pb:   pinboard.New(user.ApiKey),
	}
}
