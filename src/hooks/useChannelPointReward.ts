import { useEffect, useState } from "react";
import { doFetch, Method, RejectReason } from "../service/doFetch";
import { ChannelPointReward, RawBttvChannelPointReward, RewardTypes } from "../types/Rewards";

export function useChannelPointReward(userID: string, type: RewardTypes, defaultReward: ChannelPointReward): [ChannelPointReward, (reward: ChannelPointReward) => void, () => void, string | null, boolean] {
    const [reward, setReward] = useState<ChannelPointReward>(defaultReward);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    const fetchReward = () => {
        setLoading(true);
        const endPoint = "/api/reward";
        const searchParams = new URLSearchParams();
        searchParams.append("type", type);

        doFetch(Method.GET, endPoint, searchParams).then((response: RawBttvChannelPointReward) => ({...response, AdditionalOptionsParsed: response.AdditionalOptions !== "" ? JSON.parse(response.AdditionalOptions) : defaultReward.AdditionalOptionsParsed})).then(setReward).catch(reason => {
            if (reason !== RejectReason.NotFound) {
                throw new Error(reason);
            }
            setReward(defaultReward);
        }).then(() => setLoading(false));
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchReward, [userID, type, defaultReward]);

    const updateReward = (reward: ChannelPointReward) => {
        setLoading(true);
        const endPoint = "/api/reward";
        const searchParams = new URLSearchParams();
        searchParams.append("type", type);

        doFetch(Method.POST, endPoint, searchParams, reward).then(() => setErrorMessage(null)).then(fetchReward).catch(setErrorMessage).finally(() => setLoading(false));
    }

    const deleteReward = () => {
        setLoading(true);
        const endPoint = "/api/reward";
        const searchParams = new URLSearchParams();
        searchParams.append("type", type);

        doFetch(Method.DELETE, endPoint, searchParams, reward).then(() => setErrorMessage(null)).then(fetchReward).catch(setErrorMessage).finally(() => setLoading(false));
    }

    return [reward, updateReward, deleteReward, errorMessage, loading];
}
