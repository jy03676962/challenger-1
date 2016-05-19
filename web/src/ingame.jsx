import React from 'react'
import { render } from 'react-dom'
import { observable, computed } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/ingame.css'
import { wsAddressWithPath } from '~/js/util.jsx'

class IngameData {

  @ observable match
  @ observable connected

  constructor() {
    this._reset()
  }

  _reset() {
    this.sock = null
    this.match = null
    this.connected = false
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
        ID: 'ingame',
        TYPE: '9',
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
      case 'updateMatch':
        if (json.data != null && this.connected) {
          this.match = JSON.parse(json.data)
        }
        break
      case 'matchStop':
        this.match = null
        break
    }
  }

}

const PlayerInfo = CSSModules(observer(React.createClass({
  render() {
    let player = this.props.player
    let idx = this.props.idx
    let top = 3 + 125 * idx
    let style = {
      position: 'absolute',
      width: '872px',
      height: '122px',
      top: top + 'px',
    }
    if (player == null) {
      return (
        <div style={style}>
          <img styleName='tableImg' src={require('./assets/energy_off.png')} />
          </div>
      )
    } else {
      let name = player.cid.split(':')[1] + 'P'
      return (
        <div style={style}>
        <div styleName='tableName'>{name}</div>
        <img styleName='tableImg' src={require('./assets/energy_on.png')} />
        <div styleName='tableEnergy'>{player.energy}</div>
      </div>
      )
    }
  }
})), styles)

const IngameView = CSSModules(observer(React.createClass({
  render() {
    let data = this.props.data
    if (data.match == null) {
      return (
        <div styleName='root'>
        <img src={require('./assets/ibg.png')} />
        <div styleName='goldX'>X</div>
        <div styleName='tableBg'>
          <img styleName='tableBgImg' src={require('./assets/itb.png')}/>
        </div>
      </div>
      )
    } else {
      var content = []
      let sortedMember = data.match.member.sort((a, b) => {
        return a.cid.localeCompare(b.cid)
      })
      content = []
      for (var i = 0; i < 4; i++) {
        if (i < sortedMember.length) {
          content.push(<PlayerInfo player={sortedMember[i]} idx={i} key={i} />)
        } else {
          content.push(<PlayerInfo idx={i} key ={i} />)
        }
      }
      let min = Math.floor(data.match.elasped / 60)
      let sec = Math.floor(data.match.elasped - 60 * min)
      let pad = (i) => {
        return (i < 10 ? '0' : '') + i
      }
      let time = pad(min) + ':' + pad(sec)
      let showGold = data.match.mode == 'g' && data.match.gold > 0
      var barBg, barFront
      if (data.match.mode == 'g') {
        barBg = require('./assets/g_b.png')
        barFront = data.match.rampageTime > 0 ? require('./assets/g_r.png') : require('./assets/g_n.png')
      } else {
        barBg = require('./assets/s_b.png')
        barFront = data.match.rampageTime > 0 ? require('./assets/s_r.png') : require('./assets/s_n.png')
      }
      let r = (1 - data.match.energy / data.match.maxEnergy) * 1911 / 19.2
      let clip = `inset(0 ${r}vw 0 0)`
      let s = {
        position: 'absolute',
        width: '100%',
        left: '0',
        top: '0',
        '-webkit-clip-path': clip,
      }
      return (
        <div styleName='root'>
        <img src={require('./assets/ibg.png')} />
        {showGold ? null : <div styleName='goldX'>X</div>}
        {showGold ? <div styleName='goldValue'>{data.match.gold + 'G'}</div> : null}
        <div styleName='timeValue'>{time}</div>
        <div styleName='tableBg'>
          <img styleName='tableBgImg' src={require('./assets/itb.png')}/>
          {content}
        </div>
        <div styleName='bar'>
          <img src={barBg} />
          <img src={barFront} style={s} />
        </div>
      </div>
      )
    }
  },
  componentDidMount() {
    this.props.data.connect()
  }
})), styles)

var d = new IngameData()

render(
  (<IngameView data={d} />),
  document.getElementById('ingame'),
  function() {
    console.log('render ingame')
  }
)