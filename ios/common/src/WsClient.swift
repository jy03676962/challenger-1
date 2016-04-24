//
//  WebSocketClient.swift
//  postgame
//
//  Created by tassar on 4/14/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import Starscream
import SwiftyUserDefaults
import SwiftyJSON

class WsClient {
	static let singleton = WsClient()

	private static let ERROR_WAIT_SECOND: UInt64 = 10
	private var socket: WebSocket?
	private var address: String?

	func sendCmd(cmd: String) {
		let json = JSON([
			"cmd": cmd
		])
		sendJSON(json)
	}

	func sendJSON(json: JSON) {
		let str = json.rawString(NSUTF8StringEncoding, options: [])!
		socket!.writeString(str)
	}

	func connect(addr: String) {
		if address == addr && socket != nil && socket!.isConnected {
			return
		}
		address = addr
		if socket == nil {
			initSocket()
			doConnect()
		} else if socket!.isConnected {
			socket!.disconnect()
		}
	}

	private init() {
	}

	private func initSocket() {
		socket = WebSocket(url: NSURL(string: address!)!)
		socket?.delegate = self
	}

	private func doConnect() {
		socket!.connect()
		NSNotificationCenter.defaultCenter().postNotificationKey(.WsConnecting, object: nil)
	}
}

// MARK: websocket回调方法
extension WsClient: WebSocketDelegate {

	func websocketDidConnect(socket: WebSocket) {
		log.debug("socket connected")
		NSNotificationCenter.defaultCenter().postNotificationKey(.WsConnected, object: nil)
		sendCmd("init")
	}

	func websocketDidReceiveData(socket: WebSocket, data: NSData) {
	}

	func websocketDidDisconnect(socket: WebSocket, error: NSError?) {
		log.debug("socket disconnected:\(error?.localizedDescription)")
		NSNotificationCenter.defaultCenter().postNotificationKey(.WsDisconnected, object: nil)
		if socket.currentURL.absoluteString != address {
			initSocket()
		}
		if error == nil {
			doConnect()
		} else {
			dispatch_after(dispatch_time(DISPATCH_TIME_NOW, Int64(WsClient.ERROR_WAIT_SECOND * NSEC_PER_SEC)), dispatch_get_main_queue(), {
				self.doConnect()
			})
		}
	}

	func websocketDidReceiveMessage(socket: WebSocket, text: String) {
		log.debug("socket got:\(text)")
	}
}
