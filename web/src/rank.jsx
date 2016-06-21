import React from 'react'
import { render } from 'react-dom'
import { observable, computed } from 'mobx'
import { observer } from 'mobx-react'
import CSSModules from 'react-css-modules'
import styles from '~/styles/rank.css'
import * as util from '~/js/util.jsx'
import Enum from 'es6-enum'

const RankMode = Enum('SurvivalMode', 'GeneralMode')
const ChangeModeInterval = 10

class Rank {
  @ observable survivalModeData
  @ observable goldModeData
  @ observable error
  @ observable mode
}

const RankView = CSSModules(observer(React.createClass({
  render() {}
})), styles)

var rank = new Rank()

render(<RankView rank={rank} />, document.getElementById('rank'))
