import React from 'react'
import CSSModules from 'react-css-modules'
import { observer } from 'mobx-react'

export function wsAddressWithPath(path) {
	let loc = window.location
	let uri = `ws://${loc.host}/${path}`
	return uri
}

export function timeStr(elasped) {
	let min = Math.floor(elasped / 60)
	let sec = Math.floor(elasped - 60 * min)
	let pad = (i) => {
		return (i < 10 ? '0' : '') + i
	}
	return pad(min) + ':' + pad(sec)
}

export function cssCreate(styles, specs) {
	return CSSModules(React.createClass(specs), styles)
}

export function cssMobxCreate(styles, specs) {
	return CSSModules(observer(React.createClass(specs)), styles)
}
