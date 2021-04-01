import React from "react";

export default class MessageRecords extends React.Component {
    render() {
        return <div className="MessageRecords">
            {this.props.records.map(record =>
                <div>{record.channel}</div>
            )}
        </div>
    }
}