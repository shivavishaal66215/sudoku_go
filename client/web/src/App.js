import React, { Component } from "react";
import axios from "axios";

export default class App extends Component {
	async componentDidMount() {
		let res = await axios.get("/test");
		console.log(res);
	}

	render() {
		return <div>hello there</div>;
	}
}
