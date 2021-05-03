import { useContext, useEffect, useState } from "react";
import { useDebounce } from "react-use";
import { checkToken } from "../service/checkToken";
import { handleResponse } from "../service/fetch";
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
                .then(handleResponse)
                .then((userConfig) => setUserConfig(userConfig))
                .catch((resp) => checkToken(setScToken, resp));
        }
    };

    useEffect(fetchConfig, [scToken, apiBaseUrl, setScToken]);


    useDebounce(() => {
        if (changeCounter && userConfig && scToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { Authorization: "Bearer " + scToken }, method: "POST", body: JSON.stringify(userConfig) })
            .then(handleResponse)
            .then(onSave)
            .catch((resp) => checkToken(setScToken, resp));
        } else if (changeCounter && userConfig === null && scToken) {
            fetch(apiBaseUrl + "/api/userConfig", { headers: { Authorization: "Bearer " + scToken }, method: "DELETE"})
            .then(handleResponse)
            .then(fetchConfig)
            .catch((resp) => checkToken(setScToken, resp));
        }
    }, 500, [changeCounter]);

    const setCfg = (config: UserConfig | null) => {
        setUserConfig(config);
        setChangeCounter(changeCounter + 1);
    };

    return [userConfig, setCfg]
}