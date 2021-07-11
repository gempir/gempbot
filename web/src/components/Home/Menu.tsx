import React from "react";
import { SetUserConfig, UserConfig } from "../../hooks/useUserConfig";

export function Menu({ userConfig, setUserConfig }: { userConfig: UserConfig, setUserConfig: SetUserConfig }) {
    return <div className="flex flex-row flex-wrap gap-4">
        <BotManager userConfig={userConfig} setUserConfig={setUserConfig} />
    </div>
}
function BotManager({ userConfig, setUserConfig }: { userConfig: UserConfig, setUserConfig: SetUserConfig }) {
    const classes = "p-3 rounded shadow w-28 truncate text-center cursor-pointer".split(" ")

    if (userConfig?.BotJoin) {
        classes.push(..."bg-green-900 hover:bg-green-800 focus:bg-green:800".split(" "));
    } else {
        classes.push(..."bg-red-900 hover:bg-red-800 focus:bg-red:800".split(" "));
    }

    return <div className={classes.join(" ")} onClick={() => setUserConfig({ ...userConfig, BotJoin: !userConfig?.BotJoin })}>{userConfig?.BotJoin ? "Bot Joined" : "Join Bot"}</div>
}

