import { useState } from "react";
import { createLoginUrl } from "../../factory/createLoginUrl";
import { User } from "../../icons/User";
import { store } from "../../store";

export function Login() {
    const { apiBaseUrl, twitchClientId } = store.getRawState();
    const url = createLoginUrl(apiBaseUrl, twitchClientId);

    const [hovering, setHovering] = useState(false);

    return <a
        onMouseEnter={() => setHovering(true)}
        onMouseLeave={() => setHovering(false)}
        className={"p-3 flex justify-center rounded shadow opacity-25 bg-purple-800 hover:bg-purple-600 hover:opacity-100"}
        href={url.toString()}>
        {hovering ? <><User />&nbsp;&nbsp;Login again</> : <><User />&nbsp;&nbsp;Logged in</>}
    </a>
}