import { useEffect, useState } from 'react';
import { useSubscribtions } from '../../hooks/useSubscriptions';
import { Toggle } from './Toggle';

export function Bot() {
    const [subscribe, unsubscribe, subscriptionsStatus, loading] = useSubscribtions();

    const [predictionsAnnouncements, setPredictionAnnouncements] = useState(false);

    useEffect(() => {
        setPredictionAnnouncements(subscriptionsStatus.predictions);
    }, [subscriptionsStatus.predictions]);

    const handlePredictionAnnouncementChange = (value: boolean) => {
        setPredictionAnnouncements(value);
        if (value) {
            subscribe();
        } else {
            unsubscribe();
        }
    };

    return <div className={"p-4"}>
        <div className={"bg-gray-800 rounded shadow relative p-4 " + (loading ? "animate-pulse pointer-events-none" : "")}>
            <div className="flex items-start justify-between">
                <div>
                    <h3 className="font-bold text-xl">Prediction Announcements</h3>
                    <div className="p-2 text-gray-200 mx-0 px-0">
                        Announces when predictions
                        <ul className="list-disc pl-6 mt-2">
                            <li>are made</li>
                            <li>locked</li>
                            <li>canceled</li>
                            <li>resolved</li>
                        </ul>
                        <img src={"/images/announcement.png"} className="mt-2" />
                    </div>
                </div>
                <Toggle checked={predictionsAnnouncements} onChange={handlePredictionAnnouncementChange} />
            </div>
        </div>
    </div >;
}


