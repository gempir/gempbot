export function createLoginUrl(apiBaseUrl: string, twitchClientId: string): URL {
    const url = new URL("https://id.twitch.tv/oauth2/authorize")
    url.searchParams.set("client_id", twitchClientId);
    url.searchParams.set("redirect_uri", apiBaseUrl + "/api/callback");
    url.searchParams.set("response_type", "code");
    url.searchParams.set("scope", ["channel:read:redemptions", "channel:manage:redemptions", "channel:read:predictions", "channel:manage:predictions", "moderation:read", "channel:bot", "user:write:chat", "user:bot", "moderator:manage:announcements"].join(" "));

    return url;
}