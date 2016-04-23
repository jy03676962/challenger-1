//
//  ViewController.swift
//  postgame
//
//  Created by tassar on 3/29/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

let gameData: [String: Any] = [
	"total_time": 290,
	"rank": 10,
	"mode": GameMode.Fun.rawValue,
	"players": [
		[
			"id": "xiaoming",
			"level": 2,
			"grade": "A",
			"gold_get": 100,
			"gold_lose": 10,
			"energy": 500,
			"combo": 8
		],
		[
			"id": "xiaoqiang",
			"level": 0,
			"grade": "S",
			"gold_get": 80,
			"gold_lose": 10,
			"energy": 1000,
			"combo": 6
		],
	]
]

class StatViewController: PLViewController {
	var headerBackgroundView: UIImageView?
	var QRImageView: UIImageView?
	var timeLabel: UILabel?
	var rankLebel: UILabel?
	var tableView: UITableView?
	var data: [String: Any]?

	func loadData(data: [String: Any]) {
		self.data = data
		var img: UIImage?
		if (data["mode"] as! Int == GameMode.Fun.rawValue) {
			img = UIImage(named: "FunImage")
		} else {
			img = UIImage(named: "SurivalImage")
		}
		headerBackgroundView?.image = img
		QRImageView?.image = UIImage(named: "QRImage")
	}
}

// MARK: tableview
extension StatViewController: UITableViewDataSource, UITableViewDelegate {
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}

	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		let players = self.data!["players"] as! [Any]
		return players.count
	}

	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		return UITableViewCell()
	}
}

// MARK: UIViewController
extension StatViewController {

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		loadData(gameData)
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		let headerSize = UIImage(named: "FunImage")!.size
		let headerView = UIView()
		headerView.backgroundColor = UIColor.clearColor()
		view.addSubview(headerView)
//		headerView.mas_makeConstraints({ m in
//			m.height.equalTo()(headerSize.height)
//			m.width.equalTo()(headerSize.width)
//			m.centerX.equalTo()(self.view.mas_centerX)
//			m.top.equalTo()(120)
//		})
//		let bgView = UIImageView()
//		headerView.addSubview(bgView)
//		bgView.mas_makeConstraints({ m in
//			m.edges.equalTo()(headerView)
//		})
//		let QRImageView = UIImageView()
//		headerView.addSubview(QRImageView)
//		QRImageView.mas_makeConstraints({ m in
//			m.width.equalTo()(115)
//			m.height.equalTo()(115)
//			m.top.equalTo()(headerView.mas_top).offset()(18)
//			m.right.equalTo()(headerView.mas_right).offset()(-120)
//		})
//		let tableView = UITableView()
//		tableView.delegate = self
//		tableView.dataSource = self
//		view.addSubview(tableView)
//		tableView.mas_makeConstraints({ m in
//			m.top.equalTo()(headerView.mas_bottom).offset()(26)
//			m.width.equalTo()(headerView.mas_width)
//		})
//		tableView.rowHeight = 58
//		self.headerBackgroundView = bgView
//		self.QRImageView = QRImageView
//		self.tableView = tableView
	}
}
