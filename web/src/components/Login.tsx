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
    const {state, setState} = useContext(store);

    useEffect(() => {
        const hash = window.location.hash;
        window.location.hash = "";
        if (hash) {
            const reg = /#access_token=(\w*)&/ig;
            const match = reg.exec(hash);
            if (!match || typeof match[1] === "undefined") {
                return;
            }

            window.localStorage.setItem("accessToken", match[1]);
            setState({accessToken: match[1]});

            
        }
    }, [state.apiBaseUrl, setState]);

    useEffect(() => {
        fetch(state.apiBaseUrl + "/api/oauth", {
            method: 'post',
            body: JSON.stringify({accessToken: state.accessToken})
        })
    }, [state.accessToken, state.apiBaseUrl]);

    if (state.accessToken) {
        return <LoggedIn />;
    }

    const url = new URL("https://id.twitch.tv/oauth2/authorize")
    url.searchParams.set("client_id", state.twitchClientId);
    url.searchParams.set("redirect_uri", state.baseUrl);
    url.searchParams.set("response_type", "token");
    url.searchParams.set("scope", "channel:read:redemptions");

    return <LoginContainer href={url.toString()}>Login</LoginContainer>
}

const LoggedInContainer = styled.a`
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

function LoggedIn() {
    return <LoggedInContainer>
        Logged In
    </LoggedInContainer>;
}