import React from 'react'
import { render } from 'react-dom'
import { observable, computed } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/api.css'
import { wsAddressWithPath } from '~/js/util.jsx'

class Api {
  @ observable state
  @ observable addr
  @ observable log
  constructor() {
    this._reset()
  }

  _reset() {
    this.addr = ''
    this.output = ''
    this.state = 'not connected'
    this.sock = null
    this.log = ''
  }
  connect() {
    if (this.sock) {
      this.sock.close()
      this._reset()
      return
    }
    let uri = wsAddressWithPath('api')
    let sock = new WebSocket(uri)
    console.log('socket is ' + uri)
    this.state = 'connecting...'
    sock.onopen = () => {
      console.log('connected to ' + uri)
      this.state = 'connected'
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
      case 'addTCP':
        this.addr = data.addr
        this.log = `新连接: ${this.addr}\n` + this.log
        break
      case 'delTCP':
        this._reset()
        break
      case 'errTCP':
        this.log = `错误: ${data.msg}\n` + this.log
        break
      default:
        this.addr = data.addr
        this.log = `收到:${msg}\n` + this.log
    }
  }
  send(data) {
    let d = JSON.stringify(data)
    this.log = `发送:${d}\n` + this.log
    this.sock.send(d)
  }
}

const ApiView = CSSModules(observer(React.createClass({
  render() {
    let c = this.props.api.state == 'connected' ? '断开' : '连接'
    return (
      <div styleName='root'>
        <div styleName='left'>
          <div styleName='block'>
            <label styleName='content'>地址</label>
            <label>{this.props.api.addr}</label>
            <label styleName='content'>状态</label>
            <label>{this.props.api.state}</label>
            <button onClick={this.connect}>{c}</button>
          </div>
          <div styleName='block'>
            <label styleName='title'>灯带效果</label><br/>
            <label styleName='content'>wall</label>
            <input type='text' ref='wall'></input><br/>
            <label styleName='content'>led_t</label>
            <input type='text' ref='led_t'></input><br/>
            <label styleName='content'>mode</label>
            <input type='text' ref='mm'></input><br/>
            <button onClick={this.ledCtrl}>发送</button>
          </div>
          <div styleName='block'>
            <label styleName='title'>激光控制</label><br/>
            <label styleName='content'>laser_n</label>
            <input type='text' ref='laser_n'></input><br/>
            <label styleName='content'>laser_s</label>
            <input type='text' ref='laser_s'></input><br/>
            <button onClick={this.laserCtrl}>发送</button>
          </div>
          <div styleName='block'>
            <label styleName='title'>按键</label>
            <input type='checkbox' ref='btn'/>可用<br/>
            <button onClick={this.btnCtrl}>发送</button>
          </div>
          <div styleName='block'>
            <label styleName='title'>播放音乐</label>
            <input type='text' ref='music'></input><br/>
            <button onClick={this.musicCtrl}>发送</button>
          </div>
          <div styleName='block'>
            <label styleName='title'>指示灯</label>
            <input type='text' ref='light'></input><br/>
            <button onClick={this.lightCtrl}>发送</button>
          </div>
          <div styleName='block'>
            <label styleName='title'>模式选择</label>
            <input type='text' ref='mode'></input><br/>
            <button onClick={this.modeCtrl}>发送</button>
          </div>
          <div styleName='block'>
            <label styleName='title'>分数设置</label><br/>
            <label styleName='content'>T1</label>
            <input type='text' ref='t1'></input><br/>
            <label styleName='content'>T2</label>
            <input type='text' ref='t2'></input><br/>
            <label styleName='content'>T3</label>
            <input type='text' ref='t3'></input><br/>
            <label styleName='content'>暴走</label>
            <input type='text' ref='tr'></input><br/>
            <button onClick={this.scoreCtrl}>发送</button>
          </div>
        </div>
        <div styleName='right'>
          <textarea ref='console' readOnly value={this.props.api.log}></textarea>
        </div>
      </div>
    )
  },
  connect: function(e) {
    this.props.api.connect()

  },
  ledCtrl: function(e) {
    let d = {
      cmd: 'led_ctrl',
      led: [{
        wall: this.refs.wall.value,
        led_t: this.refs.led_t.value,
        mode: this.refs.mm.value
      }]
    }
    this.props.api.send(d)
  },
  btnCtrl: function(e) {
    let d = {
      cmd: 'btn_ctrl',
      useful: this.refs.btn.checked ? '1' : '0'
    }
    this.props.api.send(d)
  },
  musicCtrl: function(e) {
    let d = {
      cmd: 'mp3_ctrl',
      music: this.refs.music.value
    }
    this.props.api.send(d)
  },
  lightCtrl: function(e) {
    let d = {
      cmd: 'light_ctrl',
      light_mode: this.refs.light.value
    }
    this.props.api.send(d)
  },
  modeCtrl: function(e) {
    let d = {
      cmd: 'mode_change',
      mode: this.refs.mode.value
    }
    this.props.api.send(d)
  },
  scoreCtrl: function(e) {
    let d = {
      cmd: 'init_score',
      score: [{
        status: 'T1',
        'time': this.refs.t1.value
      }, {
        status: 'T2',
        'time': this.refs.t2.value
      }, {
        status: 'T3',
        'time': this.refs.t3.value
      }, {
        status: 'TR',
        'time': this.refs.tr.value
      }, ]
    }
    this.props.api.send(d)
  },
  laserCtrl: function(e) {
    let d = {
      cmd: 'laser_ctrl',
      laser: [{
        laser_n: this.refs.laser_n.value,
        laser_s: this.refs.laser_s.value,
      }]
    }
    this.props.api.send(d)
  }
})), styles)

var api = new Api()

render((
  <ApiView api={api}>
  </ApiView>
), document.getElementById('api'), function() {
  console.log('render api')
});