import { createContext, useState } from "react";
import { getCookie } from "./service/cookie";

export interface State {
    twitchClientId: string,
    baseUrl: string,
    scToken: string | null,
}

export type Action = Record<string, unknown>;

const defaultContext = {
    state: {
        twitchClientId: process.env.REACT_APP_TWITCH_CLIENT_ID,
        baseUrl: process.env.REACT_APP_BASE_URL,
        scToken: getCookie("scToken") ,
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
