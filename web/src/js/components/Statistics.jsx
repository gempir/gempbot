import React from "react";
import {connect} from "react-redux";
import ProfilePicture from "./ProfilePicture";

class Statistics extends React.Component {
    render() {
        const statistics = [];

        const stats = this.props.channelStats;
        stats.sort((a, b) => {
            if (a.value < b.value) {
                return 1;
            }
            if (a.value > b.value) {
                return -1;
            }
            return 0;
        });

        for (const stat of this.props.channelStats) {
            statistics.push(<li key={stat.id}>
                <ProfilePicture src={this.props.channels[stat.id]?.profile_image_url}/>
                <span>{this.props.channels[stat.id]?.display_name ?? ""}</span>
                <span className={"value"}>{stat.value}</span>
            </li>);
        }

        return <ul className="Statistics">
            {statistics}
        </ul>
    }
}

export default connect(state => ({
    channels: state.channels,
}))(Statistics);