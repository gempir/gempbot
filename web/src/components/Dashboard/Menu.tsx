import React, { KeyboardEvent, useRef, useState } from "react";
import { SetUserConfig, UserConfig } from "../../hooks/useUserConfig";
import { store } from "../../store";
import { Managing } from "./Managing";
import { Reset } from "./Reset";

export function Menu({ userConfig, setUserConfig }: { userConfig: UserConfig | null | undefined, setUserConfig: SetUserConfig }) {
    return <div className="flex flex-row flex-wrap gap-4 ml-4">
        <Managing userConfig={userConfig} />
        <Reset setUserConfig={setUserConfig} />
        <EditorManager setUserConfig={setUserConfig} userConfig={userConfig} />
    </div>
}

function EditorManager({ setUserConfig, userConfig }: { setUserConfig: SetUserConfig, userConfig: UserConfig | null | undefined }) {
    const managing = store.useState(s => s.managing);
    
    const loadedEditors: Record<string, string> = {};
    userConfig?.Editors.map(editor => {
        loadedEditors["e" + (Object.values(loadedEditors).length + 1)] = editor
        return editor;
    });

    const [editors, setEditors] = useState<Record<string, string>>(loadedEditors);

    const changeEditor = (id: string, value: string | null) => {
        const newEditors = { ...editors };

        if (value === null) {
            delete newEditors[id];
            setEditors(newEditors);
        } else {
            newEditors[id] = value;
            setEditors(newEditors);
        }

        // @ts-ignore
        setUserConfig({ ...userConfig, Editors: Object.values(newEditors) })
    };

    const editorInputs = [];
    for (const [key, value] of Object.entries(editors).sort()) {
        editorInputs.push(<EditorInput key={key} editor={value} id={key} onChange={changeEditor} />);
    }

    const creatingNewEditor = Object.values(editors).filter(editor => editor.length === 0).length === 0;

    if (managing !== "") {
        return null;
    }

    return <>
        {editorInputs}
        {creatingNewEditor && <div className="bg-green-700 rounded shadow p-3 cursor-pointer hover:bg-green-600" onClick={() => {
            const newEditors = { ...editors };
            newEditors["e" + (Object.values(editors).length + 1)] = "";
            setEditors(newEditors);
        }}>add editor</div>}
    </>
}

function EditorInput({ id, editor, onChange }: { id: string, editor: string, onChange: (id: string, value: string | null) => void }) {
    const [value, setValue] = useState(editor);
    const ref = useRef<HTMLInputElement>(null);

    const handleKeyPress = (e: KeyboardEvent) => {
        if (e.keyCode === 13) {
            if (ref.current !== null) {
                ref.current.blur();
            }
        }
    };

    return <div className="flex items-center">
        <input className="bg-blue-900 hover:bg-blue-800 focus:bg-blue:800 p-3 rounded shadow w-24 rounded-r-none truncate" spellCheck={false} ref={ref} type="text" value={value} onKeyDown={(e) => handleKeyPress(e)} onChange={(e) => setValue(e.target.value)} onBlur={() => onChange(id, value)} key={editor} />
        <span className="bg-red-900 p-3 px-1 font-bold opacity-25 hover:opacity-100 cursor-pointer rounded rounded-l-none" onClick={() => onChange(id, null)}>X</span>
    </div>;
}