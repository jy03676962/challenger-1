//
//  AppDelegate.swift
//  postgame
//
//  Created by tassar on 3/29/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import XCGLogger
import SwiftyUserDefaults
import ObjectMapper
import PKHUD

let log = XCGLogger.defaultInstance()

@UIApplicationMain
class AppDelegate: UIResponder, UIApplicationDelegate {

	var window: UIWindow?
	var navi: UINavigationController? {
		return window?.rootViewController as? UINavigationController
	}

	func application(application: UIApplication, didFinishLaunchingWithOptions launchOptions: [NSObject: AnyObject]?) -> Bool {
		#if DEBUG
			log.setup(.Debug, showThreadName: true, showLogLevel: true, showFileNames: true, showLineNumbers: true, writeToFile: nil)
		#else
			log.setup(.Severe, showThreadName: true, showLogLevel: true, showFileNames: true, showLineNumbers: true, writeToFile: nil)
		#endif
		Defaults[.host] = "192.168.1.5:3000"
		Defaults[.deviceID] = "2"
		Defaults[.socketType] = "4"
		Defaults[.matchID] = 0
		Defaults[.websiteHost] = "puapi.hualinfor.com"
		DataManager.singleton.subscribeData([.StopAnswer], receiver: self)
		WsClient.singleton.connect(PLConstants.getWsAddress())
		return true
	}
}

extension AppDelegate: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .StopAnswer {
			guard navi?.visibleViewController as? LoginViewController == nil else {
				return
			}
			HUD.hide()
			let sb = UIStoryboard(name: "Main", bundle: nil)
			let login = sb.instantiateViewControllerWithIdentifier("LoginViewController")
			navi?.setViewControllers([login], animated: true)
		}
	}
}
