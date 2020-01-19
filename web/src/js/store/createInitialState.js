import loadState from "../storage/loadState";

export default () => {

    const persistedState = loadState();

    return {
        loading: false,
        ...persistedState,
    }
}