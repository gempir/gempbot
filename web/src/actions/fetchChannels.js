import addChannels from "./addChannels";
import persistState from "../storage/persistState";

export default function (channelIds) {
    return function (dispatch, getState) {

        const channelIdsToFetch = [];
        const state = getState();
        for (const channelId of channelIds) {
            if (!state.channels.hasOwnProperty(channelId)) {
                channelIdsToFetch.push(channelId);
            }
        }

        if (channelIdsToFetch.length === 0) {
            return;
        }

        return new Promise((resolve, reject) => {
            fetch(`${getState().apiBaseUrl}/api/channel?channelids=${channelIdsToFetch.join(",")}`).then((response) => {
                if (response.status >= 200 && response.status < 300) {
                    return response
                } else {
                    const error = new Error(response.statusText);
                    error.response = response;
                    throw error
                }
            }).then((response) => {
                return response.json();
            }).then((json) => {
                dispatch(addChannels(json)).then(() => persistState(getState()));

                resolve();
            }).catch(err => {
                reject(err);
            });
        });
    };
}