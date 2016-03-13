import {observable, computed} from 'mobx'
import {wsAddressWithPath} from '~/js/util.jsx'

class Game {
  @observable match
  @observable playerName
  @observable options

  @computed get stage() {
    if (!this.playerName) {
      return "login"
    }
    if (!this.match || this.match.stage == "before") {
      return "hall"
    }
    if (this.match && (this.match.stage == "ongoing" || this.match.stage == "warmup")) {
      return "arena"
    }
  }

  constructor() {
    this._reset()
  }

  _reset() {
    this.playerName = ""
    this.sock = null
    this.match = null
    this.arg = null
    this.options = null
  }

  connectServer(playerName) {
    let uri = wsAddressWithPath('ws')
    let sock = new WebSocket(uri)
    console.log('socket is '+uri)
    this.playerName = playerName
    sock.onopen = () => {
      console.log("connected to " + uri)
      this.login(playerName)
    }
    sock.onclose = (e) => {
      console.log("connection closed (" + e.code + ")")
      this._reset()
    }
    sock.onmessage = (e) => {
      console.log("message received: " + e.data)
      this.onMessage(e.data)
    }
    this.sock = sock
  }

  login(playerName) {
    if (this.sock) {
      let data = {
        cmd: "login",
        name: playerName
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  joinMatch() {
    if (this.sock) {
      let data = {
        cmd: "joinMatch",
        name: this.playerName
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  createMatch() {
    if (this.sock) {
      let data = {
        cmd: "createMatch",
        name: this.playerName
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  startMatch(mode) {
    if (this.sock) {
      let data = {
        cmd: "startMatch",
        mode: mode,
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  onMessage(msg) {
    let data = JSON.parse(msg)
    switch (data.cmd) {
      case "login":
      if ("match" in data) {
        this.match = data.match
      }
      break
      case "matchChanged":
      this.match = data.match
      this.options = data.options
      break
      case "matchTick":
      this.match = data.match
      break
    }
  }
}


export default Game
