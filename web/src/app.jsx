import React from 'react'
import {render} from 'react-dom'
import Login from '~/js/login.jsx'
import Hall from '~/js/hall.jsx'
import Arena from '~/js/arena.jsx'
import Game from '~/js/game.jsx'
import {observer} from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/base.css'

const App = CSSModules(observer(React.createClass({
  render() {
    var element
    const game = this.props.game
    switch (game.stage) {
      case 'login':
      element = <Login game={game} />
      break
      case 'hall':
      element = <Hall game={game} />
      break
      case 'arena':
      element = <Arena game={game} />
      break
    }
    return (
      <div id='app' styleName='base-div'>{element}</div>
    )
  }
})), styles);

var game = new Game()

document.onkeydown = function(e) {
  game.onKeyDown(e)
};

document.onkeyup = function(e) {
  game.onKeyUp(e)
}

render((
  <App game={game}>
  </App>
), document.getElementById('root'));

