import React from "react";
import EventService from "../service/EventService";
import Statistics from "./Statistics";
import { connect } from "react-redux";
import fetchChannels from "../actions/fetchChannels";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import Record from "./Record";

class Base extends React.Component {

    state = {
        channelStats: [],
        records: [],
        activeChannels: null,
    };

    componentDidMount() {
        new EventService(this.props.apiBaseUrl, data => {
            this.setState({
                activeChannels: data.activeChannels,
                records: data.records,
                channelStats: data.channelStats,
            });
            this.props.dispatch(fetchChannels(data.channelStats.map(stat => stat.id)));
        });
    }

    render() {
        return (
            <BrowserRouter>
                <div className={"Base"}>
                    <span className="ActiveChannels">{this.state.activeChannels} channels</span>
                    <Switch>
                        <Route path="/aiden">
                            hello Aiden
                        </Route>
                        <Route path="/">
                            <div className="MessagesPerSecond">
                                <h2>Messages per Second</h2>
                                <Statistics channelStats={this.state.channelStats.map(stat => ({ value: stat.msgps, id: stat.id }))} />
                            </div>
                            {this.state.records.map(record => <Record record={record} />)}
                        </Route>
                    </Switch>
                </div>
            </BrowserRouter>
        );
    }
}

export default connect(state => ({ apiBaseUrl: state.apiBaseUrl }))(Base);