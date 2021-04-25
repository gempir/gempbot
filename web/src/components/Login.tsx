import { useContext, useEffect } from "react";
import { store } from "./../store";
import styled from "styled-components";

const LoginContainer = styled.a`
    position: absolute;
    display: block;
    color: white;
    top: 1rem;
    right: 1rem;
    padding: 1rem 2rem;
    text-decoration: none;
    font-weight: bold;
    border-radius: 3px;
    background: var(--twitch);
`;

export function Login() {
    const state = useContext(store).state;

    const url = new URL("https://id.twitch.tv/oauth2/authorize")
    url.searchParams.set("client_id", state.twitchClientId);
    url.searchParams.set("redirect_uri", state.baseUrl);
    url.searchParams.set("response_type", "token");
    url.searchParams.set("scope", "channel:read:redemptions");


    useEffect(() => {
        if (window.location.hash) {
            const reg = /#access_token=(\w*)&/ig;
            const match = reg.exec(window.location.hash);
            if (!match || typeof match[1] === "undefined") {
                return;
            }

            fetch(state.apiBaseUrl + "/api/oauth", {
                method: 'post',
                body: JSON.stringify({accessToken: match[1]})
            })
        }
    }, [state.apiBaseUrl])

    return <LoginContainer href={url.toString()}>Login</LoginContainer>
}