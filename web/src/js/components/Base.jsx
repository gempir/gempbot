import React from "react";
import EventService from "../service/EventService";
import { connect } from "react-redux";
import fetchChannels from "../actions/fetchChannels";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import ReactWordcloud from "react-wordcloud";
import Record from "./Record";

class Base extends React.Component {

    state = {
        activeChannels: null,
        records: [],
        wordcloudWords: [],
    };

    componentDidMount() {
        new EventService(this.props.apiBaseUrl, data => {
            this.setState({
                activeChannels: data.activeChannels,
                wordcloudWords: data.wordcloudWords,
                records: data.records,
            });
            this.props.dispatch(fetchChannels(new Set(data.records.map(record => record.scores).flat().map(score => score.id))));
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
                            {this.state.records.map(record => <Record record={record} key={record.title} />)}
                            {this.state.wordcloudWords.length > 0 && 
                            <div className="Record">
                                <h2>Wordcloud</h2>
                                <div className="WordCloud">
                                    <ReactWordcloud words={this.state.wordcloudWords} options={{ deterministic: true, fontSizes: [10, 40], enableTooltip: false, transitionDuration: 250, fontFamily: "monospaced"}} />
                                </div>
                            </div>}
                        </Route>
                    </Switch>
                </div>
            </BrowserRouter>
        );
    }
}

export default connect(state => ({ apiBaseUrl: state.apiBaseUrl }))(Base);