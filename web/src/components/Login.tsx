import { useContext, useEffect, useState } from "react";
import { store } from "./../store";
import styled from "styled-components";
import { Link } from "react-router-dom";

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

    &:hover {
        background: var(--twitch-dark);
    }
`;

export function Login() {
    const {state, setAccessToken} = useContext(store);

    useEffect(() => {
        const hash = window.location.hash;
        window.location.hash = "";
        if (hash) {
            const reg = /#access_token=(\w*)&/ig;
            const match = reg.exec(hash);
            if (!match || typeof match[1] === "undefined") {
                return;
            }

            if (match[1] ){
                window.localStorage.setItem("accessToken", match[1]);
                setAccessToken(match[1]);
            }
        }
    }, [state.apiBaseUrl, setAccessToken]);

    useEffect(() => {
        if (state.accessToken) {
            fetch(state.apiBaseUrl + "/api/oauth", {
                method: 'post',
                body: JSON.stringify({accessToken: state.accessToken})
            })
            // validate accessToken
        }
    }, [state.accessToken, state.apiBaseUrl]);

    if (state.accessToken) {
        return <LoggedIn />;
    }

    const url = new URL("https://id.twitch.tv/oauth2/authorize")
    url.searchParams.set("client_id", state.twitchClientId);
    url.searchParams.set("redirect_uri", state.baseUrl);
    url.searchParams.set("response_type", "token");
    url.searchParams.set("scope", "channel:read:redemptions");

    return <LoginContainer href={url.toString()}>Login</LoginContainer>;
}

const ButtonsContainer = styled.div`
    position: absolute;
    top: 1rem;
    right: 1rem;
    display: flex;

    a {
        text-decoration: none;
    }
`;

const LoggedInContainer = styled.a`
    display: block;
    color: white;
    width: 150px;
    padding: 1rem 2rem;
    text-decoration: none;
    font-weight: bold;
    border-radius: 3px;
    background: var(--twitch);
    cursor: pointer;

    &:hover {
        background: var(--twitch-dark);
    }
`;

const DashboardButton = styled.div`
    display: block;
    color: white;
    margin-right: 1rem;
    padding: 1rem 2rem;
    text-decoration: none;
    font-weight: bold;
    border-radius: 3px;
    background: var(--theme2-dark);
    cursor: pointer;

    &.dashboard {
        background: var(--theme);

        &:hover {
        background: var(--theme-bright);
    }
    }

    &:hover {
        background: var(--theme2);
    }
`;

function LoggedIn() {
    const {state} = useContext(store);
    const [loginText, setLoginText] = useState("Logged In");

    const url = new URL("https://id.twitch.tv/oauth2/authorize")
    url.searchParams.set("client_id", state.twitchClientId);
    url.searchParams.set("redirect_uri", state.baseUrl);
    url.searchParams.set("response_type", "token");
    url.searchParams.set("scope", "channel:read:redemptions");

    return <ButtonsContainer>
        <Link to="/">
            <DashboardButton>
                Home
            </DashboardButton>
        </Link>
        <Link to="/dashboard">
            <DashboardButton className="dashboard">
                Dashboard
            </DashboardButton>
        </Link>
        <LoggedInContainer href={url.toString()} onMouseEnter={() => setLoginText("force login")} onMouseLeave={() => setLoginText("Logged In")}>
            {loginText}
        </LoggedInContainer>
    </ButtonsContainer>;
}