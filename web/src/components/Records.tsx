import styled from "styled-components";
import { ProfilePicture } from "./ProfilePicture";
import { Record } from "../types/Events";

const RecordsContainer = styled.div`
    display: flex;
`;

export function Records({ records }: { records: Array<Record> }) {
    return <RecordsContainer>
        {records.map(record => <RecordComponent key={record.title} record={record} />)}

    </RecordsContainer>
}

const RecordConatiner = styled.div`
    background: var(--bg-bright);
    border: 1px solid var(--bg-brighter);
    margin: 1rem;
    margin-right: 0;
    padding: 1rem;

    h2 {
        color: var(--text);
        margin: 0;
        margin-bottom: 1rem;
        padding: 0;
    }

    ol {
        color: white;
        font-size: 1.5rem;
        font-weight: bold;
        width: 500px;
        padding: 0;
        margin: 0;
        margin-right: 15px;
        background: var(--lightBackground);
        border: 1px solid var(--lightBorder);
        border-radius: 3px;

        li {
            display: flex;
            align-items: center;
            margin-bottom: 0.25rem;

            img {
                margin-right: 10px;
            }

            .value {
                text-align: right;
                flex: 1 1 auto;
            }
        }
    }
`;

function RecordComponent({ record }: { record: Record }) {
    return <RecordConatiner>
        <h2>{record.title}</h2>
        <ol>
            {record.scores.map(score => <li key={score.user.id}>
                <ProfilePicture src={score.user.profilePicture} />
                <span>{score.user.displayName}</span>
                <span className={"value"}>{score.score}</span>
            </li>)}
        </ol>
    </RecordConatiner>
}