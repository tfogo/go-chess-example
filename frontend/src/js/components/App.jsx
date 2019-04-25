import React from "react"
import { connect } from 'react-redux'
import Header from './Header'
import Start from './Start'
import WithMoveValidation from "./Board"

class Demo extends React.Component {
  render() {

    return (
      <div>
        <div style={boardsContainer}>
          <WithMoveValidation />
          <div style={infoStyle}>
            <h5>White: {this.props.white}</h5>
            <h5>Black: {this.props.black}</h5>
            <p>{this.props.history !== undefined ? this.props.history.join(" ") : ""}</p>
          </div>
        </div>
        
      </div>
    );
  }
}

const infoStyle = {
  marginLeft: 30
};

const boardsContainer = {
  display: "flex",
  justifyContent: "space-around",
  alignItems: "center",
  flexWrap: "wrap",
  width: "100vw",
  marginTop: 30,
  marginBottom: 50
};


class AppView extends React.Component {  
  render() {
    let { gameStarted, history, white, black } = this.props
    
    return (

        <div className="h-100">
        <Header />
        
        <div className="container-fluid hero-video h-75">
          <div className="row">
            <div className="col-12">
            {
              gameStarted 
                ? <Demo history={history} black={black} white={white}></Demo>
                : <Start></Start>
            }
            </div>
          </div>
        </div>
      </div>

      
    )
  }
}

const mapStateToProps = (state) => {
  return {
    ...state
  }
}

const App = connect(
  mapStateToProps
)(AppView)

export default App
