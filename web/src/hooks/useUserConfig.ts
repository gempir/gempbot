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
    const { scToken, apiBaseUrl } = useContext(store).state;
    const { setScToken } = useContext(store);

    const [userConfig, setUserConfig] = useState<UserConfig | null | undefined>(undefined);
    const [changeCounter, setChangeCounter] = useState(0);

    const fetchConfig = () => {
        if (scToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { Authorization: "Bearer " + scToken } })
                .then(response => checkToken(setScToken, response))
                .then(response => response.json())
                .then((userConfig) => setUserConfig(userConfig))
                .catch();
        }
    };

    useEffect(fetchConfig, [scToken, apiBaseUrl, setScToken]);


    useDebounce(() => {
        if (changeCounter && userConfig && scToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { Authorization: "Bearer " + scToken }, method: "POST", body: JSON.stringify(userConfig) })
            .then(response => checkToken(setScToken, response))
            .then(onSave);
        } else if (changeCounter && userConfig === null && scToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { Authorization: "Bearer " + scToken }, method: "DELETE"})
            .then(response => checkToken(setScToken, response))
            .then(fetchConfig);
        }
    }, 500, [changeCounter]);

    const setCfg = (config: UserConfig | null) => {
        setUserConfig(config);
        setChangeCounter(changeCounter + 1);
    };

    return [userConfig, setCfg]
}

function checkToken(setScToken: (scToken: string | null) => void, response: Response) {
    if (response.status === 403) {
        setScToken(null);
    }

    return response
}