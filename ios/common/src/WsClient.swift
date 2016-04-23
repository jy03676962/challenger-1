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

class WsClient {
	static let singleton = WsClient()
	private var socket: WebSocket?
	private var currentAddress: String? {
		didSet {
			if currentAddress == oldValue || currentAddress == nil {
				return
			}
			if socket == nil {
				doConnect()
			} else if socket!.isConnected {
				socket!.disconnect()
			}
		}
	}

	@objc func onHostChanged(notif: NSNotification) {
		currentAddress = Defaults[.host]
	}

	func onConnect() {
		log.debug("socket connected")
		NSNotificationCenter.defaultCenter().postNotificationKey(.WsConnected, object: nil)
	}

	func onDisconnect(error: NSError?) {
		log.debug("socket disconnected:\(error?.localizedDescription)")
		NSNotificationCenter.defaultCenter().postNotificationKey(.WsDisconnected, object: nil)
		if error == nil {
			doConnect()
		} else {
			dispatch_after(dispatch_time(DISPATCH_TIME_NOW, Int64(10 * NSEC_PER_SEC)), dispatch_get_main_queue(), {
				self.doConnect()
			})
		}
	}

	func onText(text: String) {
		log.debug("socket got:\(text)")
	}

	private init() {
	}

	func doConnect() {
		socket = WebSocket(url: NSURL(string: currentAddress!)!)
		socket!.onConnect = {
			self.onConnect()
		}
		socket!.onDisconnect = { error in
			self.onDisconnect(error)
		}
		socket!.onText = { text in
			self.onText(text)
		}
		NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(WsClient.onHostChanged(_:)), key: .HostChanged)
		socket!.connect()
		NSNotificationCenter.defaultCenter().postNotificationKey(.WsConnecting, object: nil)
	}

	func connect(address: String) {
		currentAddress = address
	}
}
