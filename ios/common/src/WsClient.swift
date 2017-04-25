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

public protocol WsClientDelegate: class {
	func wsClientDidInit(_ client: WsClient, data: [String: Any])
	func wsClientDidReceiveMessage(_ client: WsClient, cmd: String, data: [String: Any])
	func wsClientDidDisconnect(_ client: WsClient, error: NSError?)
}

open class WsClient {
	// notification
	open static let WsInitedNotification = "WsInited"
	open static let WsDisconnectedNotification = "WsDisconnected"
	open static let WsConnectingNotification = "WsConnecting"

	open static let singleton = WsClient()
	open weak var delegate: WsClientDelegate?
	open var didInit: Bool = false

	fileprivate static let ERROR_WAIT_SECOND: UInt64 = 10
	fileprivate var socket: WebSocket?
	fileprivate var address: String?

	open func sendCmd(_ cmd: String) {
		if !didInit {
			return
		}
		let json = JSON([
			"cmd": cmd
		])
		sendJSON(json)
	}

	open func sendJSON(_ json: JSON) {
		let str = json.rawString(String.Encoding.utf8, options: [])!
        socket!.write(string: str)
	}

	open func connect(_ addr: String) {
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

	fileprivate init() {
		NotificationCenter.default.addObserver(self, selector: #selector(WsClient.appDidEnterBackground), name: NSNotification.Name.UIApplicationDidEnterBackground, object: nil)
		NotificationCenter.default.addObserver(self, selector: #selector(WsClient.appWillEnterForeground), name: NSNotification.Name.UIApplicationWillEnterForeground, object: nil)
	}

	deinit {
		NotificationCenter.default.removeObserver(self)
	}

	fileprivate func initSocket() {
		socket = WebSocket(url: URL(string: address!)!)
		socket?.delegate = self
	}

	fileprivate func doConnect() {
		socket!.connect()
		NotificationCenter.default.post(name: Notification.Name(rawValue: WsClient.WsConnectingNotification), object: nil)
	}
}

// MARK: websocket回调方法
extension WsClient: WebSocketDelegate {

	public func websocketDidConnect(socket: WebSocket) {
		log.debug("socket connected")
		let json = JSON([
			"cmd": "init",
			"ID": Defaults[.deviceID],
			"TYPE": Defaults[.socketType],
		])
		self.sendJSON(json)
	}

	public func websocketDidReceiveData(socket: WebSocket, data: Data) {
	}

	public func websocketDidDisconnect(socket: WebSocket, error: NSError?) {
		self.didInit = false
		delegate?.wsClientDidDisconnect(self, error: error)
		NotificationCenter.default.post(name: Notification.Name(rawValue: WsClient.WsDisconnectedNotification), object: nil)
		if UIApplication.shared.applicationState == .background {
			return
		}
		if socket.currentURL.absoluteString != address {
			initSocket()
		}
		if error == nil {
			doConnect()
		} else {
			DispatchQueue.main.asyncAfter(deadline: DispatchTime.now() + Double(Int64(WsClient.ERROR_WAIT_SECOND * NSEC_PER_SEC)) / Double(NSEC_PER_SEC), execute: {
				self.doConnect()
			})
		}
	}

	public func websocketDidReceiveMessage(socket: WebSocket, text: String) {
		log.debug("socket got:\(text)")
		let dataFromString = text.data(using: String.Encoding.utf8, allowLossyConversion: false)
		guard dataFromString != nil else {
			return
		}
		let json = JSON(data: dataFromString!)
        guard json.type == .dictionary else {
            return
        }
		let cmd = json["cmd"].string
		guard cmd != nil else {
			return
		}
		if cmd == "init" {
			self.didInit = true
			NotificationCenter.default.post(name: Notification.Name(rawValue: WsClient.WsInitedNotification), object: nil)
			self.delegate?.wsClientDidInit(self, data: json.dictionaryObject!)
		} else if self.didInit {
			self.delegate?.wsClientDidReceiveMessage(self, cmd: cmd!, data: json.dictionaryObject!)
		}
	}
}
