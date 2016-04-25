//
//  DataManager.swift
//  admin
//
//  Created by tassar on 4/25/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import SwiftyJSON
import Starscream

public protocol DataReceiver {
	func onReceivedData(json: JSON, type: DataType)
}

public enum DataType: String {
	case HallData = "HallData"

	var queryCmd: String {
		return "query\(self.rawValue)"
	}
}

public class DataManager {

	private var receiversMap: [DataType: [DataReceiver]] = [:]

	public static let singleton = DataManager()

	public func subscriptData(types: [DataType], receiver: DataReceiver) {
		for type in types {
			if var list = receiversMap[type] {
				list.append(receiver)
			} else {
				var list = [DataReceiver]()
				list.append(receiver)
				receiversMap[type] = list
			}
			WsClient.singleton.sendCmd(type.queryCmd)
		}
	}

	public func queryData(type: DataType) {
		WsClient.singleton.sendCmd(type.queryCmd)
	}

	private func dispatch() {
	}

	private init() {
		WsClient.singleton.delegate = self
	}
}

// MARK: websocket notificaiton
extension DataManager: WebSocketDelegate {
	public func websocketDidConnect(socket: WebSocket) {
		for (type, _) in receiversMap {
			WsClient.singleton.sendCmd(type.queryCmd)
		}
	}

	public func websocketDidReceiveData(socket: WebSocket, data: NSData) {
	}
	public func websocketDidDisconnect(socket: WebSocket, error: NSError?) {
	}
	public func websocketDidReceiveMessage(socket: WebSocket, text: String) {
		let dataFromString = text.dataUsingEncoding(NSUTF8StringEncoding, allowLossyConversion: false)
		guard dataFromString != nil else {
			return
		}
		let json = JSON(data: dataFromString!)
		let cmd = json["cmd"].string
		guard cmd != nil else {
			return
		}
		for (type, receivers) in receiversMap {
			if type.rawValue == cmd {
				for receiver in receivers {
					receiver.onReceivedData(json, type: type)
				}
			}
		}
	}
}