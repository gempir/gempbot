import { useLayoutEffect } from 'react';
import create from 'zustand';
import createContext from 'zustand/context';

export interface ScTokenContent {
    Login: string;
    UserID: string;
}

export interface Store {
    twitchClientId: string;
    apiBaseUrl: string;
    baseUrl: string;
    yjsWsUrl: string;
    scToken: string  | null;
    scTokenContent: ScTokenContent  | null;
    setScToken: (token: string) => void;
    setManaging: (managing: string | null) => void;
    managing: string | null;
}

let store: any;

const initialState = {
    twitchClientId: "",
    apiBaseUrl: "",
    yjsWsUrl: "",
    baseUrl: "",
    scToken: "",
    scTokenContent: null,
    managing: null,
}

const zustandContext = createContext<Store>()

export const StoreProvider = zustandContext.Provider;
export const useStore: typeof zustandContext.useStore = zustandContext.useStore

export const initializeStore = (preloadedState = {}) => {
    return create<Store>((set) => ({
        ...initialState,
        ...preloadedState,
        setScToken: (token: string) => set(state => ({ scToken: token })),
        setManaging: (managing: string | null) => set(state => ({ managing: managing })),
    }))
}

export function useCreateStore(serverInitialState: Store) {
    // Server side code: For SSR & SSG, always use a new store.
    if (typeof window === 'undefined') {
        return () => initializeStore(serverInitialState)
    }
    // End of server side code

    // Client side code:
    // Next.js always re-uses same store regardless of whether page is a SSR or SSG or CSR type.
    const isReusingStore = Boolean(store)
    store = store ?? initializeStore(serverInitialState)
    // When next.js re-renders _app while re-using an older store, then replace current state with
    // the new state (in the next render cycle).
    // (Why next render cycle? Because react cannot re-render while a render is already in progress.
    // i.e. we cannot do a setState() as that will initiate a re-render)
    //
    // eslint complaining "React Hooks must be called in the exact same order in every component render"
    // is ignorable as this code runs in same order in a given environment (i.e. client or server)
    // eslint-disable-next-line react-hooks/rules-of-hooks
    useLayoutEffect(() => {
        // serverInitialState is undefined for CSR pages. It is up to you if you want to reset
        // states on CSR page navigation or not. I have chosen not to, but if you choose to,
        // then add `serverInitialState = getDefaultInitialState()` here.
        if (serverInitialState && isReusingStore) {
            store.setState(
                {
                    // re-use functions from existing store
                    ...store.getState(),
                    // but reset all other properties.
                    ...serverInitialState,
                },
                true // replace states, rather than shallow merging
            )
        }
    })

    return () => store
}
