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
       new EventService(data => {
            this.setState({
                channelStats: data.channelStats
            });
            this.props.dispatch(fetchChannels(data.channelStats.map(stat => stat.id)));
        });
    }

    render() {
        return <div>
            <Statistics channelStats={this.state.channelStats} />
        </div>;
    }
}

export default connect(state => ({}))(Base);