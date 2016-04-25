//
//  WebSocketClient.swift
//  postgame
//
//  Created by tassar on 4/14/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import Starscream
import SwiftyUserDefaults
import SwiftyJSON

public class WsClient {
	// notification
	public static let WsConnectedNotification = "WsConnected"
	public static let WsDisconnectedNotification = "WsDisconnected"
	public static let WsConnectingNotification = "WsConnecting"

	public static let singleton = WsClient()

	private static let ERROR_WAIT_SECOND: UInt64 = 10
	private var socket: WebSocket?
	private var address: String?

	public func sendCmd(cmd: String) {
		let json = JSON([
			"cmd": cmd
		])
		sendJSON(json)
	}

	public func sendJSON(json: JSON) {
		let str = json.rawString(NSUTF8StringEncoding, options: [])!
		socket!.writeString(str)
	}

	public func connect(addr: String) {
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

	@objc func appDidEnterBackground() {
		if socket != nil && socket!.isConnected {
			socket!.disconnect()
		}
	}

	@objc func appWillEnterForeground() {
		if socket != nil && socket!.isConnected {
			return
		}
		guard address != nil else {
			return
		}
		initSocket()
		doConnect()
	}

	private init() {
		NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(WsClient.appDidEnterBackground), name: UIApplicationDidEnterBackgroundNotification, object: nil)
		NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(WsClient.appWillEnterForeground), name: UIApplicationWillEnterForegroundNotification, object: nil)
	}

	deinit {
		NSNotificationCenter.defaultCenter().removeObserver(self)
	}

	private func initSocket() {
		socket = WebSocket(url: NSURL(string: address!)!)
		socket?.delegate = self
	}

	private func doConnect() {
		socket!.connect()
		NSNotificationCenter.defaultCenter().postNotificationName(WsClient.WsConnectingNotification, object: nil)
	}
}

// MARK: websocket回调方法
extension WsClient: WebSocketDelegate {

	public func websocketDidConnect(socket: WebSocket) {
		log.debug("socket connected")
		NSNotificationCenter.defaultCenter().postNotificationName(WsClient.WsConnectedNotification, object: nil)
		sendCmd("init")
	}

	public func websocketDidReceiveData(socket: WebSocket, data: NSData) {
	}

	public func websocketDidDisconnect(socket: WebSocket, error: NSError?) {
		log.debug("socket disconnected:\(error?.localizedDescription)")
		NSNotificationCenter.defaultCenter().postNotificationName(WsClient.WsDisconnectedNotification, object: nil)
		if UIApplication.sharedApplication().applicationState == .Background {
			return
		}
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

	public func websocketDidReceiveMessage(socket: WebSocket, text: String) {
		log.debug("socket got:\(text)")
	}
}
