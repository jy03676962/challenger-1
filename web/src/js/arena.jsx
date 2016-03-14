import React from 'react';
import {observer} from 'mobx-react'
import Player from '~/js/player.jsx'
import Scheme from '~/js/scheme.jsx'

const Arena = observer(React.createClass({
  render() {
    let opt = this.props.game.options
    let arenaWidth = (opt.arenaCellSize + opt.arenaBorder) * opt.arenaWidth * opt.webScale
    let arenaHeight = (opt.arenaCellSize + opt.arenaBorder) * opt.arenaHeight * opt.webScale
    const infoHeight = 40
    let infoStyle = {
      width: arenaWidth + "px",
      height: infoHeight + "px",
      margin: "auto",
      textAlign: "center",
    }
    let bgStyle = {
      width: arenaWidth + "px",
      fontSize: "0",
      margin: "auto",
      position: "relative",
      border: opt.arenaBorder / 2 * opt.webScale + "px solid " + Scheme.wall
    }
    let gStyle = {
      position: "absolute",
      top: infoHeight + "px",
      left: "0",
      right: "0",
      margin: "auto",
      width: arenaWidth + "px",
      height: arenaHeight + "px",
      border: opt.arenaBorder / 2 * opt.webScale + "px solid " + Scheme.wall
    }
    return (
      <div style={{position:"relative"}}>
      <ArenaInfoBar game={this.props.game} rootStyle={infoStyle}/>
      <ArenaBackground opt={opt} rootStyle={bgStyle} />
      <ArenaButtonLayer game={this.props.game} rootStyle={gStyle} />
      <ArenaGround game={this.props.game} rootStyle={gStyle} />
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
  let elements = []
  for (let i = 0; i < opt.arenaHeight; i++) {
    for (let j = 0; j < opt.arenaWidth; j++) {
      let cellStyle = {
        width: size,
        height: size,
        display: "inline-block",
        border: opt.arenaBorder / 2 * opt.webScale + "px solid " + Scheme.border,
        backgroundColor: Scheme.normalTile
      }
      if (j == opt.arenaEntrance.X && i == opt.arenaEntrance.Y) {
        cellStyle.backgroundColor = Scheme.entranceTile
      } else if (j == opt.arenaExit.X && i == opt.arenaExit.Y) {
        cellStyle.backgroundColor = Scheme.exitTile
      } else {
        cellStyle.backgroundColor= Scheme.normalTile
      }
      elements.push(<div style={cellStyle} key={"cell:"+i * opt.arenaWidth + j}></div>)
    }
  }
  for (let [index, wall] of opt.walls.entries()) {
    let wallStyle = {
      position: "absolute",
      backgroundColor: Scheme.wall,
      left: wall.X * opt.webScale + "px",
      top: wall.Y * opt.webScale + "px",
      width: wall.W * opt.webScale + "px",
      height: wall.H * opt.webScale + "px",
    }
    elements.push(<div style={wallStyle} key={"wall:" + index}></div>)
  }
  return (
  <div style={rootStyle}>
  {elements}
  </div>
  );
}

const ArenaButtonLayer = observer(React.createClass({
  render() {
    let match = this.props.game.match
    let opt = this.props.game.options
    return (
      <div style={this.props.rootStyle}>
      {
        opt.buttons.map((button) => {
          let color = Scheme.buttonT0
          if (match.rampage) {
            color = Scheme.buttonRampage
          } else {
            let t = match.buttonState[button.id]
            if (typeof t != "undefined") {
              color = Scheme[`buttonT${t}`]
            }
          }
          let r = button.r
          let buttonStyle = {
            position: "absolute",
            backgroundColor: color,
            left: r.X * opt.webScale + "px",
            top: r.Y * opt.webScale + "px",
            width: r.W * opt.webScale + "px",
            height: r.H * opt.webScale + "px",
          }
          return <div style={buttonStyle} key={"button:" + button.id}></div>
        })
      }
      </div>
      )
  }
}))

const ArenaGround = observer(React.createClass({
  render() {
    let match = this.props.game.match
    let opt = this.props.game.options
    return (
      <div style={this.props.rootStyle}>
      {
        match.member.map((member, idx) => {
          return <Player player={member} options={opt} key={member.name} idx={idx}/>
        })
      }
      </div>
      )
  }
}))

export default Arena
