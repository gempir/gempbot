import styled from "styled-components";
import { Reward as RewardInterface } from "../hooks/useRewards";

const RewardsContainer = styled.div`
    display: grid;
    grid-gap: 1rem;
`;

const RewardContainer = styled.div`
    background: var(--bg-bright);
    border: 1px solid var(--bg-brighter);
    padding: 0.5rem;

    .container {
        display: grid;
        grid-template-columns: 56px 1fr;
        grid-gap: 1rem;
    }
`;

export function Rewards({ rewards }: { rewards: Array<RewardInterface> }) {
    return <RewardsContainer>
        {rewards.map(reward => <RewardContainer>
            <div className="container">
                <img src={reward.image?.url_2x ?? reward.default_image.url_2x} alt={reward.title} />
                <div>
                    <h5>{reward.title}</h5>
                    <p>{reward.prompt}</p>
                </div>
            </div>
        </RewardContainer>)}
    </RewardsContainer>

}