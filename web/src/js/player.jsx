import React from 'react';
import {observer} from 'mobx-react'

const Player = observer(React.createClass({
  render() {
    const arrowSize = 5
    let player = this.props.player
    let scale = this.props.options.webScale
    let bodySize = this.props.options.playerSize * scale
    let size = bodySize + 2 * arrowSize
    let degree, arrowX, arrowY
    switch (player.dir) {
      case "up":
      degree = 0
      arrowX = (size - arrowSize) / 2
      arrowY = 0
      break
      case "left":
      degree = 90
      arrowX = size - arrowSize
      arrowY = (size - arrowSize) / 2
      break
      case "down":
      degree = 180
      arrowX = (size - arrowSize) / 2
      arrowY = size -arrowSize
      break
      case "right":
      degree = 270
      arrowX = 0
      arrowY = (size - arrowSize) / 2
      break
    }
    let style = {
      position: "absolute",
      width: size + "px",
      height: size + "px",
      top: player.pos.Y * scale - size / 2 + "px",
      left: player.pos.X * scale - size / 2 + "px",
    }
    let bodyStyle = {
      textAlign: "center",
      backgroundColor: player.color,
      borderRadius: bodySize / 2 + "px",
      position: "absolute",
      width: bodySize + "px",
      height: bodySize + "px",
      top: arrowSize + "px",
      left: arrowSize + "px",
    }
    let imgStyle = {
      display: "block",
      position: "absolute",
      width: arrowSize + "px",
      height: arrowSize + "px",
      transform: `rotate(${degree})`,
      top: arrowY + "px",
      left: arrowX + "px",
    }
    return (
      <div style={style}>
      <div style={bodyStyle}>
      {player.name}
      </div>
      <img style={imgStyle} src={require('../assets/arrow.jpg')} />
      </div>
    );
  }
}))

export default Player


