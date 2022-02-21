import React, { Component } from "react";
import axios from "axios";
import querystring from "querystring";

export default class Game extends Component {
	constructor(props) {
		super(props);

		this.state = {
			sudoku: [],
			mustGenSudoku: true,
			isLoading: true,
		};

		this.selectRef = React.createRef();

		this.handleLogout = this.handleLogout.bind(this);
		this.generateNewSudoku = this.generateNewSudoku.bind(this);
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

	async generateNewSudoku() {
		this.setState(() => {
			return { ...this.state, isLoading: true };
		});

		try {
			let result = await axios({
				method: "post",
				url: "/gen-sudoku",
				data: querystring.stringify({
					difficulty: this.selectRef.current.value,
				}),
				withCredentials: true,
			});

			this.setState(() => {
				return {
					...this.state,
					sudoku: result.data,
					isLoading: false,
					mustGenSudoku: false,
				};
			});
		} catch (e) {
			this.setState(() => {
				return { ...this.state, isLoading: false };
			});
			console.log(e.response);
			alert("Something went wrong");
		}
	}

	async componentDidMount() {
		try {
			//checking for existing unsolved sudokus
			let res = await axios({
				method: "post",
				url: "/check-unsolved",
				withCredentials: true,
			});
			this.setState(() => {
				return {
					...this.state,
					mustGenSudoku: false,
					isLoading: false,
					sudoku: res.data,
				};
			});
		} catch (e) {
			const status = e.response.status;
			if (status === 404) {
				//generate new sudoku
				this.setState(() => {
					return { ...this.state, mustGenSudoku: true, isLoading: false };
				});
			} else {
				this.setState(() => {
					return { ...this.state, isLoading: false };
				});
				alert("something went wrong");
			}
		}
	}

	render() {
		let row = 0;
		let col = 0;
		return (
			<div>
				<button onClick={this.handleLogout}>Logout</button>
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
						{this.state.mustGenSudoku ? (
							<div>
								<select ref={this.selectRef} defaultValue="1">
									<option value="0">Easy</option>
									<option value="1">Medium</option>
									<option value="2">Hard</option>
								</select>
								<button onClick={this.generateNewSudoku}>Generate</button>
							</div>
						) : (
							<div>{/*TODO: Render game cells*/}</div>
						)}
					</div>
				)}
			</div>
		);
	}
}
