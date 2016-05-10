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

let log = XCGLogger.defaultInstance()

@UIApplicationMain
class AppDelegate: UIResponder, UIApplicationDelegate {

	var window: UIWindow?

	func application(application: UIApplication, didFinishLaunchingWithOptions launchOptions: [NSObject: AnyObject]?) -> Bool {
		#if DEBUG
			log.setup(.Debug, showThreadName: true, showLogLevel: true, showFileNames: true, showLineNumbers: true, writeToFile: nil)
		#else
			log.setup(.Severe, showThreadName: true, showLogLevel: true, showFileNames: true, showLineNumbers: true, writeToFile: nil)
		#endif
		UITabBar.appearance().barTintColor = UIColor.clearColor()
		UITabBar.appearance().backgroundImage = UIImage()
		UITabBar.appearance().shadowImage = UIImage()
		Defaults[.host] = "localhost:3000"
		Defaults[.deviceID] = "admin"
		Defaults[.socketType] = "1"
		Defaults[.matchID] = 0
		Defaults[.qCount] = 7
		Defaults[.websiteHost] = "puapi.hualinfor.com"
		WsClient.singleton.connect(PLConstants.getWsAddress())
		DataManager.singleton.subscribeData([.NewMatch, .QuestionCount], receiver: self)
		return true
	}
}

extension AppDelegate: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .NewMatch {
			Defaults[.matchID] = json["data"] as! Int
		} else if type == .QuestionCount {
			Defaults[.qCount] = json["data"] as! Int
		}
	}
}
