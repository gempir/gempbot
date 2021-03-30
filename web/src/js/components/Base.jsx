import React from "react";
import EventService from "../service/EventService";
import { connect } from "react-redux";
import fetchChannels from "../actions/fetchChannels";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import Record from "./Record";

class Base extends React.Component {

    state = {
        joinedChannels: null,
        activeChannels: null,
        records: [],
    };

    componentDidMount() {
        new EventService(this.props.apiBaseUrl, data => {
            this.setState({
                joinedChannels: data.joinedChannels,
                activeChannels: data.activeChannels,
                records: data.records,
            });
            this.props.dispatch(fetchChannels(new Set(data.records.map(record => record.scores).flat().map(score => score.id))));
        });
    }

    render() {
        return (
            <BrowserRouter>
                <div className={"Base"}>
                    <span className="Meta">{this.state.joinedChannels} join channels | {this.state.activeChannels} active channels</span>
                    <Switch>
                        <Route path="/aiden">
                            hello Aiden
                        </Route>
                        <Route path="/">
                            {this.state.records.map(record => <Record record={record} key={record.title} />)}
                        </Route>
                    </Switch>
                </div>
            </BrowserRouter>
        );
    }
}

export default connect(state => ({ apiBaseUrl: state.apiBaseUrl }))(Base);