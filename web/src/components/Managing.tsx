import styled from "styled-components";
import { UserConfig } from "../hooks/useUserConfig";
import { store } from "../store";

const ManagingContainer = styled.div`
    position: relative;

    select {
        -moz-appearance: none;
        -webkit-appearance: none;
        appearance: none;
        background: var(--bg-bright);
        border: 1px solid var(--bg-brighter);
        color: white;
        outline: none;
        padding: 1rem;
        padding-right: 1.5rem;
        font-size: 0.8rem;
        line-height: 1.1rem;
        position: relative;
        height: 50px;
    }

    &:after {
        content: "▼";
        color: white;
        pointer-events: none;
        position: absolute;
        right: 10px;
        top: 20px;
        margin: auto;
    }
`;

export function Managing({ userConfig }: { userConfig: UserConfig | null | undefined }) {
    if (userConfig?.Protected.EditorFor.length === 0) {
        return null;
    }
    const managing = store.useState(s => s.managing);    

    return <ManagingContainer>
        <select onChange={e => store.update(s => {s.managing = e.target.value})} value={managing} defaultValue={""}>
            {userConfig?.Protected.EditorFor.map(editorFor => <option key={editorFor} value={editorFor}>{editorFor}</option>)}
            <option value="">you</option>
        </select>
    </ManagingContainer>
}