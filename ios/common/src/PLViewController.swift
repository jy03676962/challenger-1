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
	private var timeLabel: UILabel!
	override func viewDidLoad() {
		super.viewDidLoad()
		let imageView = UIImageView()
		imageView.image = UIImage(named: "GlobalBackground")
		view.insertSubview(imageView, atIndex: 0)
		imageView <- Edges()
		timeLabel = UILabel()
		timeLabel.font = UIFont(name: PLConstants.usualFont, size: 30)
		timeLabel.textColor = UIColor.whiteColor()
		imageView.addSubview(timeLabel)
		timeLabel <- [
			CenterX(0),
			Top(10)
		]
		NSTimer.scheduledTimerWithTimeInterval(0.5, target: self, selector: #selector(PLViewController.tickTime), userInfo: nil, repeats: true)
	}

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
//		NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(PLViewController.onWsConnected), key: .WsConnected)
//		NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(PLViewController.onWsConnecting), key: .WsConnecting)
//		NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(PLViewController.onWsDisconnected), key: .WsDisconnected)
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		NSNotificationCenter.defaultCenter().removeObserver(self)
	}

	func onWsConnected() {
		timeLabel.textColor = UIColor.whiteColor()
	}

	func onWsConnecting() {
		timeLabel.textColor = UIColor.greenColor()
	}

	func onWsDisconnected() {
		timeLabel.textColor = UIColor.redColor()
	}

	func tickTime() {
		let now = NSDate()
		let fmt = NSDateFormatter()
		fmt.dateFormat = "HH:mm"
		let str = fmt.stringFromDate(now)
		timeLabel.text = "TIME \(str)"
	}
	override func prefersStatusBarHidden() -> Bool {
		return true
	}
}
