export function reducer(state, action) {
    switch (action.type) {
        case "SET_CHANNELS":
            return {...state, channels: action.channels};
        default:
            return {...state};
    }
};