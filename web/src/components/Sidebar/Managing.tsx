import React from "react";
import { UserConfig } from "../../hooks/useUserConfig";
import { setCookie } from "../../service/cookie";
import { useStore } from "../../store";
import { UserGroupIcon } from "@heroicons/react/24/solid";
import { NativeSelect } from "@mantine/core";


export function Managing({ userConfig }: { userConfig: UserConfig | null | undefined }) {
    const setManaging = useStore(state => state.setManaging);
    const managing = useStore(state => state.managing);
    const updateManaging = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setManaging(String(e.target.value).trim() !== "" ? e.target.value : null);
        setCookie("managing", e.target.value);
    };

    const data = userConfig?.Protected.EditorFor.sort().map(channel => ({label: channel, value: channel})) || [];
    data.unshift({label: "You", value: ""});

    return <div className="Managing m-3">
        <NativeSelect size="sm" className="w-full" data={data} onChange={updateManaging} value={managing ?? ""}
            rightSection={<UserGroupIcon className="text-gray-400 h-4" />}
        />
    </div>
}