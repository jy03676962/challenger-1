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
      position: "relative",
      border: arg.arenaBorder + "px solid blue"
    }
    let gStyle = {
      width: arenaWidth + "px",
      margin: "auto",
      position: "relative",
    }
    return (
      <div>
      <ArenaBackground arg={this.props.game.arg} rootStyle={bgStyle} />
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
    border: arg.arenaBorder + "px solid #CCCCCC",
    backgroundColor: "#669900",
  }
  let elements = []
  for (let i = 0; i < arg.arenaHeight; i++) {
    for (let j = 0; j < arg.arenaWidth; j++) {
      elements.push(<div style={cellStyle} key={"cell:"+i * arg.arenaWidth + j}></div>)
    }
  }
  for (let [index, wall] of arg.walls.entries()) {
    let horizontal = wall.P1.X == wall.P2.X
    let w, h, l, t
    if (horizontal) {
      w = arg.arenaCellSize
      h = 2 * arg.arenaBorder
      l = wall.P1.X * arg.arenaCellSize
      t = Math.max(wall.P1.Y, wall.P2.Y) * arg.arenaCellSize - arg.arenaBorder
    } else {
      w = 2 * arg.arenaBorder
      h = arg.arenaCellSize
      t = wall.P1.Y * arg.arenaCellSize
      l = Math.max(wall.P1.X, wall.P2.X) * arg.arenaCellSize - arg.arenaBorder
    }
    // vertical is {{top, left}, {height, width}}
    // horizontal is {left, top}, {width, height}}
    let wallStyle = {
      position: "absolute",
      backgroundColor: "blue",
      left: l + "px",
      top: t + "px",
      width: w + "px",
      height: h + "px",
    }
    elements.push(<div style={wallStyle} key={"wall:" + index}></div>)
  }
  return (
  <div style={rootStyle}>
  {elements}
  </div>
  );
}

const ArenaGround = observer(React.createClass({
  render() {
    return (<div></div>)
  }
}))

export default Arena
