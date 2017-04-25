//
//  QuickCheckViewController.swift
//  admin
//
//  Created by tassar on 6/23/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import SwiftyJSON

class QuickCheckViewController: UIViewController {
	static let CellReuseIdentifier = "CellIdentifier"
	@IBOutlet weak var unknowButton: UIButton!
	@IBOutlet weak var disableButton: UIButton!
	@IBOutlet weak var disableOnButton: UIButton!
	@IBOutlet weak var enableOffButton: UIButton!
	@IBOutlet weak var normalButton: UIButton!
	@IBOutlet weak var tableView: UITableView!

	var data: [[String]] = [[], [], [], [], []]
	var index = 0
	var timer = Timer()

	@IBAction func onButtonClicked(_ sender: UIButton) {
		switch sender {
		case unknowButton:
			index = 0
		case disableButton:
			index = 1
		case disableOnButton:
			index = 2
		case enableOffButton:
			index = 3
		case normalButton:
			index = 4
		default:
			index = 0
		}
		self.tableView.reloadData()
	}

	@IBAction func done() {
		let alert = UIAlertController(title: "是否保存检测结果", message: nil, preferredStyle: .alert)
		let cancelAction = UIAlertAction(title: "取消", style: .cancel) { action in
			self.finish(0)
		}
		alert.addAction(cancelAction)
		let doneAction = UIAlertAction(title: "确定", style: .default) { (action) in
			self.finish(1)
		}
		alert.addAction(doneAction)
		present(alert, animated: true, completion: nil)
	}

	func finish(_ save: Int) {
		let json = JSON([
			"cmd": "stopQuickCheck",
			"save": save,
		])
		WsClient.singleton.sendJSON(json)
		DataManager.singleton.unsubscribe(self)
		presentingViewController?.dismiss(animated: true, completion: nil)
	}
	override func viewDidLoad() {
		super.viewDidLoad()
		DataManager.singleton.subscribeData([.QuickCheckInfo], receiver: self)
		self.tableView.register(UITableViewCell.self, forCellReuseIdentifier: QuickCheckViewController.CellReuseIdentifier)
		WsClient.singleton.sendCmd("startQuickCheck")
		timer = Timer.scheduledTimer(timeInterval: 2, target: self, selector: #selector(query), userInfo: nil, repeats: true)
	}
	override func viewDidDisappear(_ animated: Bool) {
		super.viewDidDisappear(animated)
		timer.invalidate()
	}
	func query() {
		WsClient.singleton.sendCmd(DataType.QuickCheckInfo.queryCmd)
	}
}

extension QuickCheckViewController: DataReceiver {
	func onReceivedData(_ json: [String: Any], type: DataType) {
		if type == .QuickCheckInfo {
			let d = json["data"] as! [String: Int]
			self.data = [[], [], [], [], []]
			for (k, v) in d {
				self.data[v].append(k)
			}
			for (idx, ar) in self.data.enumerated() {
				switch idx {
				case 0:
					self.unknowButton.setTitle("未知:\(ar.count)", for: UIControlState())
				case 1:
					self.disableButton.setTitle("无效:\(ar.count)", for: UIControlState())
				case 2:
					self.disableOnButton.setTitle("无效亮:\(ar.count)", for: UIControlState())
				case 3:
					self.enableOffButton.setTitle("有效不亮:\(ar.count)", for: UIControlState())
				case 4:
					self.normalButton.setTitle("正常:\(ar.count)", for: UIControlState())
				default: break
				}
			}
			self.tableView.reloadData()
		}
	}
}

extension QuickCheckViewController: UITableViewDelegate, UITableViewDataSource {
	func numberOfSections(in tableView: UITableView) -> Int {
		return 1
	}
	func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return self.data[self.index].count
	}
	func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCell(withIdentifier: QuickCheckViewController.CellReuseIdentifier)!
		cell.textLabel?.text = self.data[self.index][indexPath.row]
		return cell
	}
}
