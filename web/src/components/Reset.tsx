import styled from "styled-components";
import { UserConfig } from "../hooks/useUserConfig";


const ResetContainer = styled.div`
    display: block;
    color: white;
    padding: 1rem 1rem;
    text-decoration: none;
    font-weight: bold;
    border-radius: 3px;
    background: var(--danger-dark);
    cursor: pointer;
    opacity: 0.25;
    transition: opacity 0.2s ease-in-out;

    &:hover {
        background: var(--danger);
        opacity: 1;
    }
`;

export function Reset({ setUserConfig }: { setUserConfig: (userConfig: UserConfig | null) => void }) {
    return <ResetContainer onClick={() => {
        if (window.confirm(`Do you really want to reset?\n- Channel Point Rewards on Twitch from spamchamp\n- Settings on spamchamp.gempir.com\n- Unsubscribes all webhooks for your channel`)) {
            setUserConfig(null);
        }
    }}>Reset</ResetContainer>;
}