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

public protocol DataReceiver: class {
	func onReceivedData(_ json: [String: Any], type: DataType)
}

public enum DataType: String {
	case HallData = "HallData"
	case ControllerData = "ControllerData"
	case NewMatch = "newMatch"
	case ArduinoList = "ArduinoList"
	case UpdateMatch = "updateMatch"
	case MatchStop = "matchStop"
	case StartAnswer = "startAnswer"
	case StopAnswer = "stopAnswer"
	case UpdatePlayerData = "updatePlayerData"
	case QuestionCount = "QuestionCount"
	case LaserInfo = "laserInfo"
	case QuickCheckInfo = "QuickCheck"
	case Error = "error"

	var queryCmd: String {
		return "query\(self.rawValue)"
	}

	var shouldQuery: Bool {
		let first: String = self.rawValue[0]
		return first.uppercased() == first
	}
}

open class DataManager {

	fileprivate var receiversMap: [DataType: [DataReceiver]] = [:]

	open static let singleton = DataManager()

	open func subscribeData(_ types: [DataType], receiver: DataReceiver) {
		for type in types {
			var list = receiversMap[type] ?? [DataReceiver]()
			if !list.contains(where: { (rcv) -> Bool in rcv === receiver }) {
				list.append(receiver)
				receiversMap[type] = list
			}
			if type.shouldQuery {
				WsClient.singleton.sendCmd(type.queryCmd)
			}
		}
	}

	open func unsubscribe(_ receiver: DataReceiver) {
		for (t, l) in receiversMap {
			var nl = [DataReceiver]()
			for r in l {
				if r !== receiver {
					nl.append(r)
				}
			}
			receiversMap[t] = nl
		}
	}
	open func unsubscribe(_ receiver: DataReceiver, type: DataType) {
		for (t, l) in receiversMap {
			if t != type {
				continue
			}
			var nl = [DataReceiver]()
			for r in l {
				if r !== receiver {
					nl.append(r)
				}
			}
			receiversMap[t] = nl
		}
	}

	open func queryData(_ type: DataType) {
		WsClient.singleton.sendCmd(type.queryCmd)
	}

	fileprivate func dispatch() {
	}

	fileprivate init() {
		WsClient.singleton.delegate = self
	}
}

// MARK: websocket notificaiton
extension DataManager: WsClientDelegate {

	public func wsClientDidInit(_ client: WsClient, data: [String: Any]) {
		for (type, _) in receiversMap {
			WsClient.singleton.sendCmd(type.queryCmd)
		}
	}

	public func wsClientDidDisconnect(_ client: WsClient, error: NSError?) {
	}

	public func wsClientDidReceiveMessage(_ client: WsClient, cmd: String, data: [String: Any]) {
		for (type, receivers) in receiversMap {
			if type.rawValue == cmd {
				for receiver in receivers {
					receiver.onReceivedData(data, type: type)
				}
			}
		}
	}
}
