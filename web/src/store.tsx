import { createContext, useState } from "react";

export interface State {
    apiBaseUrl: string,
    twitchClientId: string,
    baseUrl: string,
    accessToken: string | null,
}

export type Action = Record<string, unknown>;

const defaultContext = {
    state: {
        apiBaseUrl: process.env.REACT_APP_API_BASE_URL,
        twitchClientId: process.env.REACT_APP_TWITCH_CLIENT_ID,
        baseUrl: process.env.REACT_APP_BASE_URL,
        accessToken: window.localStorage.getItem("accessToken")
    } as State,
    setState: (state: State) => { },
};

const store = createContext(defaultContext);
const { Provider } = store;

const StateProvider = ({ children }: { children: JSX.Element }): JSX.Element => {
    const [state, setState] = useState({ ...defaultContext.state });

    return <Provider value={{ state, setState }}>{children}</Provider>;
};

export { store, StateProvider };
