import React, { Component } from "react";
import Game from "./Game";
import axios from "axios";
import LoginRegister from "./LoginRegister";
export default class App extends Component {
	constructor(props) {
		super(props);

		this.state = {
			isLoggedin: false,
			isLoading: true,
		};

		this.setLoggedInStatus = this.setLoggedInStatus.bind(this);
	}

	setLoggedInStatus(current) {
		this.setState(() => {
			return { ...this.state, isLoggedin: current };
		});
		//re-trigger the mounting process everytime loginStatus changes
		this.componentDidMount();
	}

	async componentDidMount() {
		try {
			await axios({
				method: "post",
				url: "/check-login",
				withCredentials: true,
			});

			this.setState(() => {
				return { ...this.state, isLoggedin: true, isLoading: false };
			});
		} catch (e) {
			this.setState(() => {
				return { ...this.state, isLoggedin: false, isLoading: false };
			});
		}
	}

	render() {
		return (
			<div>
				{this.state.isLoading ? (
					<div>
						<img
							src="http://i.stack.imgur.com/SBv4T.gif"
							alt="this slowpoke moves"
							width="250"
						/>
					</div>
				) : (
					<div>
						{!this.state.isLoggedin ? (
							<LoginRegister setLoggedInStatus={this.setLoggedInStatus} />
						) : (
							<div>
								<Game setLoggedInStatus={this.setLoggedInStatus} />
							</div>
						)}
					</div>
				)}
			</div>
		);
	}
}
