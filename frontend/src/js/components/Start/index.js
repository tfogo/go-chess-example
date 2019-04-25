import React from "react"
import { connect } from 'react-redux'
import { changeName, dispatchStartGame } from '../../actions'

class StartView extends React.Component {  
    render() {

        const { searching, handleStartGame, handleChangeName, username } = this.props

        const changeHandler = (e) => {
            handleChangeName(e.target.value)
        }

        const handleSubmit = (e) => {
            e.preventDefault()
            handleStartGame(username)
          }
      
      return (
        <div>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="name">Name</label>
              <input onChange={changeHandler} className="form-control" id="name" aria-describedby="nameHelp" placeholder="Enter name" value={username}></input>
              <small id="nameHelp" className="form-text text-muted">The username for your chess match</small>
            </div>
            <button disabled={searching} type="submit" className="btn btn-primary">{searching ? "Searching..." : "Start Game"}</button>
          </form>
        </div>
      )
    }
  }

const mapStateToProps = (state) => {
    return {
        ...state
    }
}


const mapDispatchToProps = (dispatch) => ({
    handleChangeName(move) {
        dispatch(changeName(move))
    },
    handleStartGame(username) {
        dispatch(dispatchStartGame(username))
    },
})
      

const Start = connect(
    mapStateToProps,
    mapDispatchToProps
)(StartView)

export default Start