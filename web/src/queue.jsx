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
  constructor() {
    this._reset()
  }

  _reset() {
    this.sock = null
    this.state = false
    this.data = null
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
        console.log(json.data)
        if (json.data != null && this.connected) {
          this.data = json.data
        }
        break
    }
  }

  send(data) {
    if (this.sock) {
      let d = JSON.stringify(data)
      this.sock.send(d)
    }
  }

}

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
