import React from 'react'
import { render } from 'react-dom'
import { observable, computed } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/ingame.css'
import Game from '~js/game.jsx'

var game = new Game()

const IngameView = CSSModules(observer(React.createClass({
  render() {}
})), styles)

render(
  (<IngameView game={game} />),
  document.getElementById('ingame'),
  function() {
    console.log('render ingame')
  }
)
