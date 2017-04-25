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
// FIXME: comparison operators with optionals were removed from the Swift Standard Libary.
// Consider refactoring the code to use the non-optional operators.
fileprivate func < <T : Comparable>(lhs: T?, rhs: T?) -> Bool {
  switch (lhs, rhs) {
  case let (l?, r?):
    return l < r
  case (nil, _?):
    return true
  default:
    return false
  }
}

// FIXME: comparison operators with optionals were removed from the Swift Standard Libary.
// Consider refactoring the code to use the non-optional operators.
fileprivate func > <T : Comparable>(lhs: T?, rhs: T?) -> Bool {
  switch (lhs, rhs) {
  case let (l?, r?):
    return l > r
  default:
    return rhs < lhs
  }
}


class ConfigController: PLViewController {

	@IBOutlet weak var wrapperView: UIView!
	@IBOutlet weak var idTextField: UITextField!
	@IBOutlet weak var hostTextField: UITextField!
	@IBOutlet weak var modeControl: UISegmentedControl!
	@IBOutlet weak var arduinoView: UIView!
	@IBOutlet weak var webHostTextField: UITextField!

	var arduinoViewMap: [String: UILabel] = [String: UILabel]()
	var timer = Timer()

	@IBAction func modeChange(_ sender: UISegmentedControl) {
		WsClient.singleton.sendJSON(JSON([
			"cmd": "arduinoModeChange",
			"mode": sender.selectedSegmentIndex
			]))
	}
	@IBAction func saveID() {
		if idTextField.text != nil && idTextField.text?.characters.count > 0 {
			Defaults[.deviceID] = idTextField.text!
			WsClient.singleton.connect(PLConstants.getWsAddress())
		}
	}

	@IBAction func saveConfig() {
		if hostTextField.text != nil && hostTextField.text?.characters.count > 0 {
			Defaults[.host] = hostTextField.text!
			WsClient.singleton.connect(PLConstants.getWsAddress())
		}
	}

	@IBAction func saveWebHost() {
		if webHostTextField.text != nil && webHostTextField.text?.characters.count > 0 {
			Defaults[.websiteHost] = webHostTextField.text!
		}
	}

	@IBAction func showLaserMenu(_ sender: UITapGestureRecognizer) {
		let alert = UIAlertController(title: "激光控制", message: nil, preferredStyle: .actionSheet)
		alert.addAction(UIAlertAction(title: "调试激光", style: .default, handler: { (action) in
			self.performSegue(withIdentifier: "ShowLaserLoop", sender: nil)
			}));
		alert.addAction(UIAlertAction(title: "检测激光", style: .default, handler: { (action) in
			self.performSegue(withIdentifier: "ShowQuickCheck", sender: nil)
			}));
		alert.addAction(UIAlertAction(title: "取消", style: .cancel, handler: nil));
		alert.popoverPresentationController?.sourceView = sender.view;
		present(alert, animated: true, completion: nil)
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		idTextField.placeholder = Defaults[.deviceID]
		hostTextField.placeholder = Defaults[.host]
		webHostTextField.placeholder = Defaults[.websiteHost]
		let scrollView = AutoKeyboardScrollView()
		scrollView.backgroundColor = UIColor.clear
		view.addSubview(scrollView)
		wrapperView.removeFromSuperview()
		scrollView.addSubview(wrapperView)
		scrollView <- Edges()
		wrapperView <- Edges()

		modeControl.tintColor = UIColor.white
	}

	override func viewWillAppear(_ animated: Bool) {
		super.viewWillAppear(animated)
		DataManager.singleton.subscribeData([.ArduinoList], receiver: self)
		timer = Timer.scheduledTimer(timeInterval: 2, target: self, selector: #selector(queryArduinoList), userInfo: nil, repeats: true)
	}

	override func viewDidDisappear(_ animated: Bool) {
		super.viewDidDisappear(animated)
		DataManager.singleton.unsubscribe(self)
		timer.invalidate()
	}

	override func onWsDisconnected() {
		super.onWsDisconnected()
		if self.presentedViewController != nil {
			self.dismiss(animated: false, completion: nil)
		}
		for (_, label) in self.arduinoViewMap {
			label.textColor = UIColor.red
		}
		modeControl.selectedSegmentIndex = 0
	}

	func queryArduinoList() {
		WsClient.singleton.sendCmd(DataType.ArduinoList.queryCmd)
	}

	func renderArduinoList(_ list: [ArduinoController]) {
		let margin: CGFloat = 10
		let row = 8
		let width: CGFloat = 100
		let height: CGFloat = 20
		let firstRender = self.arduinoViewMap.count == 0
		for (i, controller) in list.enumerated() {
			let label: UILabel
			if firstRender {
				label = UILabel()
				label.text = controller.address.id
				label.font = UIFont.systemFont(ofSize: 12)
				let top = CGFloat(i / row) * (margin + height) + margin
				let left = CGFloat(i % row) * (margin + width) + margin
				label.frame = CGRect(x: left, y: top, width: width, height: height)
				label.textAlignment = .center
				self.arduinoView.addSubview(label)
				arduinoViewMap[controller.address.id] = label
			} else {
				label = arduinoViewMap[controller.address.id]!
			}
			if (controller.online == true) {
				if controller.mode == .on {
					label.textColor = UIColor.blue
				} else if controller.mode == .off {
					label.textColor = UIColor.black
				} else if controller.mode == .free {
					label.textColor = UIColor.green
				} else if controller.mode == .scan {
					label.textColor = UIColor.orange
				} else {
					label.textColor = UIColor.purple
				}
			} else {
				label.textColor = UIColor.red
			}
			label.borderWidth = 0
			if (controller.address.type == .mainArduino) {
				if !controller.scoreUpdated {
					label.borderWidth = 1
					label.borderColor = UIColor.red
				}
			}
		}
	}
}

extension ConfigController: DataReceiver {
	func onReceivedData(_ json: [String: Any], type: DataType) {
		if type == .ArduinoList {
            let arduinoList = Mapper<ArduinoController>().mapArray(JSONObject:json["data"])
			if arduinoList != nil {
				renderArduinoList(arduinoList!)
			}
		}
	}
}
