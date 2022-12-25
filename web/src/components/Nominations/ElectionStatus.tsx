import { useEffect, useState } from "react";
import { useElection } from "../../hooks/useElection";
import { isSSR } from "../../service/isSSR";

export function ElectionStatus({channel}: {channel: string}): JSX.Element | null {
    const [election] = useElection(channel);
    const [renderAllowed, setRenderAllowed] = useState(false);

    useEffect(() => {
        if (!isSSR()) {
            setRenderAllowed(true);
        }
    }, []);

    return <div>
        <span className="text-gray-400">Ending at</span> <strong>{renderAllowed && election.StartedRunAt.add(election.Hours, 'hour').format('L LT')}</strong>
        <div>
            <span className="text-gray-400">Duration</span> <strong>{election.Hours} hours</strong>
        </div>
    </div>;
}