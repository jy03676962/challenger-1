import React from 'react';
import {observer} from 'mobx-react'
import CSSModules from 'react-css-modules'

const Arena = observer(React.createClass({
  render() {
    let arg = this.props.game.arg
    let arenaWidth = arg.arenaCellSize * arg.arenaWidth
    let bgStyle = {
      width: arenaWidth + "px",
      fontSize: "0",
      margin: "auto",
    }
    let gStyle = {
      width: arenaWidth + "px",
      margin: "auto",
      position: "relative",
    }
    return (
      <div>
      <ArenaBackground arg={this.props.game.arg} rootstyle={bgStyle} />
      <ArenaGround arg={this.props.game.arg} rootStyle={gStyle} />
      </div>
    );
  }
}))

const ArenaBackground = ({arg, rootStyle}) => {
  let size = arg.arenaCellSize - 2 * arg.arenaBorder + "px"
  let cellStyle = {
    width: size,
    height: size,
    display: "inline-block",
    border: arg.arenaBorder + "px solid black",
    backgroundColor: "green",
  }
  let rows = []
  for (let i = 0; i < arg.arenaHeight; i++) {
    for (let j = 0; j < arg.arenaWidth; j++) {
      rows.push(<div style={cellStyle} key={i * arg.arenaWidth + j}></div>)
    }
  }
  return (
  <div style={rootStyle}>
  {rows}
  </div>
  );
}

const ArenaGround = observer(React.createClass({
  render() {
  }
}))

export default Arena
