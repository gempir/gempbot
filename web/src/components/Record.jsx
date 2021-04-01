import React from "react";
import { connect } from "react-redux";
import ProfilePicture from "./ProfilePicture";

class Record extends React.Component {
    render() {
        return <div className="Record">
            <h2>{this.props.record.title}</h2>
            <ol>
                {this.props.record.scores.map(score => <li key={score.id}>
                    <ProfilePicture src={this.props.channels[score.id]?.profile_image_url} />
                    <span>{this.props.channels[score.id]?.display_name ?? ""}</span>
                    <span className={"value"}>{score.score}</span>
                </li>)}
            </ol>
        </div>
    }
}

export default connect(state => ({
    channels: state.channels,
}))(Record);