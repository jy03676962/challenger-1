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
