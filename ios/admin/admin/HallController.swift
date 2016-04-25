//
//  HallController.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import SwiftyJSON

class HallController: PLViewController {
	@IBOutlet weak var teamtableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!

	@IBOutlet weak var modeImageView: UIImageView!
	@IBOutlet weak var modeLabel: UILabel!
	@IBOutlet weak var playerNumberLabel: UILabel!
	@IBOutlet weak var readyButton: UIButton!
	@IBOutlet weak var startButton: UIButton!
	var refreshControl: UIRefreshControl!
	var teams: [JSON]?

	override func viewDidLoad() {
		super.viewDidLoad()
		refreshControl = UIRefreshControl()
		refreshControl.addTarget(self, action: #selector(HallController.refreshTeamData), forControlEvents: UIControlEvents.ValueChanged)
		teamtableView.addSubview(refreshControl)
	}
	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		DataManager.singleton.subscriptData([.HallData], receiver: self)
	}
	func refreshTeamData() {
		DataManager.singleton.refreshData(.HallData)
	}
	@IBAction func changeMode(sender: UITapGestureRecognizer) {
	}
	@IBAction func callTeam(sender: UIButton) {
	}
	@IBAction func delayTeam(sender: UIButton) {
	}
	@IBAction func addPlayer(sender: UIButton) {
	}
	@IBAction func removePlayer(sender: UIButton) {
	}
	@IBAction func ready(sender: UIButton) {
	}
	@IBAction func start(sender: AnyObject) {
	}
}

extension HallController: DataReceiver {
	func onReceivedData(json: JSON, type: DataType) {
		if type == .HallData {
			teams = json["teams"].array
			teamtableView.reloadData()
			refreshControl.endRefreshing()
		}
	}
}

extension HallController: UITableViewDataSource, UITableViewDelegate {
	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return teams != nil ? teams!.count : 0
	}
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}
	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCellWithIdentifier("HallTableViewCell")! as! HallTableViewCell
		cell.setData(teams![indexPath.row])
		return cell
	}
}
