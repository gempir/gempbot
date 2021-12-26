import React from "react";
import { UserConfig } from "../../hooks/useUserConfig";
import { setCookie } from "../../service/cookie";
import { useStore } from "../../store";
import { UserGroupIcon } from "@heroicons/react/solid";


export function Managing({ userConfig }: { userConfig: UserConfig | null | undefined }) {
    const setManaging = useStore(state => state.setManaging);
    const managing = useStore(state => state.managing);
    const updateManaging = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setManaging(e.target.value);
        setCookie("managing", e.target.value);
    };

    return <div className="Managing flex items-center my-4">
        <UserGroupIcon className="h-6" style={{width: 21}} />
        <select className="block ml-2 p-1 rounded bg-gray-900 shadow focus:outline-none w-full" style={{ maxWidth: 96 }} onChange={updateManaging} value={managing}>
            <option value="">you</option>
            {userConfig?.Protected.EditorFor.map(editorFor => <option key={editorFor} value={editorFor}>{editorFor}</option>)}
        </select>
    </div>
}