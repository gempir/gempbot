import { useState } from "react";
import { Link } from "react-router-dom";
import styled from "styled-components";
import { store } from "./../store";

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
    const scToken = store.useState(s => s.scToken);

    if (scToken) {
        return <LoggedIn />;
    }

    return <LoginContainer href={createOAuthUrl().toString()}>Login</LoginContainer>;
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
    const [loginText, setLoginText] = useState("Logged In");

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
        <LoggedInContainer href={createOAuthUrl().toString()} onMouseEnter={() => setLoginText("force login")} onMouseLeave={() => setLoginText("Logged In")}>
            {loginText}
        </LoggedInContainer>
    </ButtonsContainer>;
}

function createOAuthUrl(): URL {
    const { apiBaseUrl, twitchClientId } = store.getRawState();

    const url = new URL("https://id.twitch.tv/oauth2/authorize")
    url.searchParams.set("client_id", twitchClientId);
    url.searchParams.set("redirect_uri", apiBaseUrl + "/api/callback");
    url.searchParams.set("response_type", "code");
    url.searchParams.set("scope", "channel:read:redemptions channel:manage:redemptions");

    return url;
}