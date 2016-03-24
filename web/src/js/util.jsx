export function wsAddressWithPath(path) {
  let loc = window.location
  let uri = `ws://${loc.host}/${path}`
  return uri
}
