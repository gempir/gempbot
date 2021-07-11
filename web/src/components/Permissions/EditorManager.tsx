import React, { useState, useRef } from "react";
import { SetUserConfig, UserConfig } from "../../hooks/useUserConfig";
import { store } from "../../store";

export function EditorManager({ setUserConfig, userConfig }: { setUserConfig: SetUserConfig, userConfig: UserConfig }) {
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

    return <div className="flex mb-4 flex-row flex-wrap gap-4 items-center">
        Editors
        {editorInputs}
        {creatingNewEditor && <div className="bg-green-700 rounded shadow p-3 cursor-pointer hover:bg-green-600" onClick={() => {
            const newEditors = { ...editors };
            newEditors["e" + (Object.values(editors).length + 1)] = "";
            setEditors(newEditors);
        }}>add editor</div>}
    </div>
}

function EditorInput({ id, editor, onChange }: { id: string, editor: string, onChange: (id: string, value: string | null) => void }) {
    const [value, setValue] = useState(editor);
    const ref = useRef<HTMLInputElement>(null);

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === "Enter") {
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