//
//  AppDelegate.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import XCGLogger
import SwiftyUserDefaults

let log = XCGLogger.default

@UIApplicationMain
class AppDelegate: UIResponder, UIApplicationDelegate {

	var window: UIWindow?

	func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplicationLaunchOptionsKey: Any]?) -> Bool {
		Instabug.start(withToken: "c22cb54ceb917a9b8c01418b3ce3b9a5", invocationEvent: IBGInvocationEvent.none)
		#if DEBUG
			log.setup(level: .debug, showThreadName: true, showLevel: true, showFileNames: true, showLineNumbers: true, writeToFile: nil)
		#else
			log.setup(.severe, showThreadName: true, showLogLevel: true, showFileNames: true, showLineNumbers: true, writeToFile: nil)
		#endif
		UITabBar.appearance().barTintColor = UIColor.clear
		UITabBar.appearance().backgroundImage = UIImage()
		UITabBar.appearance().shadowImage = UIImage()
		if Defaults[.host] == "" {
			Defaults[.host] = "192.168.1.5:3000"
		}
		if Defaults[.deviceID] == "" {
			Defaults[.deviceID] = "admin"
		}
		if Defaults[.websiteHost] == "" {
			Defaults[.websiteHost] = PLConstants.defaultWebsiteHost
		}
		Defaults[.socketType] = "1"
		Defaults[.matchID] = 0
		Defaults[.qCount] = 7
		WsClient.singleton.connect(PLConstants.getWsAddress())
		DataManager.singleton.subscribeData([.NewMatch, .QuestionCount], receiver: self)
		return true
	}
}

extension AppDelegate: DataReceiver {
	func onReceivedData(_ json: [String: Any], type: DataType) {
		if type == .NewMatch {
			Defaults[.matchID] = json["data"] as! Int
		} else if type == .QuestionCount {
			Defaults[.qCount] = json["data"] as! Int
		}
	}
}
