import React, { Component } from "react";
import axios from "axios";
import querystring from "querystring";
import "./styles/LoginRegister.css";

export default class LoginRegister extends Component {
	constructor(props) {
		super(props);

		this.usernameRef = React.createRef();
		this.passwordRef = React.createRef();
		this.confirmRef = React.createRef();

		this.state = {
			isLogin: true,
		};

		this.handleSubmit = this.handleSubmit.bind(this);
		this.handleLogin = this.handleLogin.bind(this);
		this.handleRegister = this.handleRegister.bind(this);
		this.setIsLogin = this.setIsLogin.bind(this);
		this.unsetIsLogin = this.unsetIsLogin.bind(this);
	}

	async handleLogin() {
		let user = this.usernameRef.current.value;
		let pass = this.passwordRef.current.value;

		try {
			await axios({
				method: "post",
				url: "/login",
				data: querystring.stringify({
					username: user,
					password: pass,
				}),
			});
			this.props.setLoggedInStatus(true);
		} catch (e) {
			const status = e.response.status;

			if (status === 403) {
				alert("invalid username or password");
			} else {
				alert("something went wrong");
			}
			this.props.setLoggedInStatus(false);
		}
	}

	async handleRegister() {
		let user = this.usernameRef.current.value;
		let pass = this.passwordRef.current.value;
		let conf = this.confirmRef.current.value;

		if (pass !== conf) {
			alert("Passwords don't match");
			return;
		}

		try {
			await axios({
				method: "post",
				url: "/register",
				data: querystring.stringify({
					username: user,
					password: pass,
				}),
			});

			this.setState(() => {
				return { ...this.state, isLogin: true };
			});
		} catch (e) {
			alert("username already exists");
		}
	}

	handleSubmit() {
		let current = this.state.isLogin;
		if (current === true) {
			this.handleLogin();
		} else {
			this.handleRegister();
		}
	}

	setIsLogin() {
		this.setState(() => {
			return { ...this.state, isLogin: true };
		});
	}
	unsetIsLogin() {
		this.setState(() => {
			return { ...this.state, isLogin: false };
		});
	}

	render() {
		return (
			<div className="LoginRegister">
				<div className="LoginRegister-Type-Select">
					<div onClick={this.setIsLogin} className="button-regular">
						Login
					</div>
					<div onClick={this.unsetIsLogin} className="button-regular">
						Register
					</div>
				</div>
				<h1>{this.state.isLogin ? "Login" : "Register"}</h1>
				<div className="LoginRegister-Body">
					<div>
						<label className="label-regular">Username</label>
						<br />
						<input
							type="text"
							name="username"
							id="username"
							ref={this.usernameRef}
						/>
					</div>
					<div>
						<label className="label-regular">Password</label>
						<br />
						<input
							type="password"
							name="password"
							id="password"
							ref={this.passwordRef}
						/>
					</div>
					{/* display confirm password field only for register */}
					{!this.state.isLogin ? (
						<div>
							<label className="label-regular">Confirm Password</label>
							<br />
							<input
								type="password"
								name="confirm-password"
								id="confirm-password"
								ref={this.confirmRef}
							/>
						</div>
					) : null}
					<div
						onClick={this.handleSubmit}
						className="button-regular button-submit"
					>
						Submit
					</div>
				</div>
			</div>
		);
	}
}
