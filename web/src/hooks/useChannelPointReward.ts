import { useEffect, useState } from "react";
import { doFetch, Method, RejectReason } from "../service/doFetch";
import { ChannelPointReward, RewardTypes, RawBttvChannelPointReward } from "../types/Rewards";

export function useChannelPointReward(userID: string, type: RewardTypes, defaultReward: ChannelPointReward, onUpdate: () => void): [ChannelPointReward, (reward: ChannelPointReward) => void, () => void] {
    const [reward, setReward] = useState<ChannelPointReward>(defaultReward);

    const fetchReward = () => {
        doFetch(Method.GET, `/api/reward/${userID}/type/${type}`).then((response: RawBttvChannelPointReward) => ({...response, AdditionalOptionsParsed: response.AdditionalOptions !== "" ? JSON.parse(response.AdditionalOptions) : defaultReward.AdditionalOptionsParsed})).then(setReward).catch(reason => {
            if (reason !== RejectReason.NotFound) {
                throw new Error(reason);
            }
            setReward(defaultReward);
        }).then(onUpdate)
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchReward, [userID, type, defaultReward]);

    const updateReward = (reward: ChannelPointReward) => {
        doFetch(Method.POST, `/api/reward/${userID}`, undefined, reward).then(fetchReward);
    }

    const deleteReward = () => {
        doFetch(Method.DELETE, `/api/reward/${userID}/type/${type}`, undefined, reward).then(fetchReward);
    }

    return [reward, updateReward, deleteReward];
}