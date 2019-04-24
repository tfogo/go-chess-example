import React from "react"
import { connect } from 'react-redux'
import Header from './Header'
import { Switch, Route, Link } from 'react-router-dom'

class AppView extends React.Component {  
  render() {
    let { currentFrame } = this.props
    
    return (
      <div className="h-100">
        <Header />
        <div className="container-fluid hero-video h-75">
          <div className="row">
            <div className="col-12"><h5>Time: {currentFrame.Time}</h5></div>
            <div className="col"><p>GOMAXPROCS: {currentFrame.Gomaxprocs}</p></div>
            <div className="col"><p>Idle Processes: {currentFrame.Idleprocs}</p></div>
            <div className="col"><p>Idle Threads: {currentFrame.Idlethreads}</p></div>
            <div className="col"><p>Threads (Ms): {currentFrame.Threads}</p></div>
            <div className="col"><p>Garbage Collector Waiting: {currentFrame.Gcwaiting}</p></div>
            <div className="col"><p>Threads Locked: {currentFrame.Nmidlelocked}</p></div>
            <div className="col"><p>Global Queue: {currentFrame.Runqueue}</p></div>
            <div className="col"><p>Spinning Threads: {currentFrame.Spinningthreads}</p></div>
            <div className="col"><p>Stop Wait: {currentFrame.Stopwait}</p></div>
            <div className="col"><p>Sysmon Wait: {currentFrame.Sysmonwait}</p></div>
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
