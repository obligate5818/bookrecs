package discord

import (
	"context"
	"encoding/json"
	"net/http"
)

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func GetUser(ctx context.Context, tokenType, accessToken string) (*DiscordUser, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me", nil)
	req.Header.Set("Authorization", tokenType+" "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserGuilds(ctx context.Context, tokenType, accessToken string) ([]string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me/guilds", nil)
	req.Header.Set("Authorization", tokenType+" "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var guilds []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return nil, err
	}

	ids := make([]string, len(guilds))
	for i, g := range guilds {
		ids[i] = g.ID
	}
	return ids, nil
}
