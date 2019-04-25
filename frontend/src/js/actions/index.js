import Api from '../services'
import * as types from '../constants/actionTypes'

let webSocket

export const initWebSocket = (dispatch, id) => {
  webSocket = new WebSocket(`ws://localhost:9001/ws?id=${id}`)
  webSocket.onmessage = (event) => dispatch({ type: types.WS_MSG_RECEIVED, payload: JSON.parse(event.data) })
}

export const makeMove = (gameId, move) => (dispatch) => {
  webSocket.send(JSON.stringify({ID: gameId, move: move[move.length-1]}))
  dispatch({ type: types.MK_MOVE })
}

export const changeName = (name) => ({
  type: types.CH_NAME,
  payload: name
})

export const dispatchStartGame = (username) => (dispatch) => {
  dispatch({ type: types.START_GAME_PENDING })
  return Api.startGame(username).then(
    res => {
      dispatch({ type: types.START_GAME_FULFILLED, payload: res.data })
      initWebSocket(dispatch, res.data.gameId)
    },
    err => {
      dispatch({ type: types.START_GAME_FAILED, err })
    }
  )
}
