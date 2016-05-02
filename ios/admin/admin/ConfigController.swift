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

class ConfigController: PLViewController {

	@IBOutlet weak var wrapperView: UIView!
	@IBOutlet weak var hostTextField: UITextField!
	@IBOutlet weak var modeControl: UISegmentedControl!

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
		DataManager.singleton.subscriptData([.ArduinoModeChange], receiver: self)
	}
}

extension ConfigController: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .ArduinoModeChange {
			let mode = json["data"] as! Int
			modeControl.selectedSegmentIndex = mode
		}
	}
}
