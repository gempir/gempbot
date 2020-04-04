import React from "react";
import EventService from "../service/EventService";
import MessageRecords from "./MessageRecords";
import Statistics from "./Statistics";
import {connect} from "react-redux";
import fetchChannels from "../actions/fetchChannels";

class Base extends React.Component {

    state = {
        channelStats: [],
    };

    componentDidMount() {
       new EventService(this.props.apiBaseUrl, data => {
            this.setState({
                channelStats: data.channelStats
            });
            this.props.dispatch(fetchChannels(data.channelStats.map(stat => stat.id)));
        });
    }

    render() {
        return <div className={"Base"}>
            <div className="MessagesPerSecond">
                <h2>Messages per Second</h2>
                <Statistics channelStats={this.state.channelStats.map(stat => ({value: stat.msgps, id: stat.id}))} />
            </div>
            <div className="MessagesPerMinute">
                <h2>Messages per Minute</h2>
                <Statistics channelStats={this.state.channelStats.map(stat => ({value: stat.msgpm, id: stat.id}))} />
            </div>
        </div>;
    }
}

export default connect(state => ({apiBaseUrl: state.apiBaseUrl}))(Base);