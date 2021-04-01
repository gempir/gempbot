import loadState from "../storage/loadState";

export default () => {

    const persistedState = loadState();

    console.log(process.env);

    return {
        loading: false,
        channels: {},
        apiBaseUrl: process.env.REACT_APP_API_BASE_URL,
        ...persistedState,
    }
}