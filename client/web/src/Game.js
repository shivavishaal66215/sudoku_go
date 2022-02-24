import React, { Component } from "react";
import axios from "axios";
import querystring from "querystring";
import "./styles/Game.css";

export default class Game extends Component {
	constructor(props) {
		super(props);

		this.state = {
			sudoku: [],
			mustGenSudoku: true,
			isLoading: true,
			active_row: null,
			active_col: null,
			blueprint: [],
		};

		this.selectRef = React.createRef();

		this.handleLogout = this.handleLogout.bind(this);
		this.generateNewSudoku = this.generateNewSudoku.bind(this);
		this.handleCellClick = this.handleCellClick.bind(this);
		this.handleControlClick = this.handleControlClick.bind(this);
		this.handleSave = this.handleSave.bind(this);
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
					blueprint: result.data,
				};
			});
		} catch (e) {
			this.setState(() => {
				return { ...this.state, isLoading: false };
			});
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
					blueprint: res.data,
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

	handleCellClick(e) {
		const click = e.target.id.split("_");

		const row = parseInt(click[1]);
		const col = parseInt(click[2]);

		if (this.state.active_col === col && this.state.active_row === row) {
			//deactivate current cell
			this.setState(() => {
				return { ...this.state, active_col: null, active_row: null };
			});
			return;
		}

		this.setState(() => {
			return { ...this.state, active_col: col, active_row: row };
		});
	}

	handleControlClick(e) {
		const num = parseInt(e.target.id.split("_")[1]);

		if (
			this.state.blueprint[this.state.active_row][this.state.active_col] === 0
		) {
			alert("Cannot change that cell");
		} else {
			const arr = this.state.sudoku;
			arr[this.state.active_row][this.state.active_col] = num;
			this.setState(() => {
				return {
					...this.state,
					sudoku: arr,
				};
			});
		}
	}

	async handleSave() {
		//TODO
		//save sudoku progress in db
		try {
			const res = await axios({
				method: "post",
				url: "/save-sudoku",
				data: querystring.stringify({
					sudoku: JSON.stringify(this.state.sudoku),
				}),
				withCredentials: true,
			});
			console.log(res);
		} catch (e) {
			alert("Could not save progress");
		}
	}

	render() {
		let row = -1;
		let col = -1;
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
						<div>
							<select ref={this.selectRef} defaultValue="1">
								<option value="0">Easy</option>
								<option value="1">Medium</option>
								<option value="2">Hard</option>
							</select>
							<button onClick={this.generateNewSudoku}>Generate</button>
						</div>
						{this.state.mustGenSudoku ? null : (
							<div>
								<div>
									{this.state.active_col !== null &&
									this.state.active_row !== null ? (
										[1, 2, 3, 4, 5, 6, 7, 8, 9].map((ele) => {
											return (
												<span
													key={`control_${ele}`}
													id={`control_${ele}`}
													onClick={this.handleControlClick}
												>
													{ele}
												</span>
											);
										})
									) : (
										<div>Select a cell</div>
									)}
								</div>
								<div>
									{this.state.sudoku.map((ele) => {
										row += 1;
										col = -1;
										return (
											<div key={`row_${row}`}>
												{ele.map((item) => {
													col += 1;
													return (
														<span
															key={`item_${row}_${col}`}
															id={`item_${row}_${col}`}
															onClick={this.handleCellClick}
															className={
																this.state.active_col === col &&
																this.state.active_row === row
																	? "Game-Cell-Selected"
																	: null
															}
														>
															{item}
														</span>
													);
												})}
											</div>
										);
									})}
								</div>
								<button onClick={this.handleSave}>Save</button>
								<button>Submit</button>
							</div>
						)}
					</div>
				)}
			</div>
		);
	}
}
