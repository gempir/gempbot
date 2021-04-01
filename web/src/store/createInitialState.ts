import {loadState} from "../storage/loadState";

export function createInitialState() {
    const persistedState = loadState();

    return {
        loading: false,
        channels: {},
        apiBaseUrl: process.env.REACT_APP_API_BASE_URL,
        ...persistedState,
    }
}