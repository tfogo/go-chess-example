import Api from '../services'
import * as types from '../constants/actionTypes'

export const initWebSocket = (dispatch) => {
  const webSocket = new WebSocket(`ws://localhost:6061/ws`)
  webSocket.onmessage = (event) => dispatch({ type: types.WS_MSG_RECEIVED, payload: JSON.parse(event.data) })
}