import React from "react";
import EventService from "../service/EventService";
import MessageRecords from "./MessageRecords";
import Statistics from "./Statistics";

export default class Base extends React.Component {

    state = {
        channels: [],
    }

    componentDidMount() {
       new EventService(data => {
            this.setState({
                channels: data.channels
            })
        });
    }

    render() {
        return <div>
            <Statistics channels={this.state.channels} />
            <MessageRecords records={[]}/>
        </div>;
    }
}