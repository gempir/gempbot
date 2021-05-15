import { Record } from "../types/Events";
import { ProfilePicture } from "./ProfilePicture";

export function Records({ records }: { records: Array<Record> }) {
    return <div className="flex flex-row pl-5">
        {records.map(record => <RecordComponent key={record.title} record={record} />)}
    </div>
}

function RecordComponent({ record }: { record: Record }) {
    return <div className="bg-gray-600 rounded shadow p-5 w-64 mr-5 ml-0">
        <h2 className="text-1xl font-bold">{record.title}</h2>
        <ol>
            {record.scores.map(score => <li key={score.user.id} className="flex flex-row justify-between items-center m-1">
                <ProfilePicture src={score.user.profilePicture} />
                <span className="font-bold">{score.user.displayName}</span>
                <span className="font-bold">{score.score}</span>
            </li>)}
        </ol>
    </div>
}