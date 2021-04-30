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

    useEffect(() => {
        if (accessToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { accessToken } })
                .then(response => response.json())
                .then((userConfig) => setUserConfig(userConfig));
        }
    }, [accessToken, apiBaseUrl]);


    useDebounce(() => {
        if (userConfig && accessToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { accessToken }, method: "POST", body: JSON.stringify(userConfig) }).then(onSave);
        }
    }, 500, [userConfig]);

    return [userConfig, setUserConfig]
}