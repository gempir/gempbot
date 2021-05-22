import { useEffect, useState } from "react";
import { useDebounce } from "react-use";
import { doFetch, Method } from "../service/doFetch";
import { store } from "../store";

export interface UserConfig {
    Editors: Array<string>;
    Protected: Protected;
    Rewards: Rewards;
    CurrentUserID?: string;
}

export interface Rewards {
    Bttv: null | BttvReward
}

export interface BttvReward {
    title: string;
    prompt?: string;
    cost: number;
    backgroundColor?: string;
    isMaxPerStreamEnabled?: boolean;
    maxPerStream?: number;
    isUserInputRequired?: boolean;
    isMaxPerUserPerStreamEnabled?: boolean;
    maxPerUserPerStream?: number;
    isGlobalCooldownEnabled?: boolean;
    globalCooldownSeconds?: number;
    shouldRedemptionsSkipRequestQueue?: boolean;
    enabled?: boolean;
    isDefault: boolean;
    ID?: string;
}

export interface Protected {
    EditorFor: Array<string>;
}

export type SetUserConfig = (userConfig: UserConfig | null) => void;

export function useUserConfig(onUserConfigChange: () => void = () => {}): [UserConfig | null | undefined, SetUserConfig, () => void] {
    const [userConfig, setUserConfig] = useState<UserConfig | null | undefined>(undefined);
    const [changeCounter, setChangeCounter] = useState(0);
    const managing = store.useState(s => s.managing);

    const fetchConfig = () => {
        let endPoint = "/api/userConfig";
        if (managing) {
            endPoint += `?managing=${managing}`;
        }
        doFetch(Method.GET, endPoint).then((userConfig) => setUserConfig(userConfig))
    };

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(fetchConfig, [managing]);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(onUserConfigChange, [userConfig]);

    useDebounce(() => {
        if (changeCounter && userConfig) {
            let endPoint = "/api/userConfig";
            if (managing) {
                endPoint += `?managing=${managing}`;
            }
            doFetch(Method.POST, endPoint, userConfig).then(fetchConfig);
        } else if (changeCounter && userConfig === null) {
            doFetch(Method.DELETE, "/api/userConfig");
        }
    }, 500, [changeCounter]);

    const setCfg = (config: UserConfig | null) => {
        setUserConfig(config);
        setChangeCounter(changeCounter + 1);
    };

    return [userConfig, setCfg, fetchConfig]
}