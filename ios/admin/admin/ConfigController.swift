//
//  ConfigController.swift
//  admin
//
//  Created by tassar on 4/23/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import AutoKeyboardScrollView
import EasyPeasy
import SwiftyUserDefaults
import SwiftyJSON

let arenaMainSize: CGFloat = 60
let arenaSubSize: CGFloat = 15

class ConfigController: PLViewController {

	@IBOutlet weak var wrapperView: UIView!
	@IBOutlet weak var hostTextField: UITextField!
	@IBOutlet weak var modeControl: UISegmentedControl!
	var arenaWidth: Int = 0
	var arenaHeight: Int = 0

	@IBAction func modeChange(sender: UISegmentedControl) {
		WsClient.singleton.sendJSON(JSON([
			"cmd": "arduinoModeChange",
			"mode": sender.selectedSegmentIndex
			]))
	}

	@IBAction func saveConfig() {
		if hostTextField.text != nil && hostTextField.text?.characters.count > 0 {
			Defaults[.host] = hostTextField.text
			WsClient.singleton.connect(PLConstants.getWsAddress())
		}
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		let scrollView = AutoKeyboardScrollView()
		scrollView.backgroundColor = UIColor.clearColor()
		view.addSubview(scrollView)
		wrapperView.removeFromSuperview()
		scrollView.addSubview(wrapperView)
		scrollView <- Edges()
		wrapperView <- Edges()

		modeControl.tintColor = UIColor.whiteColor()
	}

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		DataManager.singleton.subscriptData([.ArenaSize, .ArduinoMode], receiver: self)
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		DataManager.singleton.removeSubscript(self)
	}
}

extension ConfigController: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .ArenaSize {
			arenaWidth = json["data"]!["width"] as! Int
			arenaHeight = json["data"]!["height"] as! Int
		}
	}
}
