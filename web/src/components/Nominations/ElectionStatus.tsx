import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { isSSR } from "../../service/isSSR";
import { Election } from "../../types/Election";

export function ElectionStatus({ election }: { election?: Election }): JSX.Element | null {
    const [renderAllowed, setRenderAllowed] = useState(false);

    useEffect(() => {
        if (!isSSR()) {
            setRenderAllowed(true);
        }
    }, []);

    let endingAt;
    if (renderAllowed && election?.StartedRunAt) {
        const endingTime = election.StartedRunAt.add(election.Hours, 'hour');
        if (election.SpecificTime) {
            if (election.SpecificTime.isAfter(endingTime)) {
                endingAt = election.SpecificTime.format("L LT");
            } else {
                endingAt = election.SpecificTime.add(1, 'day').format("L LT");
            }
        } else {
            endingAt = endingTime.format("L LT");
        }
    }
    let startedAt;
    if (renderAllowed) {
        startedAt = election?.StartedRunAt?.format("L LT");

        if (!election?.StartedRunAt && election?.SpecificTime) {
            const specificTime = dayjs().set("hour", election.SpecificTime.hour()).set("minute", election.SpecificTime.minute());
            if (specificTime.isBefore(dayjs())) {
                startedAt = specificTime.add(1, 'day').format("L LT");
            } else {
                startedAt = specificTime.format("L LT");
            }
        }
    }

    return <div>
        <div>
            <span className="text-gray-400">Start at ~</span> <strong>{startedAt}</strong>
        </div>
        <div>
            <span className="text-gray-400">End at ~</span> <strong>{endingAt}</strong>
        </div>
        <div>
            <span className="text-gray-400">Duration</span> {!!election && <strong>{election.Hours} hours</strong>}
        </div>
        <div>
            <span className="text-gray-400">Top</span> {!!election && <strong>{election.EmoteAmount}</strong>}
        </div>
    </div>;
}