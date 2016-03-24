import React from 'react'
import {render} from 'react-dom'
import {observable, computed} from 'mobx'
import {observer} from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/api.css'

class Api {
  constructor() {
    this._reset()
  }

  _reset() {
    this.addr = ""
    this.output = ""
  }
}

const ApiView = CSSModules(observer(React.createClass({
  render() {
    return (
      <div styleName="root">
        <div styleName="block">
          <label>客户端ip</label>
          <input type="text" ref="addr"></input>
        </div>
        <div styleName="block">
          <label>灯带效果</label><br/>
          <label>wall</label>
          <input type="text" ref="wall"></input><br/>
        </div>
      </div>
    )
  }
})), styles)

var api = new Api()

render((
  <ApiView api={api}>
  </ApiView>
), document.getElementById('api'), function(){
  console.log("render api")
});
