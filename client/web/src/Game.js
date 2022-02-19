import React, { Component } from "react";
import axios from "axios";

export default class Game extends Component {
	constructor(props) {
		super(props);

		this.handleLogout = this.handleLogout.bind(this);
	}

	async handleLogout() {
		try {
			await axios({
				method: "post",
				url: "/logout",
				withCredentials: true,
			});
			this.props.setLoggedInStatus(false);
		} catch (e) {
			alert("could not log you out");
		}
	}

	render() {
		return (
			<div>
				<div>Game</div>
				<button onClick={this.handleLogout}>Logout</button>
			</div>
		);
	}
}
