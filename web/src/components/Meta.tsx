import styled from "styled-components";

const MetaContainer = styled.div`
    display: inline-block;
    background: var(--bg-bright);
    border: 1px solid var(--bg-brighter);
    padding: 0.5rem;
    margin: 1rem;
    margin-bottom: 0;
`;

export function Meta({activeChannels, joinedChannels}: {[key: string]: number}) {
    return <MetaContainer>
        Joined Channels: <strong>{joinedChannels}</strong> | Active Channels: <strong>{activeChannels}</strong>
    </MetaContainer> 
}