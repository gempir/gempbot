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

    if (userConfig?.Protected?.EditorFor.length === 0) {
        return null;
    }

    return <div className="Managing flex items-center my-4">
        <UserGroupIcon className="h-6" />
        <select className="block ml-2 p-1 rounded bg-gray-900 shadow focus:outline-none" style={{maxWidth: 96}} onChange={updateManaging} value={managing}>
            {userConfig?.Protected.EditorFor.map(editorFor => <option key={editorFor} value={editorFor}>{editorFor}</option>)}
            <option value="">you</option>
        </select>
    </div>
}