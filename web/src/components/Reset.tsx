import { UserConfig } from "../hooks/useUserConfig";
import { store } from "../store";


export function Reset({ setUserConfig }: { setUserConfig: (userConfig: UserConfig | null) => void }) {
    const managing = store.useState(s => s.managing);
    if (managing !== "") {
        return null;
    }

    return <div className="p-3 opacity-25 hover:opacity-100 bg-red-900 hover:bg-red-800 shadow rounded cursor-pointer" onClick={() => {
        if (window.confirm(`Do you really want to reset?\n- Channel Point Rewards on Twitch from spamchamp\n- Settings on spamchamp.gempir.com\n- Unsubscribes all webhooks for your channel`)) {
            setUserConfig(null);
        }
    }}>Reset</div>;
}