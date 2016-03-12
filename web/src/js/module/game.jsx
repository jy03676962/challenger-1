import {observable} from 'mobx'
import {wsAddressWithPath} from '~/js/util.jsx'

class Game {
  @observable stage
  @observable room

  constructor() {
    this._reset()
  }

  _reset() {
    this.playerName = ""
    this.stage = "login"
    this.sock = null
    this.room = null
    this.arg = null
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

  joinRoom() {
    if (this.sock) {
      let data = {
        cmd: "joinRoom",
        name: this.playerName
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  createRoom() {
    if (this.sock) {
      let data = {
        cmd: "createRoom",
        name: this.playerName
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  startGame() {
    if (this.sock) {
      let data = {
        cmd: "startGame",
      }
      this.sock.send(JSON.stringify(data))
    }
  }

  onMessage(msg) {
    let data = JSON.parse(msg)
    switch (data.cmd) {
      case "login":
      if ("room" in data) {
        this.room = data.room
      }
      this.stage = "hall"
      break
      case "roomChanged":
      if ("room" in data) {
        this.room = data.room
        this.arg = data.arg
      }
      break
      case "gameStarted":
      this.stage = "arena"
      break
    }
  }
}


export default Game
