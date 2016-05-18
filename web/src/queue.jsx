import React from 'react'
import { render } from 'react-dom'
import { observable, computed } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/queue.css'

class Queue {
  @ observable data
  @ observable state
  constructor() {
    this._reset()
  }

  _reset() {
    this.state = 'not connected'
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
    this.state = 'connecting...'
    sock.onopen = () => {
      console.log('connected to ' + uri)
      let data = {
        cmd: 'init',
        ID: 'queue',
        TYPE: '8',
      }
      sock.send(JSON.stringify(data))
    }
    sock.onclose = (e) => {
      console.log('connection closed (' + e.code + ')')
      this._reset()
    }
    sock.onmessage = (e) => {
      this.onMessage(e.data)
    }
    this.sock = sock
  }

  onMessage(msg) {
    console.log('got socket message: ' + msg)
    let data = JSON.parse(msg)
    switch (data.cmd) {
      case 'init':
        this.state = 'connected'
      case 'queueData':
        this.data = data.data
    }
  }

  send(data) {
    let d = JSON.stringify(data)
    this.log = `发送:${d}\n` + this.log
    this.sock.send(d)
  }

}

const QueueView = CSSModules(observer(React.createClass({
  render() {
    return (
      <div styleName='root'>
        <img src={require('./assets/qbg.png')} />
  </div>
    )
  }
})), styles)

var queue = new Queue()

render((
  <QueueView queue={queue}>
  </QueueView>
), document.getElementById('queue'), function() {
  console.log('render queue')
});
