import React, { Component } from 'react';
import { createStore, applyMiddleware } from "redux";
import thunk from "redux-thunk";
import { Provider } from "react-redux";
import reducer from "./store/reducer";
import createInitialState from "./store/createInitialState";
import Base from './components/Base';
import persistState from './storage/persistState';


export default class App extends Component {

	constructor(props) {
		super(props);

		this.store = createStore(reducer, createInitialState(), applyMiddleware(thunk));

		this.store.subscribe(() => persistState(this.store.getState()));
	}

	render() {
		return (
			<Provider store={this.store}>
				<Base />
			</Provider>
		);
	}
}