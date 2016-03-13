import React from 'react';
import {observer} from 'mobx-react'

const Arena = observer(React.createClass({
  render() {
    let opt = this.props.game.options
    let arenaWidth = (opt.arenaCellSize + opt.arenaBorder) * opt.arenaWidth * opt.webScale
    let infoStyle = {
      width: arenaWidth + "px",
      height: "80px",
      margin: "auto",
      textAlign: "center",
    }
    let bgStyle = {
      width: arenaWidth + "px",
      fontSize: "0",
      margin: "auto",
      position: "relative",
      border: opt.arenaBorder / 2 * opt.webScale + "px solid blue"
    }
    let gStyle = {
      width: arenaWidth + "px",
      margin: "auto",
      position: "relative",
    }
    return (
      <div>
      <ArenaInfoBar game={this.props.game} rootStyle={infoStyle}/>
      <ArenaBackground opt={opt} rootStyle={bgStyle} />
      <ArenaGround opt={opt} rootStyle={gStyle} />
      </div>
    )
  }
}))

const ArenaInfoBar = observer(React.createClass({
  render() {
    let game = this.props.game
    let content
    if (game.match.stage == "warmup") {
      let second = (game.options.warmup - game.match.elasped).toFixed(1)
      content = `准备中,还剩${second}`
    } else {
      content = `游戏开始`
    }
    return (
      <div style={this.props.rootStyle}>{content}</div>
      )
  }
}))

const ArenaBackground = ({opt, rootStyle}) => {
  let size = opt.arenaCellSize * opt.webScale + "px"
  let cellStyle = {
    width: size,
    height: size,
    display: "inline-block",
    border: opt.arenaBorder / 2 * opt.webScale + "px solid #CCCCCC",
    backgroundColor: "#669900",
  }
  let elements = []
  for (let i = 0; i < opt.arenaHeight; i++) {
    for (let j = 0; j < opt.arenaWidth; j++) {
      elements.push(<div style={cellStyle} key={"cell:"+i * opt.arenaWidth + j}></div>)
    }
  }
  for (let [index, wall] of opt.walls.entries()) {
    let horizontal = wall.P1.X == wall.P2.X
    let w, h, l, t
    if (horizontal) {
      w = opt.arenaCellSize + 2 * opt.arenaBorder
      h = opt.arenaBorder
      l = wall.P1.X * (opt.arenaCellSize + opt.arenaBorder) - opt.arenaBorder / 2
      t = Math.max(wall.P1.Y, wall.P2.Y) * (opt.arenaCellSize + opt.arenaBorder) - opt.arenaBorder / 2
    } else {
      w = opt.arenaBorder
      h = opt.arenaCellSize + 2 * opt.arenaBorder
      t = wall.P1.Y * (opt.arenaCellSize + opt.arenaBorder) - opt.arenaBorder / 2
      l = Math.max(wall.P1.X, wall.P2.X) * (opt.arenaCellSize + opt.arenaBorder) - opt.arenaBorder / 2
    }
    let wallStyle = {
      position: "absolute",
      backgroundColor: "blue",
      left: l * opt.webScale + "px",
      top: t * opt.webScale + "px",
      width: w * opt.webScale + "px",
      height: h * opt.webScale + "px",
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
