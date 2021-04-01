export default function (channels) {
    return function (dispatch, getState) {
        dispatch({
            type: "SET_CHANNELS",
            channels: {...getState().channels, ...channels}
        });
    };
}