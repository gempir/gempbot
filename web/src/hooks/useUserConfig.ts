import { useEffect, useState } from "react";
import { useDebounce } from "react-use";
import { doFetch, Method } from "../service/doFetch";
import { store } from "../store";

export interface UserConfig {
    Editors: Array<string>;
    Protected: Protected;
    Rewards: Rewards;
}

export interface Rewards {
    Bttv: null | BttvReward
}

export interface BttvReward {
    title: string;
    cost: number;
    backgroundColor: string;
    maxPerStream: number;
    maxPerUserPerStream: number;
    globalCooldownSeconds: number;
}

export interface Protected {
    EditorFor: Array<string>;
}

export type SetUserConfig = (userConfig: UserConfig | null) => void;

export function useUserConfig(onSave: () => void = () => { }, onError: () => void = () => { }): [UserConfig | null | undefined, SetUserConfig] {
    const [userConfig, setUserConfig] = useState<UserConfig | null | undefined>(undefined);
    const [changeCounter, setChangeCounter] = useState(0);
    const managing = store.useState(s => s.managing);

    const fetchConfig = () => {
        let endPoint = "/api/userConfig";
        if (managing) {
            endPoint += `?managing=${managing}`;
        }
        doFetch(Method.GET, endPoint).then((userConfig) => setUserConfig(userConfig)).catch(onError)
    };

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchConfig, [managing]);

    useDebounce(() => {
        if (changeCounter && userConfig) {
            let endPoint = "/api/userConfig";
            if (managing) {
                endPoint += `?managing=${managing}`;
            }
            doFetch(Method.POST, endPoint, userConfig).then(onSave).catch(onError);
        } else if (changeCounter && userConfig === null) {
            doFetch(Method.DELETE, "/api/userConfig").then(fetchConfig).catch(onError);
        }
    }, 500, [changeCounter]);

    const setCfg = (config: UserConfig | null) => {
        setUserConfig(config);
        setChangeCounter(changeCounter + 1);
    };

    return [userConfig, setCfg]
}