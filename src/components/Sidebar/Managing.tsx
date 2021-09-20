import React from "react";
import { UserConfig } from "../../hooks/useUserConfig";
import { UserGroup } from "../../icons/UserGroup";
import { setCookie } from "../../service/cookie";
import { useStore } from "../../store";

export function Managing({ userConfig }: { userConfig: UserConfig | null | undefined }) {
    if (userConfig?.Protected.EditorFor.length === 0) {
        return null;
    }
    const managing = useStore(state => state.managing);
    const updateManaging = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const setManaging = useStore(state => state.setManaging);
        setManaging(e.target.value);
        setCookie("managing", e.target.value);
    };

    return <div className="Managing flex items-center my-4">
        <UserGroup />
        <select className="block ml-2 p-1 rounded bg-gray-900 shadow focus:outline-none" onChange={updateManaging} value={managing}>
            {userConfig?.Protected.EditorFor.map(editorFor => <option key={editorFor} value={editorFor}>{editorFor}</option>)}
            <option value="">you</option>
        </select>
    </div>
}