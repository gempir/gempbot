import { usePredictionLogs } from "../../hooks/usePredictionLogs";
import { Refresh } from "../../icons/Refresh";


export function PredictionLog() {
    const [logs, fetch, loading] = usePredictionLogs();

    return <div className="mt-4 p-4 bg-gray-800 rounded shadow w-full relative">
        <div className="absolute top-4 left-4 cursor-pointer text-2xl" onClick={fetch}><Refresh className={"hover:text-gray-400 " + (loading ? "animate-spin" : "")} /></div>
        <table className="w-full">
            <thead>
                <tr>
                    <th>Title</th>
                    <th>Status</th>
                    <th>Winner</th>
                    <th>StartedAt</th>
                    <th>LockedAt</th>
                    <th>EndedAt</th>
                </tr>
            </thead>
            <tbody>
                {logs.map((log, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
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