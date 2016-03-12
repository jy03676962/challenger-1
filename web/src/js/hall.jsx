import React from 'react';
import {observer} from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/base.css'

const Hall = observer(React.createClass({
  render: function() {
    var content
    if (this.props.game.room) {
      content = <RoomView game={this.props.game} actionFunc={this.joinOrStart} />
    } else {
      content = <button onClick={this.createRoom}>创建房间</button>
    }
    return (
      <div>
      <h1>欢迎，{this.props.game.playerName}</h1>
      {content}
      </div>
    );
  },
  createRoom: function(e) {
    this.props.game.createRoom()
  },
  joinOrStart: function(e) {
    if (this.props.game.room.hoster == this.props.game.playerName) {
      this.props.game.startGame()
    } else {
      this.props.game.joinRoom()
    }
  }
}));

const RoomView = CSSModules(observer(
  ({game, actionFunc}) => {
    let room = game.room
    let joined = room.member.indexOf(game.playerName) >= 0
    let button = null
    if (joined) {
      if (room.hoster == game.playerName) {
        button = <li><button onClick={actionFunc}>开始</button></li>
      }
    } else {
      button = <li><button onClick={actionFunc}>加入</button></li>
    }
    return(
    <ul styleName='menu'>
    <li style={{color:"red"}}>{room.hoster}</li>
    <li>{room.member.length + "/" + room.max}</li>
    {button}
    </ul>
    )
}), styles)

export default Hall
