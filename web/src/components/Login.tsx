import { useContext } from "react";
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
    url.searchParams.set("redirect_uri", state.apiBaseUrl + "/oauth");
    url.searchParams.set("response_type", "token");
    url.searchParams.set("scope", "channel:read:redemptions");

    return <LoginContainer href={url.toString()}>Login</LoginContainer>
}