import { usePredictionLogs } from "../../hooks/usePredictionLogs";


export function PredictionLog() {
    const [logs, fetch] = usePredictionLogs();

    return <div className="m-4 p-4 bg-gray-800 rounded shadow w-full overflow-y-scroll relative" style={{maxHeight: "42rem"}}>
        <div className="absolute top-4 left-4 cursor-pointer text-2xl" onClick={fetch}>🔄</div>
        <table className="w-full">
            <thead>
                <th>Title</th>
                <th>Status</th>
                <th>Winner</th>
                <th>StartedAt</th>
                <th>LockedAt</th>
                <th>EndedAt</th>
            </thead>
            <tbody>
                {logs.map((log, index) => <tr className={index % 2 ? "bg-gray-900" : ""}>
                    <th className="p-3">{log.Title}</th>
                    <th>{log.Status}</th>
                    <th>{log.getWinningOutcome()?.Title}</th>
                    <th>{log.StartedAt.toLocaleDateString()} {log.StartedAt.toLocaleTimeString()}</th>
                    <th>{log.LockedAt.toLocaleDateString()} {log.LockedAt.toLocaleTimeString()}</th>
                    <th>{log.EndedAt.toLocaleDateString()} {log.EndedAt.toLocaleTimeString()}</th>
                </tr>)}
            </tbody>
        </table>
    </div>
}