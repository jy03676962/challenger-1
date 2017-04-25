//
//  PLViewController.swift
//  admin
//
//  Created by tassar on 4/23/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import EasyPeasy

class PLViewController: UIViewController {
	fileprivate var timeLabel: UILabel!
	override func viewDidLoad() {
		super.viewDidLoad()
		let imageView = UIImageView()
		imageView.image = UIImage(named: "GlobalBackground")
		view.insertSubview(imageView, at: 0)
		imageView <- Edges()
		timeLabel = UILabel()
		timeLabel.font = UIFont(name: PLConstants.usualFont, size: 30)
		imageView.addSubview(timeLabel)
		timeLabel <- [
			CenterX(0),
			Top(10)
		]
		Timer.scheduledTimer(timeInterval: 0.5, target: self, selector: #selector(PLViewController.tickTime), userInfo: nil, repeats: true)
	}

	override func viewWillAppear(_ animated: Bool) {
		super.viewWillAppear(animated)
		timeLabel.textColor = WsClient.singleton.didInit ? UIColor.white : UIColor.red
		NotificationCenter.default.addObserver(self, selector: #selector(onWsInited), name: NSNotification.Name(rawValue: WsClient.WsInitedNotification), object: nil)
		NotificationCenter.default.addObserver(self, selector: #selector(onWsConnecting), name: NSNotification.Name(rawValue: WsClient.WsConnectingNotification), object: nil)
        NotificationCenter.default.addObserver(self, selector: #selector(onWsDisconnected), name: NSNotification.Name(rawValue: WsClient.WsDisconnectedNotification), object: nil)
	}

	override func viewDidDisappear(_ animated: Bool) {
		super.viewDidDisappear(animated)
		NotificationCenter.default.removeObserver(self)
	}

	func onWsInited() {
		timeLabel.textColor = UIColor.white
	}

	func onWsConnecting() {
		timeLabel.textColor = UIColor.green
	}

	func onWsDisconnected() {
		timeLabel.textColor = UIColor.red
	}

	func tickTime() {
		let now = Date()
		let fmt = DateFormatter()
		fmt.dateFormat = "HH:mm"
		let str = fmt.string(from: now)
		timeLabel.text = "TIME \(str)"
	}
	override var prefersStatusBarHidden : Bool {
		return true
	}
}
