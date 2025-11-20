# Go Sudoku

Terminal User Interface (TUI) sudoku game with the posibility of playing multiplayer via web.

## Components

The game has two main components, both accessible from the main menu when you run the app.

### Game server

Creates a game server on the desplayed port to play via browser with up to 3 other players.

![Game Server](https://github.com/DNLSalazar/gosudoku/raw/master/resources/server.gif)

### Terminal game

Play the game via terminal by yourself, it supports two input methods, movements (using arrows, hjkl and wasd) and input mode. It includes help page with instructions.

![TUI Game](https://github.com/DNLSalazar/gosudoku/raw/master/resources/tui.gif)

## Installation

With go installed, run `go install github.com/DNLSalazar/gosudoku@latest`.

## Usage

Run `sudoku` and select how you want to play, additionally, you can press `ctrl+H` to go to the help page and get more details about the available modes and options.

## TODO

This is a project still in development, I want to add some other features to it:

- I want to add a sudoku board generator, so every game feels different. Right now the board is always the same.
- I'm plannign to make the web version more fun and dynamic to play using the keyboard to move, making it similar to the terminal version.
- Improve the UI of the menus.
- Implement reconnect on multiplayer game
