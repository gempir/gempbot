import { useEffect, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { ChannelPointReward, RewardTypes } from "../types/Rewards";

export function useChannelPointReward(userID: string, type: RewardTypes, defaultReward: ChannelPointReward): [ChannelPointReward, (reward: ChannelPointReward) => void] {
    const [reward, setReward] = useState<ChannelPointReward>(defaultReward);

    const fetchReward = () => {
        doFetch(Method.GET, `/api/reward/${userID}/type/${type}`).then(setReward)
    }

    useEffect(fetchReward, [userID, type]);

    const updateReward = (reward: ChannelPointReward) => {
        doFetch(Method.GET, `/api/reward/${userID}/type/${type}`, reward).then(fetchReward);
    }

    return [reward, updateReward];
}