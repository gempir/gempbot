import { useContext, useEffect, useState } from "react";
import { useDebounce } from "react-use";
import { store } from "../store";

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
    const { accessToken, apiBaseUrl } = useContext(store).state;

    const [userConfig, setUserConfig] = useState<UserConfig | null | undefined>(undefined);
    const [changeCounter, setChangeCounter] = useState(0);

    const fetchConfig = () => {
        if (accessToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { accessToken } })
                .then(response => response.json())
                .then((userConfig) => setUserConfig(userConfig));
        }
    };

    useEffect(fetchConfig, [accessToken, apiBaseUrl]);


    useDebounce(() => {
        if (changeCounter && userConfig && accessToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { accessToken }, method: "POST", body: JSON.stringify(userConfig) }).then(onSave);
        } else if (changeCounter && userConfig === null && accessToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { accessToken }, method: "DELETE"}).then(fetchConfig);
        }
    }, 500, [changeCounter]);

    const setCfg = (config: UserConfig | null) => {
        setUserConfig(config);
        setChangeCounter(changeCounter + 1);
    };

    return [userConfig, setCfg]
}