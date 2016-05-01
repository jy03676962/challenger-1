import { observable, computed } from 'mobx'
import { wsAddressWithPath } from '~/js/util.jsx'

class Game {
  @observable match
  @observable playerName
  @observable options

  @computed get stage() {
    if (!this.playerName || this.match == null || this.match.member == null || this.match.member.length == 0) {
      return 'login'
    }
    if (!this.match || this.match.stage == 'before') {
      return 'hall'
    }
    if (this.match && (this.match.stage == 'ongoing' || this.match.stage == 'warmup')) {
      return 'arena'
    }
    if (this.match && this.match.stage == 'after') {
      return 'board'
    }
  }

  constructor() {
    this._reset()
  }

  _reset() {
    this.playerName = ''
    this.sock = null
    this.match = null
    this.arg = null
    this.options = null
    this.currentKey = 0
  }

  connectServer(playerName) {
    let uri = wsAddressWithPath('ws')
    let sock = new WebSocket(uri)
    console.log('socket is ' + uri)
    this.playerName = playerName
    sock.onopen = () => {
      console.log('connected to ' + uri)
      this.login(playerName)
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

  login(playerName) {
    if (this.sock) {
      let data = {
        cmd: 'login',
        name: playerName
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  joinMatch() {
    if (this.sock) {
      let data = {
        cmd: 'joinMatch',
        name: this.playerName
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  createMatch() {
    if (this.sock) {
      let data = {
        cmd: 'createMatch',
        name: this.playerName
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  startMatch(mode) {
    if (this.sock) {
      let data = {
        cmd: 'startMatch',
        mode: mode,
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  resetMatch() {
    if (this.sock) {
      let data = {
        cmd: 'resetMatch'
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  onMessage(msg) {
    let data = JSON.parse(msg)
    switch (data.cmd) {
      case 'init':
        this.options = data
        break
      case 'updateMatch':
        this.match = JSON.parse(data)
        break
    }
  }

  onKeyDown(e) {
    if (this.stage != 'arena') {
      return
    }
    var code = e.keyCode ? e.keyCode : e.which;
    let dir
    switch (code) {
      case 37: //left
        dir = 'left'
        break
      case 38:
        dir = 'up'
        break
      case 39:
        dir = 'right'
        break
      case 40:
        dir = 'down'
        break
    }
    if (dir) {
      this.currentKey = code
      let data = {
        cmd: 'playerMove',
        dir: dir,
        name: this.playerName,
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  onKeyUp(e) {
    if (this.stage != 'arena') {
      return
    }
    var code = e.keyCode ? e.keyCode : e.which;
    if (this.currentKey == code) {
      this.currentKey = 0
      let data = {
        cmd: 'playerStop',
        name: this.playerName,
      }
      this.sock.send(JSON.stringify(data))
    }
  }
}


export default Game
