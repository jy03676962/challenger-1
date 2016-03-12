export function wsAddressWithPath(path) {
  let loc = window.location
  let uri = 'ws:'
  uri += '//' + loc.host
  uri += loc.pathname + path
  return uri
}
