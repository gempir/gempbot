import { usePredictionLogs } from "../../hooks/usePredictionLogs";
import { ChevronLeft } from "../../icons/ChevronLeft";
import { ChevronRight } from "../../icons/ChevronRight";
import { Refresh } from "../../icons/Refresh";


export function PredictionLog({channel}: {channel?: string}) {
    const [logs, fetch, loading, page, increasePage, decreasePage] = usePredictionLogs(channel);

    return <div className="p-4 bg-gray-800 rounded shadow w-full relative select-none">
        <div className="text-2xl flex gap-5 w-full" onClick={fetch}>
            <div className="flex gap-2 items-center">
                <div onClick={decreasePage} className="cursor-pointer hover:text-blue-500">
                    <ChevronLeft />
                </div>
                <div className="text-base w-4 text-center">
                    {page}
                </div>
                <div onClick={increasePage} className="cursor-pointer hover:text-blue-500">
                    <ChevronRight />
                </div>
            </div>
            <Refresh className={"hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} />
        </div>
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
                    <td className="p-3">{log.Title}</td>
                    <td className="text-center">{log.Status}</td>
                    <td className="text-center">{log.getWinningOutcome()?.Title}</td>
                    <td className="text-center">{log.StartedAt.toLocaleDateString()} {log.StartedAt.toLocaleTimeString()}</td>
                    <td className="text-center">{log.LockedAt.toLocaleDateString()} {log.LockedAt.toLocaleTimeString()}</td>
                    <td className="text-center">{log.EndedAt.toLocaleDateString()} {log.EndedAt.toLocaleTimeString()}</td>
                </tr>)}
            </tbody>
        </table>
    </div>
}