import { createContext, useState } from "react";

export interface State {
    apiBaseUrl: string,
    twitchClientId: string,
    baseUrl: string,
    scToken: string | null,
}

export type Action = Record<string, unknown>;

const defaultContext = {
    state: {
        apiBaseUrl: process.env.REACT_APP_API_BASE_URL,
        twitchClientId: process.env.REACT_APP_TWITCH_CLIENT_ID,
        baseUrl: process.env.REACT_APP_BASE_URL,
        scToken: window.localStorage.getItem("scToken")
    } as State,
    setState: (state: State) => { },
    setScToken: (scToken: string) => { },
};

const store = createContext(defaultContext);
const { Provider } = store;

const StateProvider = ({ children }: { children: JSX.Element }): JSX.Element => {
    const [state, setState] = useState({ ...defaultContext.state });

    const setScToken = (scToken: string) => {
        window.localStorage.setItem("scToken", scToken);
        setState({...state, scToken});
    };

    return <Provider value={{ state, setState, setScToken }}>{children}</Provider>;
};

export { store, StateProvider };
