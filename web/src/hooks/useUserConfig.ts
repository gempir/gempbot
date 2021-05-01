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


export function useUserConfig(onSave: () => void): [UserConfig | null, (userConfig: UserConfig) => void] {
    const { accessToken, apiBaseUrl } = useContext(store).state;

    const [userConfig, setUserConfig] = useState<UserConfig | null>(null);
    const [changeCounter, setChangeCounter] = useState(0);

    useEffect(() => {
        if (accessToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { accessToken } })
                .then(response => response.json())
                .then((userConfig) => setUserConfig(userConfig));
        }
    }, [accessToken, apiBaseUrl]);


    useDebounce(() => {
        if (changeCounter && userConfig && accessToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { accessToken }, method: "POST", body: JSON.stringify(userConfig) }).then(onSave);
        }
    }, 500, [changeCounter]);

    const setCfg = (config: UserConfig) => {
        setUserConfig(config)
        setChangeCounter(changeCounter + 1);
    };

    return [userConfig, setCfg]
}