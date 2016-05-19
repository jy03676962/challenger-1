import React from 'react'
import { render } from 'react-dom'
import { observable, computed } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/queue.css'
import * as util from '~/js/util.jsx'

class Queue {
  @ observable data
  @ observable connected
  @ observable match
  constructor() {
    this._reset()
  }

  _reset() {
    this.sock = null
    this.state = false
    this.data = null
    this.match = null
  }

  connect() {
    if (this.sock) {
      this.sock.close()
      this._reset()
      return
    }
    let uri = util.wsAddressWithPath('ws')
    let sock = new WebSocket(uri)
    sock.onopen = () => {
      let data = {
        cmd: 'init',
        ID: 'queue',
        TYPE: '8',
      }
      sock.send(JSON.stringify(data))
    }
    sock.onclose = (e) => {
      this._reset()
    }
    sock.onmessage = (e) => {
      this.onMessage(e.data)
    }
    this.sock = sock
  }

  onMessage(msg) {
    let json = JSON.parse(msg)
    switch (json.cmd) {
      case 'init':
        this.connected = true
        break
      case 'matchData':
        if (json.data != null && this.connected) {
          this.data = json.data
        }
        break
      case 'matchStop':
        this.match = null
      case 'updateMatch':
        if (json.data != null && this.connected) {
          this.match = JSON.parse(json.data)
        }
    }
  }

  send(data) {
    if (this.sock) {
      let d = JSON.stringify(data)
      this.sock.send(d)
    }
  }
}

const HistoryCellView = CSSModules(React.createClass({
  render() {
    let matchData = this.props.data
    let idx = this.props.idx
    let top = (idx * 58) / 10.8 + 'vw'
    if (matchData.mode = 'g') {
      var modeImg = require('./assets/g_icon.png')
      var result = `获得：${matchData.gold}G`
    } else {
      var modeImg = require('./assets/s_icon.png')
      var result = `生存：${util.timeStr(matchData.elasped)}`
    }
    let playerStr = matchData.member.map((player, idx) => {
      if (player.name) {
        return player.name
      } else {
        return 'P' + player.cid.split(':')[1]
      }
    }).join(' ')
    let style = {
      position: 'absolute',
      width: '100%',
      height: '5.185vw',
      top: top,
    }
    return (
      <div style={style}>
        <img src={require('./assets/late0.png')} styleName='historyCellImg' />
        <div styleName='historyNumber'>{matchData.teamID}</div>
        <div styleName='historyPlayer'>{playerStr}</div>
        <img src={modeImg} styleName='historyIcon' />
        <div styleName='historyResult'>{result}</div>
      </div>
    )
  }
}), styles)

const HistoryView = CSSModules(React.createClass({
  render() {
    let history = this.props.history
    if (history == null || history.length == 0) {
      return null
    }
    return (
      <div styleName='history'>
        {
        history.sort((a, b) => {
          return a.id > b.id
        }).map((m, idx) => {
          return <HistoryCellView data={m} idx={idx} key={idx}/>
        })
        }
      </div>
    )
  }
}), styles)

const CurrentMatchView = CSSModules(React.createClass({
  render() {
    let match = this.props.match
    if (match == null) {
      return null
    }
    let bg = match.mode == 'g' ? require('./assets/g_game_bg.png') : require('./assets/s_game_bg.png')
    let color = match.mode == 'g' ? '#dc8524' : '#03dceb'
    let time = util.timeStr(match.elasped)
    return (
      <div styleName='matchInfo'>
        <img src={bg}/>
        <div style={{color:color}}>
          <div styleName='matchTimeLabel'>游戏已开始：</div>
          <div styleName='matchGoldLabel'>当前金币数：</div>
          <div styleName='matchTimeValue'>{time}</div>
          <div styleName='matchGoldValue'>{match.gold + 'G'}</div>
        </div>
      </div>
    )
  }
}), styles)

const CurrentMatchCellView = CSSModules(React.createClass({
  render() {
    let match = this.props.match
    if (match == null) {
      return null
    }
    if (match.mode == 'g') {
      var bgImg = require('./assets/g_cell_bg.png')
      var modeImg = require('./assets/g_icon.png')
      var color = '#dc8524'
    } else {
      var bgImg = require('./assets/s_cell_bg.png')
      var modeImg = require('./assets/s_icon.png')
      var color = '#03dceb'
    }
    let playerStr = match.member.map((player, idx) => {
      return 'P' + player.cid.split(':')[1]
    }).join(' ')
    return (
      <div styleName='matchCell'>
        <img src={bgImg}/>
        <div style={{color:color}}>
          <img src={require('./assets/late0.png')} styleName='historyCellImg' />
          <div styleName='historyNumber'>{match.teamID}</div>
          <div styleName='historyPlayer'>{playerStr}</div>
          <img src={modeImg} styleName='historyIcon' />
          <div styleName='historyResult'>进行中...</div>
        </div>
      </div>
    )
  }
}), styles)

const QueueView = CSSModules(observer(React.createClass({
  render() {
    if (this.props.queue.data == null) {
      return (
        <div styleName='root'>
        <div styleName='container'>
          <img styleName='rootImg' src={require('./assets/qbg.png')} />
        </div>
      </div>
      )
    }
    let history = this.props.queue.data.history
    let queue = this.props.queue.data.queue
    let match = this.props.queue.match
    var count = queue.length
    let preparing = null
    for (let team of queue) {
      if (team.status != 0) {
        count--
        if (team.status == 1) {
          preparing = team
        }
      } else {
        break
      }
    }
    return (
      <div styleName='root'>
        <div styleName='container'>
          <img styleName='rootImg' src={require('./assets/qbg.png')} />
          <div styleName='timeLabel'>最长等待时间：</div>
          <div styleName='groupLabel'>当前排队组数：</div>
          <div styleName='timeValue'>{count * 5}</div>
          <div styleName='groupValue'>{count}</div>
          <div styleName='timeUnit'>分钟</div>
          <div styleName='groupUnit'>组</div>
          <CurrentMatchView match={match} />
          <HistoryView history={history} />
          <CurrentMatchCellView match={match} />
          <PrepareCellView team={preparing} />
        </div>
      </div>
    )
  },
  componentDidMount() {
    this.props.queue.connect()
  }
})), styles)

var queue = new Queue()

render((
    <QueueView queue={queue} />),
  document.getElementById('queue'),
  function() {
    console.log('render queue')
  });
