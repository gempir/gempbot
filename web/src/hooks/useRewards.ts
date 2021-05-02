import { useContext, useEffect, useState } from "react";
import { checkToken } from "../service/checkToken";
import { store } from "../store";

export interface Reward {
    broadcaster_name: string;
    broadcaster_login: string;
    broadcaster_id: string;
    id: string;
    image: Image | null;
    background_color: string;
    is_enabled: boolean;
    cost: number;
    title: string;
    prompt: string;
    is_user_input_required: boolean;
    max_per_stream_setting: MaxPerStreamSetting;
    max_per_user_per_stream_setting: MaxPerUserPerStreamSetting;
    global_cooldown_setting: GlobalCooldownSetting;
    is_paused: boolean;
    is_in_stock: boolean;
    default_image: Image;
    should_redemptions_skip_request_queue: boolean;
    redemptions_redeemed_current_stream: null;
    cooldown_expires_at: null;
}

export interface Image {
    url_1x: string;
    url_2x: string;
    url_4x: string;
}

export interface GlobalCooldownSetting {
    is_enabled: boolean;
    global_cooldown_seconds: number;
}

export interface MaxPerStreamSetting {
    is_enabled: boolean;
    max_per_stream: number;
}

export interface MaxPerUserPerStreamSetting {
    is_enabled: boolean;
    max_per_user_per_stream: number;
}

export function useRewards() {
    const { scToken, apiBaseUrl } = useContext(store).state;
    const { setScToken } = useContext(store);

    const [rewards, setRewards] = useState<Array<Reward>>([]);

    const fetchRewards = () => {
        if (scToken) {
            fetch(apiBaseUrl + "/api/rewards", { headers: { Authorization: "Bearer " + scToken } })
                .then(response => checkToken(setScToken, response))
                .then(response => response.json())
                .then(response => setRewards(response.data))
                .catch();
        }
    };

    useEffect(fetchRewards, [scToken, apiBaseUrl, setScToken]);

    return [rewards];
}

