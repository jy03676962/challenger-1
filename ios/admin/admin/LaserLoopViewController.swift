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

class LaserLoopViewController: PLViewController {
	@IBOutlet weak var senderTextField: UITextField!
	@IBOutlet weak var startButton: UIButton!
	@IBOutlet weak var nextButton: UIButton!
	@IBOutlet weak var resultTableView: UITableView!

	var senderList: [MainArduinoInfo]?
	var laserIdx: Int = 0
	var senderIdx: Int = 0
	var checking: Bool = false
	var infoList: [LaserInfo] = [LaserInfo]()

	var currentSender: MainArduinoInfo? {
		return senderList?[senderIdx]
	}

	@IBAction func start() {
		checking = !checking
		if checking {
			startButton.setTitle("停止", forState: .Normal)
			nextButton.enabled = false
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
			startButton.setTitle("开始", forState: .Normal)
			nextButton.enabled = true
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
		if (infoList.count > 1) {
			HUD.flash(.LabeledError(title: "有多个数据，无法记录", subtitle: nil), delay: 2)
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
		start()
	}
	@IBAction func done() {
		WsClient.singleton.sendCmd("stopListenLaser")
		DataManager.singleton.unsubscribe(self)
		presentingViewController?.dismissViewControllerAnimated(true, completion: nil)
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		HUD.show(.Progress)
		DataManager.singleton.subscribeData([.LaserInfo], receiver: self)
		Alamofire.request(.GET, PLConstants.getHttpAddress("api/sender_list"))
			.validate()
			.responseArray(completionHandler: { (response: Response<[MainArduinoInfo], NSError>) in
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
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .LaserInfo {
			let info = Mapper<LaserInfo>().map(json)
			if info != nil {
				for (i, inf) in infoList.enumerate() {
					if inf.id == info!.id {
						infoList.removeAtIndex(i)
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
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}
	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return self.infoList.count
	}
	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCellWithIdentifier("LaserResultCell") as! LaserResultCell
		cell.renderData(infoList[indexPath.row])
		return cell
	}
}
