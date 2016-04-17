//
//  WebSocketClient.swift
//  postgame
//
//  Created by tassar on 4/14/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import Starscream

class WsClient {
	static let singleton = WsClient()
	private var socket: WebSocket
	private var currentAddress: String

	@objc func onHostChanged(notif: NSNotification) {
		socket.disconnect()
	}

	func onConnect() {
		log.debug("socket connected")
	}

	func onDisconnect(error: NSError?) {
		log.debug("socket disconnected:\(error?.localizedDescription)")
		if error == nil {
			socket.connect()
		} else {
			dispatch_after(dispatch_time(DISPATCH_TIME_NOW, Int64(10 * NSEC_PER_SEC)), dispatch_get_main_queue(), {
				self.socket.connect()
			})
		}
	}

	func onText(text: String) {
		log.debug("socket got:\(text)")
	}

	private init() {
		currentAddress = PLConstants.getWsAddress()
		socket = WebSocket(url: NSURL(string: currentAddress)!)
		socket.onConnect = {
			self.onConnect()
		}
		socket.onDisconnect = { error in
			self.onDisconnect(error)
		}
		socket.onText = { text in
			self.onText(text)
		}
		NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(WsClient.onHostChanged(_:)), key: .HostChanged)
	}

	func connect() {
		socket.connect()
	}
}
