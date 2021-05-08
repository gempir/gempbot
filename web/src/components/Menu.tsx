import React, { KeyboardEvent, useRef, useState } from "react";
import styled from "styled-components";
import { SetUserConfig, UserConfig } from "../hooks/useUserConfig";
import { Managing } from "./Managing";
import { Reset } from "./Reset";

const MenuContainer = styled.div`
    display: inline-grid;
    grid-auto-flow: column;
    position: absolute;
    top: 1rem;
    left: 1rem;
    border-radius: 3px;
    cursor: pointer;
    grid-gap: 0.5rem;
`;

export function Menu({ userConfig, setUserConfig }: { userConfig: UserConfig | null | undefined, setUserConfig: SetUserConfig }) {
    return <MenuContainer>
        <Reset setUserConfig={setUserConfig} />
        <Managing userConfig={userConfig} />
        <EditorManager setUserConfig={setUserConfig} userConfig={userConfig} />
    </MenuContainer>
}

const EditorManagerContainer = styled.div`
    display: inline-grid;
    grid-auto-flow: column;
    grid-gap: 0.5rem;

    .editor-input {
        background: var(--bg-bright);
        border: 1px solid var(--bg-brighter);
        color: white;
        outline: none;
        padding-right: 0.5rem;
        font-size: 0.8rem;
        line-height: 1.2rem;
        position: relative;
        height: 50px;

        input {
            outline: none;
            max-width: 6rem;
            padding: 1rem;
            padding-right: 0;
            background: transparent;
            border: 0;
            color: white;
        }

        span {
            padding: 0.5rem;
            font-weight: bold;
            user-select: none;

            &:hover {
                color: var(--danger);
            }
        }
    }

    .add-editor {
        height: 50px;
        padding: 1rem;
        font-size: 0.8rem;
        background: var(--theme-bright);
        border-radius: 3px;
        opacity: 1;
        font-weight: bold;
        line-height: 1.2rem;
        transition: opacity 0.2s ease-in-out;
        user-select: none;

        &:hover {
            opacity: 1;
        }
    }
`;

function EditorManager({ setUserConfig, userConfig }: { setUserConfig: SetUserConfig, userConfig: UserConfig | null | undefined }) {
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
    for (const [key, value] of Object.entries(editors)) {
        editorInputs.push(<EditorInput key={key} editor={value} id={key} onChange={changeEditor} />);
    }

    const creatingNewEditor = Object.values(editors).filter(editor => editor.length === 0).length === 0;

    return <EditorManagerContainer>
        {editorInputs}
        {creatingNewEditor && <div className="add-editor" onClick={() => {
            const newEditors = { ...editors };
            newEditors["e" + (Object.values(editors).length + 1)] = "";
            setEditors(newEditors);
        }}>add editor</div>}
    </EditorManagerContainer>
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

    return <div className="editor-input">
        <input ref={ref} type="text" value={value} onKeyDown={(e) => handleKeyPress(e)} onChange={(e) => setValue(e.target.value)} onBlur={() => onChange(id, value)} key={editor} />
        <span onClick={() => onChange(id, null)}>X</span>
    </div>;
}