//
//  DataManager.swift
//  admin
//
//  Created by tassar on 4/25/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import Starscream
import SwiftyJSON

public protocol DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType)
}

public enum DataType: String {
	case QueueData = "updateQueue"

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
extension DataManager: WsClientDelegate {

	public func wsClientDidInit(client: WsClient, data: [String: AnyObject]) {
		for (type, _) in receiversMap {
			WsClient.singleton.sendCmd(type.queryCmd)
		}
	}

	public func wsClientDidDisconnect(client: WsClient, error: NSError?) {
	}

	public func wsClientDidReceiveMessage(client: WsClient, cmd: String, data: [String: AnyObject]) {
		for (type, receivers) in receiversMap {
			if type.rawValue == cmd {
				for receiver in receivers {
					receiver.onReceivedData(data, type: type)
				}
			}
		}
	}
}