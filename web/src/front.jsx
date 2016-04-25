import React from 'react'
import { render } from 'react-dom'
import { observable, computed } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/front.css'


class Front {
  @observable number
}

const FrontView = CSSModules(observer(React.createClass({
  render() {
    let number = this.props.front.number
    return (
      <div styleName='root'>
        <input type='text' ref='count'></input>
        <button onClick={this.addTeam}>取号</button><br/><br/>
        <button onClick={this.resetQueue}>重置</button><br/><br/>
        <label>当前号码</label><br/>
        {
          number ? <label ref='currentNumber'>{number}</label> : null
        }
      </div>
    )
  },
  addTeam: function(e) {
    let front = this.props.front
    let param = {
      count: this.refs.count.value
    }
    $.post('/api/addteam', param, function(data) {
      if (data) {
        front.number = data.num
      }
    })
  },
  resetQueue: function(e) {
    var r = window.confirm('确定要重置吗？')
    if (r == true) {
      let front = this.props.front
      $.post('/api/resetqueue', function(data) {
        front.number = null
      })
    }
  }
})), styles)

var z = require('npm-zepto')
var front = new Front()


render((<FrontView front={front}/>), document.getElementById('front'),
  function() { console.log('render front') })