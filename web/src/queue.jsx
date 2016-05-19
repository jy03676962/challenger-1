import React from 'react'
import { render } from 'react-dom'
import { observable, computed } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/queue.css'
import { wsAddressWithPath } from '~/js/util.jsx'

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
    let uri = wsAddressWithPath('ws')
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

const CurrentMatchView = CSSModules(observer(React.createClass({
  render() {
    let match = this.props.match
    if (match == null) {
      return null
    }
    let bg = match.mode == 'g' ? require('./assets/g_game_bg.png') : require('./assets/s_game_bg.png')
    let color = match.mode == 'g' ? '#dc8524' : '#03dceb'
    let min = Math.floor(match.elasped / 60)
    let sec = Math.floor(match.elasped - 60 * min)
    let pad = (i) => {
      return (i < 10 ? '0' : '') + i
    }
    let time = pad(min) + ':' + pad(sec)
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
})), styles)

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
    for (let team of queue) {
      if (team.status != 0) {
        count--
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
