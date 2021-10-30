import { Toggle } from './Toggle';
import { useEffect, useState } from 'react';
import { useSubscribtions } from '../../hooks/useSubscriptions';

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

    const classes = ["p-4"];
    if (loading) {
        classes.push("animate-pulse pointer-events-none");
    }

    return <div className={classes.join(" ")}>
        <div className="bg-gray-800 rounded shadow relative p-4">
            <div className="flex items-start justify-between">
                <div>
                    <h3>Prediction Announcements</h3>
                    <div className="p-2 text-gray-400 mx-0 px-0">
                        Announces when predictions
                        <ul className="list-disc pl-6">
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
        {/* <div className="bg-gray-800 rounded shadow relative p-4 mt-4">
            <div className="flex items-start justify-between">
                <div>
                    <h3>Prediction Commands</h3>
                    <div className="p-2 text-gray-400 mx-0 px-0">
                        Commands to manage predictions
                        <ul className="list-disc pl-6">
                            <li>are made</li>
                            <li>locked</li>
                            <li>canceled</li>
                            <li>resolved</li>
                        </ul>
                    </div>
                </div>
                <Toggle checked={predictionsAnnouncements} onChange={setPredictionAnnouncements} />
            </div>
        </div> */}
    </div >;
}


