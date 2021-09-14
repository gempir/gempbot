import { useState } from "react";
import { createLoginUrl } from "../../factory/createLoginUrl";
import { User } from "../../icons/User";
import { useStore } from "../../store";

export function Login() {
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const twitchClientId = useStore(state => state.twitchClientId);
    const isLoggedIn = useStore(state => Boolean(state.scToken));
    const url = createLoginUrl(apiBaseUrl, twitchClientId);

    const [hovering, setHovering] = useState(false);

    const classes = "p-3 flex justify-center rounded shadow bg-purple-800 hover:bg-purple-600 hover:opacity-100 whitespace-nowrap w-36".split(" ")
    if (isLoggedIn) {
        classes.push("opacity-25")
    }

    return <a
        onMouseEnter={() => setHovering(true)}
        onMouseLeave={() => setHovering(false)}
        className={classes.join(" ")}
        href={url.toString()}>
        {isLoggedIn && <>{hovering ? <><User />&nbsp;&nbsp;Login again</> : <><User />&nbsp;&nbsp;Logged in</>}</>}
        {!isLoggedIn && <>{hovering ? <><User />&nbsp;&nbsp;Login</> : <><User />&nbsp;&nbsp;Login</>}</>}
    </a>
}