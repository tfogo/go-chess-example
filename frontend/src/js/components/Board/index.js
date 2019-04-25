import React, { Component } from "react";
import { connect } from 'react-redux'
import PropTypes from "prop-types";
import { makeMove } from '../../actions'

import Chess from "chess.js"; // import Chess from  "chess.js"(default) if recieving an error about new Chess() not being a constructor

import Chessboard from "chessboardjsx";

class HumanVsHumanView extends Component {
  static propTypes = { children: PropTypes.func };

  

  state = {
    fen: "start",
    // square styles for active drop square
    dropSquareStyle: {},
    // custom square styles
    squareStyles: {},
    // square with the currently clicked piece
    pieceSquare: "",
    // currently clicked square
    square: "",
    // array of past game moves
    history: []
  };

  componentDidMount() {
    let { game } = this.props
    this.game = game;
  }

  // keep clicked square style and remove hint squares
//   removeHighlightSquare = () => {
//     this.setState(({ pieceSquare, history }) => ({
//       squareStyles: squareStyling({ pieceSquare, history })
//     }));
//   };

  // show possible moves
//   highlightSquare = (sourceSquare, squaresToHighlight) => {
//     const highlightStyles = [sourceSquare, ...squaresToHighlight].reduce(
//       (a, c) => {
//         return {
//           ...a,
//           ...{
//             [c]: {
//               background:
//                 "radial-gradient(circle, #fffc00 36%, transparent 40%)",
//               borderRadius: "50%"
//             }
//           },
//           ...squareStyling({
//             history: this.state.history,
//             pieceSquare: this.state.pieceSquare
//           })
//         };
//       },
//       {}
//     );

//     this.setState(({ squareStyles }) => ({
//       squareStyles: { ...squareStyles, ...highlightStyles }
//     }));
//   };

  onDrop = ({ sourceSquare, targetSquare }) => {
    let { gameId,  game, username, black, handleMakeMove } = this.props

    let isBlack = true
    if (username !== black) {
        isBlack = false
    }

    if (!isBlack && game.turn() == "b") return;
    if (isBlack && game.turn() == "w") return;

    // see if the move is legal
    let move = game.move({
      from: sourceSquare,
      to: targetSquare,
      promotion: "q" // always promote to a queen for example simplicity
    });

    // illegal move
    if (move === null) return;
    // this.setState(({ history, pieceSquare }) => ({
    //   fen: this.game.fen(),
    //   history: this.game.history({ verbose: true }),
    //   squareStyles: squareStyling({ pieceSquare, history })
    // }));

    handleMakeMove(gameId, game.history())
  };

//   onMouseOverSquare = square => {
//     let { username, black } = this.props

//     let isBlack = true
//     if (username !== black) {
//         isBlack = false
//     }

//     if (!isBlack && this.game.turn() == "b") return;
//     if (isBlack && this.game.turn() == "w") return;

//     // get list of possible moves for this square
//     let moves = this.game.moves({
//       square: square,
//       verbose: true
//     });

//     // exit if there are no moves available for this square
//     if (moves.length === 0) return;

//     let squaresToHighlight = [];
//     for (var i = 0; i < moves.length; i++) {
//       squaresToHighlight.push(moves[i].to);
//     }

//     this.highlightSquare(square, squaresToHighlight);
//   };

//   onMouseOutSquare = square => this.removeHighlightSquare(square);

  // central squares get diff dropSquareStyles
//   onDragOverSquare = square => {
//     this.setState({
//       dropSquareStyle:
//         square === "e4" || square === "d4" || square === "e5" || square === "d5"
//           ? { backgroundColor: "cornFlowerBlue" }
//           : { boxShadow: "inset 0 0 1px 4px rgb(255, 255, 0)" }
//     });
//   };

  onSquareClick = square => {
    let { game, gameId, handleMakeMove } = this.props

    // this.setState(({ history }) => ({
    //   squareStyles: squareStyling({ pieceSquare: square, history }),
    //   pieceSquare: square
    // }));

    let move = game.move({
      from: this.state.pieceSquare,
      to: square,
      promotion: "q" // always promote to a queen for example simplicity
    });

    // illegal move
    if (move === null) return;

    this.setState({
    //   fen: this.game.fen(),
    //   history: this.game.history({ verbose: true }),
      pieceSquare: ""
    });

    handleMakeMove(gameId, game.history())
    // this.game = new Chess('r1k4r/p2nb1p1/2b4p/1p1n1p2/2PP4/3Q1NB1/1P3PPP/R5K1 b - c3 0 19')
  };

//   onSquareRightClick = square =>
//     this.setState({
//       squareStyles: { [square]: { backgroundColor: "deepPink" } }
//     });

  render() {
    let { username, black, game } = this.props
    const { fen, dropSquareStyle, squareStyles } = this.state;

    let orientation = "white"
    if (username === black) {
        orientation = "black"
    }

    return this.props.children({
      squareStyles,
      orientation,
      position: game.fen(),
      onMouseOverSquare: this.onMouseOverSquare,
      onMouseOutSquare: this.onMouseOutSquare,
      onDrop: this.onDrop,
      dropSquareStyle,
      onDragOverSquare: this.onDragOverSquare,
      onSquareClick: this.onSquareClick,
      onSquareRightClick: this.onSquareRightClick
    });
  }
}

const mapStateToProps = (state) => {
    return {
        ...state
    }
}

const mapDispatchToProps = (dispatch) => ({
    handleMakeMove(gameId, move) {
        dispatch(makeMove(gameId, move))
    }
})
      
const HumanVsHuman = connect(
    mapStateToProps,
    mapDispatchToProps
)(HumanVsHumanView)
    

export default function WithMoveValidation() {
  return (
    <div>
      <HumanVsHuman>
        {({
          position,
          onDrop,
          orientation,
          onMouseOverSquare,
          onMouseOutSquare,
          squareStyles,
          dropSquareStyle,
          onDragOverSquare,
          onSquareClick,
          onSquareRightClick
        }) => (
          <Chessboard
            id="humanVsHuman"
            orientation={orientation}
            
            position={position}
            onDrop={onDrop}
            onMouseOverSquare={onMouseOverSquare}
            onMouseOutSquare={onMouseOutSquare}
            boardStyle={{
              borderRadius: "5px",
            }}
            squareStyles={squareStyles}
            dropSquareStyle={dropSquareStyle}
            onDragOverSquare={onDragOverSquare}
            onSquareClick={onSquareClick}
            onSquareRightClick={onSquareRightClick}
          />
        )}
      </HumanVsHuman>
    </div>
  );
}




const squareStyling = ({ pieceSquare, history }) => {
  const sourceSquare = history.length && history[history.length - 1].from;
  const targetSquare = history.length && history[history.length - 1].to;

  return {
    [pieceSquare]: { backgroundColor: "rgba(255, 255, 0, 0.4)" },
    ...(history.length && {
      [sourceSquare]: {
        backgroundColor: "rgba(255, 255, 0, 0.4)"
      }
    }),
    ...(history.length && {
      [targetSquare]: {
        backgroundColor: "rgba(255, 255, 0, 0.4)"
      }
    })
  };
};

