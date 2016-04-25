//
//  AppDelegate.swift
//  postgame
//
//  Created by tassar on 3/29/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import XCGLogger

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
		self.window = UIWindow(frame: UIScreen.mainScreen().bounds)
		let vc = LoginViewController()
		let navi = UINavigationController(rootViewController: vc)
		navi.navigationBarHidden = true
		self.window?.rootViewController = navi
		self.window?.makeKeyAndVisible()
		WsClient.singleton.connect(PLConstants.getClientWsAddress())
		return true
	}
}
