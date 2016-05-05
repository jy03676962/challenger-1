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
import ObjectMapper

class ConfigController: PLViewController {

	@IBOutlet weak var wrapperView: UIView!
	@IBOutlet weak var hostTextField: UITextField!
	@IBOutlet weak var modeControl: UISegmentedControl!
	@IBOutlet weak var arduinoView: UIView!

	var arduinoViewMap: [String: UILabel] = [String: UILabel]()
	var timer = NSTimer()

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
		DataManager.singleton.subscriptData([.ArduinoList], receiver: self)
		timer = NSTimer.scheduledTimerWithTimeInterval(2, target: self, selector: #selector(queryArduinoList), userInfo: nil, repeats: true)
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		DataManager.singleton.removeSubscript(self)
		timer.invalidate()
	}

	func queryArduinoList() {
		WsClient.singleton.sendCmd(DataType.ArduinoList.queryCmd)
	}

	func renderArduinoList(list: [ArduinoController]) {
		let margin: CGFloat = 10
		let row = 10
		let width: CGFloat = 80
		let height: CGFloat = 20
		let firstRender = self.arduinoViewMap.count == 0
		for (i, controller) in list.enumerate() {
			let label: UILabel
			if firstRender {
				label = UILabel()
				label.text = controller.address.id
				label.font = UIFont.systemFontOfSize(12)
				let top = CGFloat(i / row) * (margin + height) + margin
				let left = CGFloat(i % row) * (margin + width) + margin
//				log.debug("\(i): \(CGFloat(i % row)), top is \(top), left is \(top)")
				label.frame = CGRect(x: left, y: top, width: width, height: height)
				label.textAlignment = .Center
				self.arduinoView.addSubview(label)
				arduinoViewMap[controller.address.id] = label
			} else {
				label = arduinoViewMap[controller.address.id]!
			}
			if (controller.online == true) {
				if controller.mode == .On {
					label.textColor = UIColor.blackColor()
				} else if controller.mode == .Off {
					label.textColor = UIColor.blueColor()
				} else if controller.mode == .Free {
					label.textColor = UIColor.greenColor()
				} else {
					label.textColor = UIColor.orangeColor()
				}
			} else {
				label.textColor = UIColor.redColor()
			}
			label.borderWidth = 0
			if (controller.address.type == .MainArduino) {
				if !controller.scoreUpdated {
					label.borderWidth = 1
					label.borderColor = UIColor.redColor()
				}
			}
		}
	}
}

extension ConfigController: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .ArduinoList {
			let arduinoList = Mapper<ArduinoController>().mapArray(json["data"])
			if arduinoList != nil {
				renderArduinoList(arduinoList!)
			}
		}
	}
}
