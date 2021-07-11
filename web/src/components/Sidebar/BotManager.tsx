import { useState } from "react";
import { SetUserConfig, UserConfig } from "../../hooks/useUserConfig";
import { Chat } from "../../icons/Chat";

export function BotManager({ userConfig, setUserConfig }: { userConfig: UserConfig, setUserConfig: SetUserConfig }) {
    const classes = "p-3 flex justify-center rounded shadow cursor-pointer mt-2 w-full hover:opacity-100".split(" ")
    const [hovering, setHovering] = useState(false);

    if (userConfig?.BotJoin) {
        classes.push(..."bg-green-900 hover:bg-red-800 focus:bg-green:800 opacity-25".split(" "));
    } else {
        classes.push(..."bg-red-900 focus:bg-red:800".split(" "));

        if (hovering) {
            classes.push(..."hover:bg-green-900".split(" "));
        }
    }


    return <div className={classes.join(" ")}
        onMouseEnter={() => setHovering(true)}
        onMouseLeave={() => setHovering(false)}
        onClick={() => setUserConfig({ ...userConfig, BotJoin: !userConfig?.BotJoin })}>
        {hovering ? userConfig?.BotJoin ? <><Chat />&nbsp;&nbsp;Depart Bot</> : <><Chat />&nbsp;&nbsp;Join Bot</> : userConfig?.BotJoin ? <><Chat />&nbsp;&nbsp;Bot joined</> : <><Chat />&nbsp;&nbsp;Not Joined</>}
    </div>
}

