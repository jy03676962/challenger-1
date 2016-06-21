import React from 'react'
import { render } from 'react-dom'
import { observable } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import Enum from 'es6-enum'

import SetIntervalMixin from '~/js/mixins/SetIntervalMixin.jsx'
import styles from '~/styles/rank.css'

const RankMode = Enum('SurvivalMode', 'GeneralMode')
const ChangeModeInterval = 10

class Rank {
  @ observable survivalModeData
  @ observable generalModeData
  @ observable mode
}

const RankDataView = CSSModules(React.createClass({
  render() {
    let rank = this.props.rank
    let data = rank.mode == RankMode.SurvivalMode ? survivalModeData : generalModeData
    if (data == null || data.error.length != 0) {
      return null
    }
    return (
      <div>
      </div>
    )
  }

}), styles)

const RankView = CSSModules(observer(React.createClass({
  mixins: [SetIntervalMixin],
  componentDidMount: () => {
    this.setInterval(this.changeMode, ChangeModeInterval * 1000)
  },
  changeMode: () => {},
  render() {
    let rank = this.props.rank
    let data = rank.mode == RankMode.SurvivalMode ? survivalModeData : generalModeData
    return (
      <div styleName='root'>
        <div styleName='container'>
          <img styleName='rootImg' src={require('./assets/qbg.png')} />
          <RankDataView {...this.props} />
        </div>
      </div>
    )
  },
})), styles)

var rank = new Rank()

render(<RankView rank={rank} />, document.getElementById('rank'))
