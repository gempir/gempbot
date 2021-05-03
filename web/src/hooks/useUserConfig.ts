import { useEffect, useState } from "react";
import { useDebounce } from "react-use";
import { doFetch, Method } from "../service/doFetch";

export interface UserConfig {
    Redemptions: Redemptions;
}

export interface Redemptions {
    Bttv: Redemption;
}

export interface Redemption {
    Title: string;
    Active: boolean;
}


export function useUserConfig(onSave: () => void): [UserConfig | null | undefined, (userConfig: UserConfig | null) => void] {
    const [userConfig, setUserConfig] = useState<UserConfig | null | undefined>(undefined);
    const [changeCounter, setChangeCounter] = useState(0);

    const fetchConfig = () => {
        doFetch(Method.GET, "/api/userConfig").then((userConfig) => setUserConfig(userConfig))
    };

    useEffect(fetchConfig, []);

    useDebounce(() => {
        if (changeCounter && userConfig) {
            doFetch(Method.POST, "/api/userConfig", userConfig).then(onSave)
        } else if (changeCounter && userConfig === null) {
            doFetch(Method.DELETE, "/api/userConfig").then(fetchConfig)
        }
    }, 500, [changeCounter]);

    const setCfg = (config: UserConfig | null) => {
        setUserConfig(config);
        setChangeCounter(changeCounter + 1);
    };

    return [userConfig, setCfg]
}