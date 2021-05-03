import { createContext, useState } from "react";
import { deleteCookie, getCookie, setCookie } from "./service/cookie";

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
        scToken: getCookie("scToken")
    } as State,
    setState: (state: State) => { },
    setScToken: (scToken: string | null) => { },
};

const store = createContext(defaultContext);
const { Provider } = store;

const StateProvider = ({ children }: { children: JSX.Element }): JSX.Element => {
    const [state, setState] = useState({ ...defaultContext.state });

    const setScToken = (scToken: string | null) => {
        if (scToken === null) {
            deleteCookie("scToken")
        } else {
            setCookie("scToken", scToken, 365);
        }

        if (scToken !== state.scToken) {
            setState({...state, scToken});
        }
    };

    return <Provider value={{ state, setState, setScToken }}>{children}</Provider>;
};

export { store, StateProvider };
