import React from "react";

export default class Statistics extends React.Component {
    render() {
        const statistics = [];

        const stats = Object.values(this.props.channels);
        stats.sort((a, b) => {
            if (a.messagesPerSecond < b.messagesPerSecond) {
                return 1;
            }
            if (a.messagesPerSecond > b.messagesPerSecond) {
                return -1;
            }
            return 0;
        });

        for (const stat of stats) {
            statistics.push(<li key={stat.channelName}>{stat.channelName}: {stat.messagesPerSecond}</li>)
        }

        return <ul className="Statistics">
            {statistics}
        </ul>
    }
}