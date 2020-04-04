import loadState from "../storage/loadState";

export default () => {

    const persistedState = loadState();

    return {
        loading: false,
        channels: {},
        apiBaseUrl: "http://localhost:8000",
        ...persistedState,
    }
}