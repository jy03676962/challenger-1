//
//  LaserLoopViewController.swift
//  admin
//
//  Created by tassar on 5/21/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import PKHUD
import Alamofire
import AlamofireObjectMapper
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
fileprivate func >= <T : Comparable>(lhs: T?, rhs: T?) -> Bool {
  switch (lhs, rhs) {
  case let (l?, r?):
    return l >= r
  default:
    return !(lhs < rhs)
  }
}


class LaserLoopViewController: PLViewController {
	@IBOutlet weak var senderTextField: UITextField!
	@IBOutlet weak var startButton: UIButton!
	@IBOutlet weak var resultTableView: UITableView!

	var senderList: [MainArduinoInfo]?
	var laserIdx: Int = 0
	var senderIdx: Int = 0
	var checking: Bool = false
	var infoList: [LaserInfo] = [LaserInfo]()

	var currentSender: MainArduinoInfo? {
		return senderList?[senderIdx]
	}

	@IBAction func manuallyChangeSender(_ sender: UITextField) {
		if let value = sender.text {
			let li = value.components(separatedBy: ":")
			if li.count == 2 {
				let laserID = li[0]
				let idx: Int? = Int(li[1])
				if idx != nil {
					for (i, info) in senderList!.enumerated() {
						if info.id == laserID && idx! < info.laserNum {
							senderIdx = i
							laserIdx = idx!
							return
						}
					}
				}
			}
			HUD.flash(.labeledError(title: "输入ID有错误", subtitle: nil), delay: 1)
		}
	}

	@IBAction func start() {
		checking = !checking
		if checking {
			startButton.setTitle("停止", for: UIControlState())
			let json = JSON([
				"cmd": "laserOn",
				"id": currentSender!.id,
				"num": laserIdx,
			])
			infoList.removeAll()
			resultTableView.reloadData()
			WsClient.singleton.sendJSON(json)
		} else {
			infoList.removeAll()
			resultTableView.reloadData()
			let json = JSON([
				"cmd": "laserOff",
				"id": currentSender!.id,
				"num": laserIdx,
			])
			WsClient.singleton.sendJSON(json)
			startButton.setTitle("开始", for: UIControlState())
		}
	}

	func next() {
		laserIdx += 1
		if laserIdx >= currentSender?.laserNum {
			senderIdx += 1
			laserIdx = 0
		}
		self.fillSenderText()
	}

	@IBAction func record() {
		if (infoList.count == 0 || infoList.count > 1 || infoList[0].err != 0) {
			HUD.flash(.labeledError(title: "有错误，无法记录", subtitle: nil), delay: 1)
			return
		}
		let json = JSON([
			"cmd": "recordLaser",
			"from": currentSender!.id,
			"from_idx": String(laserIdx),
			"to": infoList[0].id,
			"to_idx": String(infoList[0].idx),
		])
		WsClient.singleton.sendJSON(json)
		if (checking) {
			start()
		}
		next()
		// NSThread.sleepForTimeInterval(0.01)
		start()
	}
	@IBAction func done() {
		WsClient.singleton.sendCmd("stopListenLaser")
		DataManager.singleton.unsubscribe(self)
		presentingViewController?.dismiss(animated: true, completion: nil)
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		HUD.show(.progress)
		DataManager.singleton.subscribeData([.LaserInfo], receiver: self)
        request(PLConstants.getHttpAddress("api/sender_list"))
			.validate()
            .responseArray(completionHandler: { (response: DataResponse<[MainArduinoInfo]>) in
				HUD.hide()
				self.senderList = response.result.value
				if self.senderList != nil {
					self.fillSenderText()
				}
		})
	}

	func fillSenderText() {
		self.senderTextField.text = "\(self.currentSender!.id):\(self.laserIdx)"
	}
}

extension LaserLoopViewController: DataReceiver {
	func onReceivedData(_ json: [String: Any], type: DataType) {
		if type == .LaserInfo {
            let info = Mapper<LaserInfo>().map(JSONObject:json)
			if info != nil {
				for (i, inf) in infoList.enumerated() {
					if inf.id == info!.id {
						infoList.remove(at: i)
						break
					}
				}
				infoList.append(info!)
				resultTableView.reloadData()
			}
		}
	}
}

extension LaserLoopViewController: UITableViewDataSource, UITableViewDelegate {
	func numberOfSections(in tableView: UITableView) -> Int {
		return 1
	}
	func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return self.infoList.count
	}
	func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCell(withIdentifier: "LaserResultCell") as! LaserResultCell
		cell.renderData(infoList[indexPath.row])
		return cell
	}
}
