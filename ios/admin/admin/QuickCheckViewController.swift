//
//  QuickCheckViewController.swift
//  admin
//
//  Created by tassar on 6/23/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

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
	var timer = NSTimer()

	@IBAction func onButtonClicked(sender: UIButton) {
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
		WsClient.singleton.sendCmd("stopQuickCheck")
		DataManager.singleton.unsubscribe(self)
		presentingViewController?.dismissViewControllerAnimated(true, completion: nil)
	}
	override func viewDidLoad() {
		super.viewDidLoad()
		DataManager.singleton.subscribeData([.QuickCheckInfo], receiver: self)
		self.tableView.registerClass(UITableViewCell.self, forCellReuseIdentifier: QuickCheckViewController.CellReuseIdentifier)
		WsClient.singleton.sendCmd("startQuickCheck")
		timer = NSTimer.scheduledTimerWithTimeInterval(2, target: self, selector: #selector(query), userInfo: nil, repeats: true)
	}
	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		timer.invalidate()
	}
	func query() {
		WsClient.singleton.sendCmd(DataType.QuickCheckInfo.queryCmd)
	}
}

extension QuickCheckViewController: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .QuickCheckInfo {
			let d = json["data"] as! [String: Int]
			self.data = [[], [], [], [], []]
			for (k, v) in d {
				self.data[v].append(k)
			}
			self.tableView.reloadData()
		}
	}
}

extension QuickCheckViewController: UITableViewDelegate, UITableViewDataSource {
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}
	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return self.data[self.index].count
	}
	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCellWithIdentifier(QuickCheckViewController.CellReuseIdentifier)!
		cell.textLabel?.text = self.data[self.index][indexPath.row]
		return cell
	}
}
