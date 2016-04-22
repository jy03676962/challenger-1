import React from 'react';
import { observer } from 'mobx-react'

const Hall = observer(React.createClass({
  render: function() {
    var content
    if (this.props.game.match) {
      content = <MatchView game={this.props.game} actionFunc={this.joinOrStart} />
    } else {
      content = <button onClick={this.createMatch}>创建房间</button>
    }
    return (
      <div>
      <h1>欢迎，{this.props.game.playerName}</h1>
      {content}
      </div>
    );
  },
  createMatch: function(e) {
    this.props.game.createMatch()
  },
}));

const MatchView = observer(React.createClass({
  render: function() {
    let match = this.props.game.match
    let hoster = match.member[0].name
    let joined = match.member.filter((member) => {
      return member.name == this.props.game.playerName
    }).length > 0
    let actionComponent = null
    if (joined) {
      if (hoster == this.props.game.playerName) {
        actionComponent =
          <div>
        <button onClick={this.startFunMode}>开始娱乐模式</button>
        <button onClick={this.startSurvivalMode}>开始生存模式</button>
        </div>
      }
    } else {
      actionComponent = <button onClick={this.joinMatch}>加入</button>
    }
    return (
      <div>
      <div style={{color:'red'}}>{'房主:' + hoster}</div>
      {
        match.member.map((member) =>{
          return <div style={{color:'green'}} key={'player:'+member.name}>{member.name}</div>
        })
      }
      {actionComponent}
      </div>
    )
  },
  startFunMode: function(e) {
    this.props.game.startMatch(1)
  },
  startSurvivalMode: function(e) {
    this.props.game.startMatch(2)
  },
  joinMatch: function(e) {
    this.props.game.joinMatch()
  },

}))

export default Hall
