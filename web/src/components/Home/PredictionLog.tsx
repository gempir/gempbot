import { usePredictionLogs } from "../../hooks/usePredictionLogs";
import { ChevronLeft } from "../../icons/ChevronLeft";
import { ChevronRight } from "../../icons/ChevronRight";
import { Refresh } from "../../icons/Refresh";


export function PredictionLog() {
    const [logs, fetch, loading, page, increasePage, decreasePage] = usePredictionLogs();

    return <div className="p-4 bg-gray-800 rounded shadow w-full relative select-none">
        <div className="text-2xl flex justify-between w-full" onClick={fetch}>
            <Refresh className={"hover:text-blue-500 cursor-pointer " + (loading ? "animate-spin" : "")} />
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