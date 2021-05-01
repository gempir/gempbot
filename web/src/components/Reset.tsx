import styled from "styled-components";
import { UserConfig } from "../hooks/useUserConfig";


const ResetContainer = styled.div`
    position: absolute;
    top: 1rem;
    left: 1rem;
    display: block;
    color: white;
    padding: 1rem 2rem;
    text-decoration: none;
    font-weight: bold;
    border-radius: 3px;
    background: var(--danger-dark);
    cursor: pointer;

    &:hover {
        background: var(--danger);
    }
`;

export function Reset({ setUserConfig }: { setUserConfig: (userConfig: UserConfig | null) => void }) {
    return <ResetContainer onClick={() => setUserConfig(null)}>Reset</ResetContainer>;
}