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
		this.handleSubmit = this.handleSubmit.bind(this);
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
				let blueprint = [];
				for (let i = 0; i < 9; i++) {
					let temp = [];
					for (let j = 0; j < 9; j++) {
						temp.push(result.data[i][j]);
					}
					blueprint.push(temp);
				}
				return {
					...this.state,
					sudoku: result.data,
					isLoading: false,
					mustGenSudoku: false,
					blueprint: blueprint,
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
				let blueprint = [];
				for (let i = 0; i < 9; i++) {
					let temp = [];
					for (let j = 0; j < 9; j++) {
						temp.push(res.data[i][j]);
					}
					blueprint.push(temp);
				}

				return {
					...this.state,
					mustGenSudoku: false,
					isLoading: false,
					sudoku: res.data,
					blueprint: blueprint,
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

		if (this.state.blueprint[row][col] !== 0) {
			return;
		}

		this.setState(() => {
			return { ...this.state, active_col: col, active_row: row };
		});
	}

	handleControlClick(e) {
		const num = parseInt(e.target.id.split("_")[1]);

		const arr = this.state.sudoku;
		arr[this.state.active_row][this.state.active_col] = num;
		this.setState(() => {
			return {
				...this.state,
				sudoku: arr,
			};
		});
	}

	async handleSubmit() {
		try {
			await axios({
				method: "post",
				url: "submit-sudoku",
				data: querystring.stringify({
					sudoku: JSON.stringify(this.state.sudoku),
				}),
				withCredentials: true,
			});

			//reset the game
			this.setState(() => {
				return {
					...this.state,
					sudoku: [],
					mustGenSudoku: true,
					isLoading: false,
					active_row: null,
					active_col: null,
					blueprint: [],
				};
			});
		} catch (e) {
			const status = e.response.status;
			if (status === 403) {
				alert("sudoku is not valid");
			} else {
				alert("unable to submit sudoku");
			}
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

	async handleStatsClick() {
		try {
			const res = await axios({
				method: "post",
				url: "/get-stats",
				withCredentials: true,
			});
			const data = res.data;
			alert(
				`easy:${data["easy"]}\nmedium:${data["medium"]}\nhard:${data["hard"]}`
			);
		} catch (e) {
			alert("could not get stats");
		}
	}

	render() {
		let row = -1;
		let col = -1;
		return (
			<div className="Game">
				<div onClick={this.handleLogout} className="button-regular">
					Logout
				</div>
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
						<div className="Game-generate">
							<select ref={this.selectRef} defaultValue="1">
								<option value="0">Easy</option>
								<option value="1">Medium</option>
								<option value="2">Hard</option>
							</select>
							<div onClick={this.generateNewSudoku} className="button-regular">
								Generate
							</div>
							<div onClick={this.handleStatsClick} className="button-regular">
								Stats
							</div>
						</div>
						{this.state.mustGenSudoku ? null : (
							<div>
								<div className="Game-Control">
									{this.state.active_col !== null &&
									this.state.active_row !== null ? (
										[1, 2, 3, 4, 5, 6, 7, 8, 9, 0].map((ele) => {
											return (
												<div
													key={`control_${ele}`}
													id={`control_${ele}`}
													onClick={this.handleControlClick}
													className="Game-Control-Cell"
												>
													{ele !== 0 ? ele : "X"}
												</div>
											);
										})
									) : (
										<div className="Game-Hint">Select a cell</div>
									)}
								</div>
								<div className="Game-Main-Game">
									{this.state.sudoku.map((ele) => {
										row += 1;
										col = -1;
										return (
											<div key={`row_${row}`} className="Game-Main-Game-Row">
												{ele.map((item) => {
													col += 1;
													return (
														<div
															key={`item_${row}_${col}`}
															id={`item_${row}_${col}`}
															onClick={this.handleCellClick}
															className={
																this.state.active_col === col &&
																this.state.active_row === row
																	? "Game-Cell-Selected Game-Main-Game-Cell"
																	: "Game-Main-Game-Cell"
															}
														>
															{item === 0 ? "" : item}
														</div>
													);
												})}
											</div>
										);
									})}
								</div>
								<div className="Game-Submit-Save">
									<div
										onClick={this.handleSave}
										className="button-regular button-submit"
									>
										Save
									</div>
									<div
										onClick={this.handleSubmit}
										className="button-regular button-submit"
									>
										Submit
									</div>
								</div>
							</div>
						)}
					</div>
				)}
			</div>
		);
	}
}
