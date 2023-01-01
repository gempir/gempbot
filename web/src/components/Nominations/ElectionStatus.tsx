import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { useForceUpdate } from "../../hooks/useForceUpdate";
import { isSSR } from "../../service/isSSR";
import { Election } from "../../types/Election";

export function ElectionStatus({ election }: { election?: Election }): JSX.Element | null {
    const [renderAllowed, setRenderAllowed] = useState(false);
    const forceUpdate = useForceUpdate();

    useEffect(() => {
        if (!isSSR()) {
            setRenderAllowed(true);
        }

        const interval = setInterval(() => {
            forceUpdate();
        }, 60000);

        return () => {
            clearInterval(interval);
        }
    }, []);

    let endingAt;
    if (renderAllowed && election?.StartedRunAt) {
        const endingTime = election.StartedRunAt.add(election.Hours, 'hour');
        if (election.SpecificTime) {
            if (election.SpecificTime.isAfter(endingTime)) {
                endingAt = election.SpecificTime;
            } else {
                endingAt = election.SpecificTime.add(1, 'day');
            }
        } else {
            endingAt = endingTime;
        }
    }
    let startedAt;
    if (renderAllowed) {
        startedAt = election?.StartedRunAt;

        if (!election?.StartedRunAt && election?.SpecificTime) {
            const specificTime = dayjs().set("hour", election.SpecificTime.hour()).set("minute", election.SpecificTime.minute());
            if (specificTime.isBefore(dayjs())) {
                startedAt = specificTime.add(1, 'day');
            } else {
                startedAt = specificTime;
            }
        }
    }

    let countdown = "";

    if (renderAllowed && endingAt) {
        let minutes = endingAt.diff(dayjs(), 'minutes');

        const days = Math.floor(minutes / 1440);
        minutes -= days * 1440;
        const hours = Math.floor(minutes / 60);
        minutes -= hours * 60;

        if (days) {
            countdown += `${days} day${days > 1 ? "s" : ""} `;
        }
        if (hours) {
            countdown += `${hours} hour${hours > 1 ? "s" : ""} `;
        }
        if (minutes) {
            countdown += `${minutes} minute${minutes > 1 ? "s" : ""} `;
        }
    }

    return <div className="flex gap-4">
        <div className="bg-gray-800 rounded p-4 shadow">
            <span className="text-gray-400">Start at ~</span> {startedAt?.format("L LT")}
        </div>
        <div className="bg-gray-800 rounded p-4 shadow">
            <span className="text-gray-400">End at ~</span> {endingAt?.format("L LT")}
        </div>
        <div className="bg-gray-800 rounded p-4 shadow">
            <span className="text-gray-400">Top</span> {!!election && election.EmoteAmount}
        </div>
        <div className="bg-gray-800 rounded p-4 shadow">
            <span className="text-gray-400">Ending in ~</span> <strong>{countdown}</strong>
        </div>
    </div>;
}