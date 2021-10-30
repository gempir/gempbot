import { ChatIcon } from "@heroicons/react/solid";
import { useState } from "react";
import { SetUserConfig, UserConfig } from "../../hooks/useUserConfig";


export function BotManager({ userConfig, setUserConfig, userConfigLoading }: { userConfig: UserConfig, setUserConfig: SetUserConfig, userConfigLoading: boolean }) {
    const classes = "p-3 flex justify-center rounded shadow cursor-pointer mt-2 w-full hover:opacity-100 whitespace-nowrap w-36".split(" ")
    const [hovering, setHovering] = useState(false);

    
    if (userConfig?.BotJoin) {
        classes.push(..."bg-green-900 hover:bg-red-800 focus:bg-green:800 opacity-25".split(" "));
    } else {
        classes.push(..."bg-red-900 focus:bg-red:800".split(" "));

        if (hovering) {
            classes.push(..."hover:bg-green-900".split(" "));
        }
    }

    if (userConfigLoading) {
        classes.push(..."cursor-wait animate-pulse".split(" "))
    }

    return <div className={classes.join(" ")}
        onMouseEnter={() => setHovering(true)}
        onMouseLeave={() => setHovering(false)}
        onClick={() => setUserConfig({ ...userConfig, BotJoin: !userConfig?.BotJoin })}>
        {hovering ? userConfig?.BotJoin ? <><ChatIcon className="h-6" />&nbsp;&nbsp;Depart Bot</> : <><ChatIcon className="h-6" />&nbsp;&nbsp;Join Bot</> : userConfig?.BotJoin ? <><ChatIcon className="h-6" />&nbsp;&nbsp;Bot joined</> : <><ChatIcon className="h-6" />&nbsp;&nbsp;Not Joined</>}
    </div>
}

