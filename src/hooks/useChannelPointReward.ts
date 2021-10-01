import { useEffect, useState } from "react";
import { doFetch, Method, RejectReason } from "../service/doFetch";
import { useStore } from "../store";
import { ChannelPointReward, RewardTypes, RawBttvChannelPointReward } from "../types/Rewards";

export function useChannelPointReward(userID: string, type: RewardTypes, defaultReward: ChannelPointReward, onUpdate: () => void): [ChannelPointReward, (reward: ChannelPointReward) => void, () => void] {
    const [reward, setReward] = useState<ChannelPointReward>(defaultReward);

    const fetchReward = () => {
        const endPoint = "/api/reward";
        const searchParams = new URLSearchParams();
        searchParams.append("type", type);

        doFetch(Method.GET, endPoint, searchParams).then((response: RawBttvChannelPointReward) => ({...response, AdditionalOptionsParsed: response.AdditionalOptions !== "" ? JSON.parse(response.AdditionalOptions) : defaultReward.AdditionalOptionsParsed})).then(setReward).catch(reason => {
            if (reason !== RejectReason.NotFound) {
                throw new Error(reason);
            }
            setReward(defaultReward);
        }).then(onUpdate)
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchReward, [userID, type, defaultReward]);

    const updateReward = (reward: ChannelPointReward) => {
        const endPoint = "/api/reward";
        const searchParams = new URLSearchParams();
        searchParams.append("type", type);

        doFetch(Method.POST, endPoint, searchParams, reward).then(fetchReward);
    }

    const deleteReward = () => {
        const endPoint = "/api/reward";
        const searchParams = new URLSearchParams();
        searchParams.append("type", type);

        doFetch(Method.DELETE, endPoint, searchParams, reward).then(fetchReward);
    }

    return [reward, updateReward, deleteReward];
}